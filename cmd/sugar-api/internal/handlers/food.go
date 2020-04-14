package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	api "github.com/igomonov88/sugar/internal/food_data_center_api"
	storage "github.com/igomonov88/sugar/internal/food_data_storage"
	"github.com/igomonov88/sugar/internal/platform/auth"
	"github.com/igomonov88/sugar/internal/platform/cache"
	"github.com/igomonov88/sugar/internal/platform/web"
)

// Food represents the Food Data Center API method handler set.
type Food struct {
	db            *sqlx.DB
	apiClient     *api.Client
	apiCache      *cache.Cache
	authenticator *auth.Authenticator

	// ADD OTHER STATE LIKE THE LOGGER AND CONFIG HERE.
}

// Search returns
func (f *Food) Search(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Search")
	defer span.End()

	var foodSearchReq api.FoodSearchRequest
	if err := web.Decode(r, &foodSearchReq); err != nil {
		return errors.Wrap(err, "")
	}
	// TODO: Add Cache usage here

	// Check searching value in storage, if we do not get any result, we will go to external api for the resources
	result, err := storage.Search(ctx, f.db, foodSearchReq.SearchInput)
	if err != nil {
		return errors.Wrapf(err, "Search input %v", foodSearchReq.SearchInput)
	}
	if len(result) != 0 {
		var foodSearchResp api.FoodSearchResponse
		for i := 0; i < len(result)-1; i++ {
			food := api.Food{
				FDCID:       result[i].FDCID,
				Description: result[i].Description,
				BrandOwner:  result[i].BrandOwner,
			}
			foodSearchResp.Foods = append(foodSearchResp.Foods, food)
		}
		return web.Respond(ctx, w, foodSearchResp, http.StatusOK)
	}

	food, err := api.Search(ctx, f.apiClient, foodSearchReq)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return web.Respond(ctx, w, food, http.StatusOK)
}
