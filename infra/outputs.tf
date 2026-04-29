output "instance_name" {
  description = "Name of the instance"
  value       = oci_core_instance.free_server.display_name
}

output "instance_id" {
  description = "OCID of the instance"
  value       = oci_core_instance.free_server.id
}

output "public_ip" {
  description = "Public IP address of the instance"
  value       = oci_core_instance.free_server.public_ip
}

output "ssh_command" {
  description = "Command to connect to the instance"
  value       = "ssh -i /home/admin/Downloads/ssh-key-2026-04-19.key ubuntu@${oci_core_instance.free_server.public_ip}"
}
