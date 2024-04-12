terraform {
  backend "remote" {
    organization = "globe-and-citizen"
    workspaces {
      name = "layer8-influxdb-development"
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
  family                   = "influxdb2"
  task_role_arn            = data.terraform_remote_state.network.outputs.task_role_arn
  execution_role_arn       = data.terraform_remote_state.network.outputs.task_execution_role_arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["EC2"]
  skip_destroy             = true

  container_definitions = jsonencode([
    {
      name              = "influxdb2",
      image             = "influxdb:2",
      essential         = true,
      cpu               = 0,
      memoryReservation = 512,
      portMappings = [
        {
          containerPort = 8086,
          hostPort      = 8086,
          protocol      = "tcp",
        }
      ],
      mountPoints = [
      ],
      environment = [
        {
          name  = "DOCKER_INFLUXDB_INIT_MODE",
          value = "setup",
        },
        {
          name  = "DOCKER_INFLUXDB_INIT_USERNAME",
          value = "influxdbadmin",
        },
        {
          name  = "DOCKER_INFLUXDB_INIT_PASSWORD",
          value = "somethingthatyoudontknow",
        },
        {
          name  = "DOCKER_INFLUXDB_INIT_ORG",
          value = "layer8",
        },
        {
          name  = "DOCKER_INFLUXDB_INIT_BUCKET",
          value = "layer8",
        },
      ]
      environmentFiles = [],
      systemControls   = [],
      volumesFrom      = [],
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-create-group" : "true",
          "awslogs-group" : "ecs/development",
          "awslogs-region" : "${var.aws_region}",
          "awslogs-stream-prefix": "influxdb2"
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
          "awslogs-stream-prefix": "cloudflare-tunnel"
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
    Name = "influxdb2"
  }
}

resource "aws_ecs_service" "service" {
  name            = "influxdb2"
  cluster         = data.terraform_remote_state.network.outputs.ecs_cluster_id
  task_definition = aws_ecs_task_definition.task_definition.arn
  desired_count   = 1
  deployment_maximum_percent = 100
  deployment_minimum_healthy_percent = 0

  deployment_circuit_breaker {
    enable   = "true"
    rollback = "false"
  }

  network_configuration {
    assign_public_ip = false
    security_groups  = [data.terraform_remote_state.network.outputs.node_security_group_id]
    subnets          = [data.terraform_remote_state.network.outputs.private_subnets[0].id]
  }

  capacity_provider_strategy {
    capacity_provider = data.terraform_remote_state.network.outputs.db_spot_capacity_provider_name
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

