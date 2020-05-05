package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	api "github.com/igomonov88/sugar/internal/fdc_api"
	"github.com/igomonov88/sugar/internal/platform/web"
	storage "github.com/igomonov88/sugar/internal/storage"
)

// Details returns info about product with given food detail
func (f *Food) Details(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Details")
	defer span.End()

	var req api.DetailsRequest
	if err := web.Decode(r, &req); err != nil {
		return errors.Wrap(err, "")
	}

	details, err := storageGetDetails(ctx, f.db, req.FDCID)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if details.FoodNutrients != nil || len(details.FoodNutrients) > 0 {
		return web.Respond(ctx, w, details, http.StatusOK)
	}

	details, err = api.Details(ctx, f.apiClient, req)
	if err != nil {
		return errors.Wrap(err, "")
	}

	go storageAddDetails(ctx, f.db, req.FDCID, details)
	return web.Respond(ctx, w, details, http.StatusOK)
}

// storageAddDetails add values from external api, and put them to storage.
func storageAddDetails(ctx context.Context, db *sqlx.DB, fdcID int, details *api.DetailsResponse) {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Details.Storage.")
	defer span.End()

	var foodDetails storage.FoodDetails
	foodDetails.Description = details.Description

	for i := range details.FoodNutrients {
			n := storage.Nutrient{
				Name:     details.FoodNutrients[i].Nutrient.Name,
				Rank:     details.FoodNutrients[i].Nutrient.Rank,
				UnitName: details.FoodNutrients[i].Nutrient.UnitName,
			}
			fn := storage.FoodNutrient{
				Type:     details.FoodNutrients[i].Type,
				Amount:   details.FoodNutrients[i].Amount,
				Nutrient: n,
			}
			foodDetails.FoodNutrients = append(foodDetails.FoodNutrients, fn)
			storage.AddDetails(ctx, db, fdcID, foodDetails)
	}
}

// storageGetDetails returns details information in api.DetailsResponse format
// or returns error.
func storageGetDetails(ctx context.Context, db *sqlx.DB, fdcID int) (*api.DetailsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Details.Storage")
	defer span.End()

	var resp api.DetailsResponse

	details, err := storage.GetDetails(ctx, db, fdcID)
	if err != nil {
		return nil, err
	}
	resp.Description = details.Description

	for i := range details.FoodNutrients {
		n := api.Nutrient{
			Name:     details.FoodNutrients[i].Name,
			Rank:     details.FoodNutrients[i].Rank,
			UnitName: details.FoodNutrients[i].UnitName,
		}
		fn := api.FoodNutrient{
			Type:     details.FoodNutrients[i].Type,
			ID:       details.FoodNutrients[i].ID,
			Nutrient: n,
			Amount:   details.FoodNutrients[i].Amount,
		}
		resp.FoodNutrients = append(resp.FoodNutrients, fn)
	}
	return &resp, nil
}
