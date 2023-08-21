variable "region" {}
variable "environment" {}
variable "project" {}
variable "account" {}
variable "stack" {}
variable "foo" {
  type = bool
}

provider "azurerm" {
  skip_provider_registration = true
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = var.project
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = replace(var.project, "-", "")
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "example" {
  name                  = var.project
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

output "foo" {
  value = azurerm_storage_container.example.id
}

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">=3.68.0"
    }
  }
  backend "azurerm" {
    #resource_group_name  = "RG-Profishop-TEST-databricks"
    #storage_account_name = "profishoptest"
    #container_name       = "jhps-analytics-tf-state-test"
    #key                 = "test.resource_group.tfstate"
  }
}
