package provider

import (
	"context"
	"terraform-provider-windows/internal/generate/provider_windows"
	"terraform-provider-windows/internal/provider/local"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure WindowsProvider satisfies various provider interfaces.
var _ provider.Provider = &WindowsProvider{}

// WindowsProvider defines the provider implementation.
type WindowsProvider struct {
	version string
}

func (p *WindowsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "windows"
	resp.Version = p.version
}

func (p *WindowsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = provider_windows.WindowsProviderSchema(ctx)
	resp.Schema.Description = `The windows provider is used to interact remotely via winrm or ssh with a windows system.

**Important**:
Due to the limitations of the terraform-plugin-framework some attributes are listed as optionals even though a combination of certain parameters are required.
Check examples below for reference.
`
}

func (p *WindowsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		local.NewLocalGroupResource,
		local.NewLocalUserResource,
		local.NewLocalGroupMemberResource,
	}
}

func (p *WindowsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		local.NewLocalGroupDataSource,
		local.NewLocalGroupsDataSource,
		local.NewLocalUserDataSource,
		local.NewLocalUsersDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WindowsProvider{
			version: version,
		}
	}
}
