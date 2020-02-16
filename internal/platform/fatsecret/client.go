package fatsecret

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
)

// Config is the required properties to use fatsecret search api
type Config struct {
	consumerKey    string
	consumerSecret string
	apiURL         string
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, cl Client) error {
	ctx, span := trace.StartSpan(ctx, "platform.Search.StatusCheck")
	defer span.End()

	requestParams := make(map[string]interface{})
	requestParams["brand_type"] = "mars"
	requestURL := buildRequestURL(cl.config.consumerKey, cl.config.apiURL, "food_brands.get", requestParams)
	_, err := http.Get(requestURL)
	return err
}
