resource "aws_ecs_cluster" "cluster" {
  name = var.cluster_name
}

data "aws_ssm_parameter" "ecs_node_ami" {
  name = "/aws/service/ecs/optimized-ami/amazon-linux-2023/recommended/image_id"
}

resource "aws_ecs_cluster_capacity_providers" "capacity_provider_mapping" {
  cluster_name       = aws_ecs_cluster.cluster.name
  capacity_providers = [
    aws_ecs_capacity_provider.service_spot_capacity_provider.name
  ]

  default_capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_provider.service_spot_capacity_provider.name
    base              = 1
    weight            = 100
  }
}

resource "aws_cloudwatch_log_group" "log_group" {
  name = "ecs/development"
}

resource "aws_security_group" "ecs_node_sg" {
  name   = "ecs-node-sg-development"
  vpc_id = var.vpc_id

  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr_block]
    description = "vpc cidr"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
