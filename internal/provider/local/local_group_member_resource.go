package local

import (
	"context"
	"fmt"
	"terraform-provider-windows/internal/generate/resource_local_group_member"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/windows/local/accounts"
	"github.com/d-strobel/gowindows/winerror"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = (*localGroupMemberResource)(nil)

func NewLocalGroupMemberResource() resource.Resource {
	return &localGroupMemberResource{}
}

type localGroupMemberResource struct {
	client *gowindows.Client
}

func (r *localGroupMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_group_member"
}

func (r *localGroupMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_local_group_member.LocalGroupMemberResourceSchema(ctx)
	resp.Schema.Description = `Manage group member for local security groups.`
}

func (r *localGroupMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*gowindows.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *gowindows.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *localGroupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_local_group_member.LocalGroupMemberModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	params := accounts.GroupMemberCreateParams{
		SID:    data.GroupId.ValueString(),
		Member: data.MemberId.ValueString(),
	}

	if err := r.client.LocalAccounts.GroupMemberCreate(ctx, params); err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to create local group member:\n%s", err.Error()))
		return
	}

	// Create the ID for the resource
	data.Id = types.StringValue(fmt.Sprintf("%s/member/%s", data.GroupId.ValueString(), data.MemberId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localGroupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_local_group_member.LocalGroupMemberModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	params := accounts.GroupMemberReadParams{
		SID:    data.GroupId.ValueString(),
		Member: data.MemberId.ValueString(),
	}

	if _, err := r.client.LocalAccounts.GroupMemberRead(ctx, params); err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to delete local group member:\n%s", err.Error()))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localGroupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not needed in this resource
}

func (r *localGroupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_local_group_member.LocalGroupMemberModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	params := accounts.GroupMemberDeleteParams{
		SID:    data.GroupId.ValueString(),
		Member: data.MemberId.ValueString(),
	}

	if err := r.client.LocalAccounts.GroupMemberDelete(ctx, params); err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to delete local group member:\n%s", err.Error()))
		return
	}
}
