package dcim_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

const resourceNameNetboxDcimSite = "netbox_dcim_site.test"

func TestAccNetboxDcimSiteMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimSiteConfig(nameSuffix, false, false, 0),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSite),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimSite,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimSiteFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, 4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimSiteConfig(nameSuffix, true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSite),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimSite,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimSiteMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, 4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimSiteConfig(nameSuffix, false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSite),
				),
			},
			{
				Config: testAccCheckNetboxDcimSiteConfig(nameSuffix, true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSite),
				),
			},
			{
				Config: testAccCheckNetboxDcimSiteConfig(nameSuffix, false, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSite),
				),
			},
			{
				Config: testAccCheckNetboxDcimSiteConfig(nameSuffix, false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimSite),
				),
			},
		},
	})
}

func testAccCheckNetboxDcimSiteConfig(nameSuffix string, resourceFull, extraResources bool, asn int64) string {
	template := `
	{{ if eq .extraresources "true" }}
	#resource "netbox_dcim_site_group" "test" {
	#	name = "test-{{ .namesuffix }}"
	#	slug = "test-{{ .namesuffix }}"
	#}

  resource "netbox_dcim_region" "test" {
  	name = "test-{{ .namesuffix }}"
  	slug = "test-{{ .namesuffix }}"
  }

	resource "netbox_ipam_rir" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_tenant" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_ipam_asn" "test" {
		asn = {{ .asn }}
	    rir_id = netbox_ipam_rir.test.id
	}

	resource "netbox_extras_tag" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_dcim_site" "test" {
		name        = "test-{{ .namesuffix }}"
		slug        = "test-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		status           = "planned"
		description = "Test device role"
		facility    = "TestFaciliy1"
		#group_id    = netbox_dcim_site_group.test.id
		region_id   = netbox_dcim_region.test.id
		latitude    = 12.54632
		longitude   = 41.21632
		tenant_id   = netbox_tenancy_tenant.test.id
		time_zone = "Europe/Berlin"

		comments = <<-EOT
		Comments for Test device role
		Multiline
		EOT

		physical_address = <<-EOT
		multiline
		physical
		address
		EOT

		shipping_address = <<-EOT
		multiline
		shipping
		address
		EOT

		asns = [
			netbox_ipam_asn.test.id
		]

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
		"asn":            strconv.FormatInt(asn, 10),
	}
	return util.RenderTemplate(template, data)
}
