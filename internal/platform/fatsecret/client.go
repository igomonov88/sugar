package fatsecret

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
)

// Config is the required properties to use fatsecret search api
type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	APIURL         string
}

// StatusCheck returns nil if it can successfully talk to the fatsecret api. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, c *Client) error {
	ctx, span := trace.StartSpan(ctx, "platform.Search.StatusCheck")
	defer span.End()

	requestParams := make(map[string]string)
	requestParams["brand_type"] = "mars"
	reqURL := buildRequestURL(c.cfg.ConsumerKey, c.cfg.ConsumerSecret, c.cfg.APIURL, FoodsSearchMethod, requestParams)
	_, err := http.Get(reqURL)
	return err
}
