package fatsecret

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type FoodSearchResponseFoods struct {
	Foods map[string]interface{}
}

func TestFatSecretClient(t *testing.T) {
	 cfg := Config{
		ConsumerKey:    "883aa16d49bf49f8b53bb47f26a4a982",
		ConsumerSecret: "4b1be7f07d974ecab518dde15206641d",
		APIURL:         "https://platform.fatsecret.com/rest/server.api",
	}


	requestParams := make(map[string]string)
	requestParams["search_expression"] = "mars"
	reqURL := buildRequestURL(cfg.ConsumerKey, cfg.ConsumerSecret, cfg.APIURL, FoodsSearchMethod, requestParams)
	response, err := http.Get(reqURL)
	if err != nil {
		t.Log(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
