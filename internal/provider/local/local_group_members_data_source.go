package local

import (
	"context"
	"fmt"
	"terraform-provider-windows/internal/generate/datasource_local_group_members"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/windows/local/accounts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.Schema.Description = `Retrieve a list of members for a specific local security group.`
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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read local group members, got error: %s", err))
		return
	}

	var membersValueList []datasource_local_group_members.MembersValue

	for _, member := range winResp {
		memberValue := datasource_local_group_members.MembersValue{
			Name:        types.StringValue(member.Name),
			Sid:         types.StringValue(member.SID.Value),
			ObjectClass: types.StringValue(member.ObjectClass),
		}
		objVal, _ := memberValue.ToObjectValue(ctx)
		newMembersValue, _ := datasource_local_group_members.NewMembersValue(objVal.AttributeTypes(ctx), objVal.Attributes())

		membersValueList = append(membersValueList, newMembersValue)
	}

	data.Members, _ = types.ListValueFrom(ctx, datasource_local_group_members.MembersValue{}.Type(ctx), membersValueList)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
