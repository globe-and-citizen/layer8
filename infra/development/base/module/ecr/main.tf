resource "aws_ecr_repository" "layer8_ecr_repo" {
  name                 = "layer8"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "wgp_ecr_repo" {
  name                 = "wgp"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}