package razor

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/juliosueiras/razor_client/api"
	"log"
)

// Provider returns a schema.Provider for Example.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("RAZOR_API"),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"razor_repo":   resourceRepo(),
			"razor_tag":    resourceTag(),
			"razor_policy": resourcePolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"razor_task": dataSourceTask(),
			"razor_node": dataSourceNode(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		return nil, nil
	}
}

func envDefaultFuncAllowMissing(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		v := os.Getenv(k)
		return v, nil
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := razor.New()
	client.SetBaseURL(d.Get("base_url").(string))

	config := Config{
		BaseUrl: d.Get("base_url").(string),
		Client:  client,
	}

	log.Printf("[INFO] Razor Client configured for use")

	return &config, nil
}
