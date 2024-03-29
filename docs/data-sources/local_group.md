---
page_title: "windows_local_group Data Source - terraform-provider-windows"
subcategory: "Local"
description: |-
  Retrieve information about a local security group
---
# windows_local_group (Data Source)

<!-- data-source description generated from schema -->
Retrieve information about a local security group.
You can get a group by the name or the security ID of the group.
<!-- examples generated from example files -->
## Example Usage

```terraform
data "windows_local_group" "Administrator" {
  name = "Administrator"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) Define the name of the local security group.
- `sid` (String) Define the security ID of the local security group.

### Read-Only

- `description` (String) The description of the local security group.
- `id` (String) The ID of the retrieved local security group. This is the same as the SID.
