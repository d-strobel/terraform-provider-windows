package local_test

import (
	"testing"

	"github.com/d-strobel/terraform-provider-windows/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check if the data source works with name
			{
				Config: acctest.ProviderConfig() + `
          data "windows_local_group" "test" {
            name = "Administrators"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group.test", "id", "S-1-5-32-544"),
					resource.TestCheckResourceAttr("data.windows_local_group.test", "sid", "S-1-5-32-544"),
					resource.TestCheckResourceAttr("data.windows_local_group.test", "name", "Administrators"),
					resource.TestCheckResourceAttr("data.windows_local_group.test", "description", "Administrators have complete and unrestricted access to the computer/domain"),
				),
			},
			// Check if the data source works with SID
			{
				Config: acctest.ProviderConfig() + `
          data "windows_local_group" "test" {
            sid = "S-1-5-32-544"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group.test", "id", "S-1-5-32-544"),
					resource.TestCheckResourceAttr("data.windows_local_group.test", "sid", "S-1-5-32-544"),
					resource.TestCheckResourceAttr("data.windows_local_group.test", "name", "Administrators"),
					resource.TestCheckResourceAttr("data.windows_local_group.test", "description", "Administrators have complete and unrestricted access to the computer/domain"),
				),
			},
		},
	})
}
