resource "aws_dynamodb_table" "terraform_statelock" {
  name           = "terraform-lock-${var.project}-${var.region}-${data.aws_caller_identity.self.account_id}"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = local.tags
}

