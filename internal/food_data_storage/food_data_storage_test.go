package food_data_storage_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

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
			{
				err := food_data_storage.AddFood(ctx, db, food, "bounty")
				if err != nil {
					t.Fatalf("\t%s\tShould be able to add food to storage: %s", tests.Failed, err)
				}
				t.Logf("\t%s\tShould be able to add food to storage.", tests.Success)
			}
			{
				foods, err := food_data_storage.Search(ctx, db, food.Description)
				if err != nil {
					t.Fatalf("\t%s\tShould be able search food in storage: %s", tests.Failed, err)
				}
				if len(foods) != 0 {
					t.Logf("\t%s\tShould be able search food in storage.", tests.Success)
				}
			}
			{
				n1 := food_data_storage.Nutrient{
					Name:     "Protein",
					Rank:     205,
					Number:   205,
					UnitName: "Protein",
				}
				n2 := food_data_storage.Nutrient{
					Name:     "Carbs",
					Rank:     206,
					Number:   206,
					UnitName: "Carbs",
				}
				n3 := food_data_storage.Nutrient{
					Name:     "Fat",
					Rank:     207,
					Number:   207,
					UnitName: "Fat",
				}
				fn1 := food_data_storage.FoodNutrient{
					Type:     "Protein",
					Amount:   20,
					Nutrient: n1,
				}
				fn2 := food_data_storage.FoodNutrient{
					Type:     "Carbs",
					Amount:   30,
					Nutrient: n2,
				}
				fn3 := food_data_storage.FoodNutrient{
					Type:     "Fat",
					Amount:   50,
					Nutrient: n3,
				}

				fns := make([]food_data_storage.FoodNutrient, 0, 3)
				fns = append(fns, fn1)
				fns = append(fns, fn2)
				fns = append(fns, fn3)
				fd := food_data_storage.FoodDetails{
					Description:   food.Description,
					FoodNutrients: fns,
				}
				err := food_data_storage.AddDetails(ctx, db, food.FDCID, fd)
				if err != nil {
					t.Fatalf("\t%s\tShould be able to add food details to storage: %s", tests.Failed, err)
				}
				t.Logf("\t%s\tShould be able to add food details to storage.", tests.Success)
			}
			{
				foodDetails, err := food_data_storage.GetDetails(ctx, db, 1234)
				if err != nil {
					t.Fatalf("\t%s\tShould be able to get food details from storage: %s", tests.Failed, err)
				}
				if diff := cmp.Diff(food.Description, foodDetails); diff != "" {
					t.Logf("\t%s\tShould be able to get food details from storage.", tests.Success)
				}
				t.Log(foodDetails.FoodNutrients[0].Nutrient.ID)
			}

		}
	}
}
