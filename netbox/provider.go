package netbox

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netbox "github.com/netbox-community/go-netbox/v3"
	"github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/dcim"
)

const authHeaderName = "Authorization"
const authHeaderFormat = "Token %v"

// Provider exports the actual provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_URL", "127.0.0.1:8000"),
				Description: "URL and port to reach netbox application (127.0.0.1:8000 by default).",
			},
			"basepath": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_BASEPATH", client.DefaultBasePath),
				Description: "URL base path to the netbox API (/api by default).",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_TOKEN", ""),
				Description: "Token used for API operations (empty by default).",
			},
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_SCHEME", "https"),
				Description: "Scheme used to reach netbox application (https by default).",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETBOX_INSECURE", false),
				Description: "Skip TLS certificate validation (false by default).",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"netbox_dcim_device_role": dcim.ResourceNetboxDcimDeviceRole(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	// basepath := d.Get("basepath").(string)
	token := d.Get("token").(string)
	scheme := d.Get("scheme").(string)
	insecure := d.Get("insecure").(bool)

	fullurl := scheme + `://` + url

	cfg := netbox.NewConfiguration()
	cfg.Servers[0].URL = fullurl
	cfg.AddDefaultHeader(
		authHeaderName,
		fmt.Sprintf(authHeaderFormat, token),
	)
	cfg.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	return netbox.NewAPIClient(cfg), nil
}
