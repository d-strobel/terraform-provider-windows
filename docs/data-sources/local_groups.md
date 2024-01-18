---
page_title: "windows_local_groups Data Source - terraform-provider-windows"
subcategory: "Local"
description: |-
  Retrieve a list of all local security groups
---
# windows_local_groups (Data Source)

<!-- data-source description generated from schema -->
Retrieve a list of all local security groups.
<!-- examples generated from example files -->
## Example Usage

```terraform
data "windows_local_groups" "all" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `groups` (Attributes List) (see [below for nested schema](#nestedatt--groups))

<a id="nestedatt--groups"></a>
### Nested Schema for `groups`

Read-Only:

- `description` (String) The description of the local security group.
- `id` (String) The ID of the retrieved local security group. This is the same as the SID.
- `name` (String) The name of the local security group.
- `sid` (String) The security ID of the local security group.