// Copyright (c)
// SPDX-License-Identifier: MIT

package ipam

import (
	"context"
	"errors"
	"fmt"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func DataNetboxIpamAsn() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about aggregate from Netbox.",
		ReadContext: dataNetboxIpamAsnRead,

		Schema: map[string]*schema.Schema{
			"asn": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The asn number of this asn.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this aggregate.",
			},
			"rir_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The rir for this asn.",
			},
		},
	}
}

func dataNetboxIpamAsnRead(ctx context.Context, d *schema.ResourceData,
	m any) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	asn32, err := safecast.ToInt32(d.Get("asn").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	asn := []int32{asn32}

	rirID32, err := safecast.ToInt32(d.Get("rir").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	rirID := []int32{rirID32}

	resource, response, err :=
		client.IpamAPI.IpamAsnsList(ctx).Asn(asn).RirId(rirID).Execute()

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
