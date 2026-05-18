variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name used in tags and resource names"
  type        = string
  default     = "campus-connect-api"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "key_name" {
  description = "Name for the AWS key pair"
  type        = string
  default     = "campus-connect-key"
}

variable "public_key" {
  description = "SSH public key content (e.g., ssh-rsa AAAA...)"
  type        = string
}

variable "ssh_allowed_cidr" {
  description = "CIDR allowed to access SSH port 22"
  type        = string
  default     = "0.0.0.0/0"
}

variable "app_port" {
  description = "API port exposed on the EC2 instance"
  type        = number
  default     = 8080
}

variable "app_allowed_cidr" {
  description = "CIDR allowed to access API port"
  type        = string
  default     = "0.0.0.0/0"
}

variable "root_volume_size_gb" {
  description = "Root volume size in GB"
  type        = number
  default     = 20
}
