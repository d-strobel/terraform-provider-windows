// Example with a local vagrant machine via WinRM
provider "windows" {
  endpoint = "127.0.0.1"

  winrm = {
    username = "vagrant"
    password = "vagrant"
    port     = 15985
  }
}
