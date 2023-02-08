# Terraform Provider Planetscale

This is an unofficial Terraform provider for [Planetscale](https://planetscale.com/) built using the new
[Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) 🔥

## Development

This provider has been based on the awesome work Planetscale has done with their [Golang SDK](https://github.com/planetscale/planetscale-go.

Planetscale API docs can be found here: https://api-docs.planetscale.com/reference/list-databases

### Contributing

Please use this provider and test this out in all possible scenarios! The main drive behind making this Terraform 
provider was to enable the use of Planetscale in an as automated as possible manner. 

Please open [Issues](https://github.com/koslib/terraform-provider-planetscale/issues) for any problem you face or if you
got an idea for an improvement - or just need help!

Being more technical? [PRs](https://github.com/koslib/terraform-provider-planetscale/pulls) are welcome!

Interested in maintaining this Github project? [Get in touch](https://twitter.com/koslib)!

### Docs

Docs can be generated automatically with 

```
go generate ./...
``` 

Some points to remember and enhance potentially in future releases:

1. The `/examples` path contains examples for both data-sources and resources that will end up in the documentation generated.
2. The `Schema` in the resources and data-sources need to contain a `Description` so that the fields appearing in the documentation can actually make sense.
