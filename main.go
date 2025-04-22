package main

import (
	"context"
	"flag"
	"log"

	"github.com/atpons/terraform-provider-cognito-idp-link/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/atpons/cognito-idp-link",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
