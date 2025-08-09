// Copyright (c)
// SPDX-License-Identifier: MIT

package tenancy

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxTenancyContact() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a contact within Netbox.",
		CreateContext: resourceNetboxTenancyContactCreate,
		ReadContext:   resourceNetboxTenancyContactRead,
		UpdateContext: resourceNetboxTenancyContactUpdate,
		DeleteContext: resourceNetboxTenancyContactDelete,
		Exists:        resourceNetboxTenancyContactExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const200),
				Description:  "The address for this contact.",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "Comments for this contact.",
			},
			"contact_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the group where this contact belongs to.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this contact.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this contact was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				ValidateFunc: func(val any, key string) (warns []string,
					errs []error) {
					email := val.(string)
					if len(email) >= util.Const254 {
						errs = append(errs,
							fmt.Errorf("Length of %q must be lower than 254, "+
								"got: %d", key, len(email)))
					} else if !strfmt.IsEmail(email) {
						errs = append(errs, fmt.Errorf("%q is not a valid "+
							"email", key))
					}
					return warns, errs
				},
				Description: "The e-mail for this contact.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last date when this contact was updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The name for this contact.",
			},
			"phone": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const50),
				Description:  "The phone for this contact.",
			},
			"tag": &tag.TagSchema,
			"title": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The title for this contact.",
			},
		},
	}
}

func resourceNetboxTenancyContactCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	address := d.Get("address").(string)
	comments := d.Get("comments").(string)
	groupID := d.Get("contact_group_id").(int)
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	email := d.Get("email").(string)
	name := d.Get("name").(string)
	phone := d.Get("phone").(string)
	tags := d.Get("tag").(*schema.Set).List()
	title := d.Get("title").(string)

	newResource := netbox.NewContactRequestWithDefaults()
	newResource.SetAddress(address)
	newResource.SetComments(comments)
	newResource.SetCustomFields(customFields)
	newResource.SetEmail(email)
	newResource.SetName(name)
	newResource.SetPhone(phone)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	newResource.SetTitle(title)

	if groupID != 0 {
		b, err := brief.GetBriefContactGroupRequestFromID(ctx, client, groupID)
		if err != nil {
			return err
		}
		newResource.SetGroup(*b)
	}

	_, response, err :=
		client.TenancyAPI.TenancyContactsCreate(ctx).ContactRequest(
			*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxTenancyContactRead(ctx, d, m)
}

func resourceNetboxTenancyContactRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.TenancyAPI.TenancyContactsRetrieve(ctx,
		int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields,
		resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("address", resource.GetAddress()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("comments", resource.GetComments()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("contact_group_id", resource.GetGroup().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type", util.ConvertURLContentType(
		resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("email", resource.GetEmail()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("phone", resource.GetPhone()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag",
		tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("title", resource.GetTitle()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxTenancyContactUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewContactRequestWithDefaults()

	// Required parameters
	name := d.Get("name").(string)
	resource.SetName(name)

	if d.HasChange("address") {
		resource.SetAddress(d.Get("address").(string))
	}

	if d.HasChange("comments") {
		resource.SetComments(d.Get("comments").(string))
	}

	if d.HasChange("contact_group_id") {
		groupID := d.Get("contact_group_id").(int)
		if groupID == 0 {
			resource.SetGroupNil()
		} else {
			b, err := brief.GetBriefContactGroupRequestFromID(ctx, client,
				groupID)
			if err != nil {
				return err
			}
			resource.SetGroup(*b)
		}
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields :=
			customfield.ConvertCustomFieldsFromTerraformToAPI(
				stateCustomFields.(*schema.Set).List(),
				resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("email") {
		resource.SetEmail(d.Get("email").(string))
	}

	if d.HasChange("phone") {
		resource.SetPhone(d.Get("phone").(string))
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("title") {
		resource.SetTitle(d.Get("title").(string))
	}

	if _, response, err := client.TenancyAPI.TenancyContactsUpdate(ctx,
		int32(resourceID)).ContactRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxTenancyContactRead(ctx, d, m)
}

func resourceNetboxTenancyContactDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxTenancyContactExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}

	if response, err := client.TenancyAPI.TenancyContactsDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxTenancyContactExists(d *schema.ResourceData,
	m any) (b bool,
	e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.TenancyAPI.TenancyContactsRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
