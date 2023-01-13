package ipam

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

func ResourceNetboxIpamVlanGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a vlan group (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamVlanGroupCreate,
		ReadContext:   resourceNetboxIpamVlanGroupRead,
		UpdateContext: resourceNetboxIpamVlanGroupUpdate,
		DeleteContext: resourceNetboxIpamVlanGroupDelete,
		Exists:        resourceNetboxIpamVlanGroupExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this vlan group (ipam module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
				Description:  "The name for this vlan group (ipam module).",
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
				Description: "The slug for this vlan group (ipam module).",
			},
			"tag": &tag.TagSchema,
		},
	}
}

func resourceNetboxIpamVlanGroupCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	groupName := d.Get("name").(string)
	groupSlug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.VLANGroup{
		Name: &groupName,
		Slug: &groupSlug,
		Tags: tag.ConvertTagsToNestedTags(tags),
	}

	resource := ipam.NewIpamVlanGroupsCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamVlanGroupsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))
	return resourceNetboxIpamVlanGroupRead(ctx, d, m)
}

func resourceNetboxIpamVlanGroupRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamVlanGroupsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamVlanGroupsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
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

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamVlanGroupUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.VLANGroup{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	slug := d.Get("slug").(string)
	params.Slug = &slug

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = tag.ConvertTagsToNestedTags(tags)

	resource := ipam.NewIpamVlanGroupsPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamVlanGroupsPartialUpdate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamVlanGroupRead(ctx, d, m)
}

func resourceNetboxIpamVlanGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamVlanGroupExists(d, m)
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

	resource := ipam.NewIpamVlanGroupsDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamVlanGroupsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamVlanGroupExists(d *schema.ResourceData, m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := ipam.NewIpamVlanGroupsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamVlanGroupsList(params, nil)
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
