package dcim_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

const resourceNameNetboxDcimLocation = "netbox_dcim_location.test"

func TestAccNetboxDcimLocationMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimLocationConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimLocation),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimLocation,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimLocationFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimLocationConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimLocation),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimLocation,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimLocationMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimLocationConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimLocation),
				),
			},
			{
				Config: testAccCheckNetboxDcimLocationConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimLocation),
				),
			},
			{
				Config: testAccCheckNetboxDcimLocationConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimLocation),
				),
			},
			{
				Config: testAccCheckNetboxDcimLocationConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimLocation),
				),
			},
		},
	})
}

func testAccCheckNetboxDcimLocationConfig(nameSuffix string, resourceFull, extraResources bool) string {
	template := `
  {{ if eq .extraresources "true" }}
  resource "netbox_extras_tag" "test" {
    name = "test-{{ .namesuffix }}"
    slug = "test-{{ .namesuffix }}"
  }

  resource "netbox_tenancy_tenant" "test" {
    name = "test-{{ .namesuffix }}"
    slug = "test-{{ .namesuffix }}"
  }

  resource "netbox_dcim_location" "test2" {
		name        = "test2-{{ .namesuffix }}"
    site_id     = netbox_dcim_site.test.id
		slug        = "test2-{{ .namesuffix }}"
  }
  {{ end }}

  resource "netbox_dcim_site" "test" {
    name = "test-{{ .namesuffix }}"
    slug = "test-{{ .namesuffix }}"
  }

	resource "netbox_dcim_location" "test" {
		name        = "test-{{ .namesuffix }}"
    site_id     = netbox_dcim_site.test.id
		slug        = "test-{{ .namesuffix }}"

    {{ if eq .resourcefull "true" }}
    description = "Test location"
    parent_id = netbox_dcim_location.test2.id
    status = "staging"
    tag {
      name = netbox_extras_tag.test.name
      slug = netbox_extras_tag.test.slug
    }
    tenant_id = netbox_tenancy_tenant.test.id
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
