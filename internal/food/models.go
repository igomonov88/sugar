package food

// FoodSearch
type SearchResult struct {
	Foods []Food `json:"foods"`
}

// Food
type Food struct {
	FoodName string  `json:"food_name"`
	Carbs    float64 `json:"carbs"`
	Grams    int64   `json:"grams"`
}
