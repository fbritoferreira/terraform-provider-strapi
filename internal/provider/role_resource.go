package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fbritoferreira/terraform-provider-strapi/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &RoleResource{}

type RoleResource struct {
	client *client.StrapiClient
}

type RoleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func NewRoleResource() resource.Resource {
	return &RoleResource{}
}

func (r *RoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *RoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Strapi role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of role.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the role.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The description of the role.",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The type of the role (e.g., 'authenticated', 'public').",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The creation timestamp of the role.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last update timestamp of the role.",
			},
		},
	}
}

func (r *RoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.StrapiClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.StrapiClient, got: %T. Please report this issue to provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleResourceModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := client.Role{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Type:        plan.Type.ValueString(),
	}

	createdRole, err := r.client.CreateRole(role)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating role",
			fmt.Sprintf("Could not create role: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(createdRole.ID))
	plan.Name = types.StringValue(createdRole.Name)
	plan.Description = types.StringValue(createdRole.Description)
	plan.Type = types.StringValue(createdRole.Type)
	plan.CreatedAt = types.StringValue(createdRole.CreatedAt)
	plan.UpdatedAt = types.StringValue(createdRole.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Created role with ID: %d", createdRole.ID))
}

func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing role ID",
			fmt.Sprintf("Could not parse role ID '%s': %s", state.ID.ValueString(), err),
		)
		return
	}

	role, err := r.client.GetRole(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading role",
			fmt.Sprintf("Could not read role: %s", err),
		)
		return
	}

	state.ID = types.StringValue(strconv.Itoa(role.ID))
	state.Name = types.StringValue(role.Name)
	state.Description = types.StringValue(role.Description)
	state.Type = types.StringValue(role.Type)
	state.CreatedAt = types.StringValue(role.CreatedAt)
	state.UpdatedAt = types.StringValue(role.UpdatedAt)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read role with ID: %d", role.ID))
}

func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RoleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing role ID",
			fmt.Sprintf("Could not parse role ID '%s': %s", plan.ID.ValueString(), err),
		)
		return
	}

	role := client.Role{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Type:        plan.Type.ValueString(),
	}

	updatedRole, err := r.client.UpdateRole(id, role)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating role",
			fmt.Sprintf("Could not update role: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(updatedRole.ID))
	plan.Name = types.StringValue(updatedRole.Name)
	plan.Description = types.StringValue(updatedRole.Description)
	plan.Type = types.StringValue(updatedRole.Type)
	plan.UpdatedAt = types.StringValue(updatedRole.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updated role with ID: %d", updatedRole.ID))
}

func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing role ID",
			fmt.Sprintf("Could not parse role ID '%s': %s", state.ID.ValueString(), err),
		)
		return
	}

	err = r.client.DeleteRole(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting role",
			fmt.Sprintf("Could not delete role: %s", err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted role with ID: %d", id))
}

func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
