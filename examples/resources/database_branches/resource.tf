# Resource to manage database branches

# Define local variables
locals {
  organization = "my-org"
}

# Create a database
resource "planetscale_database" "this" {
  name         = "example"
  organization = local.organization
}

# Create a database branch
resource "planetscale_database_branch" "example" {
  name         = "example"
  database     = planetscale_database.this.name
  organization = local.organization
}