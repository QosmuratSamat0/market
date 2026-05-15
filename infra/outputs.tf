output "instance_name" {
  description = "Name of the instance"
  value       = yandex_compute_instance.instance-based-on-coi.name
}

output "instance_id" {
  description = "ID of the instance"
  value       = yandex_compute_instance.instance-based-on-coi.id
}

output "public_ip" {
  description = "Public IP address of the instance"
  value       = yandex_compute_instance.instance-based-on-coi.network_interface.0.nat_ip_address
}

output "ssh_command" {
  description = "Command to connect to the instance"
  value       = "ssh yc-user@${yandex_compute_instance.instance-based-on-coi.network_interface.0.nat_ip_address}"
}
