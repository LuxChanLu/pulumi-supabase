package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/LuxChanLu/pulumi-supabase/pkg/provider/client"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	pbempty "github.com/golang/protobuf/ptypes/empty"
)

//go:generate oapi-codegen --package=client -generate=client,types -o ./client/supabase.gen.go https://api.supabase.com/api/v1-json

const configServerKey = "server"
const configTokenKey = "token"

type supabaseProvider struct {
	host     *provider.HostClient
	name     string
	version  string
	schema   []byte
	supabase *client.ClientWithResponses
}

func makeProvider(host *provider.HostClient, name, version string, pulumiSchema []byte) (pulumirpc.ResourceProviderServer, error) {
	// Return the new provider
	return &supabaseProvider{
		host:    host,
		name:    name,
		version: version,
		schema:  pulumiSchema,
	}, nil
}

// Attach sends the engine address to an already running plugin.
func (p *supabaseProvider) Attach(context context.Context, req *pulumirpc.PluginAttach) (*emptypb.Empty, error) {
	host, err := provider.NewHostClient(req.GetAddress())
	if err != nil {
		return nil, err
	}
	p.host = host
	return &pbempty.Empty{}, nil
}

// Call dynamically executes a method in the provider associated with a component resource.
func (p *supabaseProvider) Call(ctx context.Context, req *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	return nil, status.Error(codes.Unimplemented, "call is not yet implemented")
}

// Configure configures the resource provider with "globals" that control its behavior.
func (p *supabaseProvider) Configure(_ context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	server, _ := os.LookupEnv("SUPABASE_SERVER")
	token, _ := os.LookupEnv("SUPABASE_TOKEN")
	for key, value := range req.GetVariables() {
		if key == "supabase:config:"+configServerKey {
			server = value
		}
		if key == "supabase:config:"+configTokenKey {
			token = value
		}
	}
	if server == "" {
		server = "https://api.supabase.com/"
	}
	supabase, err := client.NewClientWithResponses(server, client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}))
	if err != nil {
		return nil, err
	}
	p.supabase = supabase
	return &pulumirpc.ConfigureResponse{
		AcceptSecrets:   true,
		AcceptResources: true,
		AcceptOutputs:   true,
		SupportsPreview: true,
	}, nil
}

