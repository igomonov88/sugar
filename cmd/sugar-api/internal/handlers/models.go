package handlers

import "github.com/igomonov88/sugar/internal/carbohydrates"

// DetailsResponse represents response on http GET details request
type DetailsResponse struct {
	Description                 string `json:"description"`
	carbohydrates.Carbohydrates `json:"carbohydrates"`
}

// SearchResponse represents the request result of food search request
type SearchResponse struct {
	// Foods is the list of foods found matching the search criteria.
	Products []ProductInfo `json:"products"`
}

// Food represents a information of Food from the search request
type ProductInfo struct {
	// FDCID Unique ID of the food.
	FDCID int `json:"fdc_id"`
	// Description the description of the food
	Description string `json:"description"`
	// BrandOwner brand owner for the food
	BrandOwner string `json:"brand_owner"`
}
