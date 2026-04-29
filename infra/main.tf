terraform {
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = ">= 4.0.0"
    }
  }
}

provider "oci" {
  tenancy_ocid     = var.tenancy_ocid
  user_ocid        = var.user_ocid
  fingerprint      = var.fingerprint
  private_key_path = var.private_key_path
  region           = var.region
}

# Поиск доступных зон (Availability Domains)
data "oci_identity_availability_domains" "ads" {
  compartment_id = var.compartment_id
}

# Поиск образа Ubuntu 24.04
data "oci_core_images" "ubuntu" {
  compartment_id           = var.compartment_id
  operating_system         = "Canonical Ubuntu"
  operating_system_version = "24.04"
  shape                    = "VM.Standard.E2.1.Micro"
  sort_by                  = "TIMECREATED"
  sort_order               = "DESC"
}

# Поиск существующей VCN (samat-vcn-1)
data "oci_core_vcns" "existing_vcn" {
  compartment_id = var.compartment_id
  display_name   = "samat-vcn-1"
}

data "oci_core_internet_gateways" "existing" {
  compartment_id = var.compartment_id
  vcn_id         = data.oci_core_vcns.existing_vcn.virtual_networks[0].id
}

locals {
  vcn_id                        = data.oci_core_vcns.existing_vcn.virtual_networks[0].id
  existing_internet_gateway_ids = [for gateway in data.oci_core_internet_gateways.existing.gateways : gateway.id if gateway.enabled]
  internet_gateway_id           = concat(local.existing_internet_gateway_ids, oci_core_internet_gateway.server_igw[*].id)[0]
}

# Public network path for the VM. The previously configured subnet prohibits
# public IPs, so internet-facing instances need their own public subnet.
resource "oci_core_internet_gateway" "server_igw" {
  count = length(local.existing_internet_gateway_ids) == 0 ? 1 : 0

  compartment_id = var.compartment_id
  vcn_id         = local.vcn_id
  display_name   = "server-internet-gateway"
  enabled        = true
}

resource "oci_core_route_table" "server_public_route_table" {
  compartment_id = var.compartment_id
  vcn_id         = local.vcn_id
  display_name   = "server-public-route-table"

  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = local.internet_gateway_id
  }
}

resource "oci_core_subnet" "server_public_subnet" {
  compartment_id             = var.compartment_id
  vcn_id                     = local.vcn_id
  cidr_block                 = var.public_subnet_cidr
  display_name               = "server-public-subnet"
  dns_label                  = "serverpublic"
  route_table_id             = oci_core_route_table.server_public_route_table.id
  prohibit_public_ip_on_vnic = false
}

# Создание группы безопасности (Security Group) для открытия портов
resource "oci_core_network_security_group" "server_nsg" {
  compartment_id = var.compartment_id
  vcn_id         = local.vcn_id
  display_name   = "server-security-group"
}

# Правила для портов 80, 81, 3000, 9090, 22
resource "oci_core_network_security_group_security_rule" "ingress_rules" {
  for_each = toset(["80", "81", "3000", "9090", "22"])

  network_security_group_id = oci_core_network_security_group.server_nsg.id
  direction                 = "INGRESS"
  protocol                  = "6" # TCP
  source                    = "0.0.0.0/0"
  source_type               = "CIDR_BLOCK"

  tcp_options {
    destination_port_range {
      min = each.value
      max = each.value
    }
  }
}

# Создание инстанса
resource "oci_core_instance" "free_server" {
  availability_domain = data.oci_identity_availability_domains.ads.availability_domains[0].name
  compartment_id      = var.compartment_id
  shape               = "VM.Standard.E2.1.Micro"
  display_name        = "samat"

  create_vnic_details {
    subnet_id        = oci_core_subnet.server_public_subnet.id
    assign_public_ip = true
    nsg_ids          = [oci_core_network_security_group.server_nsg.id]
  }

  source_details {
    source_type = "image"
    source_id   = data.oci_core_images.ubuntu.images[0].id
  }

  metadata = {
    ssh_authorized_keys = var.ssh_public_key
    user_data           = base64encode(file("${path.module}/setup.sh"))
  }

  preserve_boot_volume = false
}
