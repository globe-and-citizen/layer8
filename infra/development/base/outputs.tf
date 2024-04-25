output "vpc_id" {
  value = module.network.vpc_id
}

output "ecs_cluster_id" {
  value = module.ecs.cluster_id
}

output "service_spot_capacity_provider_name" {
  value = module.ecs.service_spot_capacity_provider_name
}

output "private_subnets" {
  value = module.network.private_subnets
}

output "node_security_group_id" {
  value = module.ecs.security_group_id
}

output "ecr_repository_url" {
  value = module.ecr.repository_url
}