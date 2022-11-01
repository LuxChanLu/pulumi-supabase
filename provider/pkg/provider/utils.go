package provider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func propertiesMapToStruct(inputs resource.PropertyMap, output interface{}) error {
	jsonData, err := json.Marshal(inputs.Mappable())
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, output)
}

func structToOutputs(inputs interface{}, output *map[string]interface{}) error {
	// TODO: Remove ID from output
	jsonData, err := json.Marshal(inputs)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, output)
}

func checkForSupabaseError(res *http.Response, err error) error {
	if err != nil {
		return err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		// TODO: Get error content
		return fmt.Errorf("HTTP Error: %s", res.Status)
	}
	return nil
}
