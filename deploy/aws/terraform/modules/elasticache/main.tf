resource "aws_elasticache_subnet_group" "main" {
  name       = "optivor-redis-subnet-group-${var.environment}"
  subnet_ids = var.private_subnet_ids
}

resource "aws_elasticache_replication_group" "optivor" {
  replication_group_id = "optivor-redis-${var.environment}"
  description          = "Redis cache backend for Optivor Pods"
  node_type            = var.node_type
  num_cache_clusters   = 1
  port                 = 6379

  subnet_group_name  = aws_elasticache_subnet_group.main.name
  security_group_ids = [var.redis_security_group_id]

  at_rest_encryption_enabled = true
  transit_encryption_enabled = false

  tags = {
    Name        = "optivor-redis-${var.environment}"
    Environment = var.environment
  }
}
