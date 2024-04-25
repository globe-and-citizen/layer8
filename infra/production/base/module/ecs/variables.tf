variable "cluster_name" {
  type = string
}

variable "vpc_id" {
  description = "VPC ID for the ECS cluster"
}

variable "vpc_cidr_block" {
  description = "CIDR block of the VPC for the ECS cluster"
}

variable "subnets" {
  description = "List of the private subnets"
}

variable "iam_profile_arn" {
  description = "IAM profile ARN for the ECS cluster"
}

variable "task_execution_role_arn" {}

variable "task_role_arn" {}

variable "aws_region" {}