package provider

import (
	"context"

	"github.com/LuxChanLu/pulumi-supabase/pkg/provider/client"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func (p *supabaseProvider) createOrganization(ctx context.Context, inputs resource.PropertyMap, preview bool, outputs *map[string]interface{}) (string, error) {
	body := client.CreateOrganizationJSONRequestBody{}
	if err := propertiesMapToStruct(inputs, &body); err != nil {
		return "", err
	}
	if !preview {
		organization, err := p.supabase.CreateOrganizationWithResponse(ctx, body)
		if err := checkForSupabaseError(organization.HTTPResponse, err); err != nil {
			return "", err
		}
		if err := structToOutputs(organization.JSON201, outputs); err != nil {
			return "", err
		}
		return organization.JSON201.Id, nil
	}
	if err := structToOutputs(client.OrganizationResponse{Name: body.Name}, outputs); err != nil {
		return "", err
	}
	return "", nil
}

func (p *supabaseProvider) readOrganization(ctx context.Context, id string, outputs *map[string]interface{}) (string, error) {
	organizations, err := p.supabase.GetOrganizationsWithResponse(ctx)
	if err != nil || organizations.JSON200 == nil {
		return "", err
	}
	for _, organization := range *organizations.JSON200 {
		if organization.Id == id {
			if err := structToOutputs(organization, outputs); err != nil {
				return "", err
			}
			return organization.Id, nil
		}
	}
	return "", nil
}
