---
page_title: "windows_local_group Resource - terraform-provider-windows"
subcategory: "Local"
description: |-
  Manage local security groups
---
# windows_local_group (Resource)

<!-- resource description generated from schema -->
Manage local security groups.

<!-- examples generated from example files -->
## Example Usage

```terraform
resource "windows_local_group" "test" {
  name        = "Test"
  description = "This is a test group"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Define the name for the local security group.

### Optional

- `description` (String) Define a description for the local security group.

### Read-Only

- `id` (String) The ID of the retrieved local security group. This is the same as the SID.
- `sid` (String) The security ID of the local security group.