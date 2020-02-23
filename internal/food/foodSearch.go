package food

import (
	"context"
	"errors"

	"go.opencensus.io/trace"

	"github.com/igomonov88/sugar/internal/platform/fatsecret"
)

var (
	// ErrNotFound used when we get an
	ErrInternal = errors.New("failed to get response from external API")

	// ErrMethodNotSupported
	ErrMethodNotSupported = errors.New("search method is not supported from external API")
)

// Search making search request to external api and return the result.
func Search(ctx context.Context, client *fatsecret.Client, query string) (FoodsSearch, error) {
	ctx, span := trace.StartSpan(ctx, "internal.food.Search")
	defer span.End()

	fs := fatsecret.FoodsSearch{}
	err := client.Search(query, fatsecret.FoodsSearchMethod, &fs)
	if err != nil {
		switch err {
		case fatsecret.ErrMethodNotSupported:
			return FoodsSearch{}, ErrMethodNotSupported
		default:
			return FoodsSearch{}, ErrInternal
		}
	}

	// TODO: make transformation the response from the external api to the appropriate response
	return FoodsSearch{}, nil
}

// getCarbs parse the string from FoodDescription field and delete from it all noize to return only carbs value
func getCarbs(fs *FoodsSearch) {
}

// here whe should define regexp function that will be find the value we need from the string