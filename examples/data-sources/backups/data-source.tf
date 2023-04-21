# Data source for fetching database branch backups information

# Example usage of the datasource for fetching all backups for a given database branch
data "planetscale_backups" "all" {
  organization = "my-org"
  database     = "my-database"
  branch       = "my-branch"
}

# Print the output of the data source
output "backups" {
  value = data.planetscale_backups.all
}