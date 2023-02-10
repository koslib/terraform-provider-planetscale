# Data source for fetching database branches information

# List all database branches
data "planetscale_data_branches" "branches" {
  database = "my-database"
}