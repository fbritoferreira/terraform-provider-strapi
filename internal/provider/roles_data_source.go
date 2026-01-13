package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fbritoferreira/terraform-provider-strapi/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &RolesDataSource{}

type RolesDataSource struct {
	client *client.StrapiClient
}

type RolesDataSourceModel struct {
	Roles []RoleModel `tfsdk:"roles"`
}

type RoleModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func NewRolesDataSource() datasource.DataSource {
	return &RolesDataSource{}
}

func (d *RolesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roles"
}

func (d *RolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all Strapi roles available in your instance.",
		Blocks: map[string]schema.Block{
			"roles": schema.ListNestedBlock{
				Description: "List of roles",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the role.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the role.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the role.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of the role (e.g., 'authenticated', 'public').",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "The creation timestamp of the role.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "The last update timestamp of the role.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *RolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.StrapiClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.StrapiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *RolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RolesDataSourceModel

	roles, err := d.client.GetRoles()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading roles",
			fmt.Sprintf("Could not read roles: %s", err),
		)
		return
	}

	roleModels := make([]RoleModel, len(roles))
	for i, role := range roles {
		roleModels[i] = RoleModel{
			ID:          types.StringValue(strconv.Itoa(role.ID)),
			Name:        types.StringValue(role.Name),
			Description: types.StringValue(role.Description),
			Type:        types.StringValue(role.Type),
			CreatedAt:   types.StringValue(role.CreatedAt),
			UpdatedAt:   types.StringValue(role.UpdatedAt),
		}
	}

	state.Roles = roleModels

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read %d roles", len(roles)))
}
