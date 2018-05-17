package outscale

// import (
// 	"os"
// 	"strconv"
// 	"testing"

// 	"github.com/hashicorp/terraform/helper/resource"
// 	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
// )

// func TestAccOutscaleDSOAPILBU_basic(t *testing.T) {
// 	o := os.Getenv("OUTSCALE_OAPI")

// 	oapi, err := strconv.ParseBool(o)
// 	if err != nil {
// 		oapi = false
// 	}

// 	if !oapi {
// 		t.Skip()
// 	}

// 	var conf lbu.LoadBalancerDescription

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:      func() { testAccPreCheck(t) },
// 		IDRefreshName: "outscale_load_balancer.bar",
// 		Providers:     testAccProviders,
// 		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccDSOutscaleOAPILBUConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
// 					testAccCheckOutscaleOAPILBUAttributes(&conf),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "sub_region_name.#", "2"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "sub_region_name.0", "eu-west-2a"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "sub_region_name.1", "eu-west-2b"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "listener.0.backend_port", "8000"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "listener.0.backend_protocol", "http"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "listener.0.load_balancer_port", "80"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "listener.0.load_balancer_protocol", "http"),
// 				)},
// 		},
// 	})
// }

// const testAccDSOutscaleOAPILBUConfig = `
// resource "outscale_load_balancer" "bar" {
//   sub_region_name = ["eu-west-2a", "eu-west-2b"]
// 	load_balancer_name               = "foobar-terraform-elb"
//   listener {
//     backend_port = 8000
//     backend_protocol = "http"
//     load_balancer_port = 80
//     // Protocol should be case insensitive
//     load_balancer_protocol = "http"
//   }

// 	tag {
// 		bar = "baz"
// 	}
// }

// data "outscale_load_balancer" "test" {
// 	load_balancer_name = "${outscale_load_balancer.bar.id}"
// }
// `
