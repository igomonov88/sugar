package fdc_api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.opencensus.io/trace"
)

// Config is the required properties to use food data center search api
type Config struct {
	ConsumerKey string
	APIURL      string
}

// Client makes all operations with food data central external api.
type Client struct {
	cfg Config
}

// Connect knows how to connect to food data central api with provided config.
func Connect(cfg Config) (*Client, error) {
	if cfg.APIURL == "" || cfg.ConsumerKey == "" {
		return nil, ErrInvalidConfig
	}
	return &Client{cfg}, nil
}

// StatusCheck returns nil if it can successfully talk to the food data center api. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, cfg Config) error {
	ctx, span := trace.StartSpan(ctx, "platform.Search.StatusCheck")
	defer span.End()
	fs := SearchRequest{
		SearchInput: "cheese",
	}
	b, err := json.Marshal(&fs)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)
	url, err := buildRequestURL(cfg.APIURL, cfg.ConsumerKey, foodSearchMethod, nil)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("food data center api respond with status not OK")
	}
	return nil
}
