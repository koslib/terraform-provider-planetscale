# Resource for managing database objects

# Create a Planetscale database in the default region of your organization
resource "planetscale_database" "this" {
  organization = "my-awesome-org"
  name         = "my-awesome-db"
  region       = "eu-west"
}

# Create a Planetscale database in a specific region of your organization
resource "planetscale_database" "this" {
  organization = "my-awesome-org"
  name         = "my-awesome-db"
  region       = "eu-west"
}