// Copyright (c)
// SPDX-License-Identifier: MIT

package ipam

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func DataNetboxIpamIPRange() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about IP addresses from netbox.",
		ReadContext: dataNetboxIpamIPRangeRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this IP range.",
			},
			"start_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "The first address of this IP range",
			},
			"end_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "The last address of this IP range",
			},
		},
	}
}

func dataNetboxIpamIPRangeRead(ctx context.Context, d *schema.ResourceData,
	m any) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	startAddress := []string{d.Get("start_address").(string)}
	endAddress := []string{d.Get("end_address").(string)}

	resource, response, err :=
		client.IpamAPI.IpamIpRangesList(ctx).StartAddress(
			startAddress).EndAddress(endAddress).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if resource.GetCount() < 1 {
		return util.GenerateErrorMessage(nil,
			errors.New("Your query returned no results. "+
				"Please change your search criteria and try again."))

	} else if resource.GetCount() > 1 {
		return util.GenerateErrorMessage(nil,
			errors.New("Your query returned more than one result. "+
				"Please try a more specific search criteria."))
	}

	r := resource.Results[0]
	d.SetId(fmt.Sprintf("%d", r.GetId()))
	if err = d.Set("content_type",
		util.ConvertURLContentType(r.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}
