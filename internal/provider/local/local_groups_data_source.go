package local

import (
	"context"
	"fmt"
	"github.com/d-strobel/terraform-provider-windows/internal/generate/datasource_local_groups"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/winerror"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = (*localGroupsDataSource)(nil)

func NewLocalGroupsDataSource() datasource.DataSource {
	return &localGroupsDataSource{}
}

type localGroupsDataSource struct {
	client *gowindows.Client
}

func (d *localGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_groups"
}

func (d *localGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_local_groups.LocalGroupsDataSourceSchema(ctx)
	resp.Schema.Description = `Retrieve a list of all local security groups.`
}

func (d *localGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *localGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_local_groups.LocalGroupsModel
	var diags diag.Diagnostics

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	winResp, err := d.client.LocalAccounts.GroupList(ctx)
	if err != nil {
		tflog.Error(ctx, "Received unexpected error from remote windows client", map[string]interface{}{
			"command": winerror.UnwrapCommand(err),
		})
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to read local security groups:\n%s", err.Error()))
		return
	}

	// Convert the response to the expected data source schema.
	// This might be confusing but is necessary.
	// For further information, see the following issue:
	// https://github.com/hashicorp/terraform-plugin-codegen-framework/issues/80
	tflog.Trace(ctx, "Converting the windows remote client response to the expected data source schema")
	var groupsValueList []datasource_local_groups.GroupsValue

	for _, group := range winResp {
		groupsValue := datasource_local_groups.GroupsValue{
			Name:        types.StringValue(group.Name),
			Description: types.StringValue(group.Description),
			Sid:         types.StringValue(group.SID.Value),
			Id:          types.StringValue(group.SID.Value),
		}

		objVal, diags := groupsValue.ToObjectValue(ctx)
		resp.Diagnostics.Append(diags...)

		newGroupsValue, diags := datasource_local_groups.NewGroupsValue(objVal.AttributeTypes(ctx), objVal.Attributes())
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		groupsValueList = append(groupsValueList, newGroupsValue)
	}

	data.Groups, diags = types.ListValueFrom(ctx, datasource_local_groups.GroupsValue{}.Type(ctx), groupsValueList)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
