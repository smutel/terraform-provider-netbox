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

const resourceNameNetboxVirtualizationClusterType = "" +
	"netbox_virtualization_cluster_type.test"

func TestAccNetboxVirtualizationClusterTypeMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationClusterTypeConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterType),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationClusterType,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationClusterTypeFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationClusterTypeConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterType),
				),
			},
			{
				ResourceName:      resourceNameNetboxVirtualizationClusterType,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVirtualizationClusterTypeMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxVirtualizationClusterTypeConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterType),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationClusterTypeConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterType),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationClusterTypeConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterType),
				),
			},
			{
				Config: testAccCheckNetboxVirtualizationClusterTypeConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxVirtualizationClusterType),
				),
			},
		},
	})
}

func testAccCheckNetboxVirtualizationClusterTypeConfig(
	nameSuffix string, resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "virtualclustertype-{{ .namesuffix }}"
		slug = "virtualclustertype-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_virtualization_cluster_type" "test" {
		name				= "virtualclustertype-{{ .namesuffix }}"
		slug				= "virtualclustertype-{{ .namesuffix }}"
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
