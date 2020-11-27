package netbox

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxVirtualizationVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxVirtualizationVMCreate,
		Read:   resourceNetboxVirtualizationVMRead,
		Update: resourceNetboxVirtualizationVMUpdate,
		Delete: resourceNetboxVirtualizationVMDelete,
		Exists: resourceNetboxVirtualizationVMExists,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  " ",
			},
			"disk": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"local_context_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"platform_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"offline", "active",
					"planned", "staged", "failed", "decommissioning"}, false),
			},
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tenant_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// terraform default behavior sees a difference between null and an empty string
				// therefore we override the default, because null in terraform results in empty string in netbox
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// function is called for each member of map
					// including additional call on the amount of entries
					// we ignore the count, because the actual state always returns the amount of existing custom_fields and all are optional in terraform
					if k == "custom_fields.%" {
						return true
					}
					return old == new
				},
			},
		},
	}
}

func resourceNetboxVirtualizationVMCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	clusterID := int64(d.Get("cluster_id").(int))
	comments := d.Get("comments").(string)
	disk := int64(d.Get("disk").(int))
	localContextData := d.Get("local_context_data").(string)
	memory := int64(d.Get("memory").(int))
	name := d.Get("name").(string)
	platformID := int64(d.Get("platform_id").(int))
	roleID := int64(d.Get("role_id").(int))
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int64(d.Get("tenant_id").(int))
	vcpus := int64(d.Get("vcpus").(int))
	resourceCustomFields := d.Get("custom_fields").(map[string]interface{})
	customFields := make(map[string]interface{})
	for key, value := range resourceCustomFields {
		customFields[key] = value

		// special handling for booleans, as they are the only parameter not supplied as string to netbox
		if value == "true" {
			customFields[key] = 1
		} else if value == "false" {
			customFields[key] = 0
		}
	}

	newResource := &models.WritableVirtualMachineWithConfigContext{
		Cluster:          &clusterID,
		Comments:         comments,
		LocalContextData: &localContextData,
		Name:             &name,
		Status:           status,
		Tags:             convertTagsToNestedTags(tags),
	}

	if disk != 0 {
		newResource.Disk = &disk
	}

	if memory != 0 {
		newResource.Memory = &memory
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

	if vcpus != 0 {
		newResource.Vcpus = &vcpus
	}

	newResource.CustomFields = &customFields

	resource := virtualization.NewVirtualizationVirtualMachinesCreateParams().WithData(newResource)

	resourceCreated, err := client.Virtualization.VirtualizationVirtualMachinesCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxVirtualizationVMRead(d, m)
}

func resourceNetboxVirtualizationVMRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationVirtualMachinesList(
		params, nil)

	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("cluster_id", resource.Cluster.ID); err != nil {
				return err
			}

			var comments string

			if resource.Comments == "" {
				comments = " "
			} else {
				comments = resource.Comments
			}

			if err = d.Set("comments", comments); err != nil {
				return err
			}

			if err = d.Set("disk", resource.Disk); err != nil {
				return err
			}

			if err = d.Set("local_context_data",
				*resource.LocalContextData); err != nil {
				return err
			}

			if err = d.Set("memory", resource.Memory); err != nil {
				return err
			}

			if err = d.Set("name", resource.Name); err != nil {
				return err
			}

			if resource.Platform == nil {
				if err = d.Set("platform_id", 0); err != nil {
					return err
				}
			} else {
				if err = d.Set("platform_id", resource.Platform.ID); err != nil {
					return err
				}
			}

			if resource.Role == nil {
				if err = d.Set("role_id", 0); err != nil {
					return err
				}
			} else {
				if err = d.Set("role_id", resource.Role.ID); err != nil {
					return err
				}
			}

			if err = d.Set("status", resource.Status.Value); err != nil {
				return err
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return err
			}

			if resource.Tenant == nil {
				if err = d.Set("tenant_id", 0); err != nil {
					return err
				}
			} else {
				if err = d.Set("tenant_id", resource.Tenant.ID); err != nil {
					return err
				}
			}

			if err = d.Set("vcpus", resource.Vcpus); err != nil {
				return err
			}

			// custom_fields have multiple data type returns based on field type
			// but terraform only supports map[string]string, so we convert all to strings
			customFields := make(map[string]string)
			switch t := resource.CustomFields.(type) {
			case map[string]interface{}:
				for key, value := range t {
					var strValue string
					if value != nil {
						switch v := value.(type) {
						default:
							strValue = fmt.Sprintf("%v", v)
						case map[string]interface{}:
							strValue = fmt.Sprintf("%v", v["value"])
						}
					}
					customFields[key] = strValue

				}
			}

			if err = d.Set("custom_fields", customFields); err != nil {
				return err
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxVirtualizationVMUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableVirtualMachineWithConfigContext{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name
	clusterID := int64(d.Get("cluster_id").(int))
	params.Cluster = &clusterID

	if d.HasChange("comments") {
		comments := d.Get("comments").(string)
		params.Comments = comments
	}

	if d.HasChange("disk") {
		disk := int64(d.Get("disk").(int))
		params.Disk = &disk
	}

	if d.HasChange("local_context_data") {
		localContextData := d.Get("local_context_data").(string)
		params.LocalContextData = &localContextData
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
		vcpus := int64(d.Get("vcpus").(int))
		params.Vcpus = &vcpus
	}

	//
	if d.HasChange("custom_fields") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_fields")
		customFields := make(map[string]interface{})
		// netbox needs explicit empty string to remove old values
		// first we fill all existing fields from the state with an empty string
		for key, _ := range stateCustomFields.(map[string]interface{}) {
			customFields[key] = ""
		}
		// then we override the values that still exist in the terraform code with their respective value
		for key, value := range resourceCustomFields.(map[string]interface{}) {
			customFields[key] = value

			// special handling for booleans, as they are the only parameter not supplied as string to netbox
			if value == "true" {
				customFields[key] = 1
			} else if value == "false" {
				customFields[key] = 0
			}
		}
		params.CustomFields = &customFields
	}

	resource := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(
		resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxVirtualizationVMRead(d, m)
}

func resourceNetboxVirtualizationVMDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationVMExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	p := virtualization.NewVirtualizationVirtualMachinesDeleteParams().WithID(id)
	if _, err := client.Virtualization.VirtualizationVirtualMachinesDelete(
		p, nil); err != nil {
		return err
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
