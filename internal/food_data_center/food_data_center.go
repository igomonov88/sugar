package food_data_center

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

const (
	FoodSearchMethod = "foodSearch"
	FoodDetailMethod = "foodDetail"
)

var (
	// ErrInvalidConfig is used then some of the config values does not specified
	ErrInvalidConfig = errors.New("config values does not specified properly")

	// ErrMethodNotSupported is used when we try to call search with api method which is not supported
	ErrMethodNotSupported = errors.New("api method not supported")

	// ErrFailedToComposeUrl is used when we failed to build correct url
	ErrFailedToComposeUrl = errors.New("failed to compose url")

	// ErrFromExternalAPI
	ErrFromExternalAPI = errors.New("status is not 200OK in response from call to external api")
)

// Client makes all operations with fatSecret external api
type Client struct {
	cfg Config
}

// Connect knows how to connect to food data center api with provided config
func Connect(cfg Config) (*Client, error) {
	if cfg.APIURL == "" || cfg.ConsumerKey == "" {
		return nil, ErrInvalidConfig
	}
	return &Client{cfg}, nil
}

func (c *Client) FoodSearch(ctx context.Context, req FoodSearchRequest) (*FoodSearchResponse, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.Search")
	defer span.End()

	resp, err := foodSearch(c, req)
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
	if resp.StatusCode != http.StatusOK {
		return nil, ErrFromExternalAPI
	}
	fs := FoodSearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(&fs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response to food search response")
	}
	return &fs, nil
}

func (c *Client) FoodDetails(ctx context.Context, fdcID int) (*FoodDetailsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "internal.FoodDataCenter.Details")
	defer span.End()
	resp, err := foodDetails(c, fdcID)
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
	if resp.StatusCode != http.StatusOK {
		return nil, ErrFromExternalAPI
	}
	fd := FoodDetailsResponse{}
	err = json.NewDecoder(resp.Body).Decode(&fd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response to food details response")
	}
	return &fd, nil
}

// foodSearch make an external call to food data center with given client and given req parameter to get response.
//
// If we got an error during the function execution we just pull it upstears
func foodSearch(c *Client, req FoodSearchRequest) (*http.Response, error) {
	b, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrapf(err, "request value %v", req)
	}
	buf := bytes.NewBuffer(b)
	url, err := buildRequestURL(c.cfg.APIURL, c.cfg.ConsumerKey, FoodSearchMethod, nil)
	if err != nil {
		return nil, err
	}
	return http.Post(url, "application/json", buf)
}

// foodDetails make an external call to food data center with given client and fdcID parameter to get response.
//
// If we got an error during the function execution we just pull it upstears
func foodDetails(c *Client, fdcID int) (*http.Response, error) {
	url, err := buildRequestURL(c.cfg.APIURL, c.cfg.ConsumerKey, FoodDetailMethod, fdcID)
	if err != nil {
		return nil, err
	}
	return http.Get(url)
}

// buildRequestURL knows how to build url for food data center api based on given parameters
func buildRequestURL(apiURL string, consumerKey string, searchMethod string, requestParam interface{}) (string, error) {
	if apiURL == "" || consumerKey == "" {
		return "", ErrInvalidConfig
	}
	switch searchMethod {
	case FoodSearchMethod:
		return fmt.Sprintf("%s/search?api_key=%s", apiURL, consumerKey), nil
	case FoodDetailMethod:
		fdcID, ok := requestParam.(int)
		if !ok {
			return "", ErrFailedToComposeUrl
		}
		return fmt.Sprintf("%s/%v?api_key=%s", apiURL, fdcID, consumerKey), nil
	default:
		return "", ErrMethodNotSupported
	}
}
