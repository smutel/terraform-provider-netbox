data "netbox_ipam_ip_range" "iprange_test" {
  start_address = "192.168.56.1/24"
  end_address = "192.168.56.254/24"
}
