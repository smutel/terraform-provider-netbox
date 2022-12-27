terraform {
  required_providers {
    netbox = {
      source = "smutel/netbox"
      version = "~> 6.0.0"
    }
  }
}

provider netbox {
  # Environment variable NETBOX_URL
  url = "127.0.0.1:8000"

  # Environment variable NETBOX_BASEPATH
  basepath = "/api"

  # Environment variable NETBOX_TOKEN
  token = "0123456789abcdef0123456789abcdef01234567"

  # Environment variable NETBOX_SCHEME
  scheme = "http"

  # Environment variable NETBOX_INSECURE
  insecure = "true"
}
