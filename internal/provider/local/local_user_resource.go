package local

import (
	"context"
	"fmt"
	"github.com/d-strobel/terraform-provider-windows/internal/generate/resource_local_user"
	"regexp"
	"time"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/windows/local/accounts"
	"github.com/d-strobel/gowindows/winerror"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	resp.Schema = resource_local_user.LocalUserResourceSchema(ctx)
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
	var data resource_local_user.LocalUserModel

	// Add log masking for Powershell secure strings.
	ctx = tflog.MaskAllFieldValuesRegexes(ctx, regexp.MustCompile(`\$\(ConvertTo-SecureString -String '([^']*)'.+\)`))

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Read time values from the plan
	accountExpiresPlanValue, diag := data.AccountExpires.ValueRFC3339Time()
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	params := accounts.UserCreateParams{
		Name:                  data.Name.ValueString(),
		FullName:              data.FullName.ValueString(),
		Description:           data.Description.ValueString(),
		Enabled:               data.Enabled.ValueBool(),
		Password:              data.Password.ValueString(),
		PasswordNeverExpires:  data.PasswordNeverExpires.ValueBool(),
		UserMayChangePassword: data.UserMayChangePassword.ValueBool(),
		AccountExpires:        accountExpiresPlanValue,
	}

	winResp, err := r.client.LocalAccounts.UserCreate(ctx, params)
	if err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to create local user:\n%s", err.Error()))
		return
	}

	// Set data
	data.AccountExpires, diag = timetypes.NewRFC3339Value(winResp.AccountExpires.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.Description = types.StringValue(winResp.Description)
	data.Enabled = types.BoolValue(winResp.Enabled)
	data.FullName = types.StringValue(winResp.FullName)
	data.Id = types.StringValue(winResp.SID.Value)

	data.LastLogon, diag = timetypes.NewRFC3339Value(winResp.LastLogon.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.Name = types.StringValue(winResp.Name)

	data.PasswordChangeableDate, diag = timetypes.NewRFC3339Value(winResp.PasswordChangeableDate.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.PasswordExpires, diag = timetypes.NewRFC3339Value(winResp.PasswordExpires.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.PasswordLastSet, diag = timetypes.NewRFC3339Value(winResp.PasswordLastSet.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.PasswordRequired = types.BoolValue(winResp.PasswordRequired)
	data.Sid = types.StringValue(winResp.SID.Value)
	data.UserMayChangePassword = types.BoolValue(winResp.UserMayChangePassword)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_local_user.LocalUserModel
	var diag diag.Diagnostics

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	winResp, err := r.client.LocalAccounts.UserRead(ctx, accounts.UserReadParams{SID: data.Id.ValueString()})
	if err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to read local user:\n%s", err.Error()))
		return
	}

	// Set data
	data.AccountExpires, diag = timetypes.NewRFC3339Value(winResp.AccountExpires.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.Description = types.StringValue(winResp.Description)
	data.Enabled = types.BoolValue(winResp.Enabled)
	data.FullName = types.StringValue(winResp.FullName)
	data.Id = types.StringValue(winResp.SID.Value)

	data.LastLogon, diag = timetypes.NewRFC3339Value(winResp.LastLogon.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.Name = types.StringValue(winResp.Name)

	data.PasswordChangeableDate, diag = timetypes.NewRFC3339Value(winResp.PasswordChangeableDate.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.PasswordExpires, diag = timetypes.NewRFC3339Value(winResp.PasswordExpires.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.PasswordLastSet, diag = timetypes.NewRFC3339Value(winResp.PasswordLastSet.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.PasswordRequired = types.BoolValue(winResp.PasswordRequired)
	data.Sid = types.StringValue(winResp.SID.Value)
	data.UserMayChangePassword = types.BoolValue(winResp.UserMayChangePassword)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_local_user.LocalUserModel

	// Add log masking for Powershell secure strings.
	ctx = tflog.MaskAllFieldValuesRegexes(ctx, regexp.MustCompile(`\$\(ConvertTo-SecureString -String '([^']*)'.+\)`))

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Read time values from the plan
	accountExpiresValue, diag := data.AccountExpires.ValueRFC3339Time()
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	params := accounts.UserUpdateParams{
		AccountExpires:        accountExpiresValue,
		Description:           data.Description.ValueString(),
		Enabled:               data.Enabled.ValueBool(),
		FullName:              data.FullName.ValueString(),
		Password:              data.Password.ValueString(),
		PasswordNeverExpires:  data.PasswordNeverExpires.ValueBool(),
		UserMayChangePassword: data.UserMayChangePassword.ValueBool(),
		SID:                   data.Sid.ValueString(),
	}

	if err := r.client.LocalAccounts.UserUpdate(ctx, params); err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to update local user:\n%s", err.Error()))
		return
	}

	tflog.Debug(ctx, "Read local user after update to synchronize Terraform state")
	winResp, err := r.client.LocalAccounts.UserRead(ctx, accounts.UserReadParams{SID: data.Sid.ValueString()})
	if err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to read local user after update:\n%s", err.Error()))
		return
	}

	// Set data
	data.AccountExpires, diag = timetypes.NewRFC3339Value(winResp.AccountExpires.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.Description = types.StringValue(winResp.Description)
	data.Enabled = types.BoolValue(winResp.Enabled)
	data.FullName = types.StringValue(winResp.FullName)
	data.Id = types.StringValue(winResp.SID.Value)

	data.LastLogon, diag = timetypes.NewRFC3339Value(winResp.LastLogon.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.Name = types.StringValue(winResp.Name)

	data.PasswordChangeableDate, diag = timetypes.NewRFC3339Value(winResp.PasswordChangeableDate.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.PasswordExpires, diag = timetypes.NewRFC3339Value(winResp.PasswordExpires.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.PasswordLastSet, diag = timetypes.NewRFC3339Value(winResp.PasswordLastSet.Format(time.RFC3339))
	resp.Diagnostics.Append(diag...)

	data.PasswordRequired = types.BoolValue(winResp.PasswordRequired)
	data.Sid = types.StringValue(winResp.SID.Value)
	data.UserMayChangePassword = types.BoolValue(winResp.UserMayChangePassword)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *localUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_local_user.LocalUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	err := r.client.LocalAccounts.UserDelete(ctx, accounts.UserDeleteParams{SID: data.Sid.ValueString()})
	if err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to delete local user:\n%s", err.Error()))
		return
	}
}

func (r *localUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
