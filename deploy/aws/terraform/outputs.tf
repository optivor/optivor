output "alb_dns_name" {
  description = "Public DNS name of the Application Load Balancer"
  value       = module.alb.alb_dns_name
}

output "s3_bucket_name" {
  description = "Name of the created S3 storage bucket"
  value       = module.s3.bucket_id
}

output "redis_primary_endpoint" {
  description = "Primary endpoint address of the ElastiCache Redis cluster"
  value       = module.elasticache.primary_endpoint_address
}

output "ecs_cluster_name" {
  description = "Name of the deployed ECS cluster"
  value       = module.ecs.cluster_name
}

output "ecs_service_name" {
  description = "Name of the deployed ECS Fargate service"
  value       = module.ecs.service_name
}
