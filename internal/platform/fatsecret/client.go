package fatsecret

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
)

const (
	foodsSearch = "foods.search"
)

// Config is the required properties to use fatsecret search api
type Config struct {
	consumerKey    string
	consumerSecret string
	apiURL         string
}

// Client is the struct that only represents api layer of search functionality
type Client struct {
	config Config
}

// NewClient knows hot to connect to external search api based on the configuration
func NewClient(cfg Config) *Client {
	return &Client{
		config: cfg,
	}
}

// Search make search for the query string from external api
func (c *Client) Search(query string) (resp *http.Response, err error) {
	requestParams := make(map[string]string)
	requestParams["search_expression"] = query

	return http.Get(buildRequestURL(c.config.consumerKey, c.config.apiURL, foodsSearch, requestParams))
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, cl Client) error {
	ctx, span := trace.StartSpan(ctx, "platform.Search.StatusCheck")
	defer span.End()

	requestParams := make(map[string]string)
	requestParams["brand_type"] = "mars"
	requestURL := buildRequestURL(cl.config.consumerKey, cl.config.apiURL, "food_brands.get", requestParams)
	_, err := http.Get(requestURL)
	return err
}
