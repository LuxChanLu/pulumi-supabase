// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Supabase.Organization
{
    [SupabaseResourceType("supabase:organization:Project")]
    public partial class Project : Pulumi.ComponentResource
    {
        /// <summary>
        /// Project creation date
        /// </summary>
        [Output("created_at")]
        public Output<string> Created_at { get; private set; } = null!;

        /// <summary>
        /// ID of the project
        /// </summary>
        [Output("id")]
        public Output<string> Id { get; private set; } = null!;

        /// <summary>
        /// Name of the project
        /// </summary>
        [Output("name")]
        public Output<string> Name { get; private set; } = null!;

        /// <summary>
        /// Organization ID of the project
        /// </summary>
        [Output("organization_id")]
        public Output<string> Organization_id { get; private set; } = null!;

        /// <summary>
        /// Region of the project
        /// </summary>
        [Output("region")]
        public Output<Pulumi.Supabase.Organization.Region> Region { get; private set; } = null!;


        /// <summary>
        /// Create a Project resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Project(string name, ProjectArgs args, ComponentResourceOptions? options = null)
            : base("supabase:organization:Project", name, args ?? new ProjectArgs(), MakeResourceOptions(options, ""), remote: true)
        {
        }

        private static ComponentResourceOptions MakeResourceOptions(ComponentResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new ComponentResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = ComponentResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
    }

    public sealed class ProjectArgs : Pulumi.ResourceArgs
    {
        [Input("db_pass", required: true)]
        private Input<string>? _db_pass;

        /// <summary>
        /// Postgres password of the project
        /// </summary>
        public Input<string>? Db_pass
        {
            get => _db_pass;
            set
            {
                var emptySecret = Output.CreateSecret(0);
                _db_pass = Output.Tuple<Input<string>?, int>(value, emptySecret).Apply(t => t.Item1);
            }
        }

        /// <summary>
        /// KPS Enabled on the project
        /// </summary>
        [Input("kps_enabled", required: true)]
        public Input<bool> Kps_enabled { get; set; } = null!;

        /// <summary>
        /// Name of the project
        /// </summary>
        [Input("name", required: true)]
        public Input<string> Name { get; set; } = null!;

        /// <summary>
        /// Organization ID of the project
        /// </summary>
        [Input("organization_id", required: true)]
        public Input<string> Organization_id { get; set; } = null!;

        /// <summary>
        /// Plan of the project
        /// </summary>
        [Input("plan", required: true)]
        public Input<Pulumi.Supabase.Organization.Plan> Plan { get; set; } = null!;

        /// <summary>
        /// Region of the project
        /// </summary>
        [Input("region", required: true)]
        public Input<Pulumi.Supabase.Organization.Region> Region { get; set; } = null!;

        public ProjectArgs()
        {
        }
    }
}
