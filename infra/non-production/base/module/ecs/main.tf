resource "aws_ecs_cluster" "cluster" {
  name = var.cluster_name
}

data "aws_ssm_parameter" "ecs_node_ami" {
  name = "/aws/service/ecs/optimized-ami/amazon-linux-2023/recommended/image_id"
}

resource "aws_launch_template" "launch-template" {
  name          = "${aws_ecs_cluster.cluster.name}-launch-template"
  image_id      = data.aws_ssm_parameter.ecs_node_ami.value
  instance_type = "t3.micro"

  iam_instance_profile { arn = var.iam_profile_arn }
  monitoring { enabled = true }

  user_data = base64encode(<<-EOF
      #!/bin/bash
      echo ECS_CLUSTER=${aws_ecs_cluster.cluster.name} >> /etc/ecs/ecs.config;
    EOF
  )
}

resource "aws_autoscaling_group" "asg" {
  name                = "${aws_ecs_cluster.cluster.name}-asg"
  min_size            = 1
  max_size            = 10
  capacity_rebalance  = "true"
  vpc_zone_identifier = var.subnets[*].id

  mixed_instances_policy {
    instances_distribution {
      on_demand_allocation_strategy            = "prioritized"
      on_demand_base_capacity                  = "0"
      on_demand_percentage_above_base_capacity = "0"
      spot_allocation_strategy                 = "price-capacity-optimized"
      spot_instance_pools                      = "0"
    }

    launch_template {
      launch_template_specification {
        launch_template_id = aws_launch_template.launch-template.id
        version            = "$Latest"
      }

      override {
        instance_type = "t3.micro"
      }

      override {
        instance_type = "t2.micro"
      }
    }
  }

  tag {
    key                 = "AmazonECSManaged"
    value               = ""
    propagate_at_launch = true
  }
}


resource "aws_ecs_capacity_provider" "capacity_provider" {
  name = "${aws_ecs_cluster.cluster.name}-capacity-provider"
  auto_scaling_group_provider {
    auto_scaling_group_arn         = aws_autoscaling_group.asg.arn
    managed_termination_protection = "DISABLED"
    managed_scaling {
      maximum_scaling_step_size = 2
      minimum_scaling_step_size = 1
      status                    = "ENABLED"
      target_capacity           = 100
    }
  }
}

resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name       = aws_ecs_cluster.cluster.name
  capacity_providers = [aws_ecs_capacity_provider.capacity_provider.name]

  default_capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_provider.capacity_provider.name
    base              = 1
    weight            = 100
  }
}

resource "aws_security_group" "ecs_node_sg" {
  name   = "ecs-node-sg-non-prod"
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
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
