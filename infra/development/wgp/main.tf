terraform {
  backend "remote" {
    organization = "globe-and-citizen"
    workspaces {
      name = "wgp-mock-development"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

data "terraform_remote_state" "network" {
  backend = "remote"

  config = {
    organization = "globe-and-citizen"
    workspaces = {
      name = "network-development"
    }
  }
}

resource "aws_ecs_task_definition" "task_definition" {
  family                   = "wgp-mock"
  task_role_arn            = data.terraform_remote_state.network.outputs.task_role_arn
  execution_role_arn       = data.terraform_remote_state.network.outputs.task_execution_role_arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["EC2"]
  skip_destroy             = true

  container_definitions = jsonencode([
    {
      name              = "frontend",
      essential         = true,
      image             = "${var.ecr_repository_url}:frontend-latest",
      cpu               = 0,
      memoryReservation = 256,
      mountPoints       = [],
      portMappings = [
        { containerPort = 8080, hostPort = 8080, protocol = "tcp" },
      ],
      environment = [
        {
          name  = "VITE_BACKEND_URL",
          value = var.backend_url,
        },
        {
          name = "VITE_PROXY_URL",
          value = var.layer8_url
        }
      ],
      environmentFiles = [],
      systemControls = [],
      volumesFrom    = [],
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-create-group" : "true",
          "awslogs-group" : "ecs/development",
          "awslogs-region" : "${var.aws_region}",
          "awslogs-stream-prefix": "wgp-frontend"
        },
      },
      user = "0"
    },
    {
      name              = "backend",
      essential         = true,
      image             = "${var.ecr_repository_url}:backend-latest",
      cpu               = 0,
      memoryReservation = 256,
      mountPoints       = [],
      portMappings = [
        { containerPort = 8000, hostPort = 8000, protocol = "tcp" },
      ],
      environment = [
        {
          name  = "PORT",
          value = var.backend_port,
        },
        {
          name  = "FRONTEND_URL",
          value = var.frontend_url,
        },
        {
          name  = "BACKEND_URL",
          value = var.backend_url,
        },
        {
          name = "LAYER8_URL",
          value = var.layer8_url
        }
      ],
      environmentFiles = [],
      systemControls = [],
      volumesFrom    = [],
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-create-group" : "true",
          "awslogs-group" : "ecs/development",
          "awslogs-region" : "${var.aws_region}",
          "awslogs-stream-prefix": "wgp-frontend"
        },
      },
      user = "0"
    },
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
          "awslogs-stream-prefix": "layer8server"
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
    Name = "wgp"
  }
}

resource "aws_ecs_service" "service" {
  name            = "wgp"
  cluster         = data.terraform_remote_state.network.outputs.ecs_cluster_id
  task_definition = aws_ecs_task_definition.task_definition.arn
  desired_count   = 1

  deployment_circuit_breaker {
    enable   = "true"
    rollback = "false"
  }


  network_configuration {
    assign_public_ip = false
    security_groups  = [data.terraform_remote_state.network.outputs.node_security_group_id]
    subnets          = data.terraform_remote_state.network.outputs.private_subnets[*].id
  }

  capacity_provider_strategy {
    capacity_provider = data.terraform_remote_state.network.outputs.service_spot_capacity_provider_name
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

