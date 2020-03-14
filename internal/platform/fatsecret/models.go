package fatsecret

// FoodsSearch is the strange layer on abstraction from the 3rd party system called FatSecret cause it just used to pack
// Foods result.
type FoodsSearch struct {
	Foods `json:"foods"`
}

// Foods is representing the storage of Food response, since we can get multiple results in response from 3rd party
// system FatSecret.
type Foods struct {
	Food []Food `json:"food"`
}

// Food represents minimal information about food that we searching for with answer which get's from the 3rd party API
// called FatSecret. Most interesting data about product that we are searching for stored in FoodDescription field.
// Unfortunately to get information about carbs af the searched product we should parse the string value from
// FoodDescription.
type Food struct {
	BrandName       string `json:"brand_name"`
	FoodDescription string `json:"food_description"`
	FoodName        string `json:"food_name"`
}
