terraform {
  backend "remote" {
    organization = "globe-and-citizen"
    workspaces {
      name = "layer8-server-non-production"
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
      name = "network-non-production"
    }
  }
}

resource "aws_ecs_task_definition" "app" {
  family                   = "layer8-server"
  task_role_arn            = data.terraform_remote_state.network.outputs.task_role_arn
  execution_role_arn       = data.terraform_remote_state.network.outputs.task_execution_role_arn
  network_mode             = "awsvpc"
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
        { containerPort = 5001, hostPort = 80, protocol = "tcp" },
      ],
      environment = [],
      environmentFiles = [
        { type = "s3", value = var.s3_arn_env_file }
      ],
      systemControls = [],
      volumesFrom    = [],
      logConfiguration = {
        logDriver = "awsfirelens",
        options = {
          "LabelKeys" : "container_name,ecs_task_definition,source,ecs_cluster",
          "Labels" : "{job=\"firelens\"}",
          "LineFormat" : "key_value",
          "Name" : "grafana-loki",
          "RemoveKeys" : "container_id,ecs_task_arn",
          "Url" : var.loki_url
        }
      },
      user = "0"
    },
    {
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-create-group" : "true",
          "awslogs-group" : "/ecs/ecs-aws-firelens-sidecar-container",
          "awslogs-region" : "ap-southeast-1",
          "awslogs-stream-prefix" : "firelens"
        },
      },
      name              = "log-router"
      essential         = true
      image             = "grafana/fluent-bit-plugin-loki:2.9.1",
      cpu               = 0,
      memoryReservation = 128,
      mountPoints       = [],
      portMappings      = [],
      environment = [
        {
          name  = "LOKI_URL"
          value = var.loki_url
        }
      ],
      firelensConfiguration = {
        options = {
          enable-ecs-log-metadata = "true",
        },
        type = "fluentbit"
      },
      systemControls = [],
      volumesFrom    = [],
      user           = "0"
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
      environmentFiles = [],
      systemControls = [],
      volumesFrom    = [],
      logConfiguration = {
        logDriver = "awsfirelens",
        options = {
          "LabelKeys" : "container_name,ecs_task_definition,source,ecs_cluster",
          "Labels" : "{job=\"firelens\"}",
          "LineFormat" : "key_value",
          "Name" : "grafana-loki",
          "RemoveKeys" : "container_id,ecs_task_arn",
          "Url" : var.loki_url
        }
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
    Name = "layer8-server"
  }
}
