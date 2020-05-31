package carbohydrates

import (
	"strings"

	"github.com/igomonov88/sugar/internal/fdc"
)

type Carbohydrates struct {
	Amount   float64 `json:"amount"`
	UnitName string  `json:"unit_name"`
}

func Retrieve(nutrients []fdc.FoodNutrient) Carbohydrates {
	const (
		carbohydrates             = "carbohydrates"
		carbohydratesByDifference = "carbohydrate, by difference"
	)

	var (
		carbs             float64
		carbsByDifference float64
		unitName          string
	)

	for i := range nutrients {
		name := strings.ToLower(nutrients[i].Nutrient.Name)

		if strings.EqualFold(carbohydrates, name) {
			carbs = nutrients[i].Amount
		}

		if strings.EqualFold(carbohydratesByDifference, name) {
			carbsByDifference = nutrients[i].Amount
		}

		unitName = nutrients[i].Nutrient.UnitName
	}

	if carbsByDifference >= carbs {
		return Carbohydrates{
			Amount:   carbsByDifference,
			UnitName: unitName,
		}
	}

	return Carbohydrates{
		Amount:   carbs,
		UnitName: unitName,
	}
}
