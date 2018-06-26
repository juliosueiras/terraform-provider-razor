package razor

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func dataSourceNode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNodeRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"spec": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"root_password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"log": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"last_checkin": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dhcp_mac": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"hw_info": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"facts": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceNodeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Get("name").(string)
	node, err := config.Client.Node.NodeDetails(name)

	if err.Error != "" {
		return errors.New(err.Error)
	}

	hwInfo := map[string]interface{}{
		"mac": "",
	}

	if len(node.HwInfo.MAC) != 0 {
		hwInfo = map[string]interface{}{
			"mac": node.HwInfo.MAC[0],
		}
	}

	tags := make([]map[string]interface{}, len(node.Tags))

	for i := 0; i < len(node.Tags); i++ {
		for _, v := range node.Tags {
			tags[i] = map[string]interface{}{
				"id":   v.ID,
				"name": v.Name,
				"spec": v.Spec,
			}
		}
	}

	d.SetId(name)
	d.Set("hw_info", hwInfo)
	d.Set("state", node.State)
	d.Set("hostname", node.Hostname)
	d.Set("policy", node.Policy)
	d.Set("tags", tags)
	d.Set("root_password", node.RootPassword)
	d.Set("spec", node.Spec)
	d.Set("log", node.Log)
	d.Set("last_checkin", node.LastCheckin)
	d.Set("dhcp_mac", node.DHCPMAC)
	d.Set("facts", node.Facts)
	return nil
}
