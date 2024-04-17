variable "aws_region" {
  default = "ap-southeast-2"
}

variable "ecr_repository_url" {
  description = "The URL of the ECR repository"
}

variable "cloudflare_tunnel_token" {
  description = "Cloudflare tunnel token"
}

variable "backend_port" {
  default = 8000
}

variable "frontend_url" {
  default = "https://dev-fe-wgp.layer8proxy.net"
}

variable "backend_url" {
  default = "https://dev-be-wgp.layer8proxy.net"
}

variable "layer8_url" {
  default = "https://dev.layer8proxy.net"
}
