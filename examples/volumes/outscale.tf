terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = ">= 0.11.0"
    }
  }
}

provider "outscale" {
  access_key_id = var.access_key_id
  secret_key_id = var.secret_key_id
  api {
    region = var.region
  }
}
