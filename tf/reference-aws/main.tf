# Reference AWS setup for domovyk v0.0 testing.
# Deploy only into a dedicated sandbox AWS account. Not a production account.
# Creating any resource costs money — spend is the owner's decision.
# Tags: domovyk-ref = "true", owner, expires.

terraform {
  required_version = ">= 1.9"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

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

provider "aws" {
  region = var.region
}

# Common tags applied to every resource.
locals {
  common_tags = {
    domovyk-ref = "true"
    owner       = var.owner
    expires     = var.expires
  }
}
