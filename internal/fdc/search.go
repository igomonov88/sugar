package fdc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Search is searching for food given the request parameters.
func Search(ctx context.Context, client *Client, search string) (*SearchResponse, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.Search")
	defer span.End()

	// Compose internal request value.
	req := SearchInternalRequest{
		GeneralSearchInput: request.SearchInput,
	}

	// Create a context with a timeout of 120 seconds.
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Make http call to external api.
	resp, err := foodSearchHTTPRequest(ctx, client, req)
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

	// Compose response value.
	fs := SearchResponse{}
	for i := 0; i < len(fsi.Foods)-1; i++ {
		fs.Foods = append(fs.Foods, fsi.Foods[i])
	}

	return &fs, nil
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

	// Make the web call and return any error. Do will handle the
	// context level timeout.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed on making request")
	}

	return resp, nil
}
