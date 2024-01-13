package local

import (
	"context"
	"fmt"
	"terraform-provider-windows/internal/generator/local_resources"
	"time"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/windows/local"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*localUserResource)(nil)

func NewLocalUserResource() resource.Resource {
	return &localUserResource{}
}

type localUserResource struct {
	client *gowindows.Client
}

func (r *localUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_user"
}

func (r *localUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = local_resources.LocalUserResourceSchema(ctx)
	resp.Schema.Description = `Manage local users.`
}

func (r *localUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *localUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data local_resources.LocalUserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert time string to time
	accountExpires, err := time.Parse(time.DateTime, data.AccountExpires.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("account_expires"), "Config Error", fmt.Sprintf("Unable to parse time, got error: %s", err))
		return
	}

	// Create API call logic
	params := local.UserCreateParams{
		Name:                  data.Name.ValueString(),
		FullName:              data.FullName.ValueString(),
		Description:           data.Description.ValueString(),
		Enabled:               data.Enabled.ValueBool(),
		Password:              data.Password.ValueString(),
		PasswordNeverExpires:  data.PasswordNeverExpires.ValueBool(),
		UserMayChangePassword: data.UserMayChangePassword.ValueBool(),
		AccountExpires:        accountExpires,
	}

	winResp, err := r.client.Local.UserCreate(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create local user, got error: %s", err))
		return
	}

	// Set data
	data.AccountExpires = types.StringValue(winResp.AccountExpires.Format(time.DateTime))
	data.Description = types.StringValue(winResp.Description)
	data.Enabled = types.BoolValue(winResp.Enabled)
	data.FullName = types.StringValue(winResp.FullName)
	data.Id = types.StringValue(winResp.SID.Value)
	data.LastLogin = types.StringValue(winResp.LastLogon.Format(time.DateTime))
	data.Name = types.StringValue(winResp.Name)
	data.PasswordChangeableDate = types.StringValue(winResp.PasswordChangeableDate.Format(time.DateTime))
	data.PasswordExpires = types.StringValue(winResp.PasswordChangeableDate.Format(time.DateTime))
	data.PasswordLastSet = types.StringValue(winResp.PasswordLastSet.Format(time.DateTime))
	data.PasswordRequired = types.BoolValue(winResp.PasswordRequired)
	data.Sid = types.StringValue(winResp.SID.Value)
	data.UserMayChangePassword = types.BoolValue(winResp.UserMayChangePassword)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data local_resources.LocalUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	params := local.UserReadParams{
		SID: data.Id.ValueString(),
	}

	winResp, err := r.client.Local.UserRead(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read local user, got error: %s", err))
		return

	}

	// Set data
	data.Id = types.StringValue(winResp.SID.Value)
	data.Sid = types.StringValue(winResp.SID.Value)
	data.Name = types.StringValue(winResp.Name)
	data.Description = types.StringValue(winResp.Description)
	data.Enabled = types.BoolValue(winResp.Enabled)
	data.PasswordRequired = types.BoolValue(winResp.PasswordRequired)
	data.AccountExpires = types.StringValue(winResp.AccountExpires.Format(time.DateTime))
	data.FullName = types.StringValue(winResp.FullName)
	data.LastLogin = types.StringValue(winResp.LastLogon.Format(time.DateTime))
	data.PasswordChangeableDate = types.StringValue(winResp.PasswordChangeableDate.Format(time.DateTime))
	data.PasswordExpires = types.StringValue(winResp.PasswordExpires.Format(time.DateTime))
	data.PasswordLastSet = types.StringValue(winResp.PasswordLastSet.Format(time.DateTime))
	data.UserMayChangePassword = types.BoolValue(winResp.UserMayChangePassword)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data local_resources.LocalUserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data local_resources.LocalUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}
