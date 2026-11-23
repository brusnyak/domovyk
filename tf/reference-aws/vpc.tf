# VPC, subnets, internet gateway, route tables.
resource "aws_vpc" "ref" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = merge(local.common_tags, {
    Name = "domovyk-ref-vpc"
  })
}

resource "aws_subnet" "public" {
  count             = 2
  vpc_id            = aws_vpc.ref.id
  cidr_block        = "10.0.${count.index + 1}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  tags = merge(local.common_tags, {
    Name = "domovyk-ref-public-subnet-${count.index + 1}"
  })
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_internet_gateway" "ref" {
  vpc_id = aws_vpc.ref.id
  tags = merge(local.common_tags, {
    Name = "domovyk-ref-igw"
  })
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.ref.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.ref.id
  }
  tags = merge(local.common_tags, {
    Name = "domovyk-ref-public-rt"
  })
}

resource "aws_route_table_association" "public" {
  count          = 2
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_security_group" "web" {
  name        = "domovyk-ref-web"
  description = "Allow HTTP/HTTPS and SSH"
  vpc_id      = aws_vpc.ref.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "domovyk-ref-web-sg"
  })
}

resource "aws_security_group" "db" {
  name        = "domovyk-ref-db"
  description = "Allow PostgreSQL from web SG only"
  vpc_id      = aws_vpc.ref.id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.web.id]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "domovyk-ref-db-sg"
  })
}
