package dcim_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxDcimRegion = "netbox_dcim_region.test"

func TestAccNetboxDcimRegionMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
    acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimRegionConfig(nameSuffix,
					false, false, 0),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRegion),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimRegion,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimRegionFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
    acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, util.Const4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimRegionConfig(nameSuffix,
					true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRegion),
				),
			},
			{
				ResourceName:      resourceNameNetboxDcimRegion,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxDcimRegionMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
    acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, util.Const4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxDcimRegionConfig(nameSuffix,
					false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRegion),
				),
			},
			{
				Config: testAccCheckNetboxDcimRegionConfig(nameSuffix,
					true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRegion),
				),
			},
			{
				Config: testAccCheckNetboxDcimRegionConfig(nameSuffix,
					false, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRegion),
				),
			},
			{
				Config: testAccCheckNetboxDcimRegionConfig(nameSuffix,
					false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(resourceNameNetboxDcimRegion),
				),
			},
		},
	})
}

func testAccCheckNetboxDcimRegionConfig(nameSuffix string,
	resourceFull, extraResources bool, asn int64) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "dcimregion-{{ .namesuffix }}"
		slug = "dcimregion-{{ .namesuffix }}"
	}

		resource "netbox_dcim_region" "test2" {
		name				= "test2-{{ .namesuffix }}"
		slug				= "test2-{{ .namesuffix }}"
		}
		{{ end }}

	resource "netbox_dcim_region" "test" {
		name				= "dcimregion-{{ .namesuffix }}"
		slug				= "dcimregion-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test region"
				parent_id = netbox_dcim_region.test2.id

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
		"asn":            strconv.FormatInt(asn, util.Const10),
	}
	return util.RenderTemplate(template, data)
}
