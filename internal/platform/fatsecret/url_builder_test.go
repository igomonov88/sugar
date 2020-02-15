package fatsecret

import (
	"fmt"
	"testing"
)

func TestBuildURL(t *testing.T) {
	params := make(map[string]string)
	params["randomParamKey"] = "randomParamValue"
	fmt.Println((buildRequestURL("randomKEY", "randomAPIUrl", "randomAPIMethod", params)))

}
