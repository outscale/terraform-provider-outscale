terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = "0.5.32"
    }
  }
}

provider "outscale" {
  region = var.region
}
