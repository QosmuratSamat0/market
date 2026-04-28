output "container_name" {
  description = "Name of the provisioned container"
  value       = docker_container.app_server.name
}

output "container_id" {
  description = "ID of the provisioned container"
  value       = docker_container.app_server.id
}

output "exposed_ports" {
  description = "Ports exposed by the container"
  value       = "80, 3000, 9090, 2222 (SSH)"
}
