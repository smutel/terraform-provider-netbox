package ipam_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

const resourceNameNetboxIpamAggregate = "netbox_ipam_aggregate.test"

func TestAccNetboxIpamAggregateMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	prefix := "192.168.54.0/24"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, false, false, prefix),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxIpamAggregate),
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
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, true, true, prefix),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxIpamAggregate),
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
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, false, false, prefix),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, true, true, prefix),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, true, true, prefix),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
			{
				Config: testAccCheckNetboxIPAMAggregateConfig(nameSuffix, true, true, prefix),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxIpamAggregate),
				),
			},
		},
	})
}

func testAccCheckNetboxIPAMAggregateConfig(nameSuffix string, resourceFull, extraResources bool, prefix string) string {
	template := `
	resource "netbox_ipam_rir" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_tenant" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_ipam_aggregate" "test" {
		prefix = "{{ .prefix }}"
		rir_id = netbox_ipam_rir.test.id

		{{ if eq .resourcefull "true" }}
		tenant_id = netbox_tenancy_tenant.test.id
		date_added = "1971-01-02"
		description = "Test Aggregate"
		tag {
			name = netbox_extras_tag.test.name
			slug = netbox_extras_tag.test.slug
		}
		{{ end }}
	}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"prefix":         prefix,
		"extraresources": strconv.FormatBool(extraResources),
		"resourcefull":   strconv.FormatBool(resourceFull),
	}
	return util.RenderTemplate(template, data)
}
