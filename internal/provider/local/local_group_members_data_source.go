package local

import (
	"context"
	"fmt"

	"github.com/d-strobel/terraform-provider-windows/internal/generate/datasource_local_group_members"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/windows/local/accounts"
	"github.com/d-strobel/gowindows/winerror"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = (*localGroupMembersDataSource)(nil)

func NewLocalGroupMembersDataSource() datasource.DataSource {
	return &localGroupMembersDataSource{}
}

type localGroupMembersDataSource struct {
	client *gowindows.Client
}

func (d *localGroupMembersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_group_members"
}

func (d *localGroupMembersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_local_group_members.LocalGroupMembersDataSourceSchema(ctx)
}

func (d *localGroupMembersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *localGroupMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_local_group_members.LocalGroupMembersModel
	var diags diag.Diagnostics

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	params := accounts.GroupMemberListParams{
		Name: data.Name.ValueString(),
	}
	winResp, err := d.client.LocalAccounts.GroupMemberList(ctx, params)
	if err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to read local group members:\n%s", err.Error()))
		return
	}

	// Convert the response to the expected data source schema.
	// This might be confusing but is necessary.
	// For further information, see the following issue:
	// https://github.com/hashicorp/terraform-plugin-codegen-framework/issues/80
	tflog.Trace(ctx, "Converting the windows remote client response to the expected data source schema")
	memberValue := datasource_local_group_members.MembersValue{}
	var memberValues []datasource_local_group_members.MembersValue

	for _, member := range winResp {
		r := datasource_local_group_members.NewMembersValueMust(memberValue.AttributeTypes(ctx), map[string]attr.Value{
			"name":         types.StringValue(member.Name),
			"sid":          types.StringValue(member.SID.Value),
			"object_class": types.StringValue(member.ObjectClass),
		})
		memberValues = append(memberValues, r)
	}

	membersValueList, diags := types.ListValueFrom(ctx, memberValue.Type(ctx), memberValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Members, diags = types.ListValueFrom(ctx, datasource_local_group_members.MembersValue{}.Type(ctx), membersValueList)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
