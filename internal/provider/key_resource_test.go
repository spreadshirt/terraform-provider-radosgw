package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleResourceConfig(t.TempDir()), // generate new user for each test
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
	return fmt.Sprintf(`
provider "radosgw" {
	endpoint          = "http://127.0.0.1:9000"
	access_key_id     = "RMkni81ukvCYTLCjk62d"
	secret_access_key = "k8xeC8Kb62PMSXglkeuS6kLLjOHRp6y5LMntsUAR"
}

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
