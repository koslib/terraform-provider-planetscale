# Create a Planetscale database in the specified region and organization
resource "planetscale_database" "this" {
  organization = "my-awesome-org"
  name         = "my-awesome-db"
  region       = "AWS us-east-1"
}