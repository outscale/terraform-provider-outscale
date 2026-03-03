terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = ">= 1.4.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = ">= 4.0.0"
    }
    local = {
      source  = "hashicorp/local"
      version = ">= 2.0.0"
    }
  }
}

provider "outscale" {
  access_key_id = var.access_key
  secret_key_id = var.secret_key
  api {
    region = var.region
  }
}