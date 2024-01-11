provider "windows" {
  endpoint = "127.0.0.1"

  ssh = {
    username = "vagrant"
    password = "vagrant"
    port     = 1222
  }
}
