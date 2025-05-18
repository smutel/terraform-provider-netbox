package virtualization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxVirtualizationVM() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a VM resource within Netbox.",
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
				Optional:    true,
				Default:     nil,
				Description: "ID of the cluster which host this VM.",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				StateFunc:   util.TrimString,
				Description: "Comments for this VM.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this VM.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this VM was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"device_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
				Description: "Optionally pin this VM to a specific host " +
					"device within the cluster.",
			},
			"disk": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "The size in GB of the disk for this VM.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last date when this VM was updated.",
			},
			"local_context_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Local context data for this VM.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "The size in MB of the memory of this VM.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const64),
				Description:  "The name for this VM.",
			},
			"platform_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the platform for this VM.",
			},
			"primary_ip": {
				Type:     schema.TypeString,
				Computed: true,
				Default:  nil,
				Description: "Primary IP of this VM. Can be IPv4 or IPv6. " +
					"See [Netbox docs|https://docs.netbox.dev/en/stable/" +
					"models/virtualization/virtualmachine/] for more " +
					"information.",
			},
			"primary_ip4": {
				Type:        schema.TypeString,
				Computed:    true,
				Default:     nil,
				Description: "Primary IPv4 of this VM.",
			},
			"primary_ip6": {
				Type:        schema.TypeString,
				Computed:    true,
				Default:     nil,
				Description: "Primary IPv6 of this VM.",
			},
			"role_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the role for this VM.",
			},
			"site_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "ID of the site where this VM is attached. " +
					"If cluster_id is set and the cluster resides in a site, " +
					"this must be set and the same as the cluster's site",
				AtLeastOneOf: []string{"cluster_id", "site_id"},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"offline",
					"active", "planned", "staged", "failed", "decommissioning"},
					false),
				Description: "The status among offline, active, planned, " +
					"staged, failed or decommissioning (active by default).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the tenant where this VM is attached.",
			},
			"vcpus": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[0-9]+((.[0-9]){0,1}[0-9]{0,1})$"),
					"Must be like ^[0-9]+((.[0-9]){0,1}[0-9]{0,1})$"),
				//nolint:revive
				DiffSuppressFunc: func(k, old, new string,
					d *schema.ResourceData) bool {
					if old == new+".00" || old == new {
						return true
					}
					return false
				},
				Description: "The number of VCPUS for this VM.",
			},
		},
	}
}

func resourceNetboxVirtualizationVMCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	clusterID := d.Get("cluster_id").(int)
	comments := d.Get("comments").(string)
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	name := d.Get("name").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource :=
		netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()
	newResource.SetComments(comments)
	newResource.SetCustomFields(customFields)
	newResource.SetName(name)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	b, err := brief.GetBriefClusterRequestFromID(ctx, client, clusterID)
	if err != nil {
		return err
	}
	newResource.SetCluster(*b)

	s, errDiag := netbox.NewModuleStatusValueFromValue(
		d.Get("status").(string))
	if errDiag != nil {
		return util.GenerateErrorMessage(nil, errDiag)
	}
	newResource.SetStatus(*s)

	if disk := d.Get("disk").(int); disk != 0 {
		disk32, err := safecast.ToInt32(disk)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		newResource.SetDisk(disk32)
	}

	if memory := d.Get("memory").(int); memory != 0 {
		memory32, err := safecast.ToInt32(memory)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		newResource.SetMemory(memory32)
	}

	if localContextData :=
		d.Get("local_context_data").(string); localContextData != "" {
		var localContextDataMap map[string]*any
		if err := json.Unmarshal([]byte(localContextData),
			&localContextDataMap); err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		newResource.SetLocalContextData(localContextDataMap)
	}

	if deviceID := d.Get("device_id").(int); deviceID != 0 {
		b, err := brief.GetBriefDeviceRequestFromID(ctx, client, deviceID)
		if err != nil {
			return err
		}
		newResource.SetDevice(*b)
	}

	if platformID := d.Get("platform_id").(int); platformID != 0 {
		b, err := brief.GetBriefPlatformRequestFromID(ctx, client, platformID)
		if err != nil {
			return err
		}
		newResource.SetPlatform(*b)
	}

	if roleID := d.Get("role_id").(int); roleID != 0 {
		b, err := brief.GetBriefDeviceRoleRequestFromID(ctx, client, roleID)
		if err != nil {
			return err
		}
		newResource.SetRole(*b)
	}

	if siteID := d.Get("site_id").(int); siteID != 0 {
		b, err := brief.GetBriefSiteRequestFromID(ctx, client, siteID)
		if err != nil {
			return err
		}
		newResource.SetSite(*b)
	}

	if tenantID := d.Get("tenant_id").(int); tenantID != 0 {
		b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
		if err != nil {
			return err
		}
		newResource.SetTenant(*b)
	}

	if vcpus := d.Get("vcpus").(string); vcpus != "" {
		vcpusFloat, _ := strconv.ParseFloat(vcpus, util.Const32)
		newResource.SetVcpus(vcpusFloat)
	}

	_, response, errDiag :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesCreate(
			ctx).WritableVirtualMachineWithConfigContextRequest(
			*newResource).Execute()
	if response.StatusCode != util.Const201 && errDiag != nil {
		return util.GenerateErrorMessage(response, errDiag)
	}

	var resourceID int32
	if resourceID, errDiag = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, errDiag)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxVirtualizationVMRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(ctx,
			int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("cluster_id", resource.GetCluster().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("comments", resource.GetComments()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type",
		util.ConvertURLContentType(resource.GetUrl())); err != nil {
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

	if err = d.Set("device_id", resource.GetDevice().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("disk", resource.GetDisk()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	localContextDataJSON, err :=
		util.GetLocalContextData(resource.GetLocalContextData())
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("local_context_data", localContextDataJSON); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("memory", resource.GetMemory()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("platform_id", resource.GetPlatform().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("primary_ip",
		resource.GetPrimaryIp().Address); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("primary_ip4",
		resource.GetPrimaryIp4().Address); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("primary_ip6",
		resource.GetPrimaryIp6().Address); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("role_id",
		resource.GetRole().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("site_id", resource.GetSite().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("status", resource.GetStatus().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag",
		tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tenant_id", resource.GetTenant().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	vcpus := resource.GetVcpus()
	if vcpus != 0 {
		vcpusString := strconv.FormatFloat(vcpus, 'f', -1, util.Const64)
		if err = d.Set("vcpus", vcpusString); err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
	}

	return nil
}

//nolint:gocyclo
func resourceNetboxVirtualizationVMUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}
	resource :=
		netbox.NewWritableVirtualMachineWithConfigContextRequestWithDefaults()

	// Required parameters
	resource.SetName(d.Get("name").(string))

	if clusterID := d.Get("cluster_id").(int); clusterID != 0 {
		b, err := brief.GetBriefClusterRequestFromID(ctx, client, clusterID)
		if err != nil {
			return err
		}
		resource.SetCluster(*b)
	} else {
		resource.SetClusterNil()
	}

	if siteID := d.Get("site_id").(int); siteID != 0 {
		b, err := brief.GetBriefSiteRequestFromID(ctx, client, siteID)
		if err != nil {
			return err
		}
		resource.SetSite(*b)
	} else {
		resource.SetSiteNil()
	}

	if d.HasChange("comments") {
		resource.SetComments(d.Get("comments").(string))
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
			stateCustomFields.(*schema.Set).List(),
			resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("disk") {
		if disk := d.Get("disk").(int); disk != 0 {
			disk32, err := safecast.ToInt32(disk)
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			resource.SetDisk(disk32)
		} else {
			resource.SetDiskNil()
		}
	}

	if d.HasChange("local_context_data") {
		localContextData := d.Get("local_context_data").(string)
		var localContextDataMap map[string]*any
		if localContextData == "" {
			localContextDataMap = nil
		} else if err := json.Unmarshal([]byte(localContextData),
			&localContextDataMap); err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetLocalContextData(localContextDataMap)
	}

	if d.HasChange("memory") {
		if memory := d.Get("memory").(int); memory != 0 {
			memory32, err := safecast.ToInt32(memory)
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			resource.SetMemory(memory32)
		} else {
			resource.SetMemoryNil()
		}
	}

	if d.HasChange("platform_id") {
		if platformID := d.Get("platform_id").(int); platformID != 0 {
			b, err := brief.GetBriefPlatformRequestFromID(
				ctx, client, platformID)
			if err != nil {
				return err
			}
			resource.SetPlatform(*b)
		} else {
			resource.SetPlatformNil()
		}
	}

	if d.HasChange("role_id") {
		if roleID := d.Get("role_id").(int); roleID != 0 {
			b, err := brief.GetBriefDeviceRoleRequestFromID(ctx, client, roleID)
			if err != nil {
				return err
			}
			resource.SetRole(*b)
		} else {
			resource.SetRoleNil()
		}
	}

	if d.HasChange("status") {
		s, err := netbox.NewModuleStatusValueFromValue(d.Get("status").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetStatus(*s)
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("tenant_id") {
		if tenantID := d.Get("tenant_id").(int); tenantID != 0 {
			b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
			if err != nil {
				return err
			}
			resource.SetTenant(*b)
		} else {
			resource.SetTenantNil()
		}
	}

	if d.HasChange("vcpus") {
		vcpus := d.Get("vcpus").(string)
		if vcpus != "" {
			vcpusFloat, _ := strconv.ParseFloat(vcpus, util.Const32)
			resource.SetVcpus(vcpusFloat)
		} else {
			resource.SetVcpusNil()
		}
	}

	if _, response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx,
			int32(resourceID)).WritableVirtualMachineWithConfigContextRequest(
			*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxVirtualizationVMRead(ctx, d, m)
}

func resourceNetboxVirtualizationVMDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxVirtualizationVMExists(d, m)
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

	if response, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesDestroy(ctx,
			int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxVirtualizationVMExists(d *schema.ResourceData,
	m any) (b bool,
	e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err :=
		client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(nil,
			int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
