terraform {
  backend "remote" {
    organization = "globe-and-citizen"
    workspaces {
      name = "network-non-production"
    }
  }
}

provider "aws" {
  region = var.AWS_REGION
}

module "network" {
  source = "./module/network"
}

module "iam" {
  source = "./module/iam"
}

module "ecs" {
  source          = "./module/ecs"
  cluster_name    = "layer8-non-production"
  vpc_cidr_block  = module.network.vpc_cidr_block
  vpc_id          = module.network.vpc_id
  subnets         = module.network.private_subnets
  iam_profile_arn = module.iam.iam_profile_arn
}
