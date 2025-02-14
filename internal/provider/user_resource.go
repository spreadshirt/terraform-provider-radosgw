package provider

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client *admin.API
}

// Configure implements resource.ResourceWithConfigure.
func (r *userResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*admin.API)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *admin.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				Required: true,
			},
			"display_name": schema.StringAttribute{
				Required: true,
			},
			"max_buckets": schema.Int64Attribute{
				Optional: true,
			},
		},
	}
}

type userResourceModel struct {
	UserID      types.String `tfsdk:"user_id"`
	DisplayName types.String `tfsdk:"display_name"`
	MaxBuckets  types.Int64  `tfsdk:"max_buckets"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := admin.User{
		ID:          plan.UserID.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
	}

	if !plan.MaxBuckets.IsNull() {
		maxBuckets := int(plan.MaxBuckets.ValueInt64())
		user.MaxBuckets = &maxBuckets
	}

	user, err := r.client.CreateUser(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UserID = types.StringValue(user.ID)
	plan.DisplayName = types.StringValue(user.DisplayName)
	if user.MaxBuckets != nil && *user.MaxBuckets != 1000 {
		plan.MaxBuckets = types.Int64Value(int64(*user.MaxBuckets))
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(ctx, admin.User{
		ID: state.UserID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			"Could not read user "+state.UserID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.UserID = types.StringValue(user.ID)
	state.DisplayName = types.StringValue(user.DisplayName)
	if user.MaxBuckets != nil && *user.MaxBuckets != 1000 {
		state.MaxBuckets = types.Int64Value(int64(*user.MaxBuckets))
	} else {
		state.MaxBuckets = types.Int64Null()
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(ctx, admin.User{ID: plan.UserID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving user",
			"Could not retrieve user, unexpected error: "+err.Error(),
		)
		return
	}

	user.ID = plan.UserID.ValueString()
	user.DisplayName = plan.DisplayName.ValueString()

	if !plan.MaxBuckets.IsNull() {
		maxBuckets := int(plan.MaxBuckets.ValueInt64())
		user.MaxBuckets = &maxBuckets
	}

	user, err = r.client.ModifyUser(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
			"Could not update user, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UserID = types.StringValue(user.ID)
	plan.DisplayName = types.StringValue(user.DisplayName)
	if user.MaxBuckets != nil {
		plan.MaxBuckets = types.Int64Value(int64(*user.MaxBuckets))
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the user deletion request
	user := admin.User{
		ID: state.UserID.ValueString(),
	}

	// Call RemoveUser to delete the user
	err := r.client.RemoveUser(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting user",
			"Could not delete user "+state.UserID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Remove the resource from Terraform state
	resp.State.RemoveResource(ctx)
}
