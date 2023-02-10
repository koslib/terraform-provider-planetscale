# Resource for managing database branch passwords

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

# Create a database branch password
resource "planetscale_database_branch_password" "example" {
  name         = "my-password"
  database     = planetscale_database.this.name
  organization = local.organization
  branch       = planetscale_database_branch.example.name
}