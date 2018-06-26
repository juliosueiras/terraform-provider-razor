package razor

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/juliosueiras/razor_client/api/tag"
	"github.com/pkg/errors"
	"log"
	"reflect"
)

func resourceTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTagCreate,
		Read:   resourceTagRead,
		Update: resourceTagUpdate,
		Delete: resourceTagDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := AreEqualJSON(old, new)

					return equal
				},
			},
			"nodes": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policies": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTagCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	newTag := &tag.TagRequest{}

	keysBody := []byte(d.Get("rule").(string))
	keys := make([]interface{}, 0)
	json.Unmarshal(keysBody, &keys)

	newTag.Rule = keys

	newTag.Name = d.Get("name").(string)
	log.Printf("[INFO] Creating Tag: %s", newTag.Name)

	_, err := config.Client.Tag.CreateTag(newTag)

	if err.Error != "" {
		return errors.New(err.Error)
	}
	d.SetId(newTag.Name)

	return resourceTagRead(d, meta)
}

func resourceTagRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	tag, err := config.Client.Tag.TagDetails(d.Id())

	if err.Error != "" {
		return errors.New(err.Error)
	}

	d.Set("nodes", tag.Nodes.Count)
	d.Set("policies", tag.Policies.Count)
	d.Set("name", tag.Name)

	rules := new(bytes.Buffer)
	enc := json.NewEncoder(rules)
	enc.SetEscapeHTML(false)

	if err2 := enc.Encode(tag.Rule); err2 != nil {
		return errors.Wrap(err2, "Error parsing json")
	}

	d.Set("rule", rules.String())

	return nil
}

func resourceTagUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Tag Update: %s", d.Id())

	keysBody := []byte(d.Get("rule").(string))
	keys := make([]interface{}, 0)
	json.Unmarshal(keysBody, &keys)

	config := meta.(*Config)
	_, err := config.Client.Tag.UpdateTagRule(d.Id(), keys)

	if err.Error != "" {
		return errors.New(err.Error)
	}

	return resourceTagRead(d, meta)
}

func resourceTagDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Tag Delete: %s", d.Id())
	config := meta.(*Config)
	_, err := config.Client.Tag.DeleteTag(d.Id())

	if err.Error != "" {
		return errors.New(err.Error)
	}

	return nil
}

// reference: https://gist.github.com/turtlemonvh/e4f7404e28387fadb8ad275a99596f67
func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(o1, o2), nil
}
