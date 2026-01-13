package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAdminUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAdminUserResourceConfig("testadmin@example.com", "TestFirst", "TestLast"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("strapi_admin_user.test", "email", "testadmin@example.com"),
					resource.TestCheckResourceAttr("strapi_admin_user.test", "firstname", "TestFirst"),
					resource.TestCheckResourceAttr("strapi_admin_user.test", "lastname", "TestLast"),
					resource.TestCheckResourceAttrSet("strapi_admin_user.test", "id"),
				),
			},
			{
				ResourceName:            "strapi_admin_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "registration_token"},
			},
			{
				Config: testAccAdminUserResourceConfig("testadmin@example.com", "UpdatedFirst", "UpdatedLast"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("strapi_admin_user.test", "firstname", "UpdatedFirst"),
					resource.TestCheckResourceAttr("strapi_admin_user.test", "lastname", "UpdatedLast"),
				),
			},
		},
	})
}

func testAccAdminUserResourceConfig(email, firstname, lastname string) string {
	return fmt.Sprintf(`
resource "strapi_admin_user" "test" {
  email     = %[1]q
  firstname = %[2]q
  lastname  = %[3]q
  password  = "TestPass123!"
  roles     = [1]
}
`, email, firstname, lastname)
}
