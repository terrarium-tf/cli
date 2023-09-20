locals {
  # attach these tags to our resources
  tags = {
    environment = "global"
    application = "terraform"
    stack       = var.name
    project     = var.project
  }

  bucket             = "tf-state-${var.project}-${var.region}-${data.aws_caller_identity.self.account_id}"
  access_bucket_name = "${var.project}-access-${var.region}-${data.aws_caller_identity.self.account_id}"
}

