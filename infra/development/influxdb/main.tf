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

resource "aws_instance" "ec2_instance" {
  tags = {
    Name = "influxdb2-development"
  }

  ami           = "ami-023adaba598e661ac"
  instance_type = "t3.micro"
  key_name      = "key-pair-one"

  subnet_id              = data.terraform_remote_state.network.outputs.private_subnets[0].id
  vpc_security_group_ids = [data.terraform_remote_state.network.outputs.node_security_group_id]


  root_block_device {
    volume_size = 10
    volume_type = "gp3"
  }

  user_data = base64encode(<<-EOF
      #!/bin/bash
      wget -q https://repos.influxdata.com/influxdata-archive_compat.key
      echo '393e8779c89ac8d958f81f942f9ad7fb82a25e133faddaf92e15b16e6ac9ce4c influxdata-archive_compat.key' | sha256sum -c && cat influxdata-archive_compat.key | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/influxdata-archive_compat.gpg > /dev/null
      echo 'deb [signed-by=/etc/apt/trusted.gpg.d/influxdata-archive_compat.gpg] https://repos.influxdata.com/debian stable main' | sudo tee /etc/apt/sources.list.d/influxdata.list
      sudo apt-get update
      sudo apt-get install -y influxdb2
      sudo systemctl start influxdb
      sudo systemctl enable influxdb
      
      sudo apt-get install cloudflared
      curl -L --output cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb && 
      sudo dpkg -i cloudflared.deb
      sudo cloudflared service install ${var.cloudflare_tunnel_token}
    EOF
  )
}