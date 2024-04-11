data "aws_availability_zones" "available" { state = "available" }

locals {
  azs_count = 2
  azs_names = data.aws_availability_zones.available.names
}

resource "aws_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "${terraform.workspace}-vpc"
  }
}
resource "aws_subnet" "public" {
  vpc_id            = aws_vpc.vpc.id
  count             = local.azs_count
  cidr_block        = cidrsubnet(aws_vpc.vpc.cidr_block, 2, count.index)
  availability_zone = local.azs_names[count.index]
  tags = {
    Name = "${terraform.workspace}-public-subnet"
  }
}

resource "aws_subnet" "private" {
  vpc_id            = aws_vpc.vpc.id
  count             = local.azs_count
  cidr_block        = cidrsubnet(aws_vpc.vpc.cidr_block, 2, count.index + 2)
  availability_zone = local.azs_names[count.index]
  tags = {
    Name = "${terraform.workspace}-private-subnet"
  }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    Name = "${terraform.workspace}-gw"
  }
}

resource "aws_eip" "nat" {
  tags = {
    Name = "${terraform.workspace}-nat-ip"
  }
}


resource "aws_route_table" "public" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }

  tags = {
    Name = "${terraform.workspace}-public-route"
  }
}
resource "aws_route_table_association" "public" {
  count          = local.azs_count
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

module "fck-nat" {
  source = "RaJiska/fck-nat/aws"

  name               = "nat-instance"
  vpc_id             = aws_vpc.vpc.id
  subnet_id          = aws_subnet.public[0].id
  eip_allocation_ids = [aws_eip.nat.allocation_id]

  update_route_table = true
  route_table_id     = aws_route_table.private.id
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "${terraform.workspace}-private-route"
  }
}


resource "aws_route_table_association" "private" {
  count          = local.azs_count
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private.id
}
