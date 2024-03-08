// Code generated by util/generateJsonDatasources; DO NOT EDIT.
package json

import (
	"context"
	"encoding/json"
  "errors"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netbox "github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

// This file was generated by the util/generateJsonDatasources.
// Editing this file might prove futile when you re-run the util/generateJsonDatasources command

func DataNetboxJSONIpamVrfsList() *schema.Resource {
	return &schema.Resource{
		Description: "Get json output from the IpamVrfsList Netbox endpoint.",
		ReadContext: dataNetboxJSONIpamVrfsListRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filter the records returned by the query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the field to use for filtering.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the field to use for filtering.",
						},
					},
				},
			},
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

func dataNetboxJSONIpamVrfsListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

  limit := int32(d.Get("limit").(int))
  request := client.IpamAPI.IpamVrfsList(ctx)

	if filter, ok := d.GetOk("filter"); ok {
		var filterParams = filter.(*schema.Set)
		for _, f := range filterParams.List() {
			k := f.(map[string]interface{})["name"]
			v := f.(map[string]interface{})["value"]
			kString := k.(string)
			vString := v.(string)
      method := reflect.ValueOf(request).MethodByName("Set" + kString).Interface().(func(string))
      method(vString)
		}
	}

	resource, response, err := request.Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if resource.GetCount() < 1 {
		return util.GenerateErrorMessage(nil, errors.New("Your query returned no results. "+
			"Please change your search criteria and try again."))
	}

  tmp := resource.Results
	resultLength := int32(len(tmp))
	desiredLength := *resource.Count
	if limit > 0 && limit < desiredLength {
		desiredLength = limit
	}
	if limit == 0 || resultLength < desiredLength {
		limit = resultLength
	}
	offset := limit
  request = request.Offset(offset)

	for int32(len(tmp)) < desiredLength {
		offset = int32(len(tmp))
		if limit > desiredLength-offset {
			limit = desiredLength - offset
		}
    request.Execute()
		if err != nil {
		  return util.GenerateErrorMessage(response, err)
		}
		tmp = append(tmp, resource.Results...)
	}

	j, _ := json.Marshal(tmp)

	if err = d.Set("json", string(j)); err != nil {
		return util.GenerateErrorMessage(response, err)
	}
	d.SetId("NetboxJSONIpamVrfsList")

	return nil
}
