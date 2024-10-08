package virtualization_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

const resourceNameNetboxVirtualizationVM = "netbox_virtualization_vm.test"

func TestAccNetboxVirtualizationVMMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationVMConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationVM),
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

func TestAccNetboxVirtualizationVMFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationVMConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationVM),
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

func TestAccNetboxVirtualizationVMMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationVMConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationVM),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationVMConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationVM),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationVMConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationVM),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationVMConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationVM),
				),
			},
		},
	})
}

func testAccCheckNetboxVirtualizationVMConfig(nameSuffix string, resourceFull, extraResources bool) string {
	template := `
	resource "netbox_virtualization_cluster_type" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_virtualization_cluster" "test" {
		name = "test-{{ .namesuffix }}"
		type_id = netbox_virtualization_cluster_type.test.id
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_platform" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_dcim_device_role" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_extras_tag" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_tenant" "test" {
 		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_virtualization_vm" "test" {
		name            = "test-{{ .namesuffix }}"
		cluster_id      = netbox_virtualization_cluster.test.id

		{{ if eq .resourcefull "true" }}
		comments        = <<-EOT
		VM created by terraform
		Multiline
		EOT
		role_id = netbox_dcim_device_role.test.id
		platform_id = netbox_dcim_platform.test.id
		tenant_id = netbox_tenancy_tenant.test.id
		status = "planned"
		vcpus           = 2
		disk            = 50
		memory          = 16
		local_context_data = jsonencode(
		  {
			hello = "world"
			number = 1
			bool = true
		  }
		)

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
