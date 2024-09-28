package local_test

import (
	"testing"

	"github.com/d-strobel/terraform-provider-windows/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + `
          data "windows_local_groups" "test" {}
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check first group in the list
					resource.TestCheckResourceAttr("data.windows_local_groups.test", "groups.0.id", "S-1-5-32-579"),
					resource.TestCheckResourceAttr("data.windows_local_groups.test", "groups.0.sid", "S-1-5-32-579"),
					resource.TestCheckResourceAttr("data.windows_local_groups.test", "groups.0.name", "Access Control Assistance Operators"),
					resource.TestCheckResourceAttr("data.windows_local_groups.test", "groups.0.description", "Members of this group can remotely query authorization attributes and permissions for resources on this computer."),
					// Check second group in the list
					resource.TestCheckResourceAttr("data.windows_local_groups.test", "groups.1.id", "S-1-5-32-544"),
					resource.TestCheckResourceAttr("data.windows_local_groups.test", "groups.1.sid", "S-1-5-32-544"),
					resource.TestCheckResourceAttr("data.windows_local_groups.test", "groups.1.name", "Administrators"),
					resource.TestCheckResourceAttr("data.windows_local_groups.test", "groups.1.description", "Administrators have complete and unrestricted access to the computer/domain"),
				),
			},
		},
	})
}
