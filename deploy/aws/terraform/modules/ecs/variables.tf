variable "environment" { type = string }
variable "private_subnet_ids" { type = list(string) }
variable "ecs_tasks_security_group_id" { type = string }
variable "target_group_arn" { type = string }
variable "container_image" { type = string }
variable "container_port" {
  type    = number
  default = 8080
}
variable "cpu" {
  type    = number
  default = 1024
}
variable "memory" {
  type    = number
  default = 2048
}
variable "desired_count" {
  type    = number
  default = 2
}
variable "min_capacity" {
  type    = number
  default = 2
}
variable "max_capacity" {
  type    = number
  default = 10
}
variable "log_group_name" { type = string }
variable "s3_bucket_name" { type = string }
variable "s3_bucket_arn" { type = string }
variable "redis_endpoint" { type = string }
