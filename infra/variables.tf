variable "tenancy_ocid" { type = string }
variable "user_ocid" { type = string }
variable "fingerprint" { type = string }
variable "private_key_path" { type = string }
variable "region" { type = string }
variable "compartment_id" { type = string }
variable "ssh_public_key" { type = string }

variable "public_subnet_cidr" {
  description = "CIDR block for the Terraform-managed public subnet"
  type        = string
  default     = "10.0.10.0/24"
}
