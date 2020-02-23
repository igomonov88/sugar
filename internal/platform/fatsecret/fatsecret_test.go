package fatsecret

import (
	"testing"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestFatSecretClient(t *testing.T) {

	t.Log("Given we starting to test fatSecret api.")
	{
		cfg := Config{
			ConsumerKey:    "883aa16d49bf49f8b53bb47f26a4a982",
			ConsumerSecret: "4b1be7f07d974ecab518dde15206641d",
			APIURL:         "https://platform.fatsecret.com/rest/server.api",
		}
		client, err := Connect(cfg)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to connect to fatSecret client : %s.", failed, err)
		}
		t.Logf("\t%s\tShould be able to connect to fatSecret client.", success)

		type FoodsSearchResponse struct {
			Foods struct {
				Food []struct {
					BrandName       string `json:"brand_name"`
					FoodDescription string `json:"food_description"`
					FoodName        string `json:"food_name"`
				} `json:"food"`
			} `json:"foods"`
		}
		fs := FoodsSearchResponse{}
		if err := client.Search("mars", FoodsSearchMethod, &fs); err != nil {
			t.Fatalf("\t%s\tShould be able to search for given query : %s.", failed, err)
		}
		t.Logf("\t%s\tShould be able to search for given query.", success)

		t.Logf("Given we start testing negative flow.")
		{
			err := client.Search("qwerty", "qwerty", &fs)
			if err != ErrMethodNotSupported {
				t.Logf("\t%s\tUsage of unspecified search method should return: %s error but got: %s", failed, ErrMethodNotSupported, err)
			}
			t.Logf("\t%s\tUsage of unspecified search method should return %s error.", success, ErrMethodNotSupported)

			cfg.ConsumerSecret = ""
			if _, err := Connect(cfg); err != ErrInvalidConfig {
				t.Fatalf("\t%s\tUsage of inaproppriate configuration should return: %s error but got: %s.", failed, ErrInvalidConfig, err)
			}
			t.Logf("\t%s\tUsage of inaproppriate configuration should return: %s error.", success, ErrInvalidConfig)
		}
	}
}