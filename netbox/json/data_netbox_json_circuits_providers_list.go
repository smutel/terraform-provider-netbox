package json

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/circuits"
)

func DataNetboxJSONCircuitsProvidersList() *schema.Resource {
	return &schema.Resource{
		Description: "Get json output from the circuits_providers_list Netbox endpoint.",
		ReadContext: dataNetboxJSONCircuitsProvidersListRead,

		Schema: map[string]*schema.Schema{
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The max number of returned results. If 0 is specified, all records will be returned.",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON output of the list of objects for this Netbox endpoint.",
			},
		},
	}
}

func dataNetboxJSONCircuitsProvidersListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	params := circuits.NewCircuitsProvidersListParams()
	limit := int64(d.Get("limit").(int))
	params.Limit = &limit

	list, err := client.Circuits.CircuitsProvidersList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	tmp := list.Payload.Results
	resultLength := int64(len(tmp))
	desiredLength := *list.Payload.Count
	if limit > 0 && limit < desiredLength {
		desiredLength = limit
	}
	if limit == 0 || resultLength < desiredLength {
		limit = resultLength
	}
	offset := limit
	params.Offset = &offset
	for int64(len(tmp)) < desiredLength {
		offset = int64(len(tmp))
		if limit > desiredLength-offset {
			limit = desiredLength - offset
		}
		list, err = client.Circuits.CircuitsProvidersList(params, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		tmp = append(tmp, list.Payload.Results...)
	}

	j, _ := json.Marshal(tmp)

	d.Set("json", string(j))
	d.SetId("NetboxJSONCircuitsProvidersList")

	return nil
}
