package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
	"terraform-provider-planetscale/planetscale"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name terraform-provider-planetscale

func main() {
	err := providerserver.Serve(context.Background(), planetscale.New, providerserver.ServeOpts{
		// todo: fix the provider address before release
		Address: "koslib.com/tf/planetscale",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
