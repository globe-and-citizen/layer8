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
