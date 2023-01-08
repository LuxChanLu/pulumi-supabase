package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func propertiesMapToStruct(inputs resource.PropertyMap, output interface{}) error {
	jsonData, err := json.Marshal(inputs.MapRepl(nil, func(pv resource.PropertyValue) (interface{}, bool) {
		if pv.IsComputed() {
			return pv.Input().Element.Mappable(), true
		}
		if pv.IsOutput() {
			return pv.OutputValue().Element.Mappable(), true
		}
		if pv.IsSecret() {
			return pv.SecretValue().Element.Mappable(), true
		}
		if pv.IsResourceReference() {
			return pv.ResourceReferenceValue().ID.Mappable(), true
		}
		return nil, false
	}))
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
		body, err := io.ReadAll(res.Body)
		if err != nil {
			defer res.Body.Close()
			return fmt.Errorf("HTTP Error: %s, %s", res.Status, string(body))
		}
		// TODO: Get error content
		return fmt.Errorf("HTTP Error: %s", res.Status)
	}
	return nil
}
