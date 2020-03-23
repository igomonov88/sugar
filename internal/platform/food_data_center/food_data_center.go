package food_data_center

import "errors"

var (
	// ErrInvalidConfig is used then some of the config values does not specified
	ErrInvalidConfig = errors.New("config values does not specified properly")

	// ErrMethodNotSupported is used then we try to call search with api method which is not supported
	ErrMethodNotSupported = errors.New("api method not supported")

	// ErrCallExternalAPI is used then we got an http error while making external api call
	ErrCallExternalAPI = errors.New("got an http error on calling external api")
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

func buildRequestURL() {}
