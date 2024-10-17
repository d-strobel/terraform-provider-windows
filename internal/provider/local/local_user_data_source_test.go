package local_test

import (
	"testing"

	"github.com/d-strobel/terraform-provider-windows/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing with name
			{
				Config: acctest.ProviderConfig() + `
          data "windows_local_user" "test" {
            name = "Administrator"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_user.test", "name", "Administrator"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "full_name", ""),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "description", "Built-in account for administering the computer/domain"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "account_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "enabled", "true"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "user_may_change_password", "true"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "last_logon", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttrSet("data.windows_local_user.test", "password_expires"),
					resource.TestCheckResourceAttrSet("data.windows_local_user.test", "password_changeable_date"),
					resource.TestCheckResourceAttrSet("data.windows_local_user.test", "password_last_set"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "id", "S-1-5-21-153895498-367353507-3704405138-500"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "sid", "S-1-5-21-153895498-367353507-3704405138-500"),
				),
			},
			// Read testing with SID
			{
				Config: acctest.ProviderConfig() + `
          data "windows_local_user" "test" {
            sid = "S-1-5-21-153895498-367353507-3704405138-500"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_user.test", "name", "Administrator"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "full_name", ""),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "description", "Built-in account for administering the computer/domain"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "account_expires", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "enabled", "true"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "user_may_change_password", "true"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "last_logon", "0001-01-01T00:00:00Z"),
					resource.TestCheckResourceAttrSet("data.windows_local_user.test", "password_expires"),
					resource.TestCheckResourceAttrSet("data.windows_local_user.test", "password_changeable_date"),
					resource.TestCheckResourceAttrSet("data.windows_local_user.test", "password_last_set"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "id", "S-1-5-21-153895498-367353507-3704405138-500"),
					resource.TestCheckResourceAttr("data.windows_local_user.test", "sid", "S-1-5-21-153895498-367353507-3704405138-500"),
				),
			},
		},
	})
}
