package razor

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func dataSourceTask() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceTaskRead,
		Exists: dataSourceTaskExists,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"os": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"boot_seq": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Get("name").(string)
	task, err := config.Client.Task.TaskDetails(name)

	if err.Error != "" {
		return errors.New(err.Error)
	}

	d.SetId(name)
	d.Set("boot_seq", task.BootSeq)

	os := map[string]interface{}{
		"name":    task.OS.Name,
		"version": task.OS.Version,
	}

	d.Set("os", os)

	d.Set("description", task.Description)
	return nil
}

func dataSourceTaskExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	config := meta.(*Config)
	name := d.Get("name").(string)
	_, err := config.Client.Task.TaskDetails(name)

	if err.Error != "" {
		return false, errors.New(err.Error)
	}

	return true, nil
}
