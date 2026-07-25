# Optivor AWS Deployment Adapter (Production-Grade Terraform)

This deployment adapter provisions a **production-grade, multi-AZ Optivor Runtime Infrastructure on AWS** using modular Terraform, following **ADR-0002**.

## Architecture Overview

Per **ADR-0002**:
> "Deployment Adapters explicitly document whether they deploy the full runtime or a proxy in front of it."

This adapter provisions the **Full Optivor Runtime Container Cluster** on AWS ECS Fargate:

```
[ Internet ]
     │
     ▼
[ Application Load Balancer (ALB) ] (Public Subnets)
     │
     ▼
[ ECS Fargate Tasks (Optivor Containers) ] (Private Subnets, Auto Scaled 2-10 tasks)
     ├───> [ AWS S3 Bucket ] (Private IAM Access via Task IAM Role / IRSA)
     └───> [ ElastiCache Redis Cluster ] (Private Subnets, Encrypted)
```

## Infrastructure Components

- **VPC & Networking**: Multi-AZ VPC with Public & Private Subnets, Internet Gateway, NAT Gateway.
- **Security Groups**: Tight ingress/egress rules restricting traffic (ALB -> ECS Tasks -> Redis/S3).
- **Application Load Balancer**: Health checking `/healthz`, automated target group registration.
- **ECS Fargate Cluster**: Auto Scaling target tracking (CPU 70%), IAM Task Execution & Task Roles (IRSA).
- **ElastiCache Redis**: Production Redis replication group in private subnets for shared image caching.
- **S3 Bucket**: Encrypted storage bucket with public access block and IAM role policy.
- **CloudWatch Logs**: Centralized container log retention and monitoring.

## Deployment Instructions

### Prerequisites
- [Terraform](https://www.terraform.io/downloads) >= 1.5.0
- [AWS CLI](https://aws.amazon.com/cli/) configured with deployment permissions.

### Quick Start

1. Navigate to the terraform directory:
   ```bash
   cd deploy/aws/terraform
   ```

2. Create your configuration file:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```
   Edit `terraform.tfvars` with your unique `s3_bucket_name` and desired region/task sizes.

3. Initialize and apply Terraform:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

4. Verify output:
   After deployment, Terraform will output the `alb_dns_name`. Access Optivor health check at:
   ```bash
   curl http://<alb_dns_name>/healthz
   ```
