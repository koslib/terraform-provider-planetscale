---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "planetscale_database_branch_password Resource - terraform-provider-planetscale"
subcategory: ""
description: |-
  The database branch password resource allows you to manage a database branch password in Planetscale. This resource is used to create, read, update, and delete database branch passwords. For more information on database branch passwords, see the Planetscale documentation at https://planetscale.com/docs/concepts/connection-strings
---

# planetscale_database_branch_password (Resource)

The database branch password resource allows you to manage a database branch password in Planetscale. This resource is used to create, read, update, and delete database branch passwords. For more information on database branch passwords, see the Planetscale documentation at https://planetscale.com/docs/concepts/connection-strings



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `branch` (String) The name of the branch.
- `database` (String) The name of the database.
- `name` (String) The name of the database branch password.
- `organization` (String) The name of the organization.

### Optional

- `role` (String) The role of the database branch password. Defaults to admin. Once a password is created, its role cannot be changed. Supported values: admin, reader, writer, readwriter.

### Read-Only

- `plaintext` (String, Sensitive) The plaintext password of the database branch password.
- `public_id` (String) The public ID of the database branch password.
- `username` (String) The username of the database branch password.


