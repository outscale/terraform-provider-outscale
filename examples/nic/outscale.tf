terraform {
  required_providers {
    outscale = {
      source  = "outscale-dev/outscale"
      version = ">= 0.2.0"
    }
  }
}

provider "outscale" {
  access_key_id = var.access_key_id
  secret_key_id = var.secret_key_id
  region        = var.region
}
