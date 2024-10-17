package local_test

import (
	"testing"

	"github.com/d-strobel/terraform-provider-windows/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalGroupMemberResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with Guest user in Administrators group
			{
				Config: acctest.ProviderConfig() + `
          resource "windows_local_group_member" "test" {
            group_id  = "S-1-5-32-544"
            member_id = "S-1-5-21-153895498-367353507-3704405138-501"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("windows_local_group_member.test", "group_id", "S-1-5-32-544"),
					resource.TestCheckResourceAttr("windows_local_group_member.test", "member_id", "S-1-5-21-153895498-367353507-3704405138-501"),
					resource.TestCheckResourceAttr("windows_local_group_member.test", "id", "S-1-5-32-544/member/S-1-5-21-153895498-367353507-3704405138-501"),
				),
			},
			// Import testing
			{
				ResourceName:      "windows_local_group_member.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
