# RDS PostgreSQL in private subnet.
resource "aws_db_subnet_group" "ref" {
  name       = "domovyk-ref-db-subnet-group"
  subnet_ids = [aws_subnet.public[1].id]
  tags       = local.common_tags
}

resource "aws_db_instance" "ref" {
  identifier             = "domovyk-ref-db"
  engine                 = "postgres"
  engine_version         = "15.4"
  instance_class         = var.db_instance_class
  allocated_storage      = 20
  db_subnet_group_name   = aws_db_subnet_group.ref.id
  vpc_security_group_ids = [aws_security_group.db.id]
  db_name                = "refdb"
  username               = "domovyk"
  password               = var.db_password
  publicly_accessible    = false
  multi_az               = false
  skip_final_snapshot    = true
  deletion_protection    = false
  storage_encrypted      = true

  tags = merge(local.common_tags, {
    Name = "domovyk-ref-db"
  })
}

variable "db_password" {
  description = "RDS master password"
  type        = string
  sensitive   = true
}
