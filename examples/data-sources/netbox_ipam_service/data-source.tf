data "netbox_ipam_service" "service_test" {
  device_id = 5
  name      = "Mail"
  port      = 25
  protocol  = "tcp"
}
