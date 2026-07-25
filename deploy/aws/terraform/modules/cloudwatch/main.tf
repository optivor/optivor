resource "aws_cloudwatch_log_group" "optivor" {
  name              = "/ecs/optivor-${var.environment}"
  retention_in_days = var.retention_in_days

  tags = {
    Name        = "optivor-logs-${var.environment}"
    Environment = var.environment
  }
}
