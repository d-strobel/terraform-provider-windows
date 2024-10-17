package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/d-strobel/terraform-provider-windows/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccProvider(t *testing.T) {
	// Skip test if TF_ACC is not set
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC is set")
	}

	// Load and assert environment variables
	host := os.Getenv("TFWINDOWS_TEST_HOST")
	assert.NotEmpty(t, host, "Environment variable not set: TFWINDOWS_TEST_HOST")

	username := os.Getenv("TFWINDOWS_TEST_USERNAME")
	assert.NotEmpty(t, username, "Environment variable not set: TFWINDOWS_TEST_USERNAME")

	password := os.Getenv("TFWINDOWS_TEST_PASSWORD")
	assert.NotEmpty(t, password, "Environment variable not set: TFWINDOWS_TEST_PASSWORD")

	httpPort := os.Getenv("TFWINDOWS_TEST_WINRM_HTTP_PORT")
	assert.NotEmpty(t, httpPort, "Environment variable not set: TFWINDOWS_TEST_WINRM_HTTP_PORT")

	httpsPort := os.Getenv("TFWINDOWS_TEST_WINRM_HTTPS_PORT")
	assert.NotEmpty(t, httpsPort, "Environment variable not set: TFWINDOWS_TEST_WINRM_HTTPS_PORT")

	sshPort := os.Getenv("TFWINDOWS_TEST_SSH_PORT")
	assert.NotEmpty(t, sshPort, "Environment variable not set: TFWINDOWS_TEST_SSH_PORT")

	sshKeyED25519Path := os.Getenv("TFWINDOWS_TEST_SSH_PRIVATE_KEY_ED25519_PATH")
	assert.NotEmpty(t, sshKeyED25519Path, "Environment variable not set: TFWINDOWS_TEST_SSH_PRIVATE_KEY_ED25519_PATH")

	sshKeyED25519, err := os.ReadFile(sshKeyED25519Path)
	assert.NoError(t, err, "Failed to read ED25519 private key file")

	sshKeyRSAPath := os.Getenv("TFWINDOWS_TEST_SSH_PRIVATE_KEY_RSA_PATH")
	assert.NotEmpty(t, sshKeyRSAPath, "Environment variable not set: TFWINDOWS_TEST_SSH_PRIVATE_KEY_RSA_PATH")

	sshKeyRSA, err := os.ReadFile(sshKeyRSAPath)
	assert.NoError(t, err, "Failed to read RSA private key file")

	const providerTestDatasourceConfig = `
    data "windows_local_group" "test" {
      name = "Administrators"
    }
   `

	// Acceptance tests for different provider configurations
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// WinRM with HTTP
			{
				Config: providerTestDatasourceConfig + fmt.Sprintf(`
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
        `, host, username, password, httpPort),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group.test", "name", "Administrators"),
				),
			},
			// WinRM over HTTPS
			{
				Config: providerTestDatasourceConfig + fmt.Sprintf(`
          provider "windows" {
            endpoint = "%s"

            winrm = {
              username = "%s"
              password = "%s"
              port     = %s
              insecure = true
              use_tls  = true
            }
          }
        `, host, username, password, httpsPort),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group.test", "name", "Administrators"),
				),
			},
			// SSH with username and password
			{
				Config: providerTestDatasourceConfig + fmt.Sprintf(`
          provider "windows" {
            endpoint = "%s"

            ssh = {
              username = "%s"
              password = "%s"
              port     = %s
              insecure = true
            }
          }
        `, host, username, password, sshPort),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group.test", "name", "Administrators"),
				),
			},
			// SSH with ED25519 publicKeyauthentication
			{
				Config: providerTestDatasourceConfig + fmt.Sprintf(`
          provider "windows" {
            endpoint = "%s"

            ssh = {
              username         = "%s"
              private_key_path = "%s"
              port             = %s
              insecure         = true
            }
          }
        `, host, username, sshKeyED25519Path, sshPort),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group.test", "name", "Administrators"),
				),
			},
			// SSH with RSA publicKeyauthentication
			{
				Config: providerTestDatasourceConfig + fmt.Sprintf(`
          provider "windows" {
            endpoint = "%s"

            ssh = {
              username         = "%s"
              private_key_path = "%s"
              port             = %s
              insecure         = true
            }
          }
        `, host, username, sshKeyRSAPath, sshPort),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group.test", "name", "Administrators"),
				),
			},
			// SSH with ED25519 publicKeyauthentication
			{
				Config: providerTestDatasourceConfig + fmt.Sprintf(`
          provider "windows" {
            endpoint = "%s"

            ssh = {
              username    = "%s"
              port        = %s
              insecure    = true
              private_key = <<EOT
%s
EOT
            }
          }
        `, host, username, sshPort, sshKeyED25519),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group.test", "name", "Administrators"),
				),
			},
			// SSH with RSA publicKeyauthentication
			{
				Config: providerTestDatasourceConfig + fmt.Sprintf(`
          provider "windows" {
            endpoint = "%s"

            ssh = {
              username    = "%s"
              port        = %s
              insecure    = true
              private_key = <<EOT
%s
EOT
            }
          }
        `, host, username, sshPort, sshKeyRSA),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group.test", "name", "Administrators"),
				),
			},
		},
	})
}
