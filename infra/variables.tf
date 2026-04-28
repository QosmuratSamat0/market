variable "aws_region" {
  description = "AWS region to deploy in"
  default     = "us-east-1"
}

variable "instance_type" {
  description = "Type of EC2 instance"
  default     = "t3.medium"
}

variable "project_name" {
  description = "Name of the project"
  default     = "market-place"
}
