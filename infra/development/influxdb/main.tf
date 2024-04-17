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

provider "random" {
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

resource "aws_efs_file_system" "efs" {
  creation_token = "influxdb-development"

  tags = {
    Name = "influxdb-development"
  }
}

resource "aws_efs_mount_target" "efs_mount" {
  count          = 2
  file_system_id = aws_efs_file_system.efs.id
  subnet_id      = data.terraform_remote_state.network.outputs.private_subnets[count.index].id
  security_groups = [
    data.terraform_remote_state.network.outputs.node_security_group_id
  ]
}

resource "random_password" "influxdb_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_ecs_task_definition" "task_definition" {
  family                   = "influxdb2"
  task_role_arn            = data.terraform_remote_state.network.outputs.task_role_arn
  execution_role_arn       = data.terraform_remote_state.network.outputs.task_execution_role_arn
  network_mode             = "host"
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
        {
          sourceVolume  = "influxdb-data"
          containerPath = "/var/lib/influxdb2"
          readOnly      = false
        },
        {
          sourceVolume  = "influxdb-data"
          containerPath = "/etc/influxdb2"
          readOnly      = false
        },
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
          value = random_password.influxdb_password.result,
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
          "awslogs-stream-prefix" : "influxdb2"
        },
      },
      user = "0"
    },
  ])

  tags = {
    Name = "influxdb2"
  }

  volume {
    name = "influxdb-data"
    efs_volume_configuration {
      file_system_id = aws_efs_file_system.efs.id
      root_directory = "/"
    }
  }
}

resource "aws_ecs_service" "service" {
  name                               = "influxdb2"
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

