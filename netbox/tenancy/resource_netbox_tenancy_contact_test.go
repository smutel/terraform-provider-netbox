package tenancy_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxTenancyContact = "netbox_tenancy_contact.test"

func TestAccNetboxTenancyContactMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContact),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyContact,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyContactFull(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContact),
				),
			},
			{
				ResourceName:      resourceNameNetboxTenancyContact,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxTenancyContactMinimalFullMinimal(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxTenancyContactConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContact),
				),
			},
			{
				Config: testAccCheckNetboxTenancyContactConfig(
					nameSuffix, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContact),
				),
			},
			{
				Config: testAccCheckNetboxTenancyContactConfig(
					nameSuffix, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContact),
				),
			},
			{
				Config: testAccCheckNetboxTenancyContactConfig(
					nameSuffix, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxTenancyContact),
				),
			},
		},
	})
}

func testAccCheckNetboxTenancyContactConfig(nameSuffix string,
	resourceFull, extraResources bool) string {

	template := `
	{{ if eq .extraresources "true" }}
	resource "netbox_extras_tag" "test" {
		name = "tenancycontact-{{ .namesuffix }}"
		slug = "tenancycontact-{{ .namesuffix }}"
	}

	resource "netbox_tenancy_contact_group" "test2" {
		name				= "tenancycontact-{{ .namesuffix }}"
		slug				= "tenancycontact-{{ .namesuffix }}"
	}
	{{ end }}

	resource "netbox_tenancy_contact" "test" {
		name				= "tenancycontact-{{ .namesuffix }}"
		{{ if eq .resourcefull "true" }}
		title			 = "tenancycontact-{{ .namesuffix }}"
		phone			 = "0000000000"
		email			 = "test@test.com"
		address		 = "56 avenue Netbox"
		comments		= "tenancycontact-{{ .namesuffix }}"

		{{ if eq .extraresources "true" }}
		contact_group_id = netbox_tenancy_contact_group.test2.id
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
	}
	return util.RenderTemplate(template, data)
}
