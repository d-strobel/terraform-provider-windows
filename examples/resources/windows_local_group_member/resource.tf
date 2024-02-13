resource "windows_local_user" "this" {
  name = "test-user"
}

resource "windows_local_group" "this" {
  name = "test-group"
}

resource "windows_local_group_member" "this" {
  group_id  = windows_local_group.this.id
  member_id = windows_local_user.this.id
}