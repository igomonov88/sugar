package fdc

import "fmt"

// buildRequestURL knows how to build url for food data center api based on
// given parameters.
func buildRequestURL(apiURL string, consumerKey string, searchMethod string,
	requestParam interface{}) (string, error) {

	if apiURL == "" || consumerKey == "" {
		return "", ErrInvalidConfig
	}
	switch searchMethod {
	case foodSearchMethod:
		return fmt.Sprintf("%ssearch?api_key=%s", apiURL, consumerKey), nil
	case foodDetailMethod:
		fdcID, ok := requestParam.(int)
		if !ok {
			return "", ErrFailedToComposeURL
		}
		return fmt.Sprintf("%s%v?api_key=%s", apiURL, fdcID, consumerKey), nil
	default:
		return "", ErrMethodNotSupported
	}
}
