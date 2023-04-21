# Demonstrate how to create a backup of a database branch

# Create a Planetscale database in the default region of your organization
resource "planetscale_database" "this" {
  organization = "my-awesome-org"
  name         = "my-awesome-db"
  region       = "eu-west"
}

# Create a branch on the database
resource "planetscale_branch" "branch" {
  organization = planetscale_database.this.organization
  database     = planetscale_database.this.name
  name         = "my-branch"
}

# Create a backup of the database branch
# Unfortunately naming with custom names the backups is not supported yet by PlanetScale's API.
resource "planetscale_backup" "my-backup" {
    organization = planetscale_database.this.organization
    database     = planetscale_database.this.name
    branch       = planetscale_branch.branch.name
}