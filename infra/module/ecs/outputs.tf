output "cluster_id" {
  description = "ECS Cluster Name"
  value       = aws_ecs_cluster.layer8.id
}

output "capacity_provider_name" {
  description = "ECS Capacity Provider Name"
  value       = aws_ecs_capacity_provider.capacity_provider.name
}

output "security_group_id" {
  description = "ECS Security Group ID"
  value       = aws_security_group.ecs_node_sg.id
}