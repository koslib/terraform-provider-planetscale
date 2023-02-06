package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
	"terraform-provider-planetscale/planetscale"
)

func main() {
	err := providerserver.Serve(context.Background(), planetscale.New, providerserver.ServeOpts{
		// todo: fix the provider address before release
		Address: "koslib.com/tf/planetscale",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
