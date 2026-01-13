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

var _ resource.Resource = &AdminUserResource{}

type AdminUserResource struct {
	client *client.StrapiClient
}

type AdminUserResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Email             types.String `tfsdk:"email"`
	Firstname         types.String `tfsdk:"firstname"`
	Lastname          types.String `tfsdk:"lastname"`
	Password          types.String `tfsdk:"password"`
	IsActive          types.Bool   `tfsdk:"is_active"`
	Roles             types.List   `tfsdk:"roles"`
	PreferedLanguage  types.String `tfsdk:"prefered_language"`
	RegistrationToken types.String `tfsdk:"registration_token"`
}

func NewAdminUserResource() resource.Resource {
	return &AdminUserResource{}
}

func (r *AdminUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_admin_user"
}

func (r *AdminUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Strapi admin dashboard user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the admin user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The email address of the admin user.",
			},
			"firstname": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The first name of the admin user.",
			},
			"lastname": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The last name of the admin user.",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The password for the admin user. Required for create, optional for update. Must be at least 8 characters with 1 uppercase, 1 lowercase, and 1 digit.",
			},
			"is_active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the admin user account is active.",
			},
			"roles": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "List of role IDs assigned to the admin user.",
				ElementType:         types.Int64Type,
			},
			"prefered_language": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The preferred language of the admin user.",
			},
			"registration_token": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The registration token for the admin user.",
				Sensitive:           true,
			},
		},
	}
}

func (r *AdminUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AdminUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AdminUserResourceModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roles := []int{}
	diags = plan.Roles.ElementsAs(ctx, &roles, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminUser := client.AdminUser{
		Email:            plan.Email.ValueString(),
		Firstname:        plan.Firstname.ValueString(),
		Lastname:         plan.Lastname.ValueString(),
		Roles:            roles,
		PreferedLanguage: plan.PreferedLanguage.ValueString(),
	}

	if !plan.Password.IsNull() {
		adminUser.Password = plan.Password.ValueString()
	}

	createdUser, err := r.client.CreateAdminUser(adminUser)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating admin user",
			fmt.Sprintf("Could not create admin user: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(createdUser.ID))
	plan.Email = types.StringValue(createdUser.Email)
	plan.Firstname = types.StringValue(createdUser.Firstname)
	plan.Lastname = types.StringValue(createdUser.Lastname)
	plan.IsActive = types.BoolValue(true)
	plan.PreferedLanguage = types.StringValue(createdUser.PreferedLanguage)
	plan.RegistrationToken = types.StringValue(createdUser.RegistrationToken)

	roleList, diags := types.ListValueFrom(ctx, types.Int64Type, convertIntSliceToInt64Slice(createdUser.Roles))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Roles = roleList
	plan.Password = types.StringNull()

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Created admin user with ID: %d", createdUser.ID))
}

func (r *AdminUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AdminUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing admin user ID",
			fmt.Sprintf("Could not parse admin user ID '%s': %s", state.ID.ValueString(), err),
		)
		return
	}

	user, err := r.client.GetAdminUser(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading admin user",
			fmt.Sprintf("Could not read admin user: %s", err),
		)
		return
	}

	state.ID = types.StringValue(strconv.Itoa(user.ID))
	state.Email = types.StringValue(user.Email)
	state.Firstname = types.StringValue(user.Firstname)
	state.Lastname = types.StringValue(user.Lastname)
	state.IsActive = types.BoolValue(true)
	state.PreferedLanguage = types.StringValue(user.PreferedLanguage)
	state.RegistrationToken = types.StringValue(user.RegistrationToken)

	roleList, diags := types.ListValueFrom(ctx, types.Int64Type, convertIntSliceToInt64Slice(user.Roles))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Roles = roleList

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read admin user with ID: %d", user.ID))
}

func (r *AdminUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AdminUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing admin user ID",
			fmt.Sprintf("Could not parse admin user ID '%s': %s", plan.ID.ValueString(), err),
		)
		return
	}

	roles := []int{}
	diags = plan.Roles.ElementsAs(ctx, &roles, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminUser := client.AdminUser{
		Email:            plan.Email.ValueString(),
		Firstname:        plan.Firstname.ValueString(),
		Lastname:         plan.Lastname.ValueString(),
		Roles:            roles,
		PreferedLanguage: plan.PreferedLanguage.ValueString(),
	}

	if !plan.IsActive.IsNull() {
		isActive := plan.IsActive.ValueBool()
		adminUser.IsActive = &isActive
	}

	if !plan.Password.IsNull() {
		adminUser.Password = plan.Password.ValueString()
	}

	updatedUser, err := r.client.UpdateAdminUser(id, adminUser)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating admin user",
			fmt.Sprintf("Could not update admin user: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(updatedUser.ID))
	plan.Email = types.StringValue(updatedUser.Email)
	plan.Firstname = types.StringValue(updatedUser.Firstname)
	plan.Lastname = types.StringValue(updatedUser.Lastname)
	plan.PreferedLanguage = types.StringValue(updatedUser.PreferedLanguage)
	plan.RegistrationToken = types.StringValue(updatedUser.RegistrationToken)

	roleList, diags := types.ListValueFrom(ctx, types.Int64Type, convertIntSliceToInt64Slice(updatedUser.Roles))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Roles = roleList
	plan.Password = types.StringNull()

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updated admin user with ID: %d", updatedUser.ID))
}

func (r *AdminUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AdminUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing admin user ID",
			fmt.Sprintf("Could not parse admin user ID '%s': %s", state.ID.ValueString(), err),
		)
		return
	}

	err = r.client.DeleteAdminUser(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting admin user",
			fmt.Sprintf("Could not delete admin user: %s", err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted admin user with ID: %d", id))
}

func (r *AdminUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertIntSliceToInt64Slice(intSlice []int) []int64 {
	result := make([]int64, len(intSlice))
	for i, v := range intSlice {
		result[i] = int64(v)
	}
	return result
}
