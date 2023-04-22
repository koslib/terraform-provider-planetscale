---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "planetscale_backup Resource - terraform-provider-planetscale"
subcategory: ""
description: |-
  A Planetscale backup. This resource will create a new backup for a database in your Planetscale organization. The backup will be created for the specified branch.
---

# planetscale_backup (Resource)

A Planetscale backup. This resource will create a new backup for a database in your Planetscale organization. The backup will be created for the specified branch.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `branch` (String) The name of the branch to create the backup for.
- `database` (String) The name of the database to create the backup for.
- `organization` (String) The organization where the backup will be created as well as the database/branch belong to.

### Optional

- `public_id` (String) The public ID of the backup.

### Read-Only

- `completed_at` (String) If the backup is completed, this is the timestamp of when it was completed.
- `created_at` (String) The timestamp of when the backup object was created.
- `expires_at` (String) If the backup is completed, this is the timestamp of when it will expire.
- `name` (String) The name of the backup.
- `size` (Number) The size of the backup.
- `started_at` (String) The timestamp of when the backup started.
- `state` (String) The state of the backup. Options are: 'pending', 'running', 'success', 'failed', 'canceled'.
- `updated_at` (String) Last update timestamp for the backup.

