package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/smutel/terraform-provider-netbox/v4/netbox"
)

func main() {
	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return netbox.Provider()
		},
	}

	plugin.Serve(opts)
}
