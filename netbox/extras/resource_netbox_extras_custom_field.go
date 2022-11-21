package extras

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/extras"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/util"
)

func ResourceNetboxExtrasCustomField() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a custom field (extras module) within Netbox. *CAVEAT*: This module is mostly intended for testing. Be careful when changing custom fields in production.",
		CreateContext: resourceNetboxExtrasCustomFieldCreate,
		ReadContext:   resourceNetboxExtrasCustomFieldRead,
		UpdateContext: resourceNetboxExtrasCustomFieldUpdate,
		DeleteContext: resourceNetboxExtrasCustomFieldDelete,
		Exists:        resourceNetboxExtrasCustomFieldExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"choices": {
				Type:        schema.TypeSet,
				Optional:    true,
				Default:     nil,
				Description: "Avaialbe choices for selection fields.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(1, 100),
				},
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this custom field (extras module).",
			},
			"content_types": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The content types this field should be assigned to.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this custom field was created.",
			},
			"data_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The data type of this custom field.",
			},
			"default": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default value for this custom field. Encoded as JSON.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "The description of this custom field.",
			},
			"filter_logic": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "loose",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "loose", "exact"}, false),
				Description:  "The filter logic for this custom field. Allowed values: \"loose\" (default), \"exact\", \"disabled\"",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the field as displayed to users (if not provided, the field's name will be used).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this custom field was last updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
				Description:  "The name of this custom field (extras module).",
			},
			"object_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The object type for this custom field for object/multiobject fields",
			},
			"required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, this field is required when creating new objects or editing an existing object.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"text",
					"longtext",
					"integer",
					"boolean",
					"date",
					"url",
					"json",
					"select",
					"multiselect",
					"object",
					"multiobject",
				}, false),
				Description: "Type of the custom field (text, longtext, integer, boolean, url, json, select, multiselect, object, multiobject).",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this custom field (extras module).",
			},
			"validation_maximum": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum allowed value (for numeric fields)",
			},
			"validation_minimum": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum allowed value (for numeric fields)",
			},
			"validation_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Regular expression to enforce on text field values. Use ^ and $ to force matching of entire string. For example, <code>^[A-Z]{3}$</code> will limit values to exactly three uppercase letters.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Fields with higher weights appear lower in a form.",
			},
		},
	}
}

var customFieldRequiredFields = []string{
	"created",
	"last_updated",
	"name",
	"slug",
	"tags",
	"content_types",
}

func resourceNetboxExtrasCustomFieldCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	name := d.Get("name").(string)
	defaultstring := d.Get("default").(string)
	validationMaximum := int64(d.Get("validation_maximum").(int))
	var validationMaximumPtr *int64
	if validationMaximum != 0 {
		validationMaximumPtr = &validationMaximum
	} else {
		validationMaximumPtr = nil
	}
	validationMinimum := int64(d.Get("validation_minimum").(int))
	var validationMinimumPtr *int64
	if validationMinimum != 0 {
		validationMinimumPtr = &validationMinimum
	} else {
		validationMinimumPtr = nil
	}
	weight := int64(d.Get("weight").(int))

	newResource := &models.WritableCustomField{
		Choices:           util.ToListofStrings(d.Get("choices").(*schema.Set).List()),
		ContentTypes:      util.ToListofStrings(d.Get("content_types").(*schema.Set).List()),
		Default:           &defaultstring,
		Description:       d.Get("description").(string),
		FilterLogic:       d.Get("filter_logic").(string),
		Label:             d.Get("label").(string),
		Name:              &name,
		ObjectType:        d.Get("object_type").(string),
		Required:          d.Get("required").(bool),
		Type:              d.Get("type").(string),
		ValidationMaximum: validationMaximumPtr,
		ValidationMinimum: validationMinimumPtr,
		ValidationRegex:   d.Get("validation_regex").(string),
		Weight:            &weight,
	}

	resource := extras.NewExtrasCustomFieldsCreateParams().WithData(newResource)

	resourceCreated, err := client.Extras.ExtrasCustomFieldsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxExtrasCustomFieldRead(ctx, d, m)
}

func resourceNetboxExtrasCustomFieldRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := extras.NewExtrasCustomFieldsListParams().WithID(&resourceID)
	resources, err := client.Extras.ExtrasCustomFieldsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resources.Payload.Results) != 1 {
		d.SetId("")
		return nil
	}

	resource := resources.Payload.Results[0]

	if err = d.Set("choices", resource.Choices); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("content_types", resource.ContentTypes); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("created", resource.Created.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("data_type", resource.DataType); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("default", resource.Default); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", resource.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("filter_logic", resource.FilterLogic.Value); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("label", resource.Label); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("object_type", resource.ObjectType); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("required", resource.Required); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", resource.Type.Value); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("validation_maximum", resource.ValidationMaximum); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("validation_minimum", resource.ValidationMinimum); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("validation_regex", resource.ValidationRegex); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("weight", resource.Weight); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxExtrasCustomFieldUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modifiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.WritableCustomField{}

	if d.HasChange("choices") {
		params.Choices = util.ToListofStrings(d.Get("choices").(*schema.Set).List())
	}
	if d.HasChange("content_types") {
		params.ContentTypes = util.ToListofStrings(d.Get("content_types").(*schema.Set).List())
	}
	if d.HasChange("default") {
		defaultvalue := d.Get("default").(string)
		params.Default = &defaultvalue
		modifiedFields["default"] = defaultvalue
	}
	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
		modifiedFields["description"] = params.Description
	}
	if d.HasChange("filter_logic") {
		params.FilterLogic = d.Get("filter_logic").(string)
	}
	if d.HasChange("label") {
		params.Label = d.Get("label").(string)
		modifiedFields["label"] = params.Label
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("object_type") {
		params.ObjectType = d.Get("object_type").(string)
	}
	if d.HasChange("required") {
		params.Required = d.Get("required").(bool)
		modifiedFields["required"] = params.Required
	}
	if d.HasChange("type") {
		params.Type = d.Get("type").(string)
	}
	if d.HasChange("validation_maximum") {
		validationMaximum := int64(d.Get("validation_maximum").(int))
		params.ValidationMaximum = &validationMaximum
		if d.GetRawConfig().GetAttr("validation_maximum").IsNull() {
			modifiedFields["validation_maximum"] = nil
		}
	}
	if d.HasChange("validation_minimum") {
		validationMinimum := int64(d.Get("validation_minimum").(int))
		params.ValidationMinimum = &validationMinimum
		if d.GetRawConfig().GetAttr("validation_minimum").IsNull() {
			modifiedFields["validation_minimum"] = nil
		}
	}
	if d.HasChange("validation_regex") {
		params.ValidationRegex = d.Get("validation_regex").(string)
		modifiedFields["validation_regex"] = params.ValidationRegex
	}
	if d.HasChange("weight") {
		weight := int64(d.Get("weight").(int))
		params.Weight = &weight
	}

	resource := extras.NewExtrasCustomFieldsPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Extras.ExtrasCustomFieldsPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, customFieldRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxExtrasCustomFieldRead(ctx, d, m)
}

func resourceNetboxExtrasCustomFieldDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxExtrasCustomFieldExists(d, m)
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

	resource := extras.NewExtrasCustomFieldsDeleteParams().WithID(id)
	if _, err := client.Extras.ExtrasCustomFieldsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxExtrasCustomFieldExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := extras.NewExtrasCustomFieldsListParams().WithID(&resourceID)
	resources, err := client.Extras.ExtrasCustomFieldsList(params, nil)
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
