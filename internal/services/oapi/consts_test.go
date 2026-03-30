package oapi_test

import "time"

const (
	DefaultTimeout time.Duration = 2 * time.Minute

	testAccVmType         string = "tinav7.c2r2p1"
	testAccVmTypefGPU     string = "tinav5.c2r2p1"
	testAccfGPUGeneration string = "v5"
	testAccfGPUModel      string = "nvidia-p6"

	testAccCertPath string = "testdata/certificate.pem"
	testAccKeypair  string = "terraform-basic"
)
