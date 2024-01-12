package local

import (
	"context"
	"fmt"
	"terraform-provider-windows/internal/generator/local_datasources"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/windows/local"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LocalGroupDataSource{}

func NewLocalGroupDataSource() datasource.DataSource {
	return &LocalGroupDataSource{}
}

// LocalGroupDataSource defines the data source implementation.
type LocalGroupDataSource struct {
	client *gowindows.Client
}

func (d *LocalGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local" + "_group"
}

func (d *LocalGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = local_datasources.LocalGroupDataSourceSchema(ctx)
	resp.Schema.Description = `Retrieve information about a local security group.
You can get a group by the name or the security ID of the group.
`
}

func (d *LocalGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Check data source client type
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

func (d *LocalGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data local_datasources.LocalGroupModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Group params
	params := local.GroupParams{
		Name: data.Name.ValueString(),
		SID:  data.Sid.ValueString(),
	}

	// Client call
	winResp, err := d.client.Local.GroupRead(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read local security group, got error: %s", err))
		return
	}

	// Save data into Terraform state
	data.Id = types.StringValue(winResp.SID.Value)
	data.Sid = types.StringValue(winResp.SID.Value)
	data.Name = types.StringValue(winResp.Name)
	data.Description = types.StringValue(winResp.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
