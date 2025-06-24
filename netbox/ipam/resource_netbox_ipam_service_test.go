package ipam_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameIpamService = "netbox_ipam_service.test"

func TestAccNetboxIpamServiceMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamServiceConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamService),
				),
			},
			{
				ResourceName:      resourceNameIpamService,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamServiceFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamServiceConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamService),
				),
			},
			{
				ResourceName:      resourceNameIpamService,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxIpamServiceMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxIpamServiceConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamService),
				),
			},
			{
				Config: testAccCheckNetboxIpamServiceConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamService),
				),
			},
			// This step is necessary. Otherwise deleting site
			// deletes the vlan groups assigned to site.
			{
				Config: testAccCheckNetboxIpamServiceConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamService),
				),
			},
			{
				Config: testAccCheckNetboxIpamServiceConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameIpamService),
				),
			},
		},
	})
}

func testAccCheckNetboxIpamServiceConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	const template = `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "ipamservice-{{ .namesuffix }}"
		slug = "ipamservice-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_virtualization_cluster_type" "test" {
		name = "ipamservice-{{ .namesuffix }}"
		slug = "ipamservice-{{ .namesuffix }}"
	}

	resource "netbox_virtualization_cluster" "test" {
		name = "ipamservice-{{ .namesuffix }}"
		type_id = netbox_virtualization_cluster_type.test.id
	}

	resource "netbox_virtualization_vm" "test" {
		name						= "ipamservice-{{ .namesuffix }}"
		cluster_id			= netbox_virtualization_cluster.test.id
	}

	resource "netbox_ipam_service" "test" {
		name				= "ipamservice-{{ .namesuffix }}"
		ports			  = [80,443]
		protocol		= "tcp"
		virtualmachine_id = netbox_virtualization_vm.test.id
		{{ if eq .resourcefull "true" }}
		description = "Test service"

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
