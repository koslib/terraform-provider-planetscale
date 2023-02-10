# Fetching all regions enabled for you organization
data "planetscale_regions" "all" {}

output "all_regions" {
    value = data.planetscale_regions.all
}