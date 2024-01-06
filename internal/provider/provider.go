package provider

import (
	"context"
	"terraform-provider-windows/internal/generator/provider_windows"
	"terraform-provider-windows/internal/provider/local"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/connection"
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
	// Use generated schema
	resp.Schema = provider_windows.WindowsProviderSchema(ctx)
}

func (p *WindowsProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{}
}

func (p *WindowsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data provider_windows.WindowsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check for nested schema
	if data.Winrm.IsNull() && data.Ssh.IsNull() {
		resp.Diagnostics.AddError("Invalid provider configuration", "ssh or winrm must be specified")
	}
	if !data.Winrm.IsNull() && !data.Ssh.IsNull() {
		resp.Diagnostics.AddError("Invalid provider configuration", "only one of ssh or winrm must be specified")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Init connection configuration
	config := &connection.Config{}

	// Check and setup WinRM configuration
	if !data.Winrm.IsNull() {

		// Check required parameters for WinRM
		if data.Winrm.Username.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("winrm"), "Missing parameter for winrm schema", "Parameter 'username' must be set")
		}
		if data.Winrm.Password.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("winrm"), "Missing parameter for winrm schema", "Parameter 'password' must be set")
		}
		if resp.Diagnostics.HasError() {
			return
		}

		// Init WinRM config
		config.WinRM = &connection.WinRMConfig{}

		config.WinRM.WinRMHost = data.Endpoint.ValueString()
		config.WinRM.WinRMUsername = data.Winrm.Username.ValueString()
		config.WinRM.WinRMPassword = data.Winrm.Password.ValueString()

		if !data.Winrm.Port.IsNull() {
			config.WinRM.WinRMPort = int(data.Winrm.Port.ValueInt64())
		}
		if !data.Winrm.Timeout.IsNull() {
			config.WinRM.WinRMTimeout = time.Duration(data.Winrm.Timeout.ValueInt64())
		}
		if !data.Winrm.Port.IsNull() {
			config.WinRM.WinRMPort = int(data.Winrm.Port.ValueInt64())
		}
	}

	// Check and set SSH configuration
	if !data.Ssh.IsNull() {

		// Check required parameters for SSH
		if data.Ssh.Username.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("ssh"), "Missing parameter for ssh schema", "Parameter 'username' must be set")
		}
		if data.Ssh.Password.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("ssh"), "Missing parameter for ssh schema", "Parameter 'password' must be set")
		}
		if resp.Diagnostics.HasError() {
			return
		}
		// Init SSH config
		config.SSH = &connection.SSHConfig{}

		config.SSH.SSHHost = data.Endpoint.ValueString()
		config.SSH.SSHUsername = data.Ssh.Username.ValueString()
	}

	// Setup client
	client, err := gowindows.NewClient(config)
	if err != nil {
		resp.Diagnostics.AddError("Unable to setup client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *WindowsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
	// return []func() resource.Resource{
	// 	NewExampleResource,
	// }
}

func (p *WindowsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		local.NewLocalGroupDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WindowsProvider{
			version: version,
		}
	}
}
