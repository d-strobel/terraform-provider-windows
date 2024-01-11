// Example with a local vagrant machine via WinRM
provider "windows" {
  endpoint = "127.0.0.1"

  winrm = {
    username = "vagrant"
    port     = 15985
  }

  kerberos = {
    realm           = "example.local"
    krb_config_file = "/path/to/krb5.conf"
  }
}
