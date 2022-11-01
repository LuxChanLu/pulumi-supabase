// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Supabase
{
    public static class GetTypeScript
    {
        public static Task<GetTypeScriptResult> InvokeAsync(GetTypeScriptArgs? args = null, InvokeOptions? options = null)
            => Pulumi.Deployment.Instance.InvokeAsync<GetTypeScriptResult>("supabase:index:GetTypeScript", args ?? new GetTypeScriptArgs(), options.WithDefaults());

        public static Output<GetTypeScriptResult> Invoke(GetTypeScriptInvokeArgs? args = null, InvokeOptions? options = null)
            => Pulumi.Deployment.Instance.Invoke<GetTypeScriptResult>("supabase:index:GetTypeScript", args ?? new GetTypeScriptInvokeArgs(), options.WithDefaults());
    }


    public sealed class GetTypeScriptArgs : Pulumi.InvokeArgs
    {
        /// <summary>
        /// Included schemas
        /// </summary>
        [Input("includedSchemas")]
        public string? IncludedSchemas { get; set; }

        /// <summary>
        /// ID of the project
        /// </summary>
        [Input("projectId")]
        public string? ProjectId { get; set; }

        public GetTypeScriptArgs()
        {
            IncludedSchemas = "";
        }
    }

    public sealed class GetTypeScriptInvokeArgs : Pulumi.InvokeArgs
    {
        /// <summary>
        /// Included schemas
        /// </summary>
        [Input("includedSchemas")]
        public Input<string>? IncludedSchemas { get; set; }

        /// <summary>
        /// ID of the project
        /// </summary>
        [Input("projectId")]
        public Input<string>? ProjectId { get; set; }

        public GetTypeScriptInvokeArgs()
        {
            IncludedSchemas = "";
        }
    }


    [OutputType]
    public sealed class GetTypeScriptResult
    {
        /// <summary>
        /// TypeScript types of the project
        /// </summary>
        public readonly string Types;

        [OutputConstructor]
        private GetTypeScriptResult(string types)
        {
            Types = types;
        }
    }
}