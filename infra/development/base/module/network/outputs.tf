output "vpc_id" {
  description = "The ID of the VPC"
  value       = aws_vpc.vpc.id
}

output "vpc_cidr_block" {
  description = "The CIDR block of the VPC"
  value       = aws_vpc.vpc.cidr_block
}

output "public_subnets" {
  description = "List of the public subnets"
  value       = aws_subnet.public
}

output "private_subnets" {
  description = "List of the private subnets"
  value       = aws_subnet.private
}
