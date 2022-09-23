package netbox

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/virtualization"
	"github.com/smutel/go-netbox/v3/netbox/models"
)

func resourceNetboxVirtualizationVM() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a VM (virtualization module) resource within Netbox.",
		CreateContext: resourceNetboxVirtualizationVMCreate,
		ReadContext:   resourceNetboxVirtualizationVMRead,
		UpdateContext: resourceNetboxVirtualizationVMUpdate,
		DeleteContext: resourceNetboxVirtualizationVMDelete,
		Exists:        resourceNetboxVirtualizationVMExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the cluster which host this VM (virtualization module).",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "Comments for this VM (virtualization module).",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this VM (virtualization module).",
			},
			"custom_field": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the existing custom field.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{"text", "integer", "boolean",
								"date", "url", "selection", "multiple"}, false),
							Description: "Type of the existing custom field (text, integer, boolean, url, selection, multiple).",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the existing custom field.",
						},
					},
				},
				Description: "Existing custom fields to associate to this VM (virtualization module).",
			},
			"disk": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The size in GB of the disk for this VM (virtualization module).",
			},
			"local_context_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Local context data for this VM (virtualization module).",
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The size in MB of the memory of this VM (virtualization module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name for this VM (virtualization module).",
			},
			"platform_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the platform for this VM (virtualization module).",
			},
			"role_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the role for this VM (virtualization module).",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"offline", "active",
					"planned", "staged", "failed", "decommissioning"}, false),
				Description: "The status among offline, active, planned, staged, failed or decommissioning (active by default).",
			},
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the existing tag.",
						},
						"slug": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Slug of the existing tag.",
						},
					},
				},
				Description: "Existing tag to associate to this VM (virtualization module).",
			},
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the tenant where this VM (virtualization module) is attached.",
			},
			"vcpus": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[0-9]+((.[0-9]){0,1}[0-9]{0,1})$"),
					"Must be like ^[0-9]+((.[0-9]){0,1}[0-9]{0,1})$"),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == new+".00" || old == new {
						return true
					}
					return false
				},
				Description: "The number of VCPUS for this VM (virtualization module).",
			},
		},
	}
}

func resourceNetboxVirtualizationVMCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	clusterID := int64(d.Get("cluster_id").(int))
	comments := d.Get("comments").(string)
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	disk := int64(d.Get("disk").(int))
	localContextData := d.Get("local_context_data").(string)
	memory := int64(d.Get("memory").(int))
	name := d.Get("name").(string)
	platformID := int64(d.Get("platform_id").(int))
	roleID := int64(d.Get("role_id").(int))
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int64(d.Get("tenant_id").(int))
	vcpus := d.Get("vcpus").(string)

	if !strings.Contains(vcpus, ".") {
		vcpus = vcpus + ".00"
	}

	newResource := &models.WritableVirtualMachineWithConfigContext{
		Cluster:      &clusterID,
		Comments:     comments,
		CustomFields: &customFields,
		Name:         &name,
		Status:       status,
		Tags:         convertTagsToNestedTags(tags),
	}

	if disk != 0 {
		newResource.Disk = &disk
	}

	if memory != 0 {
		newResource.Memory = &memory
	}

	if localContextData != "" {
		var localContextDataMap map[string]*interface{}
		if err := json.Unmarshal([]byte(localContextData), &localContextDataMap); err != nil {
			return diag.FromErr(err)
		}
		newResource.LocalContextData = localContextDataMap
	}

	if platformID != 0 {
		newResource.Platform = &platformID
	}

	if roleID != 0 {
		newResource.Role = &roleID
	}

	if tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	if vcpus != "" {
		vcpusFloat, _ := strconv.ParseFloat(vcpus, 32)
		newResource.Vcpus = &vcpusFloat
	}

	resource := virtualization.NewVirtualizationVirtualMachinesCreateParams().WithData(newResource)

	resourceCreated, err := client.Virtualization.VirtualizationVirtualMachinesCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxVirtualizationVMRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationVirtualMachinesList(
		params, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("cluster_id", resource.Cluster.ID); err != nil {
				return diag.FromErr(err)
			}

			var comments interface{}
			if resource.Comments == "" {
				comments = nil
			} else {
				comments = resource.Comments
			}

			if err = d.Set("comments", comments); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return diag.FromErr(err)
			}

			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := updateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("disk", resource.Disk); err != nil {
				return diag.FromErr(err)
			}

			if resource.LocalContextData != nil {
				localContextDataJSON, err := json.Marshal(resource.LocalContextData)
				if err != nil {
					return diag.FromErr(err)
				}
				if err = d.Set("local_context_data",
					string(localContextDataJSON)); err != nil {
					return diag.FromErr(err)
				}
			} else {
				d.Set("local_context_data", "")
			}

			if err = d.Set("memory", resource.Memory); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("name", resource.Name); err != nil {
				return diag.FromErr(err)
			}

			if resource.Platform == nil {
				if err = d.Set("platform_id", 0); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("platform_id", resource.Platform.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if resource.Role == nil {
				if err = d.Set("role_id", 0); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("role_id", resource.Role.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if err = d.Set("status", resource.Status.Value); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return diag.FromErr(err)
			}

			if resource.Tenant == nil {
				if err = d.Set("tenant_id", 0); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("tenant_id", resource.Tenant.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if err = d.Set("vcpus", fmt.Sprintf("%v", *resource.Vcpus)); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxVirtualizationVMUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableVirtualMachineWithConfigContext{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name
	clusterID := int64(d.Get("cluster_id").(int))
	params.Cluster = &clusterID

	if d.HasChange("comments") {
		if comments, exist := d.GetOk("comments"); exist {
			params.Comments = comments.(string)
		} else {
			params.Comments = " "
		}
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := convertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("disk") {
		disk := int64(d.Get("disk").(int))
		params.Disk = &disk
	}

	if d.HasChange("local_context_data") {
		localContextData := d.Get("local_context_data").(string)
		var localContextDataMap map[string]*interface{}
		if localContextData == "" {
			localContextDataMap = nil
		} else if err := json.Unmarshal([]byte(localContextData), &localContextDataMap); err != nil {
			return diag.FromErr(err)
		}
		params.LocalContextData = localContextDataMap
	}

	if d.HasChange("memory") {
		memory := int64(d.Get("memory").(int))
		params.Memory = &memory
	}

	if d.HasChange("platform_id") {
		platformID := int64(d.Get("platform_id").(int))
		params.Platform = &platformID
	}

	if d.HasChange("role_id") {
		roleID := int64(d.Get("role_id").(int))
		params.Role = &roleID
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		params.Status = status
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	if d.HasChange("tenant_id") {
		tenantID := int64(d.Get("tenant_id").(int))
		params.Tenant = &tenantID
	}

	if d.HasChange("vcpus") {
		vcpus := d.Get("vcpus").(string)

		if !strings.Contains(vcpus, ".") {
			vcpus = vcpus + ".00"
		}

		vcpusFloat, _ := strconv.ParseFloat(vcpus, 32)
		params.Vcpus = &vcpusFloat
	}

	resource := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(
		resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxVirtualizationVMRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationVMExists(d, m)
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

	p := virtualization.NewVirtualizationVirtualMachinesDeleteParams().WithID(id)
	if _, err := client.Virtualization.VirtualizationVirtualMachinesDelete(
		p, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxVirtualizationVMExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationVirtualMachinesList(
		params, nil)
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
