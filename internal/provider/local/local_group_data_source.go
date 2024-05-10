package local

import (
	"context"
	"fmt"
	"terraform-provider-windows/internal/generate/datasource_local_group"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/windows/local/accounts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*localGroupDataSource)(nil)

func NewLocalGroupDataSource() datasource.DataSource {
	return &localGroupDataSource{}
}

type localGroupDataSource struct {
	client *gowindows.Client
}

func (d *localGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_group"
}

func (d *localGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_local_group.LocalGroupDataSourceSchema(ctx)
	resp.Schema.Description = `Retrieve information about a local security group.
You can get a group by the name or the security ID of the group.
`
}

func (d *localGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *localGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_local_group.LocalGroupModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	params := accounts.GroupReadParams{
		Name: data.Name.ValueString(),
		SID:  data.Sid.ValueString(),
	}

	winResp, err := d.client.LocalAccounts.GroupRead(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Windows Client Error", fmt.Sprintf("Unable to read local security group:\n%s", err.Error()))
		return

	}

	// Save data into Terraform state
	data.Id = types.StringValue(winResp.SID.Value)
	data.Sid = types.StringValue(winResp.SID.Value)
	data.Name = types.StringValue(winResp.Name)
	data.Description = types.StringValue(winResp.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
