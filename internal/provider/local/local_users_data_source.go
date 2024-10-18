package local

import (
	"context"
	"fmt"
	"time"

	"github.com/d-strobel/terraform-provider-windows/internal/generate/datasource_local_users"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/winerror"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = (*localUsersDataSource)(nil)

func NewLocalUsersDataSource() datasource.DataSource {
	return &localUsersDataSource{}
}

type localUsersDataSource struct {
	client *gowindows.Client
}

func (d *localUsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_users"
}

func (d *localUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_local_users.LocalUsersDataSourceSchema(ctx)
	resp.Schema.Description = `Retrieve a list of all local users.`
}

func (d *localUsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *localUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_local_users.LocalUsersModel
	var diags diag.Diagnostics

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	winResp, err := d.client.LocalAccounts.UserList(ctx)
	if err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to read local users:\n%s", err.Error()))
		return
	}

	// Convert the response to the expected data source schema.
	// This might be confusing but is necessary.
	// For further information, see the following issue:
	// https://github.com/hashicorp/terraform-plugin-codegen-framework/issues/80
	tflog.Trace(ctx, "Converting the windows remote client response to the expected data source schema")
	userValue := datasource_local_users.UsersValue{}
	var userValues []datasource_local_users.UsersValue

	for _, user := range winResp {
		r := datasource_local_users.NewUsersValueMust(userValue.AttributeTypes(ctx), map[string]attr.Value{
			"account_expires":          types.StringValue(user.AccountExpires.Format(time.RFC3339)),
			"description":              types.StringValue(user.Description),
			"enabled":                  types.BoolValue(user.Enabled),
			"full_name":                types.StringValue(user.FullName),
			"id":                       types.StringValue(user.SID.Value),
			"last_logon":               types.StringValue(user.LastLogon.Format(time.RFC3339)),
			"name":                     types.StringValue(user.Name),
			"password_changeable_date": types.StringValue(user.PasswordChangeableDate.Format(time.RFC3339)),
			"password_expires":         types.StringValue(user.PasswordExpires.Format(time.RFC3339)),
			"password_last_set":        types.StringValue(user.PasswordLastSet.Format(time.RFC3339)),
			"password_required":        types.BoolValue(user.PasswordRequired),
			"sid":                      types.StringValue(user.SID.Value),
			"user_may_change_password": types.BoolValue(user.UserMayChangePassword),
		})
		userValues = append(userValues, r)
	}

	usersValueList, diags := types.ListValueFrom(ctx, userValue.Type(ctx), userValues)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Users, diags = types.ListValueFrom(ctx, datasource_local_users.UsersValue{}.Type(ctx), usersValueList)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
