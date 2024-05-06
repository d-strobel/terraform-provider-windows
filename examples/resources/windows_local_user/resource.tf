resource "windows_local_user" "this" {
  name                     = "MyUser"
  full_name                = "My User"
  description              = "This is a test user"
  password                 = "P@ssw0rd!"
  enabled                  = true
  password_never_expires   = true
  user_may_change_password = true
}
