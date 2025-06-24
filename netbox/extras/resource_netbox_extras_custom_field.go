package extras

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

const object = "object"
const multiobject = "multiobject"

func ResourceNetboxExtrasCustomField() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a custom field within Netbox. " +
			"*CAVEAT*: This module is mostly intended for testing." +
			"Be careful when changing custom fields in production.",
		CreateContext: resourceNetboxExtrasCustomFieldCreate,
		ReadContext:   resourceNetboxExtrasCustomFieldRead,
		UpdateContext: resourceNetboxExtrasCustomFieldUpdate,
		DeleteContext: resourceNetboxExtrasCustomFieldDelete,
		Exists:        resourceNetboxExtrasCustomFieldExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"choice_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Name of the extras_custom_field_choice_set" +
					"to use for this multiselect.",
			},
			"custom_fied_choices_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the custom field choices",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this custom field.",
			},
			"content_types": {
				Type:     schema.TypeSet,
				Required: true,
				Description: "Array of content types this field should be " +
					"assigned to.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					Description: "One of the content type this field " +
						"should be assigned to.",
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
				Type:     schema.TypeString,
				Optional: true,
				Description: "The default value for this custom field. " +
					"This value must be valid Json. Strings, " +
					"List and Dicts should be wrapped in jsonencode()",
				DiffSuppressFunc: func(_, oldValue, newValue string,
					d *schema.ResourceData) bool {

					if d.Get("type").(string) == object && fmt.Sprintf(
						"\"%s\"", oldValue) == newValue {
						return true
					} else if d.Get("type").(string) == multiobject {
						oldValueUpdated := strings.Replace(
							oldValue, "[", "[\"", -1)
						oldValueUpdated = strings.Replace(
							oldValueUpdated, ",", "\",\"", -1)
						oldValueUpdated = strings.Replace(
							oldValueUpdated, "]", "\"]", -1)
						if oldValueUpdated == newValue {
							return true
						}
					} else if d.Get("type").(string) != object &&
						d.Get("type").(string) != multiobject &&
						oldValue == newValue {
						return true
					}
					return false
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const200),
				Description:  "The description of this custom field.",
			},
			"filter_logic": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "loose",
				ValidateFunc: validation.StringInSlice([]string{"disabled",
					"loose", "exact"}, false),
				Description: "The filter logic for this custom field. " +
					"Allowed values: \"loose\" (default), \"exact\", " +
					"\"disabled\"",
			},
			"group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const50),
				Description: "Custom fields within the same group will be " +
					"displayed together.",
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Name of the field as displayed to users" +
					" (if not provided, the field's name will be used).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this custom field was last updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const50),
				Description:  "The name of this custom field.",
			},
			"object_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The object type for this custom field for " +
					"object/multiobject fields",
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "If true, this field is required when " +
					"creating new objects or editing an existing object.",
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
					object,
					multiobject,
				}, false),
				Description: "Type of the custom field (text, longtext, " +
					"integer, boolean, url, json, select, multiselect, " +
					"object, multiobject).",
			},
			"ui_visibility": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "always",
				ValidateFunc: validation.StringInSlice([]string{"always",
					"if-set", "hidden"}, false),
				Description: "The filter logic for this custom field. " +
					"Allowed values: \"always\" (default), \"if-set\", " +
					"\"hidden\"",
			},
			"ui_editable": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "yes",
				ValidateFunc: validation.StringInSlice([]string{"yes", "no",
					"hidden"}, false),
				Description: "The filter logic for this custom field. " +
					"Allowed values: \"yes\" (default), \"no\", \"hidden\"",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this custom field.",
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
				Type:     schema.TypeString,
				Optional: true,
				Description: "Regular expression to enforce on " +
					"text field values. " +
					"Use ^ and $ to force matching of entire string. " +
					"For example, <code>^[A-Z]{3}$</code> will limit " +
					"values to exactly three uppercase letters.",
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  util.Const100,
				Description: "Fields with higher weights appear lower " +
					"in a form.",
			},
		},
	}
}

func resourceNetboxExtrasCustomFieldCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	newResource := netbox.NewWritableCustomFieldRequestWithDefaults()
	newResource.SetObjectTypes(util.ToListofStrings(
		d.Get("content_types").(*schema.Set).List()))
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetGroupName(d.Get("group_name").(string))
	newResource.SetLabel(d.Get("label").(string))
	newResource.SetName(d.Get("name").(string))
	newResource.SetRequired(d.Get("required").(bool))
	newResource.SetValidationMaximum(
		int64(d.Get("validation_maximum").(int)))
	newResource.SetValidationMinimum(
		int64(d.Get("validation_minimum").(int)))
	newResource.SetValidationRegex(d.Get("validation_regex").(string))
	newResource.SetRelatedObjectType(d.Get("object_type").(string))

	weight32, err := safecast.ToInt32(d.Get("weight").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetWeight(weight32)

	t, err := netbox.NewPatchedWritableCustomFieldRequestTypeFromValue(
		d.Get("type").(string))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetType(*t)

	f, err :=
		netbox.NewPatchedWritableCustomFieldRequestFilterLogicFromValue(
			d.Get("filter_logic").(string))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetFilterLogic(*f)

	if visibilityAtt, ok := d.GetOk("ui_visibility"); ok {
		v, err :=
			netbox.NewPatchedWritableCustomFieldRequestUiVisibleFromValue(
				visibilityAtt.(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		newResource.SetUiVisible(*v)
	}

	e, err :=
		netbox.NewPatchedWritableCustomFieldRequestUiEditableFromValue(
			d.Get("ui_editable").(string))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetUiEditable(*e)

	if defaultstring := d.Get("default").(string); defaultstring != "" {
		var jsonMap any
		err := json.Unmarshal([]byte(defaultstring), &jsonMap)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}

		if d.Get("type").(string) == object {
			id, err := strconv.Atoi(jsonMap.(string))
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			newResource.SetDefault(id)
		} else if d.Get("type").(string) == multiobject {
			arrayIDStr := jsonMap.([]any)
			var arrayID = make([]int, len(arrayIDStr))

			for i, IDStr := range arrayIDStr {
				id, err := strconv.Atoi(IDStr.(string))
				if err != nil {
					return util.GenerateErrorMessage(nil, err)
				}

				arrayID[i] = id
			}
			newResource.SetDefault(arrayID)
		} else {
			newResource.SetDefault(jsonMap)
		}
	}

	if csn, exist := d.GetOk("choice_set_name"); exist {
		c := netbox.NewBriefCustomFieldChoiceSetRequest(csn.(string))
		newResource.SetChoiceSet(*c)
	}

	_, response, err := client.ExtrasAPI.ExtrasCustomFieldsCreate(
		ctx).WritableCustomFieldRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxExtrasCustomFieldRead(ctx, d, m)
}

func resourceNetboxExtrasCustomFieldRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err :=
		client.ExtrasAPI.ExtrasCustomFieldsRetrieve(ctx,
			int32(resourceID)).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("choice_set_name",
		resource.GetChoiceSet().Name); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type", util.ConvertURLContentType(
		resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_types", resource.GetObjectTypes()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("data_type", resource.GetDataType()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if resource.GetDefault() != nil {
		if jsonValue, err := json.Marshal(resource.GetDefault()); err != nil {
			return util.GenerateErrorMessage(nil, err)
		} else if err = d.Set("default", string(jsonValue)); err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
	} else if err = d.Set("default", nil); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("filter_logic",
		resource.GetFilterLogic().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("group_name", resource.GetGroupName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("label", resource.GetLabel()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("object_type",
		resource.GetRelatedObjectType()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("required", resource.GetRequired()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("type", resource.GetType().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("ui_visibility",
		resource.GetUiVisible().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("ui_editable",
		resource.GetUiEditable().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("validation_maximum",
		resource.GetValidationMaximum()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("validation_minimum",
		resource.GetValidationMinimum()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("validation_regex",
		resource.GetValidationRegex()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("weight", resource.GetWeight()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

//nolint:gocyclo
func resourceNetboxExtrasCustomFieldUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritableCustomFieldRequestWithDefaults()

	// Required fields
	resource.SetName(d.Get("name").(string))
	t, err :=
		netbox.NewPatchedWritableCustomFieldRequestTypeFromValue(
			d.Get("type").(string))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	resource.SetType(*t)
	resource.SetObjectTypes(util.ToListofStrings(
		d.Get("content_types").(*schema.Set).List()))

	if d.HasChange("choice_set_name") {
		if csn, exist := d.GetOk("choice_set_name"); exist {
			c := netbox.NewBriefCustomFieldChoiceSetRequest(csn.(string))
			resource.SetChoiceSet(*c)
		} else {
			resource.SetChoiceSetNil()
		}
	}

	if d.HasChange("default") {
		if defaultstring := d.Get("default").(string); defaultstring != "" {
			var jsonMap any
			err := json.Unmarshal([]byte(defaultstring), &jsonMap)
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			if d.Get("type").(string) == object {
				id, err := strconv.Atoi(jsonMap.(string))
				if err != nil {
					return util.GenerateErrorMessage(nil, err)
				}
				resource.SetDefault(id)
			} else if d.Get("type").(string) == multiobject {
				arrayIDStr := jsonMap.([]any)
				var arrayID = make([]int, len(arrayIDStr))

				for i, IDStr := range arrayIDStr {
					id, err := strconv.Atoi(IDStr.(string))
					if err != nil {
						return util.GenerateErrorMessage(nil, err)
					}

					arrayID[i] = id
				}
				resource.SetDefault(arrayID)
			} else {
				resource.SetDefault(jsonMap)
			}
		} else {
			resource.SetDefault(nil)
		}
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("filter_logic") {
		f, err :=
			netbox.NewPatchedWritableCustomFieldRequestFilterLogicFromValue(
				d.Get("filter_logic").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetFilterLogic(*f)
	}

	if d.HasChange("group_name") {
		resource.SetGroupName(d.Get("group_name").(string))
	}

	if d.HasChange("label") {
		resource.SetLabel(d.Get("label").(string))
	}

	if d.HasChange("object_type") {
		resource.SetRelatedObjectType(d.Get("object_type").(string))
	}

	if d.HasChange("required") {
		resource.SetRequired(d.Get("required").(bool))
	}

	if visibleAtt, ok := d.GetOk("ui_visibility"); ok {
		if d.HasChange("ui_visibility") {
			v, err :=
				netbox.NewPatchedWritableCustomFieldRequestUiVisibleFromValue(
					visibleAtt.(string))
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			resource.SetUiVisible(*v)
		}
	}

	if d.HasChange("ui_editable") {
		e, err :=
			netbox.NewPatchedWritableCustomFieldRequestUiEditableFromValue(
				d.Get("ui_editable").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetUiEditable(*e)
	}

	if d.HasChange("validation_maximum") {
		if vm, exist := d.GetOk("validation_maximum"); exist {
			resource.SetValidationMaximum(int64(vm.(int)))
		} else {
			resource.SetValidationMaximumNil()
		}
	}

	if d.HasChange("validation_minimum") {
		if vm, exist := d.GetOk("validation_minimum"); exist {
			resource.SetValidationMinimum(int64(vm.(int)))
		} else {
			resource.SetValidationMinimumNil()
		}
	}

	if d.HasChange("validation_regex") {
		resource.SetValidationRegex(d.Get("validation_regex").(string))
	}

	if d.HasChange("weight") {
		weight32, err := safecast.ToInt32(d.Get("weight").(int))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetWeight(weight32)
	}

	if _, response, err := client.ExtrasAPI.ExtrasCustomFieldsUpdate(
		ctx, int32(resourceID)).WritableCustomFieldRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxExtrasCustomFieldRead(ctx, d, m)
}

func resourceNetboxExtrasCustomFieldDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxExtrasCustomFieldExists(d, m)
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

	if response, err := client.ExtrasAPI.ExtrasCustomFieldsDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxExtrasCustomFieldExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.ExtrasAPI.ExtrasCustomFieldsRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
