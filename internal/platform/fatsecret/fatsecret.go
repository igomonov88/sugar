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
	"strconv"
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
	ErrMethodNotSupported = errors.New("api method not supported")

	// ErrCallExternalAPI is used then we got an http error while making external api call
	ErrCallExternalAPI = errors.New("got an http error on calling external api")
)

// Client makes all operations with fatSecret external api
type Client struct {
	cfg Config
}

// Connect knows how to connect to fatSecret client with provided config
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
	resp, err := foodsSearchExternalCall(client.cfg.ConsumerKey, client.cfg.ConsumerSecret, client.cfg.APIURL, query)
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
// and return response as it is
func foodsSearchExternalCall(consumerKey, consumerSecret, apiURL, query string) (resp *http.Response, err error) {
	requestParams := make(map[string]string)
	requestParams["search_expression"] = query
	requestParams["max_results"] = strconv.Itoa(20)
	return http.Get(buildRequestURL(consumerKey, consumerSecret, apiURL, FoodsSearchMethod, requestParams))
}

// buildRequestURL method that compose valid url with provided parameters to use that url in Get call to fatSecret api
func buildRequestURL(consumerKey, consumerSecret, apiURL, apiMethod string, requestParams map[string]string) string {
	var (
		sigQueryStr string
		requestQuery string
		)
	requestTime := fmt.Sprintf("%d", time.Now().Unix())
	requestURL := fmt.Sprintf("%s?", apiURL)
	message := map[string]string{
		"method":                 apiMethod,
		"oauth_consumer_key":     consumerKey,
		"oauth_nonce":            fmt.Sprintf("%d", rand.NewSource(time.Now().UnixNano()).Int63()),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        requestTime,
		"oauth_version":          "1.0",
		"format":                 "json",
	}
	for requestKey, requestValue := range requestParams {
		message[requestKey] = requestValue
	}

	messageKeys := make([]string, 0, len(message))

	for messageKey, _ := range message {
		messageKeys = append(messageKeys, messageKey)
	}
	// sort keys
	sort.Strings(messageKeys)

	// build sorted k/v string for sig
	for i := range messageKeys {
		sigQueryStr = fmt.Sprintf("%s&%s=%s",sigQueryStr,  messageKeys[i], escape(message[messageKeys[i]]))
	}
	// drop initial &
	sigQueryStr = sigQueryStr[1:]

	mac := hmac.New(sha1.New, []byte(consumerSecret+"&"))
	mac.Write([]byte(fmt.Sprintf("GET&%s&%s", url.QueryEscape(apiURL), escape(sigQueryStr))))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// add sig to map
	message["oauth_signature"] = sig
	messageKeys = append(messageKeys, "oauth_signature")

	// re-sort keys after adding sig
	sort.Strings(messageKeys)
	for i := range messageKeys {
		requestQuery += fmt.Sprintf("&%s=%s", messageKeys[i], escape(message[messageKeys[i]]))
	}
	// drop initial &
	requestQuery = requestQuery[1:]

	return fmt.Sprintf("%s%s", requestURL, requestQuery)
}

// escape the given string using url-escape plus some extras
func escape(s string) string {
	return strings.Replace(strings.Replace(url.QueryEscape(s), "+", "%20", -1), "%7E", "~", -1)
}
