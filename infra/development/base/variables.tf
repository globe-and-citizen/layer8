variable "AWS_REGION" {
  description = "The region where AWS operations will take place. Examples are us-east-1, us-west-2, etc."
  default     = "ap-southeast-2"
}

variable "iam_profile_arn" {
  description = "IAM profile ARN for the ECS cluster"
}

variable "task_execution_role_arn" {}

variable "task_role_arn" {}