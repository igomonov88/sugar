package food_data_storage

// Food represents a information of Food from the search request
type Food struct {
	FDCID       int
	Description string
	BrandOwner  string
}
type FoodDetails struct {
	Description   string
	FoodNutrients []FoodNutrient
}

type FoodNutrient struct {
	Type     string  `db:"type"`
	ID       int     `db:"id"`
	Amount   float64 `db:"amount"`
	Nutrient Nutrient
}

type Nutrient struct {
	Name     string `db:"name"`
	Rank     int    `db:"rank"`
	UnitName string `db:"unit_name"`
}
