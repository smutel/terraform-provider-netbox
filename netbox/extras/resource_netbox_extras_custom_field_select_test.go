package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldSelect = "netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldSelectMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldSelect,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldSelectFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldSelect,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldSelectMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldSelectConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldSelect),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldSelectConfig(nameSuffix string, resourceFull, extraResources bool) string {
	template := `
	resource "netbox_extras_custom_field" "test" {
		name = "test_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type         = "select"
		{{ if eq .resourcefull "true" }}
		description  = "Test custom field"
		choices = [
			"test",
			"test2"
		]
		label        = "Test Label for CF"
		weight       = 50
		#required     = true
		filter_logic = "disabled"
		default      = "test"
		{{ end }}
	}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"resourcefull":   strconv.FormatBool(resourceFull),
		"extraresources": strconv.FormatBool(extraResources),
	}
	return util.RenderTemplate(template, data)
}
