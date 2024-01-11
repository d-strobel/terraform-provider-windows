---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "windows_local_group Data Source - terraform-provider-windows"
subcategory: ""
description: |-
  
---

# windows_local_group (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) Define the name of the local security group.
- `sid` (String) Define the security ID of the local security group.

### Read-Only

- `description` (String) The description of the local security group.
- `id` (String) The ID of the retrieved local security group. This is the same as the SID.