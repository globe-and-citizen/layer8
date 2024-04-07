variable "vpc_id" {
  description = "The VPC ID"
}

variable "cluster_id" {
  description = "The ECS cluster ID"
}

variable "subnets" {
  description = "The list of subnets"

}
variable "node_security_group_id" {
  description = "The security group ID for the ECS cluster"
}

variable "capacity_provider_name" {
  description = "ECS capacity provider name"
}

variable "task_execution_role_arn" {
  description = "The ARN of the task execution role that the Amazon ECS container agent and the Docker daemon can assume"
}

variable "task_role_arn" {
  description = "The ARN of IAM role that allows your Amazon ECS container task to make calls to other AWS services"
}

variable "ecr_repository_url" {
  description = "The URL of the ECR repository"
}

variable "ecr_image_tag" {
  description = "The tag of the ECR image"
}

variable "s3_arn_env_file" {
  description = "The ARN of the S3 bucket for the environment file"
}

variable "loki_url" {
  description = "The URL of the Loki server"
}

variable "cloudflare_tunnel_token" {
    description = "Cloudflare tunnel token"
}