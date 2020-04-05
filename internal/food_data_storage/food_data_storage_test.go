package food_data_storage_test

import (
	"testing"

	"github.com/igomonov88/sugar/internal/food_data_storage"
	"github.com/igomonov88/sugar/internal/tests"
)

func TestFoodDAtaStorage(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()
	ctx := tests.Context()
	t.Log("Given the need to work with Food records.")
	{
		{
			food := food_data_storage.Food{
				FDCID:       1234,
				Description: "bounty",
				BrandOwner:  "mars inc.",
			}
			err := food_data_storage.AddFood(ctx, db, food)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to add food to storage: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to add food to storage.", tests.Success)

			foods, err := food_data_storage.Search(ctx, db, food.Description)
			if err != nil {
				t.Fatalf("\t%s\tShould be able search food in storage: %s", tests.Failed, err)
			}
			t.Logf("FOODS: %v", foods)
		}


	}
}
