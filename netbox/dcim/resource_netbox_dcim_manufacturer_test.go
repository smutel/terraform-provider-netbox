package dcim_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxDcimManufacturer = "netbox_dcim_manufacturer.test"

func TestAccNetboxDcimManufacturerMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
    acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimManufacturerConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimManufacturer),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimManufacturer,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimManufacturerFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
    acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimManufacturerConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimManufacturer),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimManufacturer,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimManufacturerMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
    acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimManufacturerConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimManufacturer),
				),
			},
			{
				Config: testAccCheckNetboxDcimManufacturerConfig(nameSuffix,
					true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimManufacturer),
				),
			},
			{
				Config: testAccCheckNetboxDcimManufacturerConfig(nameSuffix,
					false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimManufacturer),
				),
			},
			{
				Config: testAccCheckNetboxDcimManufacturerConfig(nameSuffix,
					false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxDcimManufacturer),
				),
			},
		},
	})
}

func testAccCheckNetboxDcimManufacturerConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "dcimmanufacturer-{{ .namesuffix }}"
		slug = "dcimmanufacturer-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_dcim_manufacturer" "test" {
		name				= "dcimmanufacturer-{{ .namesuffix }}"
		slug				= "dcimmanufacturer-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test device role"

		tag {
			name = netbox_extras_tag.test.name
			slug = netbox_extras_tag.test.slug
		}
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
