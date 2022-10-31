package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/LuxChanLu/pulumi-supabase/pkg/provider/client"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
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

type supabaseProvider struct {
	host     *provider.HostClient
	name     string
	version  string
	schema   []byte
	supabase *client.ClientWithResponses
	ctx      context.Context
	cancel   context.CancelFunc
}

func makeProvider(host *provider.HostClient, name, version string, pulumiSchema []byte) (pulumirpc.ResourceProviderServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	// Return the new provider
	return &supabaseProvider{
		host:    host,
		name:    name,
		version: version,
		schema:  pulumiSchema,
		ctx:     ctx,
		cancel:  cancel,
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

// Construct creates a new component resource.
func (p *supabaseProvider) Construct(ctx context.Context, req *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return nil, status.Error(codes.Unimplemented, "construct is not yet implemented")
}

// Configure configures the resource provider with "globals" that control its behavior.
func (p *supabaseProvider) Configure(_ context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	server, _ := os.LookupEnv("SUPABASE_SERVER")
	token, _ := os.LookupEnv("SUPABASE_TOKEN")
	for key, value := range req.GetVariables() {
		if key == "server" {
			server = value
		}
		if key == "token" {
			token = value
		}
	}
	if server == "" {
		server = "https://api.supabase.com/v1/"
	}
	supabase, err := client.NewClientWithResponses(server, client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}))
	if err != nil {
		return nil, err
	}
	p.supabase = supabase
	return &pulumirpc.ConfigureResponse{}, nil
}

// CheckConfig validates the configuration for this provider.
func (p *supabaseProvider) CheckConfig(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	hasToken := false
	failures := []*pulumirpc.CheckFailure{}
	for key, value := range req.GetNews().GetFields() {
		if key == "server" {
			_, err := url.Parse(value.String())
			if err != nil {
				failures = append(failures, &pulumirpc.CheckFailure{Property: "server", Reason: fmt.Sprintf("error parsing supabase url: %s", err.Error())})
			}
		}
		if key == "token" {
			hasToken = true
		}
	}
	if !hasToken {
		failures = append(failures, &pulumirpc.CheckFailure{Property: "token", Reason: "missing supabase token"})
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

	if d.Changed("server") || d.Changed("token") {
		changes = pulumirpc.DiffResponse_DIFF_SOME
	}

	return &pulumirpc.DiffResponse{
		Changes:  changes,
		Replaces: []string{"server", "token"},
	}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (p *supabaseProvider) Invoke(_ context.Context, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	tok := req.GetTok()
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
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Create(%s)", p.name, urn)
	logging.V(9).Infof("%s executing", label)

	return &pulumirpc.CheckResponse{Inputs: req.News, Failures: nil}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (p *supabaseProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Diff(%s)", p.name, urn)
	logging.V(9).Infof("%s executing", label)

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

	// Replace the below condition with logic specific to your provider
	if d.Changed("length") {
		changes = pulumirpc.DiffResponse_DIFF_SOME
	}

	return &pulumirpc.DiffResponse{
		Changes:  changes,
		Replaces: []string{"length"},
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
func (p *supabaseProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())

	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	response := &pulumirpc.CreateResponse{}
	outputs := map[string]interface{}{}

	switch urn.Type() {
	case "supabase:index:Organization":
		body := client.CreateOrganizationJSONRequestBody{}
		if err := propertiesMapToStruct(inputs, &body); err != nil {
			return nil, err
		}
		organization, err := p.supabase.CreateOrganizationWithResponse(ctx, body)
		if err := checkForSupabaseError(organization.HTTPResponse, err); err != nil {
			return nil, err
		}
		response.Id = organization.JSON201.Id
		if err := structToOutputs(organization.JSON201, &outputs); err != nil {
			return nil, err
		}
	default:
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("%s does not exist", urn.Type()))
	}

	outputProperties, err := plugin.MarshalProperties(resource.NewPropertyMapFromMap(outputs), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	response.Properties = outputProperties
	return response, nil
}

// Read the current live state associated with a resource.
func (p *supabaseProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Read(%s)", p.name, urn)
	logging.V(9).Infof("%s executing", label)
	msg := fmt.Sprintf("Read is not yet implemented for %s", urn.Type())
	return nil, status.Error(codes.Unimplemented, msg)
}

// Update updates an existing resource with new values.
func (p *supabaseProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Update(%s)", p.name, urn)
	logging.V(9).Infof("%s executing", label)
	// Our example Random resource will never be updated - if there is a diff, it will be a replacement.
	msg := fmt.Sprintf("Update is not yet implemented for %s", urn.Type())
	return nil, status.Error(codes.Unimplemented, msg)
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (p *supabaseProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Update(%s)", p.name, urn)
	logging.V(9).Infof("%s executing", label)
	// Implement Delete logic specific to your provider.
	// Note that for our Random resource, we don't have to do anything on Delete.
	return &pbempty.Empty{}, nil
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (p *supabaseProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: p.version,
	}, nil
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
	p.cancel()
	// TODO: Gropu of reauest
	return &pbempty.Empty{}, nil
}
