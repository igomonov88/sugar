package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"

	"github.com/igomonov88/sugar/internal/carbohydrates"
	api "github.com/igomonov88/sugar/internal/fdc"
	"github.com/igomonov88/sugar/internal/platform/web"
	"github.com/igomonov88/sugar/internal/storage"
)

// Details returns info about product with given food detail
func (f *Food) Details(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Details")
	defer span.End()

	value, exist := f.cache.Get(strings.TrimSpace(params["fdcID"]))
	if exist {
		return web.Respond(ctx, w, value, http.StatusOK)
	}

	fdcID, err := strconv.Atoi(strings.TrimSpace(params["fdcID"]))
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	d, err := storage.RetrieveDetails(ctx, f.db, fdcID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			d, err := api.Details(ctx, f.apiClient, fdcID)
			if err != nil {
				return web.NewRequestError(err, http.StatusNotFound)
			}

			// Get information about carbohydrates from FDC API response
			carbs := carbohydrates.Retrieve(d.FoodNutrients)
			resp := DetailsResponse{
				Description:   d.Description,
				Carbohydrates: carbs,
				Portions:      make([]Portion, len(d.FoodPortions)),
			}
			for i := range d.FoodPortions {
				resp.Portions[i].GramWeight = d.FoodPortions[i].GramWeight
				resp.Portions[i].Description = d.FoodPortions[i].PortionDescription
			}

			go saveDetails(ctx, f.db, fdcID, carbs, resp.Portions)
			go f.cache.Add(strconv.Itoa(fdcID), resp)

			return web.Respond(ctx, w, resp, http.StatusOK)
		default:
			return web.NewRequestError(err, http.StatusInternalServerError)
		}
	}

	carbs := carbohydrates.Carbohydrates{
		Amount:   d.Amount,
		UnitName: d.UnitName,
	}

	resp := DetailsResponse{
		Description:   d.Description,
		Carbohydrates: carbs,
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func saveDetails(ctx context.Context, db *sqlx.DB, fdcID int, carbs carbohydrates.Carbohydrates,
	portions []Portion) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Details.Storage.SaveDetails")
	defer span.End()

	dbPortions := make([]storage.Portion, len(portions))
	dbCarbs := storage.Carbohydrates{
		FDCID:    fdcID,
		Amount:   carbs.Amount,
		UnitName: carbs.UnitName,
	}
	for i := range portions {
		dbPortions[i].FDCID = fdcID
		dbPortions[i].GramWeight = portions[i].GramWeight
		dbPortions[i].Description = portions[i].Description
	}

	return storage.SaveDetails(ctx, db, fdcID, dbCarbs, dbPortions)
}
