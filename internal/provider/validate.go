package provider

import (
	"context"
	"fmt"
	"os"
	"terraform-provider-windows/internal/generator/provider_windows"

	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// Define imperative expressions to validate the provider config.
func (p *WindowsProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{
		providervalidator.ExactlyOneOf(
			path.MatchRoot("winrm"),
			path.MatchRoot("ssh"),
		),
	}
}

// Define programmatic expressions to validate the provider config.
func (p *WindowsProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var data provider_windows.WindowsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check WinRM attributes
	if !data.Winrm.IsNull() {

		// username must be set via config or environment variable
		if data.Winrm.Username.IsNull() && os.Getenv(envWinRMUsername) == "" {
			resp.Diagnostics.AddAttributeError(path.Root("winrm"), "Missing config attribute", fmt.Sprintf("Parameter 'username' or environment variable '%s' must be set.", envWinRMUsername))
		}

		// password must be set via config or environment variable
		if data.Winrm.Password.IsNull() && os.Getenv(envWinRMPassword) == "" {
			resp.Diagnostics.AddAttributeError(path.Root("winrm"), "Missing config attribute", fmt.Sprintf("Parameter 'password' or environment variable '%s' must be set.", envWinRMPassword))
		}
	}
}
