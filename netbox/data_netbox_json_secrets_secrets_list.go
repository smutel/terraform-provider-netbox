package netbox

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/secrets"
)

func dataNetboxJSONSecretsSecretsList() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxJSONSecretsSecretsListRead,

		Schema: map[string]*schema.Schema{
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataNetboxJSONSecretsSecretsListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	params := secrets.NewSecretsSecretsListParams()
	limit := int64(d.Get("limit").(int))
	params.Limit = &limit

	list, err := client.Secrets.SecretsSecretsList(params, nil)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(list.Payload.Results)

	d.Set("json", string(j))
	d.SetId("NetboxJSONSecretsSecretsList")

	return nil
}
