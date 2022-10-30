package virtualization_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/util"
)

const resourceNameNetboxVirtualizationCluster = "netbox_virtualization_cluster.test"

func TestAccNetboxVirtualizationClusterMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationClusterConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationCluster),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationCluster,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationClusterFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationClusterConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationCluster),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationCluster,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationClusterMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationClusterConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationCluster),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationClusterConfig(nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationCluster),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationClusterConfig(nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationCluster),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationClusterConfig(nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxVirtualizationCluster),
				),
			},
		},
	})
}

func testAccCheckNetboxVirtualizationClusterConfig(nameSuffix string, resourceFull, extraResources bool) string {
	template := `
	resource "netbox_virtualization_cluster_type" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	{{ if eq .extraresources "true" }}
	resource "netbox_dcim_site" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_virtualization_cluster_group" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_tenant" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}

	resource "netbox_extras_tag" "test" {
		name = "test-{{ .namesuffix }}"
		slug = "test-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_virtualization_cluster" "test" {
		name      = "test-{{ .namesuffix }}"
		type_id   = netbox_virtualization_cluster_type.test.id
		{{ if eq .resourcefull "true" }}
		group_id  = netbox_virtualization_cluster_group.test.id
		site_id   = netbox_dcim_site.test.id
		tenant_id = netbox_tenancy_tenant.test.id

		comments = <<-EOT
		Test cluster
		EOT

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
	}
	return util.RenderTemplate(template, data)
}
