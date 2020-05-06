package storage

// Food represents a information of Food from the search request
type Food struct {
	ID          int    `db:"id"`
	FDCID       int    `db:"fdc_id"`
	Description string `db:"description"`
	BrandOwner  string `db:"brand_owner"`
}
type FoodDetails struct {
	Description   string
	FoodNutrients []FoodNutrient
}

type FoodNutrient struct {
	Type   string  `db:"type"`
	Amount float64 `db:"amount"`
	Nutrient
}

type Nutrient struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Rank     int    `db:"rank"`
	Number   int    `db:"number"`
	UnitName string `db:"unit_name"`
}