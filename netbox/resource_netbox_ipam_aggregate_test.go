package netbox_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceNameNetboxIpamAggregate = "netbox_ipam_aggregate.test"

func TestAccNetboxIpamAggregateMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	prefix := "192.168.54.0/24"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, false, false, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
			{
				ResourceName:      resourceNameNetboxIpamAggregate,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamAggregateFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	prefix := "192.168.55.0/24"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, true, true, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
			{
				ResourceName:      resourceNameNetboxIpamAggregate,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamAggregateMininmalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	prefix := "192.168.56.0/24"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, false, false, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, true, true, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, true, true, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, true, true, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
		},
	})
}

func testAccCheckNetboxIPAMAggregateConfig(nameSuffix string, resourceFull, extraResources bool, prefix string) string {
	template := `
	#resource "netbox_ipam_rir" "test" {
	#	name = "test-{{ .namesuffix }}"
	#	slug = "test-{{ .namesuffix }}"
	#}
	data "netbox_json_ipam_rirs_list" "json_rir" {
		limit = 1
	}

	{{ if eq .extraresources "true" }}
	#resource "netbox_extras_tag" "test" {
	#	name = "test-{{ .namesuffix }}"
	#	slug = "test-{{ .namesuffix }}"
	#}

	resource "netbox_tenancy_tenant" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_ipam_aggregate" "test" {
		prefix = "{{ .prefix }}"
		rir_id = jsondecode(data.netbox_json_ipam_rirs_list.json_rir.json)[0].id
		#rir_id = netbox_ipam_rir.test.id

		{{ if eq .resourcefull "true" }}
		tenant_id = netbox_tenancy_tenant.test.id
		date_added = "1971-01-02"
		description = "Test Aggregate"
		#tag {
		#	name = netbox_extras_tag.test.name
		#	slug = netbox_extras_tag.test.slug
		#}
		{{ end }}
	}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"prefix":         prefix,
		"extraresources": strconv.FormatBool(extraResources),
		"resourcefull":   strconv.FormatBool(resourceFull),
	}
	return renderTemplate(template, data)
}
