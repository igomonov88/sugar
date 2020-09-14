// Package fdc_api provide the ability to operate with Food Data Central API
// from U.S. Department Of Agriculture.
package fdc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

const (
	foodSearchMethod = "foodSearch"
	foodDetailMethod = "foodDetail"
)

var (
	// ErrInvalidConfig is used then some of config values does not specified.
	ErrInvalidConfig = errors.New("config not specified properly")

	// ErrMethodNotSupported is used when we try to call search with api method
	// which is not supported.
	ErrMethodNotSupported = errors.New("api method not supported")

	// ErrFailedToComposeURL is used when we failed to build correct url.
	ErrFailedToComposeURL = errors.New("failed to compose url")

	// ErrFromExternalAPI is used when we got not an http 200OK status in
	// response from external api.
	ErrFromExternalAPI = errors.New("error from call to external api")
)

// Details knows how to get information about product from external api.
func Details(ctx context.Context, client *Client, fdcID int) (*DetailsInternalResponse, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.Details")
	defer span.End()

	// Create a context with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()

	var fdi DetailsInternalResponse

	resp, err := foodDetailsHTTPRequest(ctx, client, fdcID)
	if err != nil {
		switch err {
		case ErrInvalidConfig:
			return nil, ErrInvalidConfig
		case ErrFailedToComposeURL:
			return nil, ErrFailedToComposeURL
		default:
			return nil, errors.Wrap(err, "got error while details request")
		}
	}

	// Check response status code. If it's not 200 return error.
	if resp.StatusCode != http.StatusOK {
		errResp := FoodDataCentralErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode error body.")
		}
		return nil, errors.Wrapf(ErrFromExternalAPI,
			"status code: %v, error: %v, message: %v, path: %v",
			resp.StatusCode, errResp.Error, errResp.Message, errResp.Path)
	}

	err = json.NewDecoder(resp.Body).Decode(&fdi)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response to food details response")
	}

	return &fdi, nil
}

// foodDetails make an external call to food data central with given client and
// fdcID parameter to get response.
//
// If we got an error during the function execution we just pull it upstears.
func foodDetailsHTTPRequest(ctx context.Context, c *Client, fdcID int) (*http.Response, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.foodDetailsHttpRequest")
	defer span.End()

	// Create request url with given client parameters.
	url, err := buildRequestURL(c.cfg.APIURL, c.cfg.ConsumerKey, foodDetailMethod, fdcID)
	if err != nil {
		return nil, err
	}

	// Create a new request.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "request url: %s", url)
	}

	//Bind the new context into the request.
	req = req.WithContext(ctx)

	// Make the web call and return any error. Do will handle the
	// context level timeout.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed on making request")
	}

	return resp, nil
}
