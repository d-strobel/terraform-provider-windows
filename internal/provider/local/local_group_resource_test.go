package local_test

import (
	"testing"

	"github.com/d-strobel/terraform-provider-windows/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + `
          resource "windows_local_group" "test" {
            name = "Test"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("windows_local_group.test", "name", "Test"),
					resource.TestCheckResourceAttr("windows_local_group.test", "description", " "),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + `
          resource "windows_local_group" "test" {
            name        = "Test"
            description = "Test description"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("windows_local_group.test", "name", "Test"),
					resource.TestCheckResourceAttr("windows_local_group.test", "description", "Test description"),
				),
			},
		},
	})
}