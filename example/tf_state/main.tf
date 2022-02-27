provider "aws" {
  region = var.region
}

# in case of an absent terraform bucket, comment this whole block
# deploy this stack to create a local tfstate file
# and enable again to move this state to the s3 bucket as well
# fix bucket_name + region accordingly
terraform {
  backend "s3" {
  }
}

