package provider

import (
	"context"
	"strings"

	"github.com/LuxChanLu/pulumi-supabase/pkg/provider/client"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func (p *supabaseProvider) createFunction(ctx context.Context, inputs resource.PropertyMap, projectId string, preview bool, outputs *map[string]interface{}) (string, error) {
	body := &client.CreateFunctionParams{}
	if err := propertiesMapToStruct(inputs, &body); err != nil {
		return "", err
	}
	if !preview {
		function, err := p.supabase.CreateFunctionWithBodyWithResponse(ctx, projectId, body, "application/json", strings.NewReader(((*outputs)["body"]).(string)))
		if err := checkForSupabaseError(function.HTTPResponse, err); err != nil {
			return "", err
		}
		if err := structToOutputs(function.JSON201, outputs); err != nil {
			return "", err
		}
		return function.JSON201.Id, nil
	}
	if err := structToOutputs(client.FunctionResponse{Name: *body.Name, Slug: *body.Slug, Status: client.FunctionResponseStatusACTIVE, VerifyJwt: body.VerifyJwt}, outputs); err != nil {
		return "", err
	}
	return "", nil
}

func (p *supabaseProvider) readFunction(ctx context.Context, projectId, slug string, outputs *map[string]interface{}) (string, error) {
	function, err := p.supabase.GetFunctionWithResponse(ctx, projectId, slug)
	if err != nil {
		return "", err
	}
	if function.JSON200 != nil {
		functionBody, err := p.supabase.GetFunctionBodyWithResponse(ctx, projectId, slug)
		if err != nil {
			return "", err
		}
		if err := structToOutputs(function.JSON200, outputs); err != nil {
			return "", err
		}
		(*outputs)["body"] = functionBody.Body
		return function.JSON200.Id, nil
	}
	return "", nil
}

func (p *supabaseProvider) updateFunction(ctx context.Context, inputs resource.PropertyMap, projectId, slug string, preview bool, outputs *map[string]interface{}) error {
	params := client.UpdateFunctionParams{}
	if err := propertiesMapToStruct(inputs, &params); err != nil {
		return err
	}
	if !preview {
		function, err := p.supabase.UpdateFunctionWithBodyWithResponse(ctx, projectId, slug, &params, "application/json", strings.NewReader(((*outputs)["body"]).(string)))
		if err := checkForSupabaseError(function.HTTPResponse, err); err != nil {
			return err
		}
		if err := structToOutputs(function.JSON200, outputs); err != nil {
			return err
		}
		return nil
	}
	if err := structToOutputs(client.FunctionResponse{Name: *params.Name, Slug: slug, VerifyJwt: params.VerifyJwt}, outputs); err != nil {
		return err
	}
	return nil
}

func (p *supabaseProvider) deleteFunction(ctx context.Context, projectId, slug string) error {
	function, err := p.supabase.DeleteFunctionWithResponse(ctx, projectId, slug)
	return checkForSupabaseError(function.HTTPResponse, err)
}

func (p *supabaseProvider) diffFunction(ctx context.Context, diff *resource.ObjectDiff) ([]string, bool) {
	changes := []string{}
	recreate := false
	for _, key := range diff.ChangedKeys() {
		if key == "name" || key == "body" || key == "verify_jwt" {
			changes = append(changes, string(key))
		}
		if key == "slug" {
			changes = append(changes, string(key))
			recreate = true
		}
	}
	return changes, recreate
}
