terraform {
  required_providers {
    outscale = {
      source  = "outscale-dev/outscale"
      version = "0.5.3"
    }
  }
}

provider "outscale" {
  region = var.region
}
