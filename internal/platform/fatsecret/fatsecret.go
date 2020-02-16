package fatsecret

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// FoodsSearchMethod is used for specifying the api request method for make search
const (
	FoodsSearchMethod = "foods.search"
)

var (
	// ErrInvalidConfig is used then some of the config values does not specified
	ErrInvalidConfig = errors.New("config values does not specified properly")

	// ErrMethodNotSupported is used then we try to call search with api method which is not supported
	ErrMethodNotSupported = errors.New("api mathod not supported")
)

// Client makes all operations with fatsecret external api
type Client struct {
	config Config
}

// Connect knows how to connect to fatsecret client with provided config
func Connect(cfg Config) (*Client, error) {
	if cfg.apiURL == "" || cfg.consumerKey == "" || cfg.consumerSecret == "" {
		return nil, ErrInvalidConfig
	}
	return &Client{config: cfg}, nil
}

// Search making call to external api for the specified query and method, and marshall response to dest value, or returns error
func (c *Client) Search(query, method string, dest interface{}) error {
	switch method {
	case FoodsSearchMethod:
		return FoodsSearch(c, query, dest)
	default:
		return ErrMethodNotSupported
	}
}

// FoodsSearch is checking for the method we trying to call from the external api and call appropriate search function
func FoodsSearch(client *Client, query string, dest interface{}) error {
	_, err := fSearch(client.config.consumerKey, client.config.apiURL, query)
	if err != nil {
		switch err {
		// TODO: think about handling an error from http.Get call for external api call
		// also we should check here status code from external API response
		}
	}

	// TODO here we should check the type of dest value and does check it for nil
	// if all checks pass, we should try to decode the response value from the api call to dest value
	return err
}

func fSearch(consumerKey, apiURL, query string) (resp *http.Response, err error) {
	requestParams := make(map[string]interface{})
	requestParams["search_expression"] = query
	requestParams["max_results"] = 20
	return http.Get(buildRequestURL(consumerKey, apiURL, FoodsSearchMethod, requestParams))
}

func buildRequestURL(consumerKey, apiURL, apiMethod string, requestParams map[string]interface{}) string {
	var (
		signatureQuery string
		signatureBase  string
		requestURL     string
		requestQuery   string
	)
	// get the oauth time parameters
	ts := fmt.Sprintf("%d", time.Now().Unix())
	nonce := fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63())

	// build the base message
	message := map[string]interface{}{
		"method":                 apiMethod,
		"oauth_consumer_key":     consumerKey,
		"oauth_nonce":            nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        ts,
		"oauth_version":          "1.0",
		"format":                 "json",
	}

	for parameterName, parameterValue := range requestParams {
		message[parameterName] = parameterValue
	}

	messageKeys := make([]string, 0, len(message))

	for messageKey := range message {
		messageKeys = append(messageKeys, messageKey)
	}

	sort.Strings(messageKeys)

	for i := range messageKeys {
		signatureQuery = fmt.Sprintf("&%s=%s", messageKeys[i], escape(messageKeys[i]))
	}

	signatureQuery = signatureQuery[1:]
	signatureBase = fmt.Sprintf("GET&%s&%s", url.QueryEscape(apiURL), escape(signatureQuery))

	mac := hmac.New(sha1.New, []byte(consumerKey+"&"))
	mac.Write([]byte(signatureBase))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	message["oauth_signature"] = signature
	messageKeys = append(messageKeys, "oauth_signature")

	sort.Strings(messageKeys)
	requestURL = fmt.Sprintf("%s?", apiURL)

	for i := range messageKeys {
		requestQuery = fmt.Sprintf("&%s=%s", messageKeys[i], escape(messageKeys[i]))
	}

	requestQuery = requestQuery[1:]

	requestURL += requestQuery

	return requestURL
}

// escape the given string using url-escape plus some extras
func escape(s string) string {
	return strings.Replace(strings.Replace(url.QueryEscape(s), "+", "%20", -1), "%7E", "~", -1)
}
