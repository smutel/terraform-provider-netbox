// Copyright (c)
// SPDX-License-Identifier: MIT

package ipam_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameIpamVlan = "netbox_ipam_vlan.test"

func TestAccNetboxIpamVlanMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamVlanConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVlan),
				),
			},
			{
				ResourceName:      resourceNameIpamVlan,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamVlanFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamVlanConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVlan),
				),
			},
			{
				ResourceName:      resourceNameIpamVlan,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamVlanMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamVlanConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVlan),
				),
			},
			{
				Config: testAccCheckNetboxIpamVlanConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVlan),
				),
			},
			// This step is necessary. Otherwise deleting site
			// deletes the vlan groups assigned to site.
			{
				Config: testAccCheckNetboxIpamVlanConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVlan),
				),
			},
			{
				Config: testAccCheckNetboxIpamVlanConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamVlan),
				),
			},
		},
	})
}

func testAccCheckNetboxIpamVlanConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	const template = `
	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test" {
		name = "ipamvlan-{{ .namesuffix }}"
		slug = "ipamvlan-{{ .namesuffix }}"
	}

	resource "netbox_extras_tag" "test" {
		name = "ipamvlan-{{ .namesuffix }}"
		slug = "ipamvlan-{{ .namesuffix }}"
	}

	resource "netbox_ipam_vlan_group" "test" {
		name				= "ipamvlan-{{ .namesuffix }}"
		slug				= "ipamvlan-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_tenant" "test" {
		 name = "ipamvlan-{{ .namesuffix }}"
		slug = "ipamvlan-{{ .namesuffix }}"
	}

	resource "netbox_ipam_role" "test" {
		 name = "ipamvlan-{{ .namesuffix }}"
		slug = "ipamvlan-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_ipam_vlan" "test" {
		name				= "ipamvlan-{{ .namesuffix }}"
		vlan_id		 = 10
		{{ if eq .resourcefull "true" }}
		description = "Test Vlan group"
		vlan_group_id = netbox_ipam_vlan_group.test.id
		role_id = netbox_ipam_role.test.id
		// site_id = netbox_dcim_site.test.id
		status = "reserved"
		tenant_id = netbox_tenancy_tenant.test.id


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
