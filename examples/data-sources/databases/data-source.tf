# List all databases
data "planetscale_databases" "all" {
  organization = "my-awesome-org"
}