package provider

import (
	"context"
	"fmt"
	"os"
	"terraform-provider-windows/internal/generate/provider_windows"

	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// ConfigValidators define imperative expressions to validate the provider config.
func (p *WindowsProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{
		providervalidator.ExactlyOneOf(
			path.MatchRoot("winrm"),
			path.MatchRoot("ssh"),
		),
	}
}

// ValidateConfig defines programmatic expressions to validate the provider config.
func (p *WindowsProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var data provider_windows.WindowsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check WinRM attributes
	if !data.Winrm.IsNull() {

		// username must be set via config or environment variable
		if data.Winrm.Username.IsNull() && os.Getenv(envWinRMUsername) == "" {
			resp.Diagnostics.AddAttributeError(path.Root("winrm"),
				"Missing config attribute",
				fmt.Sprintf("Parameter 'username' or environment variable '%s' must be set.", envWinRMUsername),
			)
		}

		// password must be set via config or environment variable
		if data.Winrm.Password.IsNull() && os.Getenv(envWinRMPassword) == "" {
			resp.Diagnostics.AddAttributeError(path.Root("winrm"),
				"Missing config attribute",
				fmt.Sprintf("Parameter 'password' or environment variable '%s' must be set.", envWinRMPassword),
			)
		}
	}

	// Check SSH attributes
	if !data.Ssh.IsNull() {

		// username must be set via config or environment variable
		if data.Ssh.Username.IsNull() && os.Getenv(envSSHUsername) == "" {
			resp.Diagnostics.AddAttributeError(path.Root("ssh"),
				"Missing config attribute",
				fmt.Sprintf("Parameter 'username' or environment variable '%s' must be set.", envSSHUsername),
			)
		}

		// password must be set via config or environment variable
		if data.Ssh.Password.IsNull() && os.Getenv(envSSHPassword) == "" && data.Ssh.PrivateKey.IsNull() && os.Getenv(envSSHPrivateKey) == "" && data.Ssh.PrivateKeyPath.IsNull() && os.Getenv(envSSHPrivateKeyPath) == "" {
			resp.Diagnostics.AddAttributeError(path.Root("ssh"),
				"Missing config attribute",
				fmt.Sprintf("Exactly one of the following parameters must be set: ['password' or environment variable '%s', 'private_key' or environment variable '%s', 'private_key_path' or environment variable '%s'].", envSSHPassword, envSSHPrivateKey, envSSHPrivateKeyPath),
			)
		}
	}
}
