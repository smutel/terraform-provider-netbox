// Copyright (c)
// SPDX-License-Identifier: MIT

package tenancy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxTenancyContactAssignment = "" +
	"netbox_tenancy_contact_assignment.test"

func TestAccNetboxTenancyContactAssignmentFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactAssignmentConfig(
					nameSuffix),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactAssignment),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyContactAssignment,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNetboxTenancyContactAssignmentConfig(
	nameSuffix string) string {
	template := `
	resource "netbox_virtualization_cluster_type" "test" {
		name = "tenancycontactassign-{{ .namesuffix }}"
		slug = "tenancycontactassign-{{ .namesuffix }}"
	}

	resource "netbox_virtualization_cluster" "test" {
		name = "tenancycontactassign-{{ .namesuffix }}"
		type_id = netbox_virtualization_cluster_type.test.id
	}

	resource "netbox_virtualization_vm" "test" {
		name						= "tenancycontactassign-{{ .namesuffix }}"
		cluster_id			= netbox_virtualization_cluster.test.id
	}

	resource "netbox_tenancy_contact" "test" {
		name				= "tenancycontactassign-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_contact_role" "test" {
		name				= "tenancycontactassign-{{ .namesuffix }}"
		slug				= "tenancycontactassign-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_contact_assignment" "test" {
		contact_id = netbox_tenancy_contact.test.id
		contact_role_id = netbox_tenancy_contact_role.test.id
		content_type = netbox_virtualization_vm.test.content_type
		object_id = netbox_virtualization_vm.test.id
		priority = "primary"
	}
	`
	data := map[string]string{
		"namesuffix": nameSuffix,
	}
	return util.RenderTemplate(template, data)
}
