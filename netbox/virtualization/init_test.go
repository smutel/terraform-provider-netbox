package virtualization_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/smutel/terraform-provider-netbox/v4/netbox"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = netbox.Provider()
	testAccProviders = map[string]*schema.Provider{
		"netbox": testAccProvider,
	}
}
