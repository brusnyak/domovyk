variable "region" {
  description = "AWS region for the reference setup (EU region preferred)"
  type        = string
  default     = "eu-central-1"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.medium"
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "owner" {
  description = "Owner tag value"
  type        = string
  default     = "domovyk"
}

variable "expires" {
  description = "Expiration tag (date the sandbox should be torn down)"
  type        = string
  default     = "2026-12-01"
}

variable "db_password" {
  description = "RDS master password"
  type        = string
  sensitive   = true
}
