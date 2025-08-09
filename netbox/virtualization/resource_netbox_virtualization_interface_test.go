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

const resourceNameNetboxVirtualizationInterface = "" +
	"netbox_virtualization_interface.test"

func TestAccNetboxVirtualizationInterfaceMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationInterfaceConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationInterface),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationInterface,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationInterfaceFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationInterfaceConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationInterface),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationInterface,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationInterfaceMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	// var resourceID *string
	// tmp := ""
	resourceID := "0"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationInterfaceConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationInterface),
					util.TestAccSaveID(
						resourceNameNetboxVirtualizationInterface, &resourceID),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationInterfaceConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationInterface),
					util.TestAccCheckID(
						resourceNameNetboxVirtualizationInterface, "id",
						&resourceID),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationInterfaceConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationInterface),
					util.TestAccCheckID(
						resourceNameNetboxVirtualizationInterface, "id",
						&resourceID),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationInterfaceConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationInterface),
					util.TestAccCheckID(
						resourceNameNetboxVirtualizationInterface, "id",
						&resourceID),
				),
			},
		},
	})
}

func testAccCheckNetboxVirtualizationInterfaceConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	resource "netbox_virtualization_cluster_type" "test" {
		name = "virtualinterface-{{ .namesuffix }}"
		slug = "virtualinterface-{{ .namesuffix }}"
	}

	resource "netbox_virtualization_cluster" "test" {
		name = "virtualinterface-{{ .namesuffix }}"
		type_id = netbox_virtualization_cluster_type.test.id
	}

	resource "netbox_virtualization_vm" "test" {
		name						= "virtualinterface-{{ .namesuffix }}"
		cluster_id			= netbox_virtualization_cluster.test.id
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_virtualization_interface" "bridge" {
		name						= "bridge-{{ .namesuffix }}"
		virtualmachine_id = netbox_virtualization_vm.test.id
	}

	resource "netbox_virtualization_interface" "parent" {
		name						= "parent-{{ .namesuffix }}"
		virtualmachine_id = netbox_virtualization_vm.test.id
	}

	resource "netbox_extras_tag" "test" {
		name = "virtualinterface-{{ .namesuffix }}"
		slug = "virtualinterface-{{ .namesuffix }}"
	}

	resource "netbox_ipam_vlan" "tagged" {
		name = "virtualinterface-{{ .namesuffix }}"
		vlan_id = 1501
	}

	resource "netbox_ipam_vlan" "untagged" {
		name = "virtualinterface-{{ .namesuffix }}"
		vlan_id = 1101
	}

	resource "netbox_ipam_vrf" "test" {
		name = "rd-{{ .namesuffix }}"
		rd	 = "rd-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_virtualization_interface" "test" {
		name						= "virtualinterface-{{ .namesuffix }}"
		virtualmachine_id = netbox_virtualization_vm.test.id

		{{ if eq .resourcefull "true" }}
		description = "Test interface"
		mac_address = "AA:AA:AA:AA:AA:AA"
		mtu = 1300
		enabled = false
		mode = "tagged"
		tagged_vlans = [
			netbox_ipam_vlan.tagged.id,
		]
		untagged_vlan = netbox_ipam_vlan.untagged.id
		parent_id = netbox_virtualization_interface.parent.id
		bridge_id = netbox_virtualization_interface.bridge.id
		vrf_id = netbox_ipam_vrf.test.id

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
	}
	return util.RenderTemplate(template, data)
}
