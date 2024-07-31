terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = "1.0.0-rc.1"
    }
  }
}

provider "outscale" {
  access_key_id = var.access_key_id
  secret_key_id = var.secret_key_id
  region        = var.region
}
