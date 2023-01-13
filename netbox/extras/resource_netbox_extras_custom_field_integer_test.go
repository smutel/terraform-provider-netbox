package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldInteger = "netbox_extras_custom_field.test"

func TestAccNetboxExtrasCustomFieldIntegerMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldIntegerConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldInteger),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldInteger,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldIntegerFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldIntegerConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldInteger),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldInteger,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldIntegerMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldIntegerConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldInteger),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldIntegerConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldInteger),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldIntegerConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldInteger),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldIntegerConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxExtrasCustomFieldInteger),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldIntegerConfig(nameSuffix string, resourceFull, extraResources bool) string {
	template := `
	resource "netbox_extras_custom_field" "test" {
		name = "test_{{ .namesuffix }}"
		content_types = [
			"dcim.site",
		]

		type               = "integer"
		{{ if eq .resourcefull "true" }}
		description        = "Test custom field"
		label              = "Test Label for CF"
		weight             = 50
		#required           = true
		filter_logic       = "disabled"
		# Fixed in Netbox 3.3
		#default            = 50
		validation_minimum = 1
		validation_maximum = 500
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
