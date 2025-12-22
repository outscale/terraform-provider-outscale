resource "outscale_oks_project" "project" {
  name                    = "test-oks-project-${random_string.suffix[0].result}"
  cidr                    = "10.50.0.0/18"
  region                  = "eu-west-2"
  disable_api_termination = false
  description             = "TF209 OKS project"
}

resource "outscale_oks_cluster" "cluster" {
  project_id      = outscale_oks_project.project.id
  admin_whitelist = ["0.0.0.0/0"]
  cidr_pods       = "10.91.0.0/16"
  cidr_service    = "10.92.0.0/16"
  version         = "1.31"
  name            = "test-oks-cluster-${random_string.suffix[0].result}"
  control_planes  = "cp.mono.master"
  tags = {
    test = "TF-209"
  }
}
