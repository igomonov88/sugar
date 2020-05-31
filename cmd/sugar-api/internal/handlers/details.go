package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

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
			}

			go storage.SaveDetails(ctx, f.db, fdcID, carbs.Amount, carbs.UnitName)
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
