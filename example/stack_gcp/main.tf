provider "aws" {
  region = "eu-central-1"
}

data "aws_caller_identity" "self" {}

variable "region" {}
variable "environment" {}
variable "project" {}
variable "account" {}
variable "stack" {}
variable "foo" {
  type = bool
}

resource "aws_s3_bucket" "test" {
  bucket = "test-${var.environment}-${data.aws_caller_identity.self.account_id}"
}

resource "aws_s3_bucket" "state" {}

terraform {
  backend "gcs"  {
  }
}

output "foo" {
  value = "test-${var.environment}-${data.aws_caller_identity.self.account_id}"
}
