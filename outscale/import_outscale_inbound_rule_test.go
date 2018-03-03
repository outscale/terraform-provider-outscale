package outscale

// func TestAccOutscaleInboundRule_importBasic(t *testing.T) {
// 	resourceName := "outscale_inbound_rule.ingress_1"

// 	rInt := acctest.RandInt()

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckOutscaleSGRuleDestroy,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: testAccOutscaleSecurityGroupRuleIngressConfig(rInt),
// 			},

// 			resource.TestStep{
// 				ResourceName:            resourceName,
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"associate_public_ip_address", "user_data", "security_group"},
// 			},
// 		},
// 	})
// }
