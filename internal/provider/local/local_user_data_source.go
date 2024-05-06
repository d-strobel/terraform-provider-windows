package local

import (
	"context"
	"fmt"
	"terraform-provider-windows/internal/generate/datasource_local_user"
	"time"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/windows/local/accounts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*localUserDataSource)(nil)

func NewLocalUserDataSource() datasource.DataSource {
	return &localUserDataSource{}
}

type localUserDataSource struct {
	client *gowindows.Client
}

func (d *localUserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_user"
}

func (d *localUserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_local_user.LocalUserDataSourceSchema(ctx)
	resp.Schema.Description = `Get information about a local user.`
}

func (d *localUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*gowindows.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *gowindows.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *localUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_local_user.LocalUserModel
	var diag diag.Diagnostics

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	params := accounts.UserReadParams{
		Name: data.Name.ValueString(),
		SID:  data.Sid.ValueString(),
	}

	winResp, err := d.client.LocalAccounts.UserRead(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read local user, got error: %s", err))
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
