terraform {
  backend "remote" {
    organization = "globe-and-citizen"
    workspaces {
      name = "layer8-server-development"
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
  family                   = "layer8-server"
  task_role_arn            = data.terraform_remote_state.network.outputs.task_role_arn
  execution_role_arn       = data.terraform_remote_state.network.outputs.task_execution_role_arn
  network_mode             = "host"
  requires_compatibilities = ["EC2"]
  skip_destroy             = true

  container_definitions = jsonencode([
    {
      name              = "layer8-server",
      essential         = true,
      image             = "${var.ecr_repository_url}:${var.ecr_image_tag}",
      cpu               = 0,
      memoryReservation = 512,
      mountPoints       = [],
      portMappings = [
        { containerPort = 5001, hostPort = 5001, protocol = "tcp" },
      ],
      environment = [],
      environmentFiles = [
        { type = "s3", value = var.s3_arn_env_file }
      ],
      systemControls = [],
      volumesFrom    = [],
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-create-group" : "true",
          "awslogs-group" : "ecs/development",
          "awslogs-region" : "${var.aws_region}",
          "awslogs-stream-prefix" : "layer8server"
        },
      },
      user = "0"
    },
    {
      name              = "telegraf",
      essential         = true,
      image             = "hannmuhammadd/telegraf-opentelemetry-influxdb",
      cpu               = 0,
      memoryReservation = 128,
      mountPoints       = [],
      portMappings = [
        { containerPort = 4317, hostPort = 4317, protocol = "tcp" },
      ],
      environment = [
        {
          name  = "INFLUXDB_URL",
          value = "${var.influxdb_url}",
        },
        {
          name  = "INFLUXDB_TOKEN",
          value = "${var.influxdb_token}",
        },
        {
          name  = "INFLUXDB_ORGANIZATION",
          value = "layer8",
        },
        {
          name  = "INFLUXDB_BUCKET_NAME",
          value = "layer8",
        }
      ],
      environmentFiles = [],
      systemControls   = [],
      volumesFrom      = [],
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-create-group" : "true",
          "awslogs-group" : "ecs/development",
          "awslogs-region" : "${var.aws_region}",
          "awslogs-stream-prefix" : "layer8server"
        },
      },
      user    = "0",
      command = []
    }
  ])

  tags = {
    Name = "layer8-server"
  }
}

resource "aws_ecs_service" "service" {
  name                               = "layer8-server"
  cluster                            = data.terraform_remote_state.network.outputs.ecs_cluster_id
  task_definition                    = aws_ecs_task_definition.task_definition.arn
  desired_count                      = 1
  deployment_maximum_percent         = 100
  deployment_minimum_healthy_percent = 0

  deployment_circuit_breaker {
    enable   = "true"
    rollback = "false"
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

