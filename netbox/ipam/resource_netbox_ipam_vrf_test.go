package ipam_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

const resourceNameIpamVrf = "netbox_ipam_vrf.test"

func TestAccNetboxIpamVrfMinimal(t *testing.T) {

	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIPAMVrfConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVrf),
				),
			},
			{
				ResourceName:      resourceNameIpamVrf,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamVrfFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIPAMVrfConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVrf),
				),
			},
			{
				ResourceName:      resourceNameIpamVrf,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamVrfMininmalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIPAMVrfConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVrf),
				),
			},
			{
				Config: testAccCheckNetboxIPAMVrfConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVrf),
				),
			},
			{
				Config: testAccCheckNetboxIPAMVrfConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVrf),
				),
			},
			{
				Config: testAccCheckNetboxIPAMVrfConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVrf),
				),
			},
		},
	})
}

func testAccCheckNetboxIPAMVrfConfig(nameSuffix string, resourceFull, extraResources bool) string {
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

  resource "netbox_ipam_route_targets" "rt_export_test" {
    name = "test1-{{ .namesuffix }}"
  }

  resource "netbox_ipam_route_targets" "rt_import_test" {
    name = "test2-{{ .namesuffix }}"
  }
	{{ end }}

	resource "netbox_ipam_vrf" "test" {
		name        = "test-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
    # enforce_unique = false
    export_targets = [ netbox_ipam_route_targets.rt_export_test.id ]
    import_targets = [ netbox_ipam_route_targets.rt_import_test.id ]
		rd = "test-{{ .namesuffix }}"
    description = "Test Vrf"
		comments = <<-EOT
		Test Vrf
		EOT

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
