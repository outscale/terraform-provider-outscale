resource "outscale_oks_project" "project" {
  name                    = "oks-project-tf-207"
  cidr                    = "10.50.0.0/18"
  region                  = "eu-west-2"
  disable_api_termination = false
  description             = "TF207 OKS project"
  tags = {
    test = "TF-207"
    key  = "value"
  }
}
