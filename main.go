package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/juliosueiras/terraform-provider-razor/razor"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: razor.Provider,
	})
}
