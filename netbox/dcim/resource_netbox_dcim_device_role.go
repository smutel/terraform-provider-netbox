// Copyright (c)
// SPDX-License-Identifier: MIT

package dcim

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxDcimDeviceRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a device role within Netbox.",
		CreateContext: resourceNetboxDcimDeviceRoleCreate,
		ReadContext:   resourceNetboxDcimDeviceRoleRead,
		UpdateContext: resourceNetboxDcimDeviceRoleUpdate,
		DeleteContext: resourceNetboxDcimDeviceRoleDelete,
		Exists:        resourceNetboxDcimDeviceRoleExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "9e9e9e",
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 6),
					validation.StringMatch(
						regexp.MustCompile("^[0-9a-f]{1,6}$"),
						"^[0-9a-f]{1,6})$")),
				Description: "The color of this device role. " +
					"Default is grey (#9e9e9e).",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this device role.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this device role was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const200),
				Description:  "The description of this device role.",
			},
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of devices with this device role.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this device role was last updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The name of this device role.",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The slug of this device role.",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this device role.",
			},
			"virtualmachine_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "The number of virtual machines with " +
					"this device role.",
			},
			"vm_role": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow this device role for virtual machines",
			},
		},
	}
}

func resourceNetboxDcimDeviceRoleCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	color := d.Get("color").(string)
	vmRole := d.Get("vm_role").(bool)
	description := d.Get("description").(string)
	tags := d.Get("tag").(*schema.Set).List()
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields :=
		customfield.ConvertCustomFieldsFromTerraformToAPI(
			nil, resourceCustomFields)

	newResource := netbox.NewDeviceRoleRequestWithDefaults()
	newResource.SetName(name)
	newResource.SetSlug(slug)
	newResource.SetColor(color)
	newResource.SetVmRole(vmRole)
	newResource.SetDescription(description)
	newResource.SetCustomFields(customFields)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	_, response, err := client.DcimAPI.DcimDeviceRolesCreate(
		ctx).DeviceRoleRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxDcimDeviceRoleRead(ctx, d, m)
}

func resourceNetboxDcimDeviceRoleRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.DcimAPI.DcimDeviceRolesRetrieve(
		ctx, int32(resourceID)).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("color", resource.GetColor()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type", util.ConvertURLContentType(
		resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(
		resourceCustomFields, resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("device_count", resource.GetDeviceCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("slug", resource.GetSlug()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(
		resource.GetTags())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("virtualmachine_count",
		resource.GetVirtualmachineCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("vm_role", resource.GetVmRole()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxDcimDeviceRoleUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewDeviceRoleRequestWithDefaults()

	// Required fields
	resource.SetName(d.Get("name").(string))
	resource.SetSlug(d.Get("slug").(string))

	if d.HasChange("color") {
		resource.SetColor(d.Get("color").(string))
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
			stateCustomFields.(*schema.Set).List(),
			resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("vm_role") {
		resource.SetVmRole(d.Get("vm_role").(bool))
	}

	if _, response, err := client.DcimAPI.DcimDeviceRolesUpdate(
		ctx, int32(resourceID)).DeviceRoleRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxDcimDeviceRoleRead(ctx, d, m)
}

func resourceNetboxDcimDeviceRoleDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxDcimDeviceRoleExists(d, m)
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

	if response, err := client.DcimAPI.DcimDeviceRolesDestroy(
		ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxDcimDeviceRoleExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.DcimAPI.DcimDeviceRolesRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
