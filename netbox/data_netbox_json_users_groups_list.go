package netbox

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/users"
)

func dataNetboxJSONUsersGroupsList() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxJSONUsersGroupsListRead,

		Schema: map[string]*schema.Schema{
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataNetboxJSONUsersGroupsListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	params := users.NewUsersGroupsListParams()
	limit := int64(d.Get("limit").(int))
	params.Limit = &limit

	list, err := client.Users.UsersGroupsList(params, nil)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(list.Payload.Results)

	d.Set("json", string(j))
	d.SetId("NetboxJSONUsersGroupsList")

	return nil
}
