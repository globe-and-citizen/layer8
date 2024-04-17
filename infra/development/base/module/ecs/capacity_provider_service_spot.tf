resource "aws_launch_template" "service_spot_launch_template" {
  name          = "${aws_ecs_cluster.cluster.name}-service-spot-launch-template"
  image_id      = data.aws_ssm_parameter.ecs_node_ami.value
  instance_type = "t3.medium"

  iam_instance_profile { arn = var.iam_profile_arn }
  monitoring { enabled = false }

  user_data = base64encode(<<-EOF
      #!/bin/bash
      echo ECS_CLUSTER=${aws_ecs_cluster.cluster.name} >> /etc/ecs/ecs.config;
    EOF
  )
}

resource "aws_autoscaling_group" "service_spot_asg" {
  name                = "${aws_ecs_cluster.cluster.name}-service-spot-asg"
  min_size            = 1
  max_size            = 1
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
        launch_template_id = aws_launch_template.service_spot_launch_template.id
        version            = "$Latest"
      }

      override {
        instance_type = "t3.medium"
      }

      override {
        instance_type = "t3a.medium"
      }
    }
  }

  tag {
    key                 = "AmazonECSManaged"
    value               = ""
    propagate_at_launch = true
  }
}


resource "aws_ecs_capacity_provider" "service_spot_capacity_provider" {
  name = "${aws_ecs_cluster.cluster.name}-service-spot-capacity-provider"
  auto_scaling_group_provider {
    auto_scaling_group_arn         = aws_autoscaling_group.service_spot_asg.arn
    managed_termination_protection = "DISABLED"
  }
}