package provider

import (
	"context"
	"fmt"

	"github.com/LuxChanLu/pulumi-supabase/pkg/provider/client"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func (p *supabaseProvider) createProject(ctx context.Context, inputs resource.PropertyMap, preview bool, outputs *map[string]interface{}) (string, error) {
	body := client.CreateProjectJSONRequestBody{}
	if err := propertiesMapToStruct(inputs, &body); err != nil {
		return "", err
	}
	if !preview {
		project, err := p.supabase.CreateProjectWithResponse(ctx, body)
		if err := checkForSupabaseError(project.HTTPResponse, err); err != nil {
			return "", err
		}
		if err := structToOutputs(project.JSON201, outputs); err != nil {
			return "", err
		}
		decorateProject(project.JSON201, *outputs)
		return project.JSON201.Id, nil
	}
	if err := structToOutputs(client.ProjectResponse{Name: body.Name, Region: string(body.Region), OrganizationId: body.OrganizationId}, outputs); err != nil {
		return "", err
	}
	return "", nil
}

func (p *supabaseProvider) readProject(ctx context.Context, id string, outputs *map[string]interface{}) (string, error) {
	projects, err := p.supabase.GetProjectsWithResponse(ctx)
	if err != nil || projects.JSON200 == nil {
		return "", err
	}
	for _, project := range *projects.JSON200 {
		if project.Id == id {
			outputs := map[string]interface{}{}
			if err := structToOutputs(project, &outputs); err != nil {
				return "", err
			}
			decorateProject(&project, outputs)
			return project.Id, nil
		}
	}
	return "", nil
}

// TODO: From api when available
func decorateProject(project *client.ProjectResponse, outputs map[string]interface{}) {
	outputs["dbUsername"] = "postgres"
	outputs["dbHost"] = fmt.Sprintf("db.%s.supabase.co", project.Id)
	outputs["dbPort"] = 5432
	outputs["dbName"] = "postgres"

	outputs["dbPoolingPort"] = 6543

	outputs["endpoint"] = fmt.Sprintf("https://%s.supabase.co", project.Id)
}
