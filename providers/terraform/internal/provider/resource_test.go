package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "flagmanagment_project" "test" {
  name        = "Test Project"
  description = "Test Description"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagmanagment_project.test", "name", "Test Project"),
					resource.TestCheckResourceAttr("flagmanagment_project.test", "description", "Test Description"),
				),
			},
		},
	})
}
