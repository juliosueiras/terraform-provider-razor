package razor

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/juliosueiras/razor_client/api/policy"
	"github.com/pkg/errors"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Exists: resourcePolicyExists,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"broker": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"before": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"after"},
			},
			"after": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"before"},
			},
			"max_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"root_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"task": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"repo": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	newPolicy := &policy.PolicyRequest{}

	newPolicy.After = d.Get("after").(string)
	newPolicy.Before = d.Get("before").(string)
	newPolicy.Broker = d.Get("broker").(string)
	newPolicy.MaxCount = d.Get("max_count").(int)
	newPolicy.Enabled = d.Get("enabled").(bool)
	newPolicy.Hostname = d.Get("hostname").(string)
	newPolicy.NodeMetadata = d.Get("node_metadata").(map[string]interface{})
	newPolicy.Repo = d.Get("repo").(string)
	newPolicy.RootPassword = d.Get("root_password").(string)
	newPolicy.Name = d.Get("name").(string)

	tags := d.Get("tags").([]interface{})
	new_tags := make([]string, len(tags))

	for _, v := range tags {
		new_tags = append(new_tags, v.(string))
	}

	newPolicy.Tags = new_tags
	newPolicy.Task = d.Get("task").(string)

	log.Printf("[INFO] Creating Policy: %s", newPolicy.Name)

	d.SetId(newPolicy.Name)
	_, err := config.Client.Policy.CreatePolicy(newPolicy)

	if err.Error != "" {
		return errors.New(err.Error)
	}

	return resourcePolicyRead(d, meta)
}

func resourcePolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	config := meta.(*Config)

	_, err := config.Client.Policy.PolicyDetails(d.Id())
	if err.Error == "no policy matched id="+d.Id() {
		d.MarkNewResource()
		return false, nil
	}

	if err.Error != "" {
		return false, errors.New(err.Error)
	}

	return true, nil
}
func resourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	newPolicy, err := config.Client.Policy.PolicyDetails(d.Id())

	if err.Error == "no policy matched id="+d.Id() {
		d.MarkNewResource()
		return nil
	}

	if err.Error != "" {
		return errors.New(err.Error)
	}

	d.Set("broker", newPolicy.Broker.Name)
	d.Set("max_count", newPolicy.MaxCount)
	d.Set("enabled", newPolicy.Enabled)
	d.Set("hostname", newPolicy.Configuration.HostnamePattern)
	d.Set("node_metadata", newPolicy.NodeMetadata)
	d.Set("repo", newPolicy.Repo.Name)
	d.Set("root_password", newPolicy.Configuration.RootPassword)
	d.Set("tags", newPolicy.Tags)

	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Policy Update: %s", d.Id())

	config := meta.(*Config)
	_, err := config.Client.Policy.UpdatePolicyTask(d.Id(), d.Get("task").(string))

	if err.Error != "" {
		return errors.New(err.Error)
	}

	return resourcePolicyRead(d, meta)
}

func resourcePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Policy Delete: %s", d.Id())
	config := meta.(*Config)
	_, err := config.Client.Policy.DeletePolicy(d.Id())

	if err.Error != "" {
		return errors.New(err.Error)
	}

	return nil
}
