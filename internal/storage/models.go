package storage

// Food represents a information of Food from the search request.
type Food struct {
	ID          int    `db:"id"`
	FDCID       int    `db:"fdc_id"`
	Description string `db:"description"`
	BrandOwner  string `db:"brand_owner"`
}

// Details represents the food details with it's nutritions.
type Details struct {
	Description string
	Nutrients   []FoodNutrient
}

// Details represents the food details with it's carbohydrate amount
type DetailsRef struct {
	Description string `db:"description"`
	Carbohydrates
	Portions []Portion
}

// Carbohydrates in specified food with provided fdcID
type Carbohydrates struct {
	ID       int     `db:"id"`
	FDCID    int     `db:"fdc_id"`
	Amount   float64 `db:"amount"`
	UnitName string  `db:"unit_name"`
}

type Portion struct {
	ID          int     `db:"id"`
	FDCID       int     `db:"fdc_id"`
	GramWeight  float64 `db:"gram_weight"`
	Description string  `db:"description"`
}

// FoodNutrient represents nutrients with amunt and type.
type FoodNutrient struct {
	Type   string  `db:"type"`
	Amount float64 `db:"amount"`
	Nutrient
}

// Nutrient is the nutrient such as Fat, Ferum and etc.
type Nutrient struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Rank     int    `db:"rank"`
	Number   int    `db:"number"`
	UnitName string `db:"unit_name"`
}
