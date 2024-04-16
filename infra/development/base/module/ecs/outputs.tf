output "cluster_id" {
  description = "ECS Cluster ID"
  value       = aws_ecs_cluster.cluster.id
}

output "service_spot_capacity_provider_name" {
  description = "ECS Capacity Provider Name for Service"
  value       = aws_ecs_capacity_provider.service_spot_capacity_provider.name
}

output "security_group_id" {
  description = "ECS Security Group ID"
  value       = aws_security_group.ecs_node_sg.id
}