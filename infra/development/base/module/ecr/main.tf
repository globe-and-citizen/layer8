resource "aws_ecr_repository" "layer8_ecr_repo" {
  name                 = "layer8"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}
