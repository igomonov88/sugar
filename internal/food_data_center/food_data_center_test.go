package food_data_center

import (
	"testing"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestFoodDataCenterClient(t *testing.T) {

	t.Log("Given we starting to test food data center api.")
	{
		cfg := Config{
			ConsumerKey: "07qblbARRNts5zU45YOPyC8NDQc1iuHQgTqLwbTL",
			APIURL:      "https://api.nal.usda.gov/fdc/v1/",
		}
		_, err := Connect(cfg)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to connect to Food Data Center client : %s.", failed, err)
		}
		t.Logf("\t%s\tShould be able to connect to Food Data Center client.", success)

	}
}
