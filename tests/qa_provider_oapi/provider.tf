terraform {
  required_providers {
    outscale = {
      source  = "outscale-dev/outscale"
      version = "0.5.1"
    }
  }
}

provider "outscale" {}
