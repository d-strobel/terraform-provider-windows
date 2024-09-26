package acctest

import (
	"github.com/d-strobel/terraform-provider-windows/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during acceptance testing.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"windows": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// ProviderConfig returns a provider configuration for use in acceptance tests.
func ProviderConfig() string {
	return `
provider "windows" {
  endpoint = "127.0.0.1"

  winrm = {
    username = "vagrant"
    password = "vagrant"
    port     = 15985
    insecure = true
    use_tls  = false
  }
}
`
}
