package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestKeyCreateMultiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleResourceConfig(acctest.RandString(8)), // generate new user for each test
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"radosgw_key.test",
						tfjsonpath.New("access_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"radosgw_key.test",
						tfjsonpath.New("secret_key"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"radosgw_key.test-second",
						tfjsonpath.New("access_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"radosgw_key.test-second",
						tfjsonpath.New("secret_key"),
						knownvalue.NotNull(),
					),

					statecheck.CompareValuePairs(
						"radosgw_key.test", tfjsonpath.New("access_key"),
						"radosgw_key.test-second", tfjsonpath.New("access_key"),
						compare.ValuesDiffer(),
					),
				},
			},
		},
	})
}

func testAccExampleResourceConfig(user string) string {
	return fmt.Sprintf(testAccProviderSetup+`
resource "radosgw_user" "test" {
	user_id      = %q
	display_name = %q
}

resource "radosgw_key" "test" {
	user = %q

	depends_on = [radosgw_user.test]
}

resource "radosgw_key" "test-second" {
	user = %q

	depends_on = [radosgw_user.test]
}`, user, user, user, user)
}

func TestKeyCreateStatic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderSetup + `
variable "test_user_id" {}

resource "radosgw_user" "test" {
	user_id      = var.test_user_id
	display_name = var.test_user_id
}

resource "radosgw_key" "test-static" {
	user = var.test_user_id

	access_key = "static-access-key"
	secret_key = "static-secret-key"

	depends_on = [radosgw_user.test]
}
					`,
				ConfigVariables: config.Variables{
					"test_user_id": config.StringVariable(acctest.RandString(8)),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"radosgw_key.test-static",
						tfjsonpath.New("access_key"),
						knownvalue.StringExact("static-access-key"),
					),
					statecheck.ExpectKnownValue(
						"radosgw_key.test-static",
						tfjsonpath.New("secret_key"),
						knownvalue.StringExact("static-secret-key"),
					),
				},
			},
		},
	})
}
