package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/igomonov88/sugar/cmd/sugar-api/internal/handlers"
	fdcAPI "github.com/igomonov88/sugar/internal/fdc"
	"github.com/igomonov88/sugar/internal/platform/cache"
	"github.com/igomonov88/sugar/internal/tests"
)

// TestFoodAPI is the entry point for testing user management functions.
func TestFoodAPI(t *testing.T) {
	test := tests.NewIntegration(t)
	defer test.Teardown()

	shutdown := make(chan os.Signal, 1)

	// Creating config for external api client
	fdcConfig := fdcAPI.Config{
		ConsumerKey: "07qblbARRNts5zU45YOPyC8NDQc1iuHQgTqLwbTL",
		APIURL:      "https://api.nal.usda.gov/fdc/v1/",
	}
	// Connect to external api
	fdcClient, err := fdcAPI.Connect(fdcConfig)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to connect to Food Data Center api", tests.Failed)
	}
	cacheCfg := cache.Config{
		DefaultDuration: 1 * time.Millisecond,
		Size:            10,
	}
	cacheClient, err := cache.New(cacheCfg)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to create cache instance", tests.Failed)
	}
	tests := FoodAPITests{
		app: handlers.API("develop", shutdown, test.Log, test.DB, fdcClient, cacheClient),
	}

	t.Run("postSearch200", tests.postSearch200)

}

func (ft *FoodAPITests) postSearch200(t *testing.T) {
	r := httptest.NewRequest("GET", "/v1/search/product=mars", nil)
	w := httptest.NewRecorder()

	ft.app.ServeHTTP(w, r)

	// f is the value we will return.
	var f fdcAPI.FoodSearchResponse

	t.Log("Given the need to search for new food.")
	{
		t.Log("\tTest 0:\tWhen using the declared foodSearch value.")
		if w.Code != http.StatusOK {
			t.Fatalf("\t%s\tShould receive a status code of 20 for the response : %v", tests.Failed, w.Code)
		}
		t.Logf("\t%s\tShould receive a status code of 200 for the response.", tests.Success)

		if err := json.NewDecoder(w.Body).Decode(&f); err != nil {
			t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", tests.Failed, err)
		}
	}
}

type FoodAPITests struct {
	app http.Handler
}
