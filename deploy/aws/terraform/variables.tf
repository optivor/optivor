variable "aws_region" {
  type        = string
  default     = "us-east-1"
  description = "AWS region for infrastructure deployment"
}

variable "environment" {
  type        = string
  default     = "prod"
  description = "Deployment environment (e.g. prod, staging, dev)"
}

variable "vpc_cidr" {
  type        = string
  default     = "10.0.0.0/16"
  description = "CIDR block for the VPC"
}

variable "public_subnet_cidrs" {
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
  description = "Public subnet CIDR blocks across Multi-AZ"
}

variable "private_subnet_cidrs" {
  type        = list(string)
  default     = ["10.0.10.0/24", "10.0.11.0/24"]
  description = "Private subnet CIDR blocks across Multi-AZ"
}

variable "availability_zones" {
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
  description = "AWS Availability Zones"
}

variable "s3_bucket_name" {
  type        = string
  description = "Name of the S3 bucket for Optivor image storage"
}

variable "container_image" {
  type        = string
  default     = "optivor/optivor:latest"
  description = "Docker container image repository and tag for Optivor"
}

variable "cpu" {
  type        = number
  default     = 1024
  description = "Fargate task CPU units (1024 = 1 vCPU)"
}

variable "memory" {
  type        = number
  default     = 2048
  description = "Fargate task Memory in MB"
}

variable "desired_count" {
  type        = number
  default     = 2
  description = "Desired number of ECS task instances"
}

variable "min_capacity" {
  type        = number
  default     = 2
  description = "Minimum instances for ECS Auto Scaling"
}

variable "max_capacity" {
  type        = number
  default     = 10
  description = "Maximum instances for ECS Auto Scaling"
}
