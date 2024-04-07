terraform {
  backend "remote" {
    organization = "globe-and-citizen"
    workspaces {
      name = "base-non-production"
    }
  }
}

provider "aws" {
  region = var.region
}

module "vpc" {
  source = "../../module/network"
}
