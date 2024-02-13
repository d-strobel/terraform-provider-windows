package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"terraform-provider-windows/internal/generate/provider_windows"
	"terraform-provider-windows/internal/provider/local"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/connection"
)

const (
	// SSH environment variables.
	envSSHUsername       string = "WIN_SSH_USERNAME"
	envSSHPassword       string = "WIN_SSH_PASSWORD"
	envSSHPrivateKey     string = "WIN_SSH_PRIVATE_KEY"
	envSSHPrivateKeyPath string = "WIN_SSH_PRIVATE_KEY_PATH"
	envSSHKnownHostsPath string = "WIN_SSH_KNOWN_HOSTS_PATH"
	envSSHPort           string = "WIN_SSH_PORT"
	envSSHInsecure       string = "WIN_SSH_INSECURE"

	// WinRM environment variables.
	envWinRMUsername string = "WIN_WINRM_USERNAME"
	envWinRMPassword string = "WIN_WINRM_PASSWORD"
	envWinRMPort     string = "WIN_WINRM_PORT"
	envWinRMTimeout  string = "WIN_WINRM_TIMEOUT"
	envWinRMInsecure string = "WIN_WINRM_INSECURE"
	envWinRMUseTLS   string = "WIN_WINRM_USE_TLS"

	// WinRM Kerberos environment variables.
	envKerberosRealm      string = "WIN_KRB_REALM"
	envKerberosConfigFile string = "WIN_KRB_CONFIG_FILE"

	// SSH default values.
	defaultSSHPort     int  = 22
	defaultSSHInsecure bool = false

	// WinRM default values.
	defaultWinRMPort     int           = 5986
	defaultWinRMTimeout  time.Duration = 0
	defaultWinRMInsecure bool          = false
	defaultWinRMUseTLS   bool          = true
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

func (p *WindowsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data provider_windows.WindowsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Init connection configuration
	config := &connection.Config{}

	// Check and setup WinRM configuration
	if !data.Winrm.IsNull() {
		config.WinRM = &connection.WinRMConfig{}

		// Endpoint
		config.WinRM.WinRMHost = data.Endpoint.ValueString()

		// Username
		config.WinRM.WinRMUsername = data.Winrm.Password.ValueString()
		if data.Winrm.Password.IsNull() {
			config.WinRM.WinRMUsername = os.Getenv(envWinRMUsername)
		}

		// Password
		config.WinRM.WinRMPassword = data.Winrm.Password.ValueString()
		if data.Winrm.Password.IsNull() {
			config.WinRM.WinRMPassword = os.Getenv(envWinRMPassword)
		}

		// Port
		config.WinRM.WinRMPort = defaultWinRMPort
		if !data.Winrm.Port.IsNull() {
			config.WinRM.WinRMPort = int(data.Winrm.Port.ValueInt64())
		} else if os.Getenv(envWinRMPort) != "" {
			winrmPort, err := strconv.Atoi(os.Getenv(envWinRMPort))
			if err != nil {
				resp.Diagnostics.AddError("Environment variable conversion error", fmt.Sprintf("Failed to convert environment variable '%s' to integer. Error: %s", envWinRMPort, err))
			}
			config.WinRM.WinRMPort = winrmPort
		}

		// Timeout
		config.WinRM.WinRMTimeout = defaultWinRMTimeout
		if !data.Winrm.Timeout.IsNull() {
			config.WinRM.WinRMTimeout = time.Duration(data.Winrm.Timeout.ValueInt64())
		} else if os.Getenv(envWinRMTimeout) != "" {
			winrmTimeout, err := strconv.Atoi(os.Getenv(envWinRMTimeout))
			if err != nil {
				resp.Diagnostics.AddError("Environment variable conversion error", fmt.Sprintf("Failed to convert environment variable '%s' to integer. Error: %s", envWinRMTimeout, err))
			}
			config.WinRM.WinRMTimeout = time.Duration(winrmTimeout)
		}

		// UseTLS (https)
		config.WinRM.WinRMUseTLS = defaultWinRMUseTLS
		if !data.Winrm.UseTls.IsNull() {
			config.WinRM.WinRMUseTLS = data.Winrm.UseTls.ValueBool()
		} else if os.Getenv(envWinRMUseTLS) != "" {
			winrmUseTls, err := strconv.ParseBool(os.Getenv(envWinRMUseTLS))
			if err != nil {
				resp.Diagnostics.AddError("Environment variable conversion error", fmt.Sprintf("Failed to convert environment variable '%s' to integer. Error: %s", envWinRMUseTLS, err))
			}
			config.WinRM.WinRMUseTLS = winrmUseTls
		}

		// Insecure
		config.WinRM.WinRMInsecure = defaultWinRMInsecure
		if !data.Winrm.Insecure.IsNull() {
			config.WinRM.WinRMInsecure = data.Winrm.Insecure.ValueBool()
		} else if os.Getenv(envWinRMInsecure) != "" {
			winrmInsecure, err := strconv.ParseBool(os.Getenv(envWinRMInsecure))
			if err != nil {
				resp.Diagnostics.AddError("Environment variable conversion error", fmt.Sprintf("Failed to convert environment variable '%s' to bool. Error: %s", envWinRMInsecure, err))
			}
			config.WinRM.WinRMUseTLS = winrmInsecure
		}

		// Kerberos
		if !data.Kerberos.IsNull() {
			config.WinRM.WinRMKerberos = &connection.KerberosConfig{}

			// Realm
			if !data.Kerberos.Realm.IsNull() {
				config.WinRM.WinRMKerberos.Realm = data.Kerberos.Realm.ValueString()
			} else if os.Getenv(envKerberosRealm) != "" {
				config.WinRM.WinRMKerberos.Realm = os.Getenv(envKerberosRealm)
			}

			// Krb config file
			if !data.Kerberos.KrbConfigFile.IsNull() {
				config.WinRM.WinRMKerberos.KrbConfigFile = data.Kerberos.KrbConfigFile.ValueString()
			} else if os.Getenv(envKerberosConfigFile) != "" {
				config.WinRM.WinRMKerberos.KrbConfigFile = os.Getenv(envKerberosConfigFile)
			}
		}
	}

	// Check and set SSH configuration
	if !data.Ssh.IsNull() {
		config.SSH = &connection.SSHConfig{}

		// Endpoint
		config.SSH.SSHHost = data.Endpoint.ValueString()

		// Username
		if !data.Ssh.Username.IsNull() {
			config.SSH.SSHUsername = data.Ssh.Username.ValueString()
		} else if os.Getenv(envSSHUsername) != "" {
			config.SSH.SSHUsername = os.Getenv(envSSHUsername)
		}

		// Password
		if !data.Ssh.Password.IsNull() {
			config.SSH.SSHPassword = data.Ssh.Password.ValueString()
		} else if os.Getenv(envSSHPassword) != "" {
			config.SSH.SSHPassword = os.Getenv(envSSHPassword)
		}

		// Private Key
		if !data.Ssh.PrivateKey.IsNull() {
			config.SSH.SSHPrivateKey = data.Ssh.PrivateKey.ValueString()
		} else if os.Getenv(envSSHPrivateKey) != "" {
			config.SSH.SSHPrivateKey = os.Getenv(envSSHPrivateKey)
		}

		// Private Key path
		if !data.Ssh.PrivateKeyPath.IsNull() {
			config.SSH.SSHPrivateKeyPath = data.Ssh.PrivateKeyPath.ValueString()
		} else if os.Getenv(envSSHPrivateKeyPath) != "" {
			config.SSH.SSHPrivateKeyPath = os.Getenv(envSSHPrivateKeyPath)
		}

		// Port
		config.SSH.SSHPort = defaultSSHPort
		if !data.Ssh.Port.IsNull() {
			config.SSH.SSHPort = int(data.Ssh.Port.ValueInt64())
		} else if os.Getenv(envSSHPort) != "" {
			sshPort, err := strconv.Atoi(os.Getenv(envSSHPort))
			if err != nil {
				resp.Diagnostics.AddError("Environment variable conversion error", fmt.Sprintf("Failed to convert environment variable '%s' to integer. Error: %s", envSSHPort, err))
			}
			config.SSH.SSHPort = sshPort
		}

		// Insecure
		config.SSH.SSHInsecureIgnoreHostKey = defaultSSHInsecure
		if !data.Ssh.Insecure.IsNull() {
			config.SSH.SSHInsecureIgnoreHostKey = data.Ssh.Insecure.ValueBool()
		} else if os.Getenv(envSSHInsecure) != "" {
			sshInsecure, err := strconv.ParseBool(os.Getenv(envSSHInsecure))
			if err != nil {
				resp.Diagnostics.AddError("Environment variable conversion error", fmt.Sprintf("Failed to convert environment variable '%s' to bool. Error: %s", envSSHInsecure, err))
			}
			config.SSH.SSHInsecureIgnoreHostKey = sshInsecure
		}

		// Known hosts path
		if !data.Ssh.KnownHostsPath.IsNull() {
			config.SSH.SSHKnownHostsPath = data.Ssh.KnownHostsPath.ValueString()
		} else if os.Getenv(envSSHKnownHostsPath) != "" {
			config.SSH.SSHKnownHostsPath = os.Getenv(envSSHKnownHostsPath)
		}
	}

	// Do not setup client if any errors occur
	if resp.Diagnostics.HasError() {
		return
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
	return []func() resource.Resource{
		local.NewLocalGroupResource,
		local.NewLocalUserResource,
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
