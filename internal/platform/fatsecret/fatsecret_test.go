package fatsecret

import (
	"fmt"
	"testing"
)

type FoodSearchResponseFoods struct {
    Foods map[string]interface{}
}

func TestFatsecretClient(t *testing.T) {
	clientConfig := Config{
		ConsumerKey:    "883aa16d49bf49f8b53bb47f26a4a982",
		ConsumerSecret: "4b1be7f07d974ecab518dde15206641d",
		APIURL:         "https://platform.fatsecret.com/rest/server.api",
	}

	client, _ := Connect(clientConfig)
	resp := FoodSearchResponseFoods{}

	err := client.Search("mars", FoodsSearchMethod, &resp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}
