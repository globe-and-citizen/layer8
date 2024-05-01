output "iam_profile_arn" {
  description = "IAM profile ARN for the ECS cluster"
  value       = aws_iam_instance_profile.ecs_node.arn
}

output "task_role_arn" {
  description = "value"
  value       = aws_iam_role.ecs_task_role.arn
}

output "task_execution_role_arn" {
  description = "value"
  value       = aws_iam_role.ecs_exec_role.arn
}