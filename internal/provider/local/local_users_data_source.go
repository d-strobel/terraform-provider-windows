package local

import (
	"context"
	"fmt"
	"terraform-provider-windows/internal/generate/datasource_local_users"
	"time"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/winerror"
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
	var usersValueList []datasource_local_users.UsersValue

	for _, user := range winResp {
		usersValue := datasource_local_users.UsersValue{
			AccountExpires:         types.StringValue(user.AccountExpires.Format(time.RFC3339)),
			Description:            types.StringValue(user.Description),
			Enabled:                types.BoolValue(user.Enabled),
			FullName:               types.StringValue(user.FullName),
			Id:                     types.StringValue(user.SID.Value),
			LastLogon:              types.StringValue(user.LastLogon.Format(time.RFC3339)),
			Name:                   types.StringValue(user.Name),
			PasswordChangeableDate: types.StringValue(user.PasswordChangeableDate.Format(time.RFC3339)),
			PasswordExpires:        types.StringValue(user.PasswordExpires.Format(time.RFC3339)),
			PasswordLastSet:        types.StringValue(user.PasswordLastSet.Format(time.RFC3339)),
			PasswordRequired:       types.BoolValue(user.PasswordRequired),
			Sid:                    types.StringValue(user.SID.Value),
			UserMayChangePassword:  types.BoolValue(user.UserMayChangePassword),
		}

		objVal, diags := usersValue.ToObjectValue(ctx)
		resp.Diagnostics.Append(diags...)

		newUsersValue, diags := datasource_local_users.NewUsersValue(objVal.AttributeTypes(ctx), objVal.Attributes())
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		usersValueList = append(usersValueList, newUsersValue)
	}

	data.Users, diags = types.ListValueFrom(ctx, datasource_local_users.UsersValue{}.Type(ctx), usersValueList)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
