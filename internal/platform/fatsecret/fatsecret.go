package fatsecret

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
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

	// ErrCallExternalAPI is used then we got an http error while making external api call
	ErrCallExternalAPI = errors.New("got an http error on calling external api")
)

// Client makes all operations with fatsecret external api
type Client struct {
	Config
}

// Connect knows how to connect to fatsecret client with provided config
func Connect(cfg Config) (*Client, error) {
	if cfg.APIURL == "" || cfg.ConsumerKey == "" || cfg.ConsumerSecret == "" {
		return nil, ErrInvalidConfig
	}
	return &Client{cfg}, nil
}

// Search making call to external api for the specified query and method, and marshall response to dest value, or returns error
func (c *Client) Search(query, method string, dest interface{}) error {
	switch method {
	case FoodsSearchMethod:
		return foodsSearch(c, query, dest)
	default:
		return ErrMethodNotSupported
	}
}

// foodsSearch is checking for the method we trying to call from the external api and call appropriate search function
func foodsSearch(client *Client, query string, dest interface{}) error {
	resp, err := foodsSearchExternalCall(client.ConsumerKey, client.APIURL, query)
	if err != nil {
		switch e := err.(type) {
		case *url.Error:
			if !e.Temporary() {
				return ErrCallExternalAPI
			}
		}
	}
	defer resp.Body.Close()

	value := reflect.ValueOf(dest)
	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(dest)
}

// foodsSearchExternalCall only knows compose a request query to external api, call it with http.Get method
// and return reposne as it is
func foodsSearchExternalCall(consumerKey, apiURL, query string) (resp *http.Response, err error) {
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
		"oauth_version":          "2.0",
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
	fmt.Printf("REQUEST URL: %v", requestQuery)
	return requestURL
}

// escape the given string using url-escape plus some extras
func escape(s string) string {
	return strings.Replace(strings.Replace(url.QueryEscape(s), "+", "%20", -1), "%7E", "~", -1)
}
