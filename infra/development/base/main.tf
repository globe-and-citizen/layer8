terraform {
  backend "remote" {
    organization = "globe-and-citizen"
    workspaces {
      name = "network-development"
    }
  }
}

provider "aws" {
  region = var.AWS_REGION
}

module "network" {
  source = "./module/network"
}

module "ecs" {
  source                  = "./module/ecs"
  cluster_name            = "layer8-development"
  vpc_cidr_block          = module.network.vpc_cidr_block
  vpc_id                  = module.network.vpc_id
  subnets                 = module.network.private_subnets
  iam_profile_arn         = var.iam_profile_arn
  task_execution_role_arn = var.task_execution_role_arn
  task_role_arn           = var.task_role_arn
  aws_region              = var.AWS_REGION
}
