resource "outscale_oks_project" "project" {
  name                    = "oks-project-tf-208"
  cidr                    = "10.50.0.0/18"
  region                  = "eu-west-2"
  disable_api_termination = false
  description             = "TF208 OKS project"
}

resource "outscale_oks_cluster" "cluster" {
  project_id              = outscale_oks_project.project.id
  admin_whitelist         = ["0.0.0.0/0", "100.0.0.0"]
  cidr_pods               = "10.91.0.0/16"
  cidr_service            = "10.92.0.0/16"
  version                 = "1.31"
  name                    = "oks-cluster-tf-208"
  control_planes          = "cp.mono.master"
  description             = "TF208 OKS cluster"
  disable_api_termination = false
  admission_flags = {
    enable_admission_plugins = ["EventRateLimit"]
  }
  tags = {
    test = "TF-208"
    key  = "value"
  }
}
