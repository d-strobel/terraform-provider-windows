package local_test

import (
	"testing"

	"github.com/d-strobel/terraform-provider-windows/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing only name
			{
				Config: acctest.ProviderConfig() + `
          resource "windows_local_user" "test_1" {
            name = "test-1"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("windows_local_user.test_1", "name", "test-1"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "full_name", ""),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "description", ""),
					resource.TestCheckNoResourceAttr("windows_local_user.test_1", "password"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "enabled", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "account_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "password_never_expires", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "user_may_change_password", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "last_logon", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "password_changeable_date", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "password_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "password_last_set", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_1", "id"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_1", "sid"),
				),
			},
			// Update and Read testing only name
			{
				Config: acctest.ProviderConfig() + `
          resource "windows_local_user" "test_1" {
            name    = "test-1"
            enabled = false
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("windows_local_user.test_1", "name", "test-1"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "full_name", ""),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "description", ""),
					resource.TestCheckNoResourceAttr("windows_local_user.test_1", "password"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "enabled", "false"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "account_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "password_never_expires", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "user_may_change_password", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "last_logon", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "password_changeable_date", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "password_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_1", "password_last_set", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_1", "id"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_1", "sid"),
				),
			},
			// Create and Read testing with name + password
			{
				Config: acctest.ProviderConfig() + `
          resource "windows_local_user" "test_2" {
            name     = "test-2"
            password = "Supersecretpassword1234"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("windows_local_user.test_2", "name", "test-2"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "full_name", ""),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "description", ""),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "enabled", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "password", "Supersecretpassword1234"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "account_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "password_never_expires", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "user_may_change_password", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "last_logon", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "password_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_2", "password_changeable_date"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_2", "password_last_set"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_2", "id"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_2", "sid"),
				),
			},
			// Update and Read testing with name + password
			{
				Config: acctest.ProviderConfig() + `
          resource "windows_local_user" "test_2" {
            name     = "test-2"
            password = "NEWSupersecretpassword1234"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("windows_local_user.test_2", "name", "test-2"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "full_name", ""),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "description", ""),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "enabled", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "password", "NEWSupersecretpassword1234"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "account_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "password_never_expires", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "user_may_change_password", "true"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "last_logon", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_2", "password_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_2", "password_changeable_date"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_2", "password_last_set"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_2", "id"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_2", "sid"),
				),
			},
			// Create and Read testing with all possible parameters
			{
				Config: acctest.ProviderConfig() + `
          resource "windows_local_user" "test_3" {
            name                     = "test-3"
            full_name                = "Test User 3"
            description              = "Test user for Terraform test"
            password                 = "SuperSecretPassword123!"
            account_expires          = "2072-12-31T23:59:59Z"
            enabled                  = false
            password_never_expires   = false
            user_may_change_password = false
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("windows_local_user.test_3", "name", "test-3"),
					resource.TestCheckResourceAttr("windows_local_user.test_3", "full_name", "Test User 3"),
					resource.TestCheckResourceAttr("windows_local_user.test_3", "description", "Test user for Terraform test"),
					resource.TestCheckResourceAttr("windows_local_user.test_3", "password", "SuperSecretPassword123!"),
					resource.TestCheckResourceAttr("windows_local_user.test_3", "account_expires", "2072-12-31T23:59:59Z"),
					resource.TestCheckResourceAttr("windows_local_user.test_3", "enabled", "false"),
					resource.TestCheckResourceAttr("windows_local_user.test_3", "password_never_expires", "false"),
					resource.TestCheckResourceAttr("windows_local_user.test_3", "user_may_change_password", "false"),
					resource.TestCheckResourceAttr("windows_local_user.test_3", "last_logon", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_3", "password_expires"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_3", "password_changeable_date"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_3", "password_last_set"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_3", "id"),
					resource.TestCheckResourceAttrSet("windows_local_user.test_3", "sid"),
				),
			},
		},
	})
}
