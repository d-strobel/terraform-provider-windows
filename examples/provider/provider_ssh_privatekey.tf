provider "windows" {
  endpoint = "127.0.0.1"

  ssh = {
    username    = "vagrant"
    port        = 1222
    private_key = <<EOT
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACCNKqLkDaaa4KGp+xaT0X94XVxGiwG6RHsymEc9/m39hwAAAJjpeDkr6Xg5
KwAAAAtzc2gtZWQyNTUxOQAAACCNKqLkDaaa4KGp+xaT0X94XVxGiwG6RHsymEc9/m39hw
AAAEAMT15+Ut2N+m9HW9wXgIeVR+qKeoT3UlVCxxnPsnoA5o0qouQNpprgoan7FpPRf3hd
XEaLAbpEezKYRz3+bf2HAAAAD2RzdHJvYmVsQE5CMDc4NAECAwQFBg==
-----END OPENSSH PRIVATE KEY-----
EOT
  }
}
