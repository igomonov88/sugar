package storage_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/igomonov88/sugar/internal/storage"
	"github.com/igomonov88/sugar/internal/tests"
)

func TestFoodDataStorage(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()
	ctx := tests.Context()
	t.Log("Given the need to work with food records.")
	{
		{
			// Prepare data for test cases
			food := storage.Food{
				FDCID:       1234,
				Description: "bounty",
				BrandOwner:  "mars inc.",
			}
			n1 := storage.Nutrient{
				Name:     "Protein",
				Rank:     205,
				Number:   205,
				UnitName: "Protein",
			}
			n2 := storage.Nutrient{
				Name:     "Carbohydrates",
				Rank:     206,
				Number:   206,
				UnitName: "g",
			}
			n3 := storage.Nutrient{
				Name:     "Fat",
				Rank:     207,
				Number:   207,
				UnitName: "Fat",
			}
			fn1 := storage.FoodNutrient{
				Type:     "Protein",
				Amount:   20,
				Nutrient: n1,
			}
			fn2 := storage.FoodNutrient{
				Type:     "Carbs",
				Amount:   30,
				Nutrient: n2,
			}
			fn3 := storage.FoodNutrient{
				Type:     "Fat",
				Amount:   50,
				Nutrient: n3,
			}
			fns := make([]storage.FoodNutrient, 0, 3)
			fns = append(fns, fn1)
			fns = append(fns, fn2)
			fns = append(fns, fn3)

			// Add Food Item to storage and check that everything is OK
			{
				err := storage.SaveSearchInput(ctx, db, food, "bounty")
				if err != nil {
					t.Fatalf("\t%s\tShould be able to add food to storage: %s", tests.Failed, err)
				}
				t.Logf("\t%s\tShould be able to add food to storage.", tests.Success)
			}

			// Search for Food item in storage and check that everything is OK
			{
				foods, err := storage.List(ctx, db, food.Description)
				if err != nil {
					t.Fatalf("\t%s\tShould be able search food in storage: %s", tests.Failed, err)
				}
				if len(foods) == 0 {
					t.Fatalf("\t%s\tShould be able search food in storage.", tests.Failed)
				}
				t.Logf("\t%s\tShould be able search food in storage.", tests.Success)
			}

			// Add Food details to storage and check that everything is OK
			{
				cs := storage.Carbohydrates{
					FDCID:    food.FDCID,
					Amount:   66.5,
					UnitName: "g",
				}
				var ps []storage.Portion

				err := storage.SaveDetails(ctx, db, food.FDCID, cs, ps)
				if err != nil {
					t.Fatalf("\t%s\tShould be able to add food details to storage: %s", tests.Failed, err)
				}
				t.Logf("\t%s\tShould be able to add food details to storage.", tests.Success)
			}

			// Get Food details from storage, compare then and check that everything is correct
			{
				foodDetails, err := storage.RetrieveDetails(ctx, db, 1234)
				if err != nil {
					t.Fatalf("\t%s\tShould be able to get food details from storage: %s", tests.Failed, err)
				}
				if diff := cmp.Diff(food.Description, foodDetails); diff != "" {
					t.Logf("%s\tShould be able to get food details from storage.", tests.Success)
				}

				if diff := cmp.Diff(foodDetails.FDCID, 1234); diff != "" {
					t.Logf("%s\tShould be able to get the same food details from storage.", tests.Failed)
				}

				if diff := cmp.Diff(foodDetails.Amount, 100.0); diff != "" {
					t.Logf("%s\tShould be able to get the same food details from storage.", tests.Failed)
				}

				if diff := cmp.Diff(foodDetails.UnitName, "g"); diff != "" {
					t.Logf("%s\tShould be able to get the same food details from storage.", tests.Failed)
				}

				t.Logf("%s\tShould be able to get the same food details from storage.", tests.Success)

			}
		}
	}
}
