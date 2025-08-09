// Copyright (c)
// SPDX-License-Identifier: MIT

package virtualization_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxVirtualizationVMPrimaryIP = "" +
	"netbox_virtualization_vm_primary_ip.test"

func TestAccNetboxVirtualizationVMPrimaryIPMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	ipnum := int64(acctest.RandIntRange(1, util.Const16384))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationVMPrimaryIPConfig(
					nameSuffix, false, false, ipnum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationVMPrimaryIP),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationVM,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationVMPrimaryIPFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	ipnum := int64(acctest.RandIntRange(1, util.Const16384))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationVMPrimaryIPConfig(
					nameSuffix, true, true, ipnum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationVMPrimaryIP),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationVM,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationVMPrimaryIPMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	ipnum := int64(acctest.RandIntRange(1, util.Const16384))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationVMPrimaryIPConfig(
					nameSuffix, false, false, ipnum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationVMPrimaryIP),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationVMPrimaryIPConfig(
					nameSuffix, true, true, ipnum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationVMPrimaryIP),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationVMPrimaryIPConfig(
					nameSuffix, false, true, ipnum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationVMPrimaryIP),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationVMPrimaryIPConfig(
					nameSuffix, false, false, ipnum),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationVMPrimaryIP),
				),
			},
		},
	})
}

func testAccCheckNetboxVirtualizationVMPrimaryIPConfig(
	nameSuffix string, resourceFull,
	extraResources bool, ipnum int64) string {

	template := `
	resource "netbox_virtualization_cluster_type" "test" {
		name = "virtualprimaryip-{{ .namesuffix }}"
		slug = "virtualprimaryip-{{ .namesuffix }}"
	}

	resource "netbox_virtualization_cluster" "test" {
		name = "virtualprimaryip-{{ .namesuffix }}"
		type_id = netbox_virtualization_cluster_type.test.id
	}

	resource "netbox_virtualization_vm" "test" {
		name						= "virtualprimaryip-{{ .namesuffix }}"
		cluster_id			= netbox_virtualization_cluster.test.id
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_virtualization_interface" "test" {
		name							= "virtualprimaryip-{{ .namesuffix }}"
		virtualmachine_id = netbox_virtualization_vm.test.id
	}

	resource "netbox_virtualization_interface" "test2" {
		name							= "test2-{{ .namesuffix }}"
		virtualmachine_id = netbox_virtualization_vm.test.id
	}

	resource "netbox_ipam_ip_addresses" "test4" {
		address		 = "${cidrhost("10.0.0.0/8", {{ .ipnum }})}/24"
		object_id	 = netbox_virtualization_interface.test.id
		object_type = "virtualization.vminterface"
	}

	resource "netbox_ipam_ip_addresses" "test6" {
		address		 = "${cidrhost("2001:db8::/32", {{ .ipnum }})}/24"
		object_id	 = netbox_virtualization_interface.test.id
		object_type = "virtualization.vminterface"
	}
	{{ end }}

	resource "netbox_virtualization_vm_primary_ip" "test" {
		virtualmachine_id = netbox_virtualization_vm.test.id
		{{ if eq .resourcefull "true" }}
		primary_ip4_id		= netbox_ipam_ip_addresses.test4.id
		primary_ip6_id		= netbox_ipam_ip_addresses.test6.id
		{{ end }}
	}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"extraresources": strconv.FormatBool(extraResources),
		"resourcefull":   strconv.FormatBool(resourceFull),
		"ipnum":          strconv.FormatInt(ipnum, util.Const10),
	}
	return util.RenderTemplate(template, data)
}
