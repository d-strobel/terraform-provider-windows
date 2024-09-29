package local_test

import (
	"regexp"
	"testing"

	"github.com/d-strobel/terraform-provider-windows/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalGroupMembersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + `
          data "windows_local_group_members" "test" {
            name = "Administrators"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.windows_local_group_members.test", "name", "Administrators"),
					// Check first group in the list
					resource.TestCheckResourceAttr("data.windows_local_group_members.test", "members.0.sid", "S-1-5-21-153895498-367353507-3704405138-500"),
					resource.TestCheckResourceAttr("data.windows_local_group_members.test", "members.0.object_class", "User"),
					resource.TestMatchResourceAttr("data.windows_local_group_members.test", "members.0.name", regexp.MustCompile(`Administrator$`)),
					// Check second group in the list
					resource.TestCheckResourceAttr("data.windows_local_group_members.test", "members.1.sid", "S-1-5-21-153895498-367353507-3704405138-1000"),
					resource.TestCheckResourceAttr("data.windows_local_group_members.test", "members.1.object_class", "User"),
					resource.TestMatchResourceAttr("data.windows_local_group_members.test", "members.1.name", regexp.MustCompile(`vagrant$`)),
				),
			},
		},
	})
}
