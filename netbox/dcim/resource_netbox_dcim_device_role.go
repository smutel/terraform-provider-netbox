package dcim

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/dcim"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

func ResourceNetboxDcimDeviceRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a device role (dcim module) within Netbox.",
		CreateContext: resourceNetboxDcimDeviceRoleCreate,
		ReadContext:   resourceNetboxDcimDeviceRoleRead,
		UpdateContext: resourceNetboxDcimDeviceRoleUpdate,
		DeleteContext: resourceNetboxDcimDeviceRoleDelete,
		Exists:        resourceNetboxDcimDeviceRoleExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this device role (dcim module).",
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "9e9e9e",
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 6),
					validation.StringMatch(
						regexp.MustCompile("^[0-9a-f]{1,6}$"),
						"^[0-9a-f]{1,6})$")),
				Description: "The color of this device role. Default is grey (#9e9e9e).",
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
				ValidateFunc: validation.StringLenBetween(1, 200),
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
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this device role (dcim module).",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The slug of this device role (dcim module).",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this device role (dcim module).",
			},
			"virtualmachine_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of virtual machines with this device role.",
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

var deviceRoleRequiredFields = []string{
	"created",
	"last_updated",
	"name",
	"slug",
	"tags",
}

func resourceNetboxDcimDeviceRoleCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.DeviceRole{
		Color:        d.Get("color").(string),
		CustomFields: customFields,
		Description:  d.Get("description").(string),
		Name:         &name,
		Slug:         &slug,
		Tags:         tag.ConvertTagsToNestedTags(tags),
		VMRole:       d.Get("vm_role").(bool),
	}

	dropFields := []string{}
	emptyFields := make(map[string]interface{})
	if !newResource.VMRole {
		emptyFields["vm_role"] = false
	}

	resource := dcim.NewDcimDeviceRolesCreateParams().WithData(newResource)

	resourceCreated, err := client.Dcim.DcimDeviceRolesCreate(resource, nil, requestmodifier.NewRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxDcimDeviceRoleRead(ctx, d, m)
}

func resourceNetboxDcimDeviceRoleRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := dcim.NewDcimDeviceRolesListParams().WithID(&resourceID)
	resources, err := client.Dcim.DcimDeviceRolesList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resources.Payload.Results) != 1 {
		d.SetId("")
		return nil
	}

	resource := resources.Payload.Results[0]

	if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("color", resource.Color); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("created", resource.Created.String()); err != nil {
		return diag.FromErr(err)
	}
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

	if err = d.Set("custom_field", customFields); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", resource.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("device_count", resource.DeviceCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("slug", resource.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("virtualmachine_count", resource.VirtualmachineCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("vm_role", resource.VMRole); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxDcimDeviceRoleUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modifiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	params := &models.DeviceRole{}

	if d.HasChange("color") {
		params.Color = d.Get("color").(string)
	}
	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
		modifiedFields["description"] = params.Description
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("slug") {
		slug := d.Get("slug").(string)
		params.Slug = &slug
	}
	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	}
	if d.HasChange("vm_role") {
		params.VMRole = d.Get("vm_role").(bool)
		modifiedFields["vm_role"] = d.Get("vm_role").(bool)
	}

	resource := dcim.NewDcimDeviceRolesPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Dcim.DcimDeviceRolesPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, deviceRoleRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxDcimDeviceRoleRead(ctx, d, m)
}

func resourceNetboxDcimDeviceRoleDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxDcimDeviceRoleExists(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource := dcim.NewDcimDeviceRolesDeleteParams().WithID(id)
	if _, err := client.Dcim.DcimDeviceRolesDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxDcimDeviceRoleExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := dcim.NewDcimDeviceRolesListParams().WithID(&resourceID)
	resources, err := client.Dcim.DcimDeviceRolesList(params, nil)
	if err != nil {
		return resourceExist, err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			resourceExist = true
		}
	}

	return resourceExist, nil
}
