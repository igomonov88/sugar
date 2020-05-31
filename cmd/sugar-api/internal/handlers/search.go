package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"

	api "github.com/igomonov88/sugar/internal/fdc"
	"github.com/igomonov88/sugar/internal/platform/web"
	"github.com/igomonov88/sugar/internal/storage"
)

// Search returns result of the food with food ids from given search query.
func (f *Food) Search(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Search")
	defer span.End()

	si := strings.TrimSpace(params["product"])

	foods , err := storage.List(ctx, f.db, si)
	if err != nil {
		return web.NewRequestError(err, http.StatusInternalServerError)
	}

	if len(foods) != 0 {
		resp := SearchResponse{Products:make([]ProductInfo, len(foods))}
		for i := range foods {
			product := ProductInfo{
				FDCID:       foods[i].FDCID,
				Description: foods[i].Description,
				BrandOwner:  foods[i].BrandOwner,
			}
			resp.Products[i] = product
		}

		return web.Respond(ctx, w, &resp, http.StatusOK)
	}

	sr, err := api.SearchOutput(ctx, f.apiClient, si)
	if err != nil {
		return web.NewRequestError(err, http.StatusInternalServerError)
	}

	resp := SearchResponse{Products:make([]ProductInfo, len(sr.Foods))}
	if len (sr.Foods) != 0 {
		for i := range sr.Foods {
			product := ProductInfo{
				FDCID:       sr.Foods[i].FDCID,
				Description: sr.Foods[i].Description,
				BrandOwner:  sr.Foods[i].BrandOwner,
			}
			resp.Products[i] = product
		}

		go saveSearchInput(ctx, f.db, si, &resp)
		return web.Respond(ctx, w, &resp, http.StatusOK)
	}

	return web.Respond(ctx, w, &resp, http.StatusOK)
}

// addToStorage is add value to the storage.
func saveSearchInput(ctx context.Context, db *sqlx.DB, searchInput string, resp *SearchResponse) {
	for i := range resp.Products {
		f := storage.Food{
			FDCID:       resp.Products[i].FDCID,
			Description: resp.Products[i].Description,
			BrandOwner:  resp.Products[i].BrandOwner,
		}
		storage.SaveSearchInput(ctx, db, f, searchInput)
	}
}
