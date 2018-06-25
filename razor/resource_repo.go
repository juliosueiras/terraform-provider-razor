package razor

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/juliosueiras/razor_client/api/repo"
	"github.com/pkg/errors"
)

func resourceRepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepoCreate,
		Read:   resourceRepoRead,
		Update: resourceRepoUpdate,
		Delete: resourceRepoDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"iso_url", "no_content"},
				ForceNew:      true,
			},
			"iso_url": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"url", "no_content"},
				ForceNew:      true,
			},
			"no_content": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"iso_url", "url"},
				ForceNew:      true,
			},
			"task": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"spec": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRepoCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	newRepo := &repo.RepoRequest{}

	exist, taskErr := config.Client.Task.CheckIfTaskExist(d.Get("task").(string))

	if taskErr != nil {
		return errors.Wrap(taskErr, "issue with checking if task exist")
	}

	if !exist {
		return errors.New("there is no task with name: " + d.Get("task").(string))
	}

	newRepo.Task = d.Get("task").(string)

	if d.Get("iso_url") != nil {
		newRepo.IsoURL = d.Get("iso_url").(string)
	}

	if d.Get("no_content") != nil {
		newRepo.NoContent = d.Get("no_content").(bool)
	}

	if d.Get("url") != nil {
		newRepo.URL = d.Get("url").(string)
	}

	newRepo.Name = d.Get("name").(string)
	log.Printf("[INFO] Creating Repo: %s", newRepo.Name)

	d.SetId(newRepo.Name)
	_, err := config.Client.Repo.CreateRepo(newRepo)

	if err.Error != "" {
		return errors.New(err.Error)
	}

	return resourceRepoRead(d, meta)
}

func resourceRepoRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	repo, err := config.Client.Repo.RepoDetails(d.Id())

	if err.Error != "" {
		return errors.New(err.Error)
	}

	if repo.IsoURL == "" && repo.URL == "" {
		d.Set("no_content", true)
	}

	d.Set("iso_url", repo.IsoURL)
	d.Set("url", repo.URL)
	d.Set("spec", repo.Spec)
	d.Set("task", repo.Task.Name)

	return nil
}

func resourceRepoUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Repo Update: %s", d.Id())

	config := meta.(*Config)
	_, err := config.Client.Repo.UpdateRepoTask(d.Id(), d.Get("task").(string))

	if err.Error != "" {
		return errors.New(err.Error)
	}

	return resourceRepoRead(d, meta)
}

func resourceRepoDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Repo Delete: %s", d.Id())
	config := meta.(*Config)
	_, err := config.Client.Repo.DeleteRepo(d.Id())

	if err.Error != "" {
		return errors.New(err.Error)
	}

	return nil
}

//func connection(d *schema.ResourceData) *ThingConnection {
//	log.Printf("[INFO] setting connection: %s", d.Get("connection.0.name"))
//
//	return &ThingConnection{
//		Name: d.Get("connection.0.name").(string),
//	}
//}

//func name(params) {
//func connection(d *schema.ResourceData,) error {
//
//}
