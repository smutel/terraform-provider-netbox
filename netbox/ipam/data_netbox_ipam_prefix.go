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

func DataNetboxIpamPrefix() *schema.Resource {
	return &schema.Resource{
		Description: "Get info about prefix from netbox.",
		ReadContext: dataNetboxIpamPrefixRead,

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this prefix.",
			},
			"prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, util.Const256),
				Description: "The prefix (IP address/mask) used " +
					"for this prefix.",
			},
		},
	}
}

func dataNetboxIpamPrefixRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	prefix := []string{d.Get("prefix").(string)}

	resource, response, err :=
		client.IpamAPI.IpamPrefixesList(ctx).Prefix(prefix).Execute()

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
