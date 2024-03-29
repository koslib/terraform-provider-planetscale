# Terraform Provider for Planetscale

This is an unofficial Terraform provider for [Planetscale](https://planetscale.com/) built using the new
[Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) 🔥

Provider documentation: https://registry.terraform.io/providers/koslib/planetscale/latest/docs

## Getting started

Add this provider in your terraform configuration block:

```terraform
terraform {
  required_providers {
    planetscale = {
      source = "koslib/planetscale"
      version = "~> 0.5"
    }
  }
}

# Configure the planetscale provider
provider "planetscale" {
  service_token_id = "<my-service-token-id>"
  service_token    = "<my-service-token>"
}
```

See the following section for more quick examples as well as the `/examples` directory for more detailed demonstrations.

## Examples

This provider focuses in being efficient and getting-the-job-done as easy as possible.

The following is a very simple demonstration of an example use of this provider for typical use-cases with Planetscale:

```terraform
# Create a database
resource "planetscale_database" "this" {
  organization = "my-awesome-org"
  name         = "test-from-tf"
}

# List databases in your organization
data "planetscale_databases" "all" {
  organization = "my-awesome-org"
}

# Output the values of the databases fetched in the data source above
output "list_databases" {
  value = data.planetscale_databases.all
}

# A useful data source for fetching all regions enabled for you organization
data "planetscale_regions" "all" {}

# Create a database branch
resource "planetscale_database_branch" "this" {
  organization = planetscale_database.this.organization
  database     = planetscale_database.this.name
  name         = "my-tf-branch"
}

# Create a database branch password
resource "planetscale_database_branch_password" "my-user" {
  organization = planetscale_database.this.organization
  database     = planetscale_database.this.name
  branch       = planetscale_database_branch.this.name
  name         = "my-staging-env"
}

```

More examples can be found in the `/examples` directory.

## Development

This provider has been based on the awesome work Planetscale has done with their [Golang SDK](https://github.com/planetscale/planetscale-go).

Planetscale API docs can be found [here](https://api-docs.planetscale.com/).

### Contributing

Please use this provider and test this out in all possible scenarios! The main drive behind making this Terraform 
provider was to enable the use of Planetscale in an as automated as possible manner. 

Please open [Issues](https://github.com/koslib/terraform-provider-planetscale/issues) for any problem you face or if you
got an idea for an improvement - or just need help!

Being more technical? [PRs](https://github.com/koslib/terraform-provider-planetscale/pulls) are welcome!

Interested in becoming a maintainer in this Github repository? [Get in touch](https://twitter.com/koslib)!

### Docs

Docs can be generated automatically with 

```
go generate ./...
``` 

Some points to remember and enhance potentially in future releases:

1. The `/examples` path contains examples for both data-sources and resources that will end up in the documentation generated.
2. The `Schema` in the resources and data-sources need to contain a `Description` so that the fields appearing in the documentation can actually make sense.

## Known limitations/issues

1. Resources updates: the Planetscale Golang SDK, on which this Terraform provider heavily relies on, does not support update operations everywhere. This means configuration of resources is not always successful.
2. Data sources filtering: the filters supported are the filters supported by the Planetscale Golang SDK. More filters will be added as soon as the SDK offers support for them.
3. No `import` functionality yet.

## Licence

MIT License

Copyright (c) 2023 Konstantinos Livieratos

See the full licence [here](https://github.com/koslib/terraform-provider-planetscale/blob/main/LICENSE.md).
