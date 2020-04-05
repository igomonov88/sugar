package food_data_center_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	// ErrInvalidConfig is used then some of the config values does not specified.
	ErrInvalidConfig = errors.New("config values does not specified properly")

	// ErrMethodNotSupported is used when we try to call search with api method which is not supported.
	ErrMethodNotSupported = errors.New("api method not supported")

	// ErrFailedToComposeUrl is used when we failed to build correct url.
	ErrFailedToComposeUrl = errors.New("failed to compose url")

	// ErrFromExternalAPI is used when we got not an http 200OK status in response from external api.
	ErrFromExternalAPI = errors.New("error from call to external api")
)

// Client makes all operations with fatSecret external api.
type Client struct {
	cfg Config
}

// Connect knows how to connect to food data center api with provided config.
func Connect(cfg Config) (*Client, error) {
	if cfg.APIURL == "" || cfg.ConsumerKey == "" {
		return nil, ErrInvalidConfig
	}
	return &Client{cfg}, nil
}

// FoodSearch is searching for food given the request parameters.
func FoodSearch(ctx context.Context, client *Client, request FoodSearchRequest) (*FoodSearchResponse, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.Search")
	defer span.End()

	// Compose internal request value.
	req := FoodSearchInternalRequest{
		GeneralSearchInput: request.SearchInput,
		BrandOwner:         request.BrandOwner,
	}

	// Create a context with a timeout of 120 seconds.
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Make http call to external api.
	resp, err := foodSearchHttpRequest(ctx, client, req)
	if err != nil {
		switch err {
		case ErrInvalidConfig:
			return nil, ErrInvalidConfig
		case ErrFailedToComposeUrl:
			return nil, ErrFailedToComposeUrl
		default:
			return nil, errors.Wrap(err, "got error while trying to make food search request")
		}
	}

	// Check response status code. If it's not 200 return error.
	if resp.StatusCode != http.StatusOK {
		errResp := FoodDataCenterErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode error response body.")
		}
		return nil, errors.Wrapf(ErrFromExternalAPI, "status code: %v, error: %v, message: %v, path: %v", resp.StatusCode, errResp.Error, errResp.Message, errResp.Path)
	}

	// Compose internal response value.
	fsi := FoodSearchInternalResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&fsi); err != nil {
		return nil, errors.Wrap(err, "failed to decode response from external api")
	}

	// Compose response value.
	fs := FoodSearchResponse{}
	for i := 0; i < len(fsi.Foods)-1; i++ {
		fs.Foods = append(fs.Foods, fsi.Foods[i])
	}
	fs.TotalHits = fsi.TotalHits

	return &fs, nil
}

// FoodDetails knows how to get information about product from external api.
func FoodDetails(ctx context.Context, client *Client, req FoodDetailsRequest) (*FoodDetailsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.Details")
	defer span.End()

	// Create a context with a timeout of 120 seconds.
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	resp, err := foodDetailsHttpRequest(ctx, client, req.FDCID)
	if err != nil {
		switch err {
		case ErrInvalidConfig:
			return nil, ErrInvalidConfig
		case ErrFailedToComposeUrl:
			return nil, ErrFailedToComposeUrl
		default:
			return nil, errors.Wrap(err, "got error while trying to make food details request")
		}
	}

	// Check response status code. If it's not 200 return error.
	if resp.StatusCode != http.StatusOK {
		errResp := FoodDataCenterErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode error response body.")
		}
		return nil, errors.Wrapf(ErrFromExternalAPI, "status code: %v, error: %v, message: %v, path: %v", resp.StatusCode, errResp.Error, errResp.Message, errResp.Path)
	}

	fdi := FoodDetailsInternalResponse{}
	err = json.NewDecoder(resp.Body).Decode(&fdi)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response to food details response")
	}
	fd := FoodDetailsResponse{}
	for i := 0; i < len(fdi.FoodNutrients); i++ {
		fd.FoodNutrients = append(fd.FoodNutrients, fdi.FoodNutrients[i])
	}
	for i := 0; i < len(fdi.FoodPortions); i++ {
		fd.FoodPortions = append(fd.FoodPortions, fdi.FoodPortions[i])
	}
	fd.Description = fdi.Description
	return &fd, nil
}

// foodSearch make an external call to food data center with given client and given req parameter to get response.
//
// If we got an error during the function execution we just pull it upstears.
func foodSearchHttpRequest(ctx context.Context, c *Client, request FoodSearchInternalRequest) (*http.Response, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.foodSearchHttpRequest")
	defer span.End()

	// Create request url with given client parameters.
	url, err := buildRequestURL(c.cfg.APIURL, c.cfg.ConsumerKey, foodSearchMethod, nil)
	if err != nil {
		return nil, err
	}

	// Marshall incoming request to json.
	b, err := json.Marshal(&request)
	if err != nil {
		return nil, errors.Wrapf(err, "request value %v", request)
	}

	// Creating new http request.
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

// foodDetails make an external call to food data center with given client and fdcID parameter to get response.
//
// If we got an error during the function execution we just pull it upstears.
func foodDetailsHttpRequest(ctx context.Context, c *Client, fdcID int) (*http.Response, error) {
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

// buildRequestURL knows how to build url for food data center api based on given parameters.
func buildRequestURL(apiURL string, consumerKey string, searchMethod string, requestParam interface{}) (string, error) {
	if apiURL == "" || consumerKey == "" {
		return "", ErrInvalidConfig
	}
	switch searchMethod {
	case foodSearchMethod:
		return fmt.Sprintf("%ssearch?api_key=%s", apiURL, consumerKey), nil
	case foodDetailMethod:
		fdcID, ok := requestParam.(int)
		if !ok {
			return "", ErrFailedToComposeUrl
		}
		return fmt.Sprintf("%s%v?api_key=%s", apiURL, fdcID, consumerKey), nil
	default:
		return "", ErrMethodNotSupported
	}
}
