# Optivor AWS Deployment Adapter (ECS Fargate)

This deployment adapter provisions a **Full Optivor Runtime Container Service** on AWS ECS Fargate with AWS Application Load Balancer (ALB) and native AWS IRSA (IAM Roles for Service Accounts) integration, following **ADR-0002**.

## Architecture

Per **ADR-0002**:
> "Deployment Adapters explicitly document whether they deploy the full runtime or a proxy in front of it."

This adapter deploys the **Full Runtime** binary inside Docker containers on AWS ECS Fargate, scaling horizontally behind an AWS ALB.

```
[ Client Request ] ---> [ AWS Application Load Balancer ] ---> [ ECS Fargate Pod (Full Optivor Runtime) ] ---> [ AWS S3 / ElastiCache Redis ]
```

## Setup & Deployment

```bash
cd deploy/aws
sam build
sam deploy --guided
```
