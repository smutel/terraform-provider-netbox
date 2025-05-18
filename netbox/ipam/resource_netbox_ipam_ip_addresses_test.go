package ipam_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameIPAddress = "netbox_ipam_ip_addresses.test"

func TestAccNetboxIpamIPAddressesMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	ipNum := int64(acctest.RandIntRange(1, util.Const16384))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamIPAddressConfig(
					nameSuffix, false, false, ipNum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameIPAddress),
				),
			},
			{
				ResourceName:      resourceNameIPAddress,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamIPAddressesFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	ipNum := int64(acctest.RandIntRange(1, util.Const16384))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamIPAddressConfig(
					nameSuffix, true, true, ipNum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIPAddress),
				),
			},
			{
				ResourceName:      resourceNameIPAddress,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamIPAddressesMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	ipNum := int64(acctest.RandIntRange(1, util.Const16384))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamIPAddressConfig(
					nameSuffix, false, false, ipNum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIPAddress),
				),
			},
			{
				Config: testAccCheckNetboxIpamIPAddressConfig(
					nameSuffix, true, true, ipNum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIPAddress),
				),
			},
			{
				Config: testAccCheckNetboxIpamIPAddressConfig(
					nameSuffix, false, true, ipNum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIPAddress),
				),
			},
			{
				Config: testAccCheckNetboxIpamIPAddressConfig(
					nameSuffix, false, false, ipNum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIPAddress),
				),
			},
		},
	})
}

func testAccCheckNetboxIpamIPAddressConfig(nameSuffix string,
	resourceFull, extraResources bool, ipNum int64) string {

	const template = `
	{{ if eq .extraresources "true" }}
	resource "netbox_virtualization_cluster_type" "test" {
		name = "ipamipaddress-{{ .namesuffix }}"
		slug = "ipamipaddress-{{ .namesuffix }}"
	}

	resource "netbox_virtualization_cluster" "test" {
		name = "ipamipaddress-{{ .namesuffix }}"
		type_id = netbox_virtualization_cluster_type.test.id
	}

	resource "netbox_virtualization_vm" "test" {
		name			 = "ipamipaddress-{{ .namesuffix }}"
		cluster_id = netbox_virtualization_cluster.test.id
	}

	resource "netbox_virtualization_interface" "test" {
		name							= "ipamipaddress-{{ .namesuffix }}"
		virtualmachine_id = netbox_virtualization_vm.test.id
	}

	resource "netbox_extras_tag" "test" {
		name = "ipamipaddress-{{ .namesuffix }}"
		slug = "ipamipaddress-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_tenant" "test" {
		name = "ipamipaddress-{{ .namesuffix }}"
		slug = "ipamipaddress-{{ .namesuffix }}"
	}

	resource "netbox_ipam_ip_addresses" "nat" {
		address = "${cidrhost("10.0.0.0/8", {{ .ipnum }} + 2 )}/24"
	}

	resource "netbox_ipam_vrf" "test" {
		name = "ipamipaddress-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_ipam_ip_addresses" "test" {
		address		 = "${cidrhost("10.0.0.0/8", {{ .ipnum }})}/24"

		{{ if eq .resourcefull "true" }}
		description	 = "Test ip address"
		dns_name			= "test.example.local"
		role					= "vip"
		status				= "reserved"
		vrf_id				= netbox_ipam_vrf.test.id
		tenant_id		 = netbox_tenancy_tenant.test.id
		nat_inside_id = netbox_ipam_ip_addresses.nat.id

		object_id	 = netbox_virtualization_interface.test.id
		object_type = "virtualization.vminterface"

		tag {
			name = netbox_extras_tag.test.name
			slug = netbox_extras_tag.test.slug
		}
		{{ end }}
	}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"extraresources": strconv.FormatBool(extraResources),
		"resourcefull":   strconv.FormatBool(resourceFull),
		"ipnum":          strconv.FormatInt(ipNum, util.Const10),
	}
	return util.RenderTemplate(template, data)
}
