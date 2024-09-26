---
page_title: "windows_local_group_members Data Source - terraform-provider-windows"
subcategory: "Local"
description: |-
  Retrieve a list of members for a specific local security group
---
# windows_local_group_members (Data Source)

<!-- data-source description generated from schema -->
Retrieve a list of members for a specific local security group.
<!-- examples generated from example files -->
## Example Usage

```terraform
data "windows_local_group_members" "administrators" {
  name = "Administrators"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the local group.

### Read-Only

- `members` (Attributes List) (see [below for nested schema](#nestedatt--members))

<a id="nestedatt--members"></a>
### Nested Schema for `members`

Read-Only:

- `name` (String) The name of the local group member.
- `object_class` (String) The ObjectClass of the local group member.
- `sid` (String) The security ID of the local group member.