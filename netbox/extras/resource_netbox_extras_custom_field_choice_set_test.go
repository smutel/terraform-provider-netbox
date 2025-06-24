package extras_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const resourceNameNetboxExtrasCustomFieldChoiceSet = "" +
	"netbox_extras_custom_field_choice_set.test"

func TestAccNetboxExtrasCustomFieldChoiceSetMinimalBase(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldChoiceSet,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldChoiceSetFullBase(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldChoiceSet,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldChoiceSetMinimalFullMinimalBase(
	t *testing.T) {

	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, false, true, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldChoiceSetMinimalExtra(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldChoiceSet,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldChoiceSetFullExtra(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				ResourceName:      resourceNameNetboxExtrasCustomFieldChoiceSet,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxExtrasCustomFieldChoiceSetMinimalFullMinimalExtra(
	t *testing.T) {

	nameSuffix := acctest.RandStringFromCharSet(util.Const10,
		acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { util.TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
			{
				Config: testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
					nameSuffix, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					util.TestAccResourceExists(
						resourceNameNetboxExtrasCustomFieldChoiceSet),
				),
			},
		},
	})
}

func testAccCheckNetboxExtrasCustomFieldChoiceSetConfig(
	nameSuffix string, resourceFull,
	extraResources, baseChoices bool) string {

	template := `
	resource "netbox_extras_custom_field_choice_set" "test" {
		name = "extrascfchoiceset_{{ .namesuffix }}"
		{{ if eq .base_choices "true" }}
		base_choices = "IATA"
		{{ else }}

		extra_choices {
			value = "zz_extra_choice_1"
			label = "zz_Extra choice 1"
		}

		extra_choices {
			value = "extra_choice_1"
			label = "Extra choice 1"
		}
		{{ end }}

		{{ if eq .resourcefull "true" }}
		description	 = "Test custom field"
		order				 = true
		{{ end }}
	}
	`
	data := map[string]string{
		"namesuffix":     nameSuffix,
		"resourcefull":   strconv.FormatBool(resourceFull),
		"extraresources": strconv.FormatBool(extraResources),
		"base_choices":   strconv.FormatBool(baseChoices),
	}
	return util.RenderTemplate(template, data)
}
