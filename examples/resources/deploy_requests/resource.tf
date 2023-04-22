# Display how to use the deploy requests resource

locals {
    organization = "my-org"
    database     = "my-db"
}

# Create a new development branch for the database so that we can point the deploy request into it
resource "planetscale_database_branch" "destination_branch" {
  organization = local.organization
  database     = local.database
  name         = "my-deploy-request-branch"
}

# Create a new deploy request from our existing branch into the new branch we created
resource "planetscale_deploy_request" "new_deploy_reqest" {
  organization = local.organization
  database     = local.database
  branch       = "my-existing-tf-branch"
  into_branch  = planetscale_database_branch.destination_branch.name
}