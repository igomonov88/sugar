package fdc

import (
	"testing"

	"github.com/igomonov88/sugar/internal/tests"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestFoodDataCenterClient(t *testing.T) {

	t.Log("Given we starting to test food data center api.")
	{
		t.Log("\tWhen handling creation of the client.")
		{
			ctx := tests.Context()

			cfg := Config{
				ConsumerKey: "07qblbARRNts5zU45YOPyC8NDQc1iuHQgTqLwbTL",
				APIURL:      "https://api.nal.usda.gov/fdc/v1/",
			}

			client, err := Connect(cfg)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to connect to Food Data Center client : %s.", failed, err)
			}
			t.Logf("\t%s\tShould be able to connect to Food Data Center client.", success)

			t.Log("\tWhen handling the request to Food Data Center.")
			{
				resp, err := SearchOutput(ctx, client, "mc donalds cheeseburger")
				if err != nil {
					t.Fatalf("\t%s\tShould be able make search request to Food Data Center: %s.", failed, err)
				}
				t.Logf("\t%s\tShould be able make search request to Food Data Center.", success)
				t.Log("\tWhen try to get FDCID value from response for making another method call.")
				{
					if len(resp.Foods) == 0 {
						t.Skipf("\tResponse result is zero, please try to use another example of product in method above.")
					}
					t.Logf("\t%s\tShould be able to get FDCID from response.", success)

					t.Log("\tWhen try to make Details method call.")
					{
						req := DetailsRequest{FDCID: resp.Foods[0].FDCID}
						resp, err := Details(ctx, client, req)
						if err != nil {
							t.Fatalf("\t%s\tShould be able make details request to Food Data Center: %s.", failed, err)
						}
						if len(resp.FoodNutrients) == 0 {
							t.Logf("\t%s\tShould be able make details request to Food Data Center", success)
						}
						t.Logf("\t%s\tShould be able make details request to Food Data Center.", success)
					}
				}
			}
		}
	}
}