// CheckConfig validates the configuration for this provider.
func (p *supabaseProvider) CheckConfig(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	hasToken := false
	failures := []*pulumirpc.CheckFailure{}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	for key, value := range news {
		if key == configServerKey {
			_, err := url.Parse(value.String())
			if err != nil {
				failures = append(failures, &pulumirpc.CheckFailure{Property: configServerKey, Reason: fmt.Sprintf("error parsing supabase url: %s", err.Error())})
			}
		}
		if key == configTokenKey {
			hasToken = true
		}
	}
	if !hasToken {
		failures = append(failures, &pulumirpc.CheckFailure{Property: configTokenKey, Reason: "missing supabase token"})
	}
	if len(failures) > 0 {
		return &pulumirpc.CheckResponse{Inputs: req.GetNews(), Failures: failures}, nil
	}
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// DiffConfig diffs the configuration for this provider.
func (p *supabaseProvider) DiffConfig(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	d := olds.Diff(news)
	changes := pulumirpc.DiffResponse_DIFF_NONE

	if d != nil && (d.Changed(configServerKey) || d.Changed(configTokenKey)) {
		changes = pulumirpc.DiffResponse_DIFF_SOME
	}

	return &pulumirpc.DiffResponse{
		Changes:  changes,
		Replaces: []string{configServerKey, configTokenKey},
	}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (p *supabaseProvider) Invoke(ctx context.Context, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	tok := req.GetTok()
	inputs, err := plugin.UnmarshalProperties(req.GetArgs(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	switch tok {
	case "supabase:index:GetTypeScript":
		schema, err := p.supabase.GetTypescriptTypesWithResponse(ctx, inputs["projectId"].String(), &client.GetTypescriptTypesParams{IncludedSchemas: pulumi.StringRef(inputs["includedSchemas"].String())})
		if err != nil {
			return nil, err
		}
		if schema.JSON200 != nil {
			outputs := map[string]interface{}{}

			if err := structToOutputs(schema.JSON200, &outputs); err != nil {
				return nil, err
			}

			outputProperties, err := plugin.MarshalProperties(resource.NewPropertyMapFromMap(outputs), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
			if err != nil {
				return nil, err
			}

			return &pulumirpc.InvokeResponse{Return: outputProperties}, nil
		}
		return &pulumirpc.InvokeResponse{Failures: []*pulumirpc.CheckFailure{{Property: "types", Reason: "Types not found"}}}, nil
	}
	return nil, fmt.Errorf("unknown Invoke token '%s'", tok)
}

// StreamInvoke dynamically executes a built-in function in the provider. The result is streamed
// back as a series of messages.
func (p *supabaseProvider) StreamInvoke(req *pulumirpc.InvokeRequest, server pulumirpc.ResourceProvider_StreamInvokeServer) error {
	tok := req.GetTok()
	return fmt.Errorf("unknown StreamInvoke token '%s'", tok)
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (p *supabaseProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{Inputs: req.News, Failures: nil}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (p *supabaseProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	d := olds.Diff(news)
	changes := pulumirpc.DiffResponse_DIFF_NONE

	if d.AnyChanges() {
		changes = pulumirpc.DiffResponse_DIFF_SOME
	}

	return &pulumirpc.DiffResponse{
		Changes:  changes,
		Replaces: []string{},
	}, nil
}

// Construct creates a new component resource.
func (p *supabaseProvider) Construct(ctx context.Context, req *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return nil, status.Error(codes.Unimplemented, "construct is not implemented")

}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
func (p *supabaseProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	id := ""

	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	outputs := map[string]interface{}{}

	switch urn.Type() {
	case "supabase:index:Organization":
		id, err = p.createOrganization(ctx, inputs, req.GetPreview(), &outputs)
		if err != nil {
			return nil, err
		}
	case "supabase:index:Project":
		id, err = p.createProject(ctx, inputs, req.GetPreview(), &outputs)
		if err != nil {
			return nil, err
		}
	case "supabase:index:Function":
		id, err = p.createFunction(ctx, inputs, inputs["projectId"].String(), req.GetPreview(), &outputs)
		if err != nil {
			return nil, err
		}
	case "supabase:index:Secret":
		id, err = p.createSecret(ctx, inputs, inputs["projectId"].String(), req.GetPreview(), &outputs)
		if err != nil {
			return nil, err
		}
	default:
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("%s does not exist", urn.Type()))
	}

	outputProperties, err := plugin.MarshalProperties(resource.NewPropertyMapFromMap(outputs), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	return &pulumirpc.CreateResponse{Id: id, Properties: outputProperties}, nil
}

// Read the current live state associated with a resource.
func (p *supabaseProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	var err error
	urn := resource.URN(req.GetUrn())
	id := ""

	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true, KeepOutputValues: false})
	if err != nil {
		return nil, err
	}

	outputs := map[string]interface{}{}

	switch urn.Type() {
	case "supabase:index:Organization":
		id, err = p.readOrganization(ctx, req.Id, &outputs)
		if err != nil {
			return nil, err
		}
	case "supabase:index:Project":
		id, err = p.readProject(ctx, req.Id, &outputs)
		if err != nil {
			return nil, err
		}
	case "supabase:index:Function":
		id, err = p.readFunction(ctx, inputs["projectId"].String(), inputs["slug"].String(), &outputs)
		if err != nil {
			return nil, err
		}
	case "supabase:index:Secret":
		id, err = p.readSecret(ctx, inputs["projectId"].String(), inputs["name"].String(), &outputs)
		if err != nil {
			return nil, err
		}
	default:
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("%s does not exist", urn.Type()))
	}
	outputProperties, err := plugin.MarshalProperties(resource.NewPropertyMapFromMap(outputs), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	return &pulumirpc.ReadResponse{Id: id, Properties: outputProperties}, nil
}

// Update updates an existing resource with new values.
func (p *supabaseProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())

	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	outputs := map[string]interface{}{}

	switch urn.Type() {
	case "supabase:index:Organization":
		return nil, status.Error(codes.Unimplemented, "no update available for organization (update manually and refresh)")
	case "supabase:index:Project":
		return nil, status.Error(codes.Unimplemented, "no update available for organization project (update manually and refresh)")
	case "supabase:index:Function":
		if err := p.updateFunction(ctx, news, olds["projectId"].String(), olds["slug"].String(), req.GetPreview(), &outputs); err != nil {
			return nil, err
		}
	case "supabase:index:Secret":
		return nil, status.Error(codes.Unimplemented, "no update available for project secret (update manually and refresh)")
	default:
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("%s does not exist", urn.Type()))
	}
	outputProperties, err := plugin.MarshalProperties(resource.NewPropertyMapFromMap(outputs), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	return &pulumirpc.UpdateResponse{Properties: outputProperties}, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (p *supabaseProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())

	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	switch urn.Type() {
	case "supabase:index:Organization":
		return nil, status.Error(codes.Unimplemented, "no delete available for organization (delete manually and refresh)")
	case "supabase:index:Project":
		return nil, status.Error(codes.Unimplemented, "no delete available for organization project (delete manually and refresh)")
	case "supabase:index:Function":
		return &pbempty.Empty{}, p.deleteFunction(ctx, inputs["projectId"].String(), inputs["slug"].String())
	case "supabase:index:Secret":
		return &pbempty.Empty{}, p.deleteSecret(ctx, inputs["projectId"].String(), inputs["name"].String())
	default:
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("%s does not exist", urn.Type()))
	}
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (p *supabaseProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{Version: p.version}, nil
}

// GetSchema returns the JSON-serialized schema for the provider.
func (p *supabaseProvider) GetSchema(ctx context.Context, req *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	if v := req.GetVersion(); v != 0 {
		return nil, fmt.Errorf("unsupported schema version %d", v)
	}
	return &pulumirpc.GetSchemaResponse{Schema: string(p.schema)}, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (p *supabaseProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	return &pbempty.Empty{}, nil
}
