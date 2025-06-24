package tenancy_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxTenancyContactGroup = "" +
	"netbox_tenancy_contact_group.test"

func TestAccNetboxTenancyContactGroupMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactGroupConfig(
					nameSuffix, false, false, 0),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactGroup),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyContactGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyContactGroupFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, util.Const4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactGroupConfig(
					nameSuffix, true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactGroup),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyContactGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyContactGroupMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)
	asn := int64(acctest.RandIntRange(1, util.Const4294967295))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactGroupConfig(
					nameSuffix, false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactGroup),
				),
			},
			{
				Config: testAccCheckNetboxTenancyContactGroupConfig(
					nameSuffix, true, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactGroup),
				),
			},
			{
				Config: testAccCheckNetboxTenancyContactGroupConfig(
					nameSuffix, false, true, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactGroup),
				),
			},
			{
				Config: testAccCheckNetboxTenancyContactGroupConfig(
					nameSuffix, false, false, asn),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContactGroup),
				),
			},
		},
	})
}

func testAccCheckNetboxTenancyContactGroupConfig(
	nameSuffix string, resourceFull,
	extraResources bool, asn int64) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "tenancycontactgroup-{{ .namesuffix }}"
		slug = "tenancycontactgroup-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_contact_group" "test2" {
		name				= "test2-{{ .namesuffix }}"
		slug				= "test2-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_tenancy_contact_group" "test" {
		name				= "tenancycontactgroup-{{ .namesuffix }}"
		slug				= "tenancycontactgroup-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		description = "Test contact role description"
		parent_id	 = netbox_tenancy_contact_group.test2.id

		{{ if eq .extraresources "true" }}
		tag {
			name = netbox_extras_tag.test.name
			slug = netbox_extras_tag.test.slug
		}
		{{ end }}
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
