terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.1"
    }
  }
}

provider "docker" {}

# Image for the app server (simulating a VM)
resource "docker_image" "ubuntu" {
  name         = "ubuntu:22.04"
  keep_locally = true
}

# Network for the infrastructure
resource "docker_network" "infra_network" {
  name = "market_place_infra"
}

# Container simulating the production VM
resource "docker_container" "app_server" {
  image = docker_image.ubuntu.image_id
  name  = "production-server-vm"
  
  networks_advanced {
    name = docker_network.infra_network.name
  }

  # Simulating open ports
  ports {
    internal = 80
    external = 80
  }
  
  ports {
    internal = 3000
    external = 3000
  }

  ports {
    internal = 9090
    external = 9090
  }

  ports {
    internal = 22
    external = 2222
  }

  # Command to keep the container running
  command = ["tail", "-f", "/dev/null"]
}
