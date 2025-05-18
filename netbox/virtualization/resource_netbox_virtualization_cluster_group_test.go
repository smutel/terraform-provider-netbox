package virtualization_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxVirtualizationClusterGroup = "" +
	"netbox_virtualization_cluster_group.test"

func TestAccNetboxVirtualizationClusterGroupMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationClusterGroupConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterGroup),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationClusterGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationClusterGroupFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationClusterGroupConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterGroup),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationClusterGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationClusterGroupMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationClusterGroupConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterGroup),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationClusterGroupConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterGroup),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationClusterGroupConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterGroup),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationClusterGroupConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterGroup),
				),
			},
		},
	})
}

func testAccCheckNetboxVirtualizationClusterGroupConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "virtualclustergroup-{{ .namesuffix }}"
		slug = "virtualclustergroup-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_virtualization_cluster_group" "test" {
		name				= "virtualclustergroup-{{ .namesuffix }}"
		slug				= "virtualclustergroup-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test device role"

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
