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

var _ resource.Resource = &UserResource{}

type UserResource struct {
	client *client.StrapiClient
}

type UserResourceModel struct {
	ID         types.String `tfsdk:"id"`
	DocumentID types.String `tfsdk:"document_id"`
	Username   types.String `tfsdk:"username"`
	Email      types.String `tfsdk:"email"`
	Confirmed  types.Bool   `tfsdk:"confirmed"`
	Blocked    types.Bool   `tfsdk:"blocked"`
	RoleName   types.String `tfsdk:"role_name"`
	RoleID     types.Int64  `tfsdk:"role_id"`
}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Strapi user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"document_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The document ID of user.",
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The username of the user.",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The email address of the user.",
			},
			"confirmed": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether that user account is confirmed.",
			},
			"blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether that user account is blocked.",
			},
			"role_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name of the role to assign to the user.",
			},
			"role_id": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the role assigned to the user.",
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.StrapiClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.StrapiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := client.User{
		Username:  plan.Username.ValueString(),
		Email:     plan.Email.ValueString(),
		Confirmed: plan.Confirmed.ValueBool(),
		Blocked:   plan.Blocked.ValueBool(),
	}

	if !plan.RoleName.IsNull() {
		role, err := r.client.FindRoleByName(plan.RoleName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error finding role",
				fmt.Sprintf("Could not find role '%s': %s", plan.RoleName.ValueString(), err),
			)
			return
		}
		user.Role = map[string]interface{}{
			"connect": []map[string]interface{}{
				{"id": role.ID},
			},
		}
		plan.RoleID = types.Int64Value(int64(role.ID))
	} else if !plan.RoleID.IsNull() {
		user.Role = map[string]interface{}{
			"connect": []map[string]interface{}{
				{"id": int(plan.RoleID.ValueInt64())},
			},
		}
	}

	createdUser, err := r.client.CreateUser(user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			fmt.Sprintf("Could not create user: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(createdUser.ID))
	plan.DocumentID = types.StringValue(createdUser.DocumentID)
	plan.Username = types.StringValue(createdUser.Username)
	plan.Email = types.StringValue(createdUser.Email)
	plan.Confirmed = types.BoolValue(createdUser.Confirmed)
	plan.Blocked = types.BoolValue(createdUser.Blocked)

	if createdUser.Role != nil {
		if id, ok := createdUser.Role["id"].(float64); ok {
			plan.RoleID = types.Int64Value(int64(id))
		} else if id, ok := createdUser.Role["id"].(int); ok {
			plan.RoleID = types.Int64Value(int64(id))
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Created user with ID: %d", createdUser.ID))
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing user ID",
			fmt.Sprintf("Could not parse user ID '%s': %s", state.ID.ValueString(), err),
		)
		return
	}

	user, err := r.client.GetUser(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("Could not read user: %s", err),
		)
		return
	}

	state.ID = types.StringValue(strconv.Itoa(user.ID))
	state.DocumentID = types.StringValue(user.DocumentID)
	state.Username = types.StringValue(user.Username)
	state.Email = types.StringValue(user.Email)
	state.Confirmed = types.BoolValue(user.Confirmed)
	state.Blocked = types.BoolValue(user.Blocked)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read user with ID: %d", user.ID))
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing user ID",
			fmt.Sprintf("Could not parse user ID '%s': %s", plan.ID.ValueString(), err),
		)
		return
	}

	user := client.User{
		Username:  plan.Username.ValueString(),
		Email:     plan.Email.ValueString(),
		Confirmed: plan.Confirmed.ValueBool(),
		Blocked:   plan.Blocked.ValueBool(),
	}

	if !plan.RoleName.IsNull() && plan.RoleName.ValueString() != "" {
		role, err := r.client.FindRoleByName(plan.RoleName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error finding role",
				fmt.Sprintf("Could not find role '%s': %s", plan.RoleName.ValueString(), err),
			)
			return
		}
		user.Role = map[string]interface{}{
			"connect": []map[string]interface{}{
				{"id": role.ID},
			},
		}
		plan.RoleID = types.Int64Value(int64(role.ID))
	} else if !plan.RoleID.IsNull() {
		user.Role = map[string]interface{}{
			"connect": []map[string]interface{}{
				{"id": int(plan.RoleID.ValueInt64())},
			},
		}
	}

	updatedUser, err := r.client.UpdateUser(id, user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
			fmt.Sprintf("Could not update user: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(updatedUser.ID))
	plan.DocumentID = types.StringValue(updatedUser.DocumentID)
	plan.Username = types.StringValue(updatedUser.Username)
	plan.Email = types.StringValue(updatedUser.Email)
	plan.Confirmed = types.BoolValue(updatedUser.Confirmed)
	plan.Blocked = types.BoolValue(updatedUser.Blocked)

	if updatedUser.Role != nil {
		if id, ok := updatedUser.Role["id"].(float64); ok {
			plan.RoleID = types.Int64Value(int64(id))
		} else if id, ok := updatedUser.Role["id"].(int); ok {
			plan.RoleID = types.Int64Value(int64(id))
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updated user with ID: %d", updatedUser.ID))
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing user ID",
			fmt.Sprintf("Could not parse user ID '%s': %s", state.ID.ValueString(), err),
		)
		return
	}

	err = r.client.DeleteUser(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting user",
			fmt.Sprintf("Could not delete user: %s", err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted user with ID: %d", id))
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
