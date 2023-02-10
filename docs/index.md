---
page_title: "PlanetScale Provider"
subcategory: "databases"
description: |-
The PlanetScale provider provides resources and data sources for interacting with PlanetScale.

---

# PlanetScale Provider

The [PlanetScale](https://planetscale.com) provider provides resources and data sources to facilitate the use of
various PlanetScale objects using Terraform.

As PlanetScale is the most advanced database system, having an IaC way to manage object programmatically is mission
critical. PlanetScale can cover various use-cases, from simple to quite sophisticated. While it's essentially MySQL
under the hood, the fact that it is serverless offers multiple advantages, most notably the branching capabilities.

Developers can have multiple branches, eg. for preview environments, and make schema migrations, have pre-production
environments of their choice and do access control on a branch-level.

This provider supports the majority of operations currently supported by the PlanetScale [Go SDK](https://github.com/planetscale/planetscale-go).


## Example Usage

```terraform

terraform {
  required_providers {
    aws = {
      source  = "koslib/planetscale"
      version = "0.3"
    }
  }
}

# Service-token based configuration for the PlanetScale provider
provider "planetscale" {
  service_token_id = "my-token-id"
  service_token    = "my-token"
}
```

## Authentication and Configuration

Configuration for the PlanetScale Provider authentication can be derived only using Service Tokens at the moment.

To learn more about service tokens, please check out the relevant docs [here](https://planetscale.com/docs/concepts/service-tokens).

## Schema

### Required

- `service_token` (String)
- `service_token_id` (String)