variable "aws_region" {
  default = "ap-southeast-2"
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

variable "cloudflare_tunnel_token" {
  description = "Cloudflare tunnel token"
}