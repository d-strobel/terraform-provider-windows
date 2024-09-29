package local_test

import (
	"testing"

	"github.com/d-strobel/terraform-provider-windows/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalUsersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing with name
			{
				Config: acctest.ProviderConfig() + `
          data "windows_local_users" "test" {}
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check first user in the list
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.0.name", "Administrator"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.0.full_name", ""),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.0.description", "Built-in account for administering the computer/domain"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.0.account_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.0.user_may_change_password", "true"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.0.last_logon", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttrSet("data.windows_local_users.test", "users.0.password_expires"),
					resource.TestCheckResourceAttrSet("data.windows_local_users.test", "users.0.password_changeable_date"),
					resource.TestCheckResourceAttrSet("data.windows_local_users.test", "users.0.password_last_set"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.0.id", "S-1-5-21-153895498-367353507-3704405138-500"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.0.sid", "S-1-5-21-153895498-367353507-3704405138-500"),
					// Check second user in the list
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.1.name", "DefaultAccount"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.1.full_name", ""),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.1.description", "A user account managed by the system."),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.1.account_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.1.enabled", "false"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.1.user_may_change_password", "true"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.1.last_logon", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttrSet("data.windows_local_users.test", "users.1.password_expires"),
					resource.TestCheckResourceAttrSet("data.windows_local_users.test", "users.1.password_changeable_date"),
					resource.TestCheckResourceAttrSet("data.windows_local_users.test", "users.1.password_last_set"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.1.id", "S-1-5-21-153895498-367353507-3704405138-503"),
					resource.TestCheckResourceAttr("data.windows_local_users.test", "users.1.sid", "S-1-5-21-153895498-367353507-3704405138-503"),
				),
			},
		},
	})
}
