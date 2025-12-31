resource "outscale_oks_project" "project" {
  name                    = "test-oks-project-${random_string.suffix[0].result}"
  cidr                    = "10.50.0.0/18"
  region                  = "eu-west-2"
  disable_api_termination = false
  tags = {
    test = "TF-207"
  }
}
