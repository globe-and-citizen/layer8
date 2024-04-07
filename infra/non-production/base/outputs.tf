output "vpc_id" {
  value = module.network.vpc_id
}

output "ecs_cluster_id" {
  value = module.ecs.cluster_id
}

output "capacity_provider_name" {
  value = module.ecs.capacity_provider_name
}

output "private_subnets" {
  value = module.network.private_subnets
}

output "node_security_group_id" {
  value = module.ecs.security_group_id
}

output "task_execution_role_arn" {
  value = module.iam.task_execution_role_arn
}

output "task_role_arn" {
  value = module.iam.task_role_arn
}

output "ecr_repository_url" {
  value = module.ecr.repository_url
}