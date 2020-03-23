package food_data_center

// FoodSearchRequest represents the request data which send to food data center api in http.Post request
type FoodSearchRequest struct {
	// Search query (general text) to query food
	GeneralSearchInput string `json:"generalSearchInput"`
	// Specific data types to include in search e.g. ["Survey (FNDDS)", "Foundation", "Branded"]
	IncludeDataTypeList []string `json:"includeDataTypeList"`
	// Ingredients The list of ingredients (as it appears on the product label)
	Ingredients string `json:"ingredients"`
	// Brand owner for the food
	BrandOwner string `json:"brandOwner"`
	// RequireAllWords bool flag, used to include all words from general search input to search query.
	// When true, the search will only return foods that contain all of the words that were entered in the search field.
	// Should be converted from bool to string.
	RequireAllWords string `json:"requireAllWords"`
	// PageNumber the page of results to return. Should be converted to string from int.
	PageNumber string `json:"pageNumber"`
	// SortField is name of the field by which to sort. Possible sorting options: lowercaseDescription.keyword,
	// dataType.keyword, publishedDate, fdcId. E.g. "sortField":"publishedDate"
	SortField string `json:"sortField"`
	// SortDirection the direction of the sorting, either "asc" or "desc".
	SortDirection string `json:"sortDirection"`
}

// FoodSearchResponse represents the request result from food data center api in http.Post request
type FoodSearchResponse struct {
	// FoodSearchCriteria is a copy of the criteria that were used in the search
	FoodSearchCriteria `json:"foodSearchCriteria"`
	// TotalHits the total number of foods found matching the search criteria.
	TotalHits int `json:"totalHits"`
	// CurrentPage the current page of results being returned.
	CurrentPage int `json:"currentPage"`
	// TotalPages represents total number of pages found matching the search criteria.
	TotalPages int `json:"totalPages"`
	// Foods is the list of foods found matching the search criteria.
	Foods []Food `json:"foods"`
}

// FoodSearchCriteria is a copy of the criteria that were used in the search.
type FoodSearchCriteria struct {
	// GeneralSearchInput search query (general text)
	GeneralSearchInput string `json:"generalSearchInput"`
	// PageNumber represents current page of results
	PageNumber int `json:"pageNumber"`
	// RequireAllWords represents does search result require all words rule
	RequireAllWords bool `json:"requireAllWords"`
}

// Food represents a information of Food from the search request
type Food struct {
	// FDCID Unique ID of the food.
	FDCID int `json:"fdcId"`
	// Description the description of the food
	Description string `json:"description"`
	// DataType the type of the food data.
	DataType string `json:"dataType"`
	// PublishedDate date the item was published to FDC.
	PublishedDate string `json:"publishedDate"`
	// BrandOwner brand owner for the food
	BrandOwner string `json:"brandOwner"`
	// Ingredients the list of ingredients (as it appears on the product label).
	Ingredients string `json:"ingredients"`
	// Score is relative score indicating how well the food matches the search criteria.
	Score float64 `json:"score"`
	// ScientificName the scientific name of the food.
	ScientificName string `json:"scientificName, omitempty"`
	// CommonNames contains any other names of the food
	CommonNames []string `json:"commonNames, omitempty"`
	// AdditionalDescriptions contains any additional descriptions of the food
	AdditionalDescriptions string `json:"additionalDescriptions, omitempty"`
	// FoodCode any A unique ID identifying the food within FNDDS.
	FoodCode string `json:"foodCode, omitempty"`
}
