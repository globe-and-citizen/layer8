resource "aws_ecs_task_definition" "task_definition" {
  family                   = "cloudflare-tunnel"
  task_role_arn            = var.task_role_arn
  execution_role_arn       = var.task_execution_role_arn
  network_mode             = "host"
  requires_compatibilities = ["EC2"]
  skip_destroy             = true

  container_definitions = jsonencode([
    {
      name              = "cloudflared-tunnel",
      essential         = true,
      image             = "cloudflare/cloudflared:latest",
      cpu               = 0,
      memoryReservation = 128,
      mountPoints       = [],
      portMappings      = [],
      environment       = [],
      environmentFiles  = [],
      systemControls    = [],
      volumesFrom       = [],
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-create-group" : "true",
          "awslogs-group" : "ecs/development",
          "awslogs-region" : "${var.aws_region}",
          "awslogs-stream-prefix" : "cloudflare-tunnel"
        },
      },
      user = "0",
      command = [
        "tunnel",
        "--no-autoupdate",
        "run",
        "--token",
        "${var.cloudflare_tunnel_token}"
      ]
    }
  ])

  tags = {
    Name = "cloudflare-tunnel"
  }
}

resource "aws_ecs_service" "service" {
  name                               = "cloudflare-tunnel"
  cluster                            = aws_ecs_cluster.cluster.id
  task_definition                    = aws_ecs_task_definition.task_definition.arn
  desired_count                      = 1
  deployment_maximum_percent         = 100
  deployment_minimum_healthy_percent = 0

  deployment_circuit_breaker {
    enable   = "true"
    rollback = "false"
  }

  capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_provider.service_spot_capacity_provider.name
    base              = 1
    weight            = 100
  }

  ordered_placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }

  lifecycle {
    ignore_changes = [desired_count]
  }

  force_new_deployment = true

  triggers = {
    redeployment = plantimestamp()
  }

  depends_on = [aws_ecs_task_definition.task_definition]
}

