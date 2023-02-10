# Data source for fetching database branch passwords information

# Fetch passwords for a single database
data "planetscale_database_branch_passwords" "branch" {
  database = "my-database"
}

# Fetch passwords for all databases in an org
data "planetscale_database_branch_passwords" "branches" {
  org = "my-org"
}

# Fetch passwords for all databases in an org, filtered by name
data "planetscale_database_branch_passwords" "branches" {
  org = "my-org"
  name = "my-branch"
}