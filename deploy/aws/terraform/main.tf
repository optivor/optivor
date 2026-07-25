terraform {
  required_version = ">= 1.5.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

module "vpc" {
  source               = "./modules/vpc"
  environment          = var.environment
  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  availability_zones   = var.availability_zones
}

module "security" {
  source      = "./modules/security"
  environment = var.environment
  vpc_id      = module.vpc.vpc_id
}

module "alb" {
  source                = "./modules/alb"
  environment           = var.environment
  vpc_id                = module.vpc.vpc_id
  public_subnet_ids     = module.vpc.public_subnet_ids
  alb_security_group_id = module.security.alb_security_group_id
}

module "s3" {
  source      = "./modules/s3"
  environment = var.environment
  bucket_name = var.s3_bucket_name
}

module "elasticache" {
  source                  = "./modules/elasticache"
  environment             = var.environment
  vpc_id                  = module.vpc.vpc_id
  private_subnet_ids      = module.vpc.private_subnet_ids
  redis_security_group_id = module.security.redis_security_group_id
}

module "cloudwatch" {
  source      = "./modules/cloudwatch"
  environment = var.environment
}

module "ecs" {
  source                      = "./modules/ecs"
  environment                 = var.environment
  private_subnet_ids          = module.vpc.private_subnet_ids
  ecs_tasks_security_group_id = module.security.ecs_tasks_security_group_id
  target_group_arn            = module.alb.target_group_arn
  container_image             = var.container_image
  cpu                         = var.cpu
  memory                      = var.memory
  desired_count               = var.desired_count
  min_capacity                = var.min_capacity
  max_capacity                = var.max_capacity
  log_group_name              = module.cloudwatch.log_group_name
  s3_bucket_name              = module.s3.bucket_id
  s3_bucket_arn               = module.s3.bucket_arn
  redis_endpoint              = module.elasticache.primary_endpoint_address
}
