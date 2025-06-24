package netbox_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/smutel/terraform-provider-netbox/v8/netbox"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = netbox.Provider()
	testAccProviders = map[string]*schema.Provider{
		"netbox": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := netbox.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

//nolint:revive
func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = netbox.Provider()
}
