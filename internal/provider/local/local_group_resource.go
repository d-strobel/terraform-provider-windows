package local

import (
	"context"
	"fmt"
	"terraform-provider-windows/internal/generate/resource_local_group"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/windows/local/accounts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*localGroupResource)(nil)

func NewLocalGroupResource() resource.Resource {
	return &localGroupResource{}
}

type localGroupResource struct {
	client *gowindows.Client
}

func (r *localGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_group"
}

func (r *localGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_local_group.LocalGroupResourceSchema(ctx)
	resp.Schema.Description = `Manage local security groups.

**Note:** The description default is a string with a space.
This is necessary because the powershell function Set-LocalGroup does not allow an empty string.
`
}

func (r *localGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *localGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_local_group.LocalGroupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	params := accounts.GroupCreateParams{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	winResp, err := r.client.LocalAccounts.GroupCreate(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create local security group, got error: %s", err))
		return
	}

	// Example data value setting
	data.Name = types.StringValue(winResp.Name)
	data.Description = types.StringValue(winResp.Description)
	data.Sid = types.StringValue(winResp.SID.Value)
	data.Id = types.StringValue(winResp.SID.Value)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_local_group.LocalGroupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	winResp, err := r.client.LocalAccounts.GroupRead(ctx, accounts.GroupReadParams{SID: data.Id.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read local security group, got error: %s", err))
		return
	}

	// Set read values
	data.Name = types.StringValue(winResp.Name)
	data.Description = types.StringValue(winResp.Description)
	data.Sid = types.StringValue(winResp.SID.Value)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_local_group.LocalGroupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	params := accounts.GroupUpdateParams{
		SID:         data.Sid.ValueString(),
		Description: data.Description.ValueString(),
	}

	err := r.client.LocalAccounts.GroupUpdate(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update local security group, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_local_group.LocalGroupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	err := r.client.LocalAccounts.GroupDelete(ctx, accounts.GroupDeleteParams{SID: data.Sid.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete local security group, got error: %s", err))
		return
	}
}

func (r *localGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
