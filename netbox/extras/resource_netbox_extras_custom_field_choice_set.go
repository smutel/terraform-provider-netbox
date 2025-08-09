// Copyright (c)
// SPDX-License-Identifier: MIT

package extras

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxExtrasCustomFieldChoiceSet() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a custom field within Netbox. " +
			"*CAVEAT*: This module is mostly intended for testing. " +
			"Be careful when changing custom fields in production.",
		CreateContext: resourceNetboxExtrasCustomFieldChoiceSetCreate,
		ReadContext:   resourceNetboxExtrasCustomFieldChoiceSetRead,
		UpdateContext: resourceNetboxExtrasCustomFieldChoiceSetUpdate,
		DeleteContext: resourceNetboxExtrasCustomFieldChoiceSetDelete,
		Exists:        resourceNetboxExtrasCustomFieldChoiceSetExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"base_choices": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IATA",
					"ISO_3166",
					"UN_LOCODE",
				}, false),
				ExactlyOneOf: []string{"base_choices", "extra_choices"},
				Description:  "Base available choices for selection fields.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this custom field.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this custom field was created.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const200),
				Description:  "The description of this custom field.",
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Date when this custom field choice was " +
					"last updated.",
			},
			"extra_choices": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Extra available choices for selection fields.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value for this extra choice.",
						},
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Lable for this extra choice.",
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const50),
				Description:  "The name of this custom field.",
			},
			"order": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Choices are automatically ordered alphabetically",
			},
		},
	}
}

func resourceNetboxExtrasCustomFieldChoiceSetCreate(
	ctx context.Context, d *schema.ResourceData,
	m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	newResource :=
		netbox.NewWritableCustomFieldChoiceSetRequestWithDefaults()
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetName(d.Get("name").(string))
	newResource.SetOrderAlphabetically(d.Get("order").(bool))

	if _, ok := d.GetOk("extra_choices"); ok {
		newResource.SetExtraChoices(
			util.ToInterfaceArrayToInteraceArrayArray(
				d.Get("extra_choices").(*schema.Set).List()))
	} else {
		empty := make([][]any, 0)
		newResource.SetExtraChoices(empty)
	}

	if _, ok := d.GetOk("base_choices"); ok {
		b, err :=
			netbox.NewPatchedWritableCustomFieldChoiceSetRequestBaseChoicesFromValue( //nolint:revive
				d.Get("base_choices").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		newResource.SetBaseChoices(*b)
	}

	_, response, err :=
		client.ExtrasAPI.ExtrasCustomFieldChoiceSetsCreate(
			ctx).WritableCustomFieldChoiceSetRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxExtrasCustomFieldChoiceSetRead(ctx, d, m)
}

func resourceNetboxExtrasCustomFieldChoiceSetRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err :=
		client.ExtrasAPI.ExtrasCustomFieldChoiceSetsRetrieve(
			ctx, int32(resourceID)).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("base_choices",
		resource.GetBaseChoices().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type",
		util.ConvertURLContentType(resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created",
		resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("extra_choices",
		util.ToInterfaceArrayArrayToInteraceArray(
			resource.GetExtraChoices())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("order", resource.GetOrderAlphabetically()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxExtrasCustomFieldChoiceSetUpdate(
	ctx context.Context, d *schema.ResourceData,
	m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritableCustomFieldChoiceSetRequestWithDefaults()

	// Required fields
	resource.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("base_choices"); ok {
		b, err :=
			netbox.NewPatchedWritableCustomFieldChoiceSetRequestBaseChoicesFromValue( //nolint:revive
				d.Get("base_choices").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetBaseChoices(*b)
	}

	if _, ok := d.GetOk("extra_choices"); ok {
		resource.SetExtraChoices(
			util.ToInterfaceArrayToInteraceArrayArray(
				d.Get("extra_choices").(*schema.Set).List()))
	} else {
		empty := make([][]any, 0)
		resource.SetExtraChoices(empty)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("order") {
		resource.SetOrderAlphabetically(d.Get("order").(bool))
	}

	if _, response, err :=
		client.ExtrasAPI.ExtrasCustomFieldChoiceSetsUpdate(
			ctx, int32(resourceID)).WritableCustomFieldChoiceSetRequest(
			*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxExtrasCustomFieldChoiceSetRead(ctx, d, m)
}

func resourceNetboxExtrasCustomFieldChoiceSetDelete(
	ctx context.Context, d *schema.ResourceData,
	m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxExtrasCustomFieldChoiceSetExists(
		d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}

	if response, err := client.ExtrasAPI.ExtrasCustomFieldChoiceSetsDestroy(
		ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxExtrasCustomFieldChoiceSetExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.ExtrasAPI.ExtrasCustomFieldChoiceSetsRetrieve(
		nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}
	return false, err
}
