package provider

import (
	"context"

	"github.com/LuxChanLu/pulumi-supabase/pkg/provider/client"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (p *supabaseProvider) createProjectFunction(ctx context.Context, inputs resource.PropertyMap, projectId string, preview bool, outputs *map[string]interface{}) (string, error) {
	body := client.CreateFunctionJSONRequestBody{}
	if err := propertiesMapToStruct(inputs, &body); err != nil {
		return "", err
	}
	if !preview {
		function, err := p.supabase.CreateFunctionWithResponse(ctx, projectId, body)
		if err := checkForSupabaseError(function.HTTPResponse, err); err != nil {
			return "", err
		}
		if err := structToOutputs(function.JSON201, outputs); err != nil {
			return "", err
		}
		return function.JSON201.Id, nil
	}
	if err := structToOutputs(client.FunctionResponse{Name: body.Name, Slug: body.Slug, Status: client.FunctionResponseStatusACTIVE, VerifyJwt: body.VerifyJwt}, outputs); err != nil {
		return "", err
	}
	return "", nil
}

func (p *supabaseProvider) readProjectFunction(ctx context.Context, projectId, slug string, outputs *map[string]interface{}) (string, error) {
	function, err := p.supabase.GetFunctionWithResponse(ctx, projectId, slug, &client.GetFunctionParams{IncludeBody: pulumi.BoolRef(true)})
	if err != nil {
		return "", err
	}
	if function.JSON200 != nil {
		if err := structToOutputs(function.JSON200, outputs); err != nil {
			return "", err
		}
		return function.JSON200.Id, nil
	}
	return "", nil
}

func (p *supabaseProvider) updateProjectFunction(ctx context.Context, inputs resource.PropertyMap, projectId, slug string, preview bool, outputs *map[string]interface{}) error {
	body := client.UpdateFunctionJSONRequestBody{}
	if err := propertiesMapToStruct(inputs, &body); err != nil {
		return err
	}
	if !preview {
		function, err := p.supabase.UpdateFunctionWithResponse(ctx, projectId, slug, body)
		if err := checkForSupabaseError(function.HTTPResponse, err); err != nil {
			return err
		}
		if err := structToOutputs(function.JSON200, outputs); err != nil {
			return err
		}
		return nil
	}
	if err := structToOutputs(client.FunctionResponse{Name: *body.Name, Slug: slug, VerifyJwt: body.VerifyJwt}, outputs); err != nil {
		return err
	}
	return nil
}

func (p *supabaseProvider) deleteProjectFunction(ctx context.Context, projectId, slug string) error {
	function, err := p.supabase.DeleteFunctionWithResponse(ctx, projectId, slug)
	return checkForSupabaseError(function.HTTPResponse, err)
}
