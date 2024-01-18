package local

import (
	"context"
	"fmt"
	"terraform-provider-windows/internal/generate/datasource_local_groups"

	"github.com/d-strobel/gowindows"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	winResp, err := d.client.Local.GroupList(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read local security groups, got error: %s", err))
		return
	}

	// Set data
	for _, v := range winResp {
		groupsData := datasource_local_groups.GroupsValue{
			Name:        types.StringValue(v.Name),
			Description: types.StringValue(v.Description),
			Sid:         types.StringValue(v.SID.Value),
			Id:          types.StringValue(v.SID.Value),
		}

		list, _ := types.ListValueFrom(ctx, datasource_local_groups.GroupsType{}, groupsData)
		data.Groups = list
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
