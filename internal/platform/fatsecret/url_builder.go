package fatsecret

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

func buildRequestURL(consumerKey, apiURL, apiMethod string, requestParams map[string]string) string {
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
	message := map[string]string{
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
