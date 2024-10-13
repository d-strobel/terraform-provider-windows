// package acctest contains utilities for acceptance testing.
package acctest

import (
	"fmt"
	"os"

	"github.com/d-strobel/terraform-provider-windows/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during acceptance testing
// and should therefore be imported in the test files.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"windows": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// ProviderConfig returns a default WinRM provider configuration for use in acceptance tests.
// Use this function in the test files to set up the provider configuration for non domain joined windows machines.
func ProviderConfig() string {
	// Load environment variables are present
	host := os.Getenv("TFWINDOWS_TEST_HOST")
	username := os.Getenv("TFWINDOWS_TEST_USERNAME")
	password := os.Getenv("TFWINDOWS_TEST_PASSWORD")
	port := os.Getenv("TFWINDOWS_TEST_WINRM_HTTP_PORT")

	return fmt.Sprintf(`
    provider "windows" {
      endpoint = "%s"

      winrm = {
        username = "%s"
        password = "%s"
        port     = %s
        insecure = true
        use_tls  = false
      }
    }
  `, host, username, password, port)
}
