package main

import (
	"context"
	"flag"
	"log"

	"github.com/fbritoferreira/terraform-provider-strapi/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run the docs generation tool, follow the instructions below.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name strapi

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/fbritoferreira/strapi",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
