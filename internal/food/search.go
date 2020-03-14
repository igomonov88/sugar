package food

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

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
func Search(ctx context.Context, client *fatsecret.Client, query string) (SearchResult, error) {
	ctx, span := trace.StartSpan(ctx, "internal.food.Search")
	defer span.End()

	var fs fatsecret.FoodsSearch
	var sr SearchResult

	err := client.Search(query, fatsecret.FoodsSearchMethod, &fs)
	if err != nil {
		switch err {
		case fatsecret.ErrMethodNotSupported:
			return SearchResult{}, ErrMethodNotSupported
		default:
			return SearchResult{}, ErrInternal
		}
	}
	if err := sanitizeFoodDescription(&fs); err != nil {
		return SearchResult{}, ErrInternal
	}

	sr.Foods = make([]Food, len(fs.Food))
	for i := 0; i <= len(fs.Food); i++ {
		sr.Foods[i].FoodName = fs.Food[i].FoodName
		carbs, err := strconv.ParseFloat(fs.Food[i].FoodDescription, 64)
		if err != nil {
			continue
		}
		sr.Foods[i].Carbs = carbs
	}

	return sr, nil
}

// sanitizeFoodDescription parse the string from FoodDescription field and delete from it all values except Carbs value
//func sanitizeFoodDescription(fs []fatsecret.Foods) (string error) {
//}

func getCarbs() (float64, error) {
	var cabs float64

	return cabs, nil
}

func getValueFromString(str, searchValue string) string {
	exp := `(\d*\.)?\d+`
	v := fmt.Sprintf("%v: %v", searchValue, exp)
	re, _ := regexp.Compile(v)
	rre, _ := regexp.Compile(`(\d*\.)?\d+`)
	sv := re.FindString(str)
	return rre.FindString(sv)
}

func parseFoodDescription(fs string) (string, error) {
	const (
		expr  = `(\d*\.)?\d+`
		carbs = "Carbs"
		grams = "Grams"
	)
	var carbsValue string

	carbsValue = getValueFromString(fs)

}
