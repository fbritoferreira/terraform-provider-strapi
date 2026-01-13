package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/fbritoferreira/terraform-provider-strapi/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &StrapiProvider{}
var _ provider.ProviderWithFunctions = &StrapiProvider{}
var _ provider.ProviderWithEphemeralResources = &StrapiProvider{}
var _ provider.ProviderWithActions = &StrapiProvider{}

type StrapiProvider struct {
	version string
}

type StrapiProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIToken types.String `tfsdk:"api_token"`
}

func (p *StrapiProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "strapi"
	resp.Version = p.version
}

func (p *StrapiProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The Strapi provider allows you to manage Strapi CMS resources.
		
Configure the provider with the endpoint and API token for your Strapi instance.`,
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The Strapi API endpoint URL. Can also be provided via STRAPI_ENDPOINT environment variable.",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The Strapi API token for authentication. Can also be provided via STRAPI_API_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *StrapiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Strapi provider")

	var config StrapiProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("STRAPI_ENDPOINT")
	apiToken := os.Getenv("STRAPI_API_TOKEN")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.APIToken.IsNull() {
		apiToken = config.APIToken.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Strapi Endpoint",
			"The provider cannot create the Strapi client as there is a missing or empty value for the Strapi endpoint. "+
				"Set the endpoint value in the configuration or use the STRAPI_ENDPOINT environment variable.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Strapi API Token",
			"The provider cannot create the Strapi client as there is a missing or empty value for the Strapi API token. "+
				"Set the api_token value in the configuration or use the STRAPI_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	strapiClient := client.New(endpoint, apiToken)
	resp.DataSourceData = strapiClient
	resp.ResourceData = strapiClient

	tflog.Info(ctx, fmt.Sprintf("Configured Strapi provider with endpoint: %s", endpoint))
}

func (p *StrapiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewRoleResource,
		NewAdminUserResource,
	}
}

func (p *StrapiProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *StrapiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRolesDataSource,
	}
}

func (p *StrapiProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *StrapiProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &StrapiProvider{
			version: version,
		}
	}
}
