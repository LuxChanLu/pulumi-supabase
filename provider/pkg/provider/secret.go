package provider

import (
	"context"

	"github.com/LuxChanLu/pulumi-supabase/pkg/provider/client"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func (p *supabaseProvider) createSecret(ctx context.Context, inputs resource.PropertyMap, projectId string, preview bool, outputs *map[string]interface{}) (string, error) {
	body := client.CreateSecretBody{}
	if err := propertiesMapToStruct(inputs, &body); err != nil {
		return "", err
	}
	if !preview {
		secret, err := p.supabase.CreateSecretsWithResponse(ctx, projectId, client.CreateSecretsJSONRequestBody{body})
		if err := checkForSupabaseError(secret.HTTPResponse, err); err != nil {
			return "", err
		}
		if err := structToOutputs(client.SecretResponse{Name: body.Name}, outputs); err != nil {
			return "", err
		}
		return body.Name, nil
	}
	if err := structToOutputs(client.SecretResponse{Name: body.Name}, outputs); err != nil {
		return "", err
	}
	return "", nil
}

func (p *supabaseProvider) readSecret(ctx context.Context, projectId, name string, outputs *map[string]interface{}) (string, error) {
	secrets, err := p.supabase.GetSecretsWithResponse(ctx, projectId)
	if err != nil || secrets.JSON200 == nil {
		return "", err
	}
	for _, secret := range *secrets.JSON200 {
		if secret.Name == name {
			if err := structToOutputs(secret, outputs); err != nil {
				return "", err
			}
			return secret.Name, nil
		}
	}
	return "", nil
}

func (p *supabaseProvider) deleteSecret(ctx context.Context, projectId, name string) error {
	function, err := p.supabase.DeleteSecretsWithResponse(ctx, projectId, client.DeleteSecretsJSONRequestBody{name})
	return checkForSupabaseError(function.HTTPResponse, err)
}
