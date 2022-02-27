data "aws_caller_identity" "self" {
}
data "aws_availability_zones" "self" {
}

/*data "aws_iam_policy_document" "state_bucket" {
  statement {
    sid     = "FullAccess"
    actions = ["s3:*"]
    principals {
      type        = "AWS"
      identifiers = [data.aws_iam_user.deploy.arn]
    }

    resources = [
      "arn:aws:s3:::${local.bucket}",
      "arn:aws:s3:::${local.bucket}/*",
    ]
  }
}
*/
