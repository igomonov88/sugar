package fdc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// SearchOutput is returning food with given request parameters.
func SearchOutput(ctx context.Context, client *Client, search string) (*SearchInternalResponse, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.Search")
	defer span.End()

	// Create a context with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Make http call to external api.
	resp, err := foodSearchHTTPRequest(ctx, client, search)
	if err != nil {
		switch err {
		case ErrInvalidConfig:
			return nil, ErrInvalidConfig
		case ErrFailedToComposeURL:
			return nil, ErrFailedToComposeURL
		default:
			return nil, errors.Wrap(err, "error while search request")
		}
	}

	// Check response status code. If it's not 200 return error.
	if resp.StatusCode != http.StatusOK {
		errResp := FoodDataCentralErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode response body.")
		}
		return nil, errors.Wrapf(
			ErrFromExternalAPI,
			"status code: %v, error: %v, message: %v, path: %v",
			resp.StatusCode, errResp.Error, errResp.Message, errResp.Path)
	}

	// Compose internal response value.
	fsi := SearchInternalResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&fsi); err != nil {
		return nil, errors.Wrap(err, "failed to decode response from external api")
	}

	return &fsi, nil
}

// foodSearchHTTPRequest make an external call to food data center with given client and
// given req parameter to get response.
//
// If we got an error during the function execution we just pull it upstears.
func foodSearchHTTPRequest(ctx context.Context, c *Client, search string) (*http.Response, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.foodSearchHttpRequest")
	defer span.End()

	// Create request url with given client parameters.
	url, err := buildRequestURL(c.cfg.APIURL, c.cfg.ConsumerKey, foodSearchMethod, search)
	if err != nil {
		return nil, err
	}

	request := SearchInternalRequest{
		GeneralSearchInput: search,
	}
	// Marshall incoming request to json.
	b, err := json.Marshal(&request)
	if err != nil {
		return nil, errors.Wrapf(err, "request value %v", request)
	}

	//// Creating new http request.
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, errors.Wrapf(err, "request url: %s", url)
	}

	// Bind the new context into the request.
	req = req.WithContext(ctx)

	// Set appropriate Content-Type to request
	req.Header.Add("Content-Type", "application/json")

	// Make the web call and return any error. Do will handle the
	// context level timeout.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed on making request")
	}

	return resp, nil
}
