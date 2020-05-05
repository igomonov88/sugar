package handlers

import (
	"context"
	"database/sql"
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

	var sreq api.SearchRequest
	if err := web.Decode(r, &sreq); err != nil {
		return errors.Wrap(err, "")
	}

	cachedResp, exist := searchInCache(ctx, f.cache, sreq)
	if exist {
		return web.Respond(ctx, w, cachedResp, http.StatusOK)
	}

	// Check searching value in storage, if we do not get any result,
	// we will go to external api for the resources.
	storageResp, err := searchInStorage(ctx, f.db, sreq.SearchInput)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// If we cant get the food items from database, we will go to external
			// api for the data.
			apiResp, err := api.Search(ctx, f.apiClient, sreq)
			if err != nil {
				return errors.Wrap(err, "")
			}
			go func() {
				addToCache(ctx, f.cache, sreq.SearchInput, apiResp)
			}()
			return web.Respond(ctx, w, apiResp, http.StatusOK)
		}
	}
	if storageResp != nil {
		return web.Respond(ctx, w, storageResp, http.StatusOK)
	}

	apiResp, err := api.Search(ctx, f.apiClient, sreq)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return web.Respond(ctx, w, apiResp, http.StatusOK)
}

// searchInStorage is searching for the given result in storage.
func searchInStorage(ctx context.Context, db *sqlx.DB, searchInput string) (*api.SearchResponse, error) {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Search.Storage")
	defer span.End()

	var sr api.SearchResponse
	result, err := storage.Search(ctx, db, searchInput)
	if err != nil {
		return nil, err
	}
	switch result {
	case nil:
		return &sr, sql.ErrNoRows
	}

	if len(result) != 0 {
		for i := range result {
			food := api.Food{
				FDCID:       result[i].FDCID,
				Description: result[i].Description,
				BrandOwner:  result[i].BrandOwner,
			}
			sr.Foods = append(sr.Foods, food)
		}
	}

	return &sr, err
}

// searchInCache is searching for the item in cache and return value and bool
// parameter. return false in bool parameter if item not exists in cache.
func searchInCache(ctx context.Context, cache *cache.Cache, r api.SearchRequest) (*api.SearchResponse, bool) {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Search.Cache.Get")
	defer span.End()

	value, exist := cache.Get(r.SearchInput)
	item, ok := value.(api.SearchResponse)
	if !ok {
		return nil, false
	}
	return &item, exist
}

// addToCache add value to cache with the given key.
func addToCache(ctx context.Context, cache *cache.Cache, key string, value *api.SearchResponse) {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Search.Cache.Add")
	defer span.End()
	cache.Add(key, value)
}
