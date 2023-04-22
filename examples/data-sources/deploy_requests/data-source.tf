# Display how getting information for a specific deploy request works.
# At this time, it is not possible to list deploy requests for a database. Therefore a data-source requires providing
# the deploy request number.

data "planetscale_deploy_requests" "this" {
  organization = "my-org"
  database     = "my-database"
  number       = "test-deploy-request-number-here"
}

output "deploy_request_info" {
  value = data.planetscale_deploy_requests.this
}