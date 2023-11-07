provider "google" {
  project     = var.project
  region      = var.region
}

variable "region" {}
variable "environment" {}
variable "project" {}
variable "account" {}
variable "stack" {}
variable "foo" {
  type = bool
}


resource "google_storage_bucket" "state" {
  location = "EU"
  name = "${var.project}-${var.account}-test"
}

terraform {
  backend "gcs"  {
  }
}

output "foo" {
  value = google_storage_bucket.state.id
}
