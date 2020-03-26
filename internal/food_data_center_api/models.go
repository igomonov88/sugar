package food_data_center_api

// FoodSearchRequest represents a request query to our api
type FoodSearchRequest struct {
	// SearchInput is the search string for given food
	SearchInput string `json:"search_input"`
	// Brand owner for the food
	BrandOwner string `json:"brand_owner"`
}

// FoodSearchResponse represents the request result of food search request
type FoodSearchResponse struct {
	// Foods is the list of foods found matching the search criteria.
	Foods []Food `json:"foods"`
	// TotalHits the total number of foods found matching the search criteria.
	TotalHits int `json:"total_hits"`
}

// FoodDetailsRequest represents a request query to our api
type FoodDetailsRequest struct {
	// FDCID Unique ID of the food.
	FDCID int `json:"fdc_id"`
	// SearchInput word which was used before get FDCID
	SearchInput string `json:"search_input, omitempty"`
}

// FoodDetailsResponse
type FoodDetailsResponse struct {
	// Description is product description
	Description string `json:"description"`
	// FoodNutrients represents nutrients of given product
	FoodNutrients []FoodNutrients `json:"food_nutrients"`
	// FoodPortions represents portion of given product
	FoodPortions []FoodPortion `json:"food_portions"`
}

// FoodSearchRequest represents the request data which send to food data center api in http.Post request
type FoodSearchInternalRequest struct {
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

// FoodSearchInternalResponse represents the request result from food data center api in http.Post request
type FoodSearchInternalResponse struct {
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
	CommonNames string `json:"commonNames, omitempty"`
	// AdditionalDescriptions contains any additional descriptions of the food
	AdditionalDescriptions string `json:"additionalDescriptions, omitempty"`
	// FoodCode any A unique ID identifying the food within FNDDS.
	FoodCode string `json:"foodCode, omitempty"`
}

type FoodDetailsInternalRequest struct {
	// FDCID Unique ID of the food.
	FDCID int `json:"fdcId"`
}

type FoodDetailsInternalResponse struct {
	FoodClass     string          `json:"foodClass"`
	Description   string          `json:"description"`
	FoodNutrients []FoodNutrients `json:"foodNutrients"`
	FoodPortions  []FoodPortion   `json:"foodPortions"`
}

type FoodNutrients struct {
	Type     string   `json:"type"`
	ID       int      `json:"id"`
	Nutrient Nutrient `json:"nutrient"`
	Amount   float64  `json:"amount"`
}

type Nutrient struct {
	ID       int    `json:"id"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	Rank     int    `json:"rank"`
	UnitName string `json:"unitName"`
}

type FoodPortion struct {
	ID       int    `json:"id"`
	Modifier string `json:"modifier"`
	// GramWeight represents total gram amount in portion
	GramWeight float64 `json:"gramWeight"`
	// PortionDescription represents information about portion 1 bra/ 1 snack etc.
	PortionDescription string `json:"portionDescription"`
	// SequenceNumber represents sequence number of the element. Can be useful for iteration, but REMEMBER that this
	// parameter starts from 1 not from 0
	SequenceNumber int `json:"sequenceNumber"`
}
