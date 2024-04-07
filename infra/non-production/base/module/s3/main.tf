resource "aws_s3_bucket" "ecs_env_bucket" {
  bucket = "ecsenv"

  tags = {
    Name        = "ECS Environment Variable"
  }
}