package dcim_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

const resourceNameNetboxDcimRack = "netbox_dcim_rack.test"

func TestAccNetboxDcimRackMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimRackConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRack),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimRack,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimRackFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimRackConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRack),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimRack,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimRackMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimRackConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRack),
				),
			},
			{
				Config: testAccCheckNetboxDcimRackConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRack),
				),
			},
			{
				Config: testAccCheckNetboxDcimRackConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRack),
				),
			},
			{
				Config: testAccCheckNetboxDcimRackConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRack),
				),
			},
		},
	})
}

func testAccCheckNetboxDcimRackConfig(nameSuffix string, resourceFull, extraResources bool) string {
	template := `
  {{ if eq .extraresources "true" }}
  resource "netbox_extras_tag" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

  resource "netbox_dcim_rack_role" "test" {
    name = "test-{{ .namesuffix }}"
    slug = "test-{{ .namesuffix }}"
  }

  resource "netbox_tenancy_tenant" "test" {
    name = "test-{{ .namesuffix }}"
    slug = "test-{{ .namesuffix }}"
  }

  resource "netbox_dcim_location" "test" {
		name        = "test-{{ .namesuffix }}"
    site_id     = netbox_dcim_site.test.id
		slug        = "test-{{ .namesuffix }}"
  }
	{{ end }}

  resource "netbox_dcim_site" "test" {
    name = "test-{{ .namesuffix }}"
    slug = "test-{{ .namesuffix }}"
  }

	resource "netbox_dcim_rack" "test" {
		name        = "test-{{ .namesuffix }}"
    site_id     = netbox_dcim_site.test.id
    height      = 10
    width       = 10

    {{ if eq .resourcefull "true" }}
		asset_tag = "test-{{ .namesuffix }}"
    comments = <<-EOT
    Comments for Test Rack
    Multiline
    EOT
    desc_units = true
    facility = "test-{{ .namesuffix }}"
    // location_id = netbox_dcim_location.test.id
    outer_depth = 1
    outer_unit = "mm"
    outer_width = 2
    role_id = netbox_dcim_rack_role.test.id
    serial = "test-{{ .namesuffix }}"
    status = "reserved"
    tag {
      name = netbox_extras_tag.test.name
      slug = netbox_extras_tag.test.slug
    }
    tenant_id = netbox_tenancy_tenant.test.id
    type = "4-post-frame"
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
