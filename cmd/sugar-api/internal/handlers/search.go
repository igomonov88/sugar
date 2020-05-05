package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	api "github.com/igomonov88/sugar/internal/fdc_api"
	"github.com/igomonov88/sugar/internal/platform/cache"
	"github.com/igomonov88/sugar/internal/platform/web"
	"github.com/igomonov88/sugar/internal/storage"
)

// Search returns result of the food with food ids from given search query.
func (f *Food) Search(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Search")
	defer span.End()

	var req api.SearchRequest
	if err := web.Decode(r, &req); err != nil {
		return errors.Wrap(err, "")
	}

	// Check searching value in cache, if value exist return it from the cache
	cachedResp, exist := searchInCache(ctx, f.cache, req)
	if exist {
		return web.Respond(ctx, w, cachedResp, http.StatusOK)
	}

	// Check searching value in storage, if we do not get any result,
	// we will go to external api for the resources.
	storageResp, err := searchInStorage(ctx, f.db, req.SearchInput)
	if err != nil {
		return errors.Wrap(err, "")
	}

	// If we get values from storage, we can return it as a response and also
	// put result into the cache.
	if len(storageResp.Foods) != 0 {
		go addToCache(ctx, f.cache, req.SearchInput, storageResp)
		return web.Respond(ctx, w, storageResp, http.StatusOK)
	}

	// If we does not have searched values in cache or in database we need to
	// get information in external api.
	apiResp, err := api.Search(ctx, f.apiClient, req)
	if err != nil {
		return errors.Wrap(err, "")
	}

	// If we successfully get the result from the external api we add this
	// result into our cache and storage, and return the result to the client.
	go addToStorage(ctx, f.db, req.SearchInput, apiResp)
	go addToCache(ctx, f.cache, req.SearchInput, apiResp)

	return web.Respond(ctx, w, apiResp, http.StatusOK)
}

// searchInStorage is searching for the given result in database, put result to
// appropriate response structure and return it, or return error.
func searchInStorage(ctx context.Context, db *sqlx.DB, searchInput string) (*api.SearchResponse, error) {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Search.Storage.Search")
	defer span.End()

	var sr api.SearchResponse
	foods, err := storage.Search(ctx, db, searchInput)
	if err != nil {
		return nil, err
	}

	if len(foods) != 0 {
		for i := range foods {
			food := api.Food{
				FDCID:       foods[i].FDCID,
				Description: foods[i].Description,
				BrandOwner:  foods[i].BrandOwner,
			}
			sr.Foods = append(sr.Foods, food)
		}
	}

	return &sr, err
}

// addToStorage is add value to the storage.
func addToStorage(ctx context.Context, db *sqlx.DB, searchInput string, resp *api.SearchResponse) {
	for i := range resp.Foods {
		f := storage.Food{
			FDCID:       resp.Foods[i].FDCID,
			Description: resp.Foods[i].Description,
			BrandOwner:  resp.Foods[i].BrandOwner,
		}
		storage.AddFood(ctx, db, f, searchInput)
	}
}

// searchInCache is searching for the item in cache and return value and bool
// parameter. return false in bool parameter if item not exists in cache.
func searchInCache(ctx context.Context, cache *cache.Cache, r api.SearchRequest) (*api.SearchResponse, bool) {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Search.Cache.Get")
	defer span.End()

	value, exist := cache.Get(r.SearchInput)
	item, ok := value.(*api.SearchResponse)
	if !ok {
		return nil, false
	}
	return item, exist
}

// addToCache add value to cache with the given key.
func addToCache(ctx context.Context, cache *cache.Cache, searchInput string, value *api.SearchResponse) {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Search.Cache.Add")
	defer span.End()
	cache.Add(searchInput, value)
}
