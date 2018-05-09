package outscale

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleLBU_basic(t *testing.T) {
	var conf lbu.LoadBalancerDescription

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testAccCheckOutscaleLBUAttributes(&conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.#", "2"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.0", "eu-west-2a"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.1", "eu-west-2b"),
					// resource.TestCheckResourceAttr(
					// 	"outscale_load_balancer.bar", "subnets.#", "3"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.0.instance_port", "8000"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.0.instance_protocol", "http"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.0.load_balancer_port", "80"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.0.protocol", "http"),
				)},
		},
	})
}

func TestAccOutscaleLBU_availabilityZones(t *testing.T) {
	var conf lbu.LoadBalancerDescription

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.#", "3"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.2487133097", "eu-west-2a"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.221770259", "eu-west-2b"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.2050015877", "eu-west-2c"),
				),
			},

			{
				Config: testAccOutscaleLBUConfig_AvailabilityZonesUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.#", "2"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.2487133097", "eu-west-2a"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.221770259", "eu-west-2b"),
				),
			},
		},
	})
}

func TestAccOutscaleLBU_swap_subnets(t *testing.T) {
	var conf lbu.LoadBalancerDescription

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.ourapp",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfig_subnets,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.ourapp", &conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.ourapp", "subnets.#", "2"),
				),
			},

			{
				Config: testAccOutscaleLBUConfig_subnet_swap,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.ourapp", &conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.ourapp", "subnets.#", "2"),
				),
			},
		},
	})
}

func TestAccOutscaleLBU_InstanceAttaching(t *testing.T) {
	var conf lbu.LoadBalancerDescription

	testCheckInstanceAttached := func(count int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(conf.Instances) != count {
				return fmt.Errorf("instance count does not match")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testAccCheckOutscaleLBUAttributes(&conf),
				),
			},

			{
				Config: testAccOutscaleLBUConfigNewInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testCheckInstanceAttached(1),
				),
			},
		},
	})
}

func TestAccOutscaleLBUUpdate_Listener(t *testing.T) {
	var conf lbu.LoadBalancerDescription

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testAccCheckOutscaleLBUAttributes(&conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.206423021.instance_port", "8000"),
				),
			},

			{
				Config: testAccOutscaleLBUConfigListener_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.3931999347.instance_port", "8080"),
				),
			},
		},
	})
}

func TestAccOutscaleLBU_HealthCheck(t *testing.T) {
	var conf lbu.LoadBalancerDescription

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfigHealthCheck,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testAccCheckOutscaleLBUAttributesHealthCheck(&conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "health_check.0.healthy_threshold", "5"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "health_check.0.unhealthy_threshold", "5"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "health_check.0.target", "HTTP:8000/"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "health_check.0.timeout", "30"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "health_check.0.interval", "60"),
				),
			},
		},
	})
}

func TestAccOutscaleLBUUpdate_HealthCheck(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfigHealthCheck,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "health_check.0.healthy_threshold", "5"),
				),
			},
			{
				Config: testAccOutscaleLBUConfigHealthCheck_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "health_check.0.healthy_threshold", "10"),
				),
			},
		},
	})
}

func TestAccOutscaleLBU_SecurityGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfig,
				Check: resource.ComposeTestCheckFunc(
					// ELBs get a default security group
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "security_groups.#", "1",
					),
				),
			},
			{
				Config: testAccOutscaleLBUConfigSecurityGroups,
				Check: resource.ComposeTestCheckFunc(
					// Count should still be one as we swap in a ceutom security group
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "security_groups.#", "1",
					),
				),
			},
		},
	})
}

func testAccCheckOutscaleLBUDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).LBU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_load_balancer" {
			continue
		}

		var err error
		var describe *lbu.DescribeLoadBalancersOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			describe, err = conn.API.DescribeLoadBalancers(&lbu.DescribeLoadBalancersInput{
				LoadBalancerNames: []*string{aws.String(rs.Primary.ID)},
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "CertificateNotFound") {
					return resource.RetryableError(
						fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err == nil {
			if len(describe.LoadBalancerDescriptions) != 0 &&
				*describe.LoadBalancerDescriptions[0].LoadBalancerName == rs.Primary.ID {
				return fmt.Errorf("LBU still exists")
			}
		}

		// Verify the error
		providerErr, ok := err.(awserr.Error)
		if !ok {
			return err
		}

		if providerErr.Code() != "LoadBalancerNotFound" {
			return fmt.Errorf("Unexpected error: %s", err)
		}
	}

	return nil
}

func testAccCheckOutscaleLBUAttributes(conf *lbu.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		zones := []string{"eu-west-2a", "eu-west-2b", "eu-west-2c"}
		azs := make([]string, 0, len(conf.AvailabilityZones))
		for _, x := range conf.AvailabilityZones {
			azs = append(azs, *x)
		}
		sort.StringSlice(azs).Sort()
		if !reflect.DeepEqual(azs, zones) {
			return fmt.Errorf("bad availability_zones_member")
		}

		l := lbu.Listener{
			InstancePort:     aws.Int64(int64(8000)),
			InstanceProtocol: aws.String("HTTP"),
			LoadBalancerPort: aws.Int64(int64(80)),
			Protocol:         aws.String("HTTP"),
		}

		if !reflect.DeepEqual(conf.ListenerDescriptions[0].Listener, &l) {
			return fmt.Errorf(
				"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
				conf.ListenerDescriptions[0].Listener,
				l)
		}

		if *conf.DNSName == "" {
			return fmt.Errorf("empty dns_name")
		}

		return nil
	}
}

func testAccCheckOutscaleLBUAttributesHealthCheck(conf *lbu.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		zones := []string{"eu-west-2a", "eu-west-2b", "eu-west-2c"}
		azs := make([]string, 0, len(conf.AvailabilityZones))
		for _, x := range conf.AvailabilityZones {
			azs = append(azs, *x)
		}
		sort.StringSlice(azs).Sort()
		if !reflect.DeepEqual(azs, zones) {
			return fmt.Errorf("bad availability_zones_member")
		}

		check := &lbu.HealthCheck{
			Timeout:            aws.Int64(int64(30)),
			UnhealthyThreshold: aws.Int64(int64(5)),
			HealthyThreshold:   aws.Int64(int64(5)),
			Interval:           aws.Int64(int64(60)),
			Target:             aws.String("HTTP:8000/"),
		}

		if !reflect.DeepEqual(conf.HealthCheck, check) {
			return fmt.Errorf(
				"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
				conf.HealthCheck,
				check)
		}

		if *conf.DNSName == "" {
			return fmt.Errorf("empty dns_name")
		}

		return nil
	}
}

func testAccCheckOutscaleLBUExists(n string, res *lbu.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LBU ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).LBU

		var err error
		var describe *lbu.DescribeLoadBalancersOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			describe, err = conn.API.DescribeLoadBalancers(&lbu.DescribeLoadBalancersInput{
				LoadBalancerNames: []*string{aws.String(rs.Primary.ID)},
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "CertificateNotFound") {
					return resource.RetryableError(
						fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		if len(describe.LoadBalancerDescriptions) != 1 ||
			*describe.LoadBalancerDescriptions[0].LoadBalancerName != rs.Primary.ID {
			return fmt.Errorf("LBU not found")
		}

		*res = *describe.LoadBalancerDescriptions[0]

		// Confirm source_security_group_id for ELBs in a VPC
		// 	See https://github.com/hashicorp/terraform/pull/3780
		if res.VPCId != nil {
			sgid := rs.Primary.Attributes["source_security_group_id"]
			if sgid == "" {
				return fmt.Errorf("Expected to find source_security_group_id for LBU, but was empty")
			}
		}

		return nil
	}
}

const testAccOutscaleLBUConfig = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b"]
	load_balancer_name               = "foobar-terraform-elb"
  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    // Protocol should be case insensitive
    protocol = "http"
  }

	tag {
		bar = "baz"
	}

}
`

const testAccOutscaleLBUFullRangeOfCharacters = `
resource "outscale_load_balancer" "foo" {
  name = "%s"
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }
}
`

const testAccOutscaleLBUAccessLogs = `
resource "outscale_load_balancer" "foo" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }
}
`

func testAccOutscaleLBUAccessLogsOn(r string) string {
	return fmt.Sprintf(`
# an S3 bucket configured for Access logs
# The 797873946194 is the AWS ID for eu-west-2, so this test
# meut be ran in eu-west-2
resource "aws_s3_bucket" "acceslogs_bucket" {
  bucket = "%s"
  acl = "private"
  force_destroy = true
  policy = <<EOF
{
  "Id": "Policy1446577137248",
  "Statement": [
    {
      "Action": "s3:PutObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::797873946194:root"
      },
      "Resource": "arn:aws:s3:::%s/*",
      "Sid": "Stmt1446575236270"
    }
  ],
  "Version": "2012-10-17"
}
EOF
}

resource "outscale_load_balancer" "foo" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }

	access_logs {
		interval = 5
		bucket = "${aws_s3_bucket.acceslogs_bucket.bucket}"
	}
}
`, r, r)
}

func testAccOutscaleLBUAccessLogsDisabled(r string) string {
	return fmt.Sprintf(`
# an S3 bucket configured for Access logs
# The 797873946194 is the AWS ID for eu-west-2, so this test
# meut be ran in eu-west-2
resource "aws_s3_bucket" "acceslogs_bucket" {
  bucket = "%s"
  acl = "private"
  force_destroy = true
  policy = <<EOF
{
  "Id": "Policy1446577137248",
  "Statement": [
    {
      "Action": "s3:PutObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::797873946194:root"
      },
      "Resource": "arn:aws:s3:::%s/*",
      "Sid": "Stmt1446575236270"
    }
  ],
  "Version": "2012-10-17"
}
EOF
}

resource "outscale_load_balancer" "foo" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }

	access_logs {
		interval = 5
		bucket = "${aws_s3_bucket.acceslogs_bucket.bucket}"
		enabled = false
	}
}
`, r, r)
}

const testAccOutscaleLBU_namePrefix = `
resource "outscale_load_balancer" "test" {
  name_prefix = "test-"
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }
}
`

const testAccOutscaleLBUGeneratedName = `
resource "outscale_load_balancer" "foo" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }
}
`

const testAccOutscaleLBU_zeroValueName = `
resource "outscale_load_balancer" "foo" {
  name               = ""
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }
}
`

const testAccOutscaleLBUConfig_AvailabilityZonesUpdate = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }
}
`

const testAccOutscaleLBUConfig_TagUpdate = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }

	tag {
		foo = "bar"
		new = "type"
	}

}
`

const testAccOutscaleLBUConfigNewInstance = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }

  instances = ["${aws_instance.foo.id}"]
}

resource "aws_instance" "foo" {
	# eu-west-2
	ami = "ami-043a5034"
	instance_type = "t1.micro"
}
`

const testAccOutscaleLBUConfigListenerSSLCertificateId = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    ssl_certificate_id = "%s"
    load_balancer_port = 443
    protocol = "https"
  }
}
`

const testAccOutscaleLBUConfigHealthCheck = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }

  health_check {
    healthy_threshold = 5
    unhealthy_threshold = 5
    target = "HTTP:8000/"
    interval = 60
    timeout = 30
  }
}
`

const testAccOutscaleLBUConfigHealthCheck_update = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }

  health_check {
    healthy_threshold = 10
    unhealthy_threshold = 5
    target = "HTTP:8000/"
    interval = 60
    timeout = 30
  }
}
`

const testAccOutscaleLBUConfigListener_update = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8080
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }
}
`

const testAccOutscaleLBUConfigIdleTimeout = `
resource "outscale_load_balancer" "bar" {
	availability_zones_member = ["eu-west-2a"]

	listeners_member {
		instance_port = 8000
		instance_protocol = "http"
		load_balancer_port = 80
		protocol = "http"
	}

	idle_timeout = 200
}
`

const testAccOutscaleLBUConfigIdleTimeout_update = `
resource "outscale_load_balancer" "bar" {
	availability_zones_member = ["eu-west-2a"]

	listeners_member {
		instance_port = 8000
		instance_protocol = "http"
		load_balancer_port = 80
		protocol = "http"
	}

	idle_timeout = 400
}
`

const testAccOutscaleLBUConfigConnectionDraining = `
resource "outscale_load_balancer" "bar" {
	availability_zones_member = ["eu-west-2a"]

	listeners_member {
		instance_port = 8000
		instance_protocol = "http"
		load_balancer_port = 80
		protocol = "http"
	}

	connection_draining = true
	connection_draining_timeout = 400
}
`

const testAccOutscaleLBUConfigConnectionDraining_update_timeout = `
resource "outscale_load_balancer" "bar" {
	availability_zones_member = ["eu-west-2a"]

	listeners_member {
		instance_port = 8000
		instance_protocol = "http"
		load_balancer_port = 80
		protocol = "http"
	}

	connection_draining = true
	connection_draining_timeout = 600
}
`

const testAccOutscaleLBUConfigConnectionDraining_update_disable = `
resource "outscale_load_balancer" "bar" {
	availability_zones_member = ["eu-west-2a"]

	listeners_member {
		instance_port = 8000
		instance_protocol = "http"
		load_balancer_port = 80
		protocol = "http"
	}

	connection_draining = false
}
`

const testAccOutscaleLBUConfigSecurityGroups = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    protocol = "http"
  }

  security_groups = ["${aws_security_group.bar.id}"]
}

resource "aws_security_group" "bar" {
  ingress {
    protocol = "tcp"
    from_port = 80
    to_port = 80
    cidr_blocks = ["0.0.0.0/0"]
  }

	tag {
		Name = "tf_elb_sg_test"
	}
}
`

// This IAM Server config is lifted from
// builtin/providers/aws/resource_aws_iam_server_certificate_test.go
func testAccELBIAMServerCertConfig(certName string) string {
	return fmt.Sprintf(`
resource "aws_iam_server_certificate" "test_cert" {
  name = "%s"
  certificate_body = <<EOF
-----BEGIN CERTIFICATE-----
MIIDBjCCAe4CCQCGWwBmOiHQdTANBgkqhkiG9w0BAQUFADBFMQswCQYDVQQGEwJB
VTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50ZXJuZXQgV2lkZ2l0
cyBQdHkgTHRkMB4XDTE2MDYyMTE2MzM0MVoXDTE3MDYyMTE2MzM0MVowRTELMAkG
A1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAfBgNVBAoTGEludGVybmV0
IFdpZGdpdHMgUHR5IEx0ZDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AL+LFlsCJG5txZp4yuu+lQnuUrgBXRG+irQqcTXlV91Bp5hpmRIyhnGCtWxxDBUL
xrh4WN3VV/0jDzKT976oLgOy3hj56Cdqf+JlZ1qgMN5bHB3mm3aVWnrnsLbBsfwZ
SEbk3Kht/cE1nK2toNVW+rznS3m+eoV3Zn/DUNwGlZr42hGNs6ETn2jURY78ETqR
mW47xvjf86eIo7vULHJaY6xyarPqkL8DZazOmvY06hUGvGwGBny7gugfXqDG+I8n
cPBsGJGSAmHmVV8o0RCB9UjY+TvSMQRpEDoVlvyrGuglsD8to/4+7UcsuDGlRYN6
jmIOC37mOi/jwRfWL1YUa4MCAwEAATANBgkqhkiG9w0BAQUFAAOCAQEAPDxTH0oQ
JjKXoJgkmQxurB81RfnK/NrswJVzWbOv6ejcbhwh+/ZgJTMc15BrYcxU6vUW1V/i
Z7APU0qJ0icECACML+a2fRI7YdLCTiPIOmY66HY8MZHAn3dGjU5TeiUflC0n0zkP
mxKJe43kcYLNDItbfvUDo/GoxTXrC3EFVZyU0RhFzoVJdODlTHXMVFCzcbQEBrBJ
xKdShCEc8nFMneZcGFeEU488ntZoWzzms8/QpYrKa5S0Sd7umEU2Kwu4HTkvUFg/
CqDUFjhydXxYRsxXBBrEiLOE5BdtJR1sH/QHxIJe23C9iHI2nS1NbLziNEApLwC4
GnSud83VUo9G9w==
-----END CERTIFICATE-----
EOF

	private_key =  <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAv4sWWwIkbm3FmnjK676VCe5SuAFdEb6KtCpxNeVX3UGnmGmZ
EjKGcYK1bHEMFQvGuHhY3dVX/SMPMpP3vqguA7LeGPnoJ2p/4mVnWqAw3lscHeab
dpVaeuewtsGx/BlIRuTcqG39wTWcra2g1Vb6vOdLeb56hXdmf8NQ3AaVmvjaEY2z
oROfaNRFjvwROpGZbjvG+N/zp4iju9QsclpjrHJqs+qQvwNlrM6a9jTqFQa8bAYG
fLuC6B9eoMb4jydw8GwYkZICYeZVXyjREIH1SNj5O9IxBGkQOhWW/Ksa6CWwPy2j
/j7tRyy4MaVFg3qOYg4LfuY6L+PBF9YvVhRrgwIDAQABAoIBAFqJ4h1Om+3e0WK8
6h4YzdYN4ue7LUTv7hxPW4gASlH5cMDoWURywX3yLNN/dBiWom4b5NWmvJqY8dwU
eSyTznxNFhJ0PjozaxOWnw4FXlQceOPhV2bsHgKudadNU1Y4lSN9lpe+tg2Xy+GE
ituM66RTKCf502w3DioiJpx6OEkxuhrnsQAWNcGB0MnTukm2f+629V+04R5MT5V1
nY+5Phx2BpHgYzWBKh6Px1puu7xFv5SMQda1ndlPIKb4cNp0yYn+1lHNjbOE7QL/
oEpWgrauS5Zk/APK33v/p3wVYHrKocIFHlPiCW0uIJJLsOZDY8pQXpTlc+/xGLLy
WBu4boECgYEA6xO+1UNh6ndJ3xGuNippH+ucTi/uq1+0tG1bd63v+75tn5l4LyY2
CWHRaWVlVn+WnDslkQTJzFD68X+9M7Cc4oP6WnhTyPamG7HlGv5JxfFHTC9GOKmz
sSc624BDmqYJ7Xzyhe5kc3iHzqG/L72ZF1aijZdrodQMSY1634UX6aECgYEA0Jdr
cBPSN+mgmEY6ogN5h7sO5uNV3TQQtW2IslfWZn6JhSRF4Rf7IReng48CMy9ZhFBy
Q7H2I1pDGjEC9gQHhgVfm+FyMSVqXfCHEW/97pvvu9ougHA0MhPep1twzTGrqg+K
f3PLW8hVkGyCrTfWgbDlPsHgsocA/wTaQOheaqMCgYBat5z+WemQfQZh8kXDm2xE
KD2Cota9BcsLkeQpdFNXWC6f167cqydRSZFx1fJchhJOKjkeFLX3hgzBY6VVLEPu
2jWj8imLNTv3Fhiu6RD5NVppWRkFRuAUbmo1SPNN2+Oa5YwGCXB0a0Alip/oQYex
zPogIB4mLlmrjNCtL4SB4QKBgCEHKMrZSJrz0irqS9RlanPUaZqjenAJE3A2xMNA
Z0FZXdsIEEyA6JGn1i1dkoKaR7lMp5sSbZ/RZfiatBZSMwLEjQv4mYUwoHP5Ztma
+wEyDbaX6G8L1Sfsv3+OWgETkVPfHBXsNtH0mZ/BnrtgsQVeBh52wmZiPAUlNo26
fWCzAoGBAJOjqovLelLWzyQGqPFx/MwuI56UFXd1CmFlCIvF2WxCFmk3tlExoCN1
HqSpt92vsgYgV7+lAb4U7Uy/v012gwiU1LK+vyAE9geo3pTjG73BNzG4H547xtbY
dg+Sd4Wjm89UQoUUoiIcstY7FPbqfBtYKfh4RYHAHV2BwDFqzZCM
-----END RSA PRIVATE KEY-----
EOF
}

resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]

  listeners_member {
    instance_port = 8000
    instance_protocol = "https"
    load_balancer_port = 80
    // Protocol should be case insensitive
    protocol = "HttPs"
    ssl_certificate_id = "${aws_iam_server_certificate.test_cert.arn}"
  }

	tag {
		bar = "baz"
	}

}
`, certName)
}

const testAccOutscaleLBUConfig_subnets = `
provider "aws" {
  region = "eu-west-2"
}

resource "aws_vpc" "azelb" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_hostnames = true

  tag {
    Name = "subnet-vpc"
  }
}

resource "aws_subnet" "public_a_one" {
  vpc_id = "${aws_vpc.azelb.id}"

  cidr_block        = "10.1.1.0/24"
  availability_zone = "eu-west-2a"
}

resource "aws_subnet" "public_b_one" {
  vpc_id = "${aws_vpc.azelb.id}"

  cidr_block        = "10.1.7.0/24"
  availability_zone = "eu-west-2b"
}

resource "aws_subnet" "public_a_two" {
  vpc_id = "${aws_vpc.azelb.id}"

  cidr_block        = "10.1.2.0/24"
  availability_zone = "eu-west-2a"
}

resource "outscale_load_balancer" "ourapp" {
  name = "terraform-asg-deployment-example"

  subnets = [
    "${aws_subnet.public_a_one.id}",
    "${aws_subnet.public_b_one.id}",
  ]

  listeners_member {
    instance_port     = 80
    instance_protocol = "http"
    load_balancer_port           = 80
    protocol       = "http"
  }

  depends_on = ["aws_internet_gateway.gw"]
}

resource "aws_internet_gateway" "gw" {
  vpc_id = "${aws_vpc.azelb.id}"

  tag {
    Name = "main"
  }
}
`

const testAccOutscaleLBUConfig_subnet_swap = `
provider "aws" {
  region = "eu-west-2"
}

resource "aws_vpc" "azelb" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_hostnames = true

  tag {
    Name = "subnet-vpc"
  }
}

resource "aws_subnet" "public_a_one" {
  vpc_id = "${aws_vpc.azelb.id}"

  cidr_block        = "10.1.1.0/24"
  availability_zone = "eu-west-2a"
}

resource "aws_subnet" "public_b_one" {
  vpc_id = "${aws_vpc.azelb.id}"

  cidr_block        = "10.1.7.0/24"
  availability_zone = "eu-west-2b"
}

resource "aws_subnet" "public_a_two" {
  vpc_id = "${aws_vpc.azelb.id}"

  cidr_block        = "10.1.2.0/24"
  availability_zone = "eu-west-2a"
}

resource "outscale_load_balancer" "ourapp" {
  name = "terraform-asg-deployment-example"

  subnets = [
    "${aws_subnet.public_a_two.id}",
    "${aws_subnet.public_b_one.id}",
  ]

  listeners_member {
    instance_port     = 80
    instance_protocol = "http"
    load_balancer_port           = 80
    protocol       = "http"
  }

  depends_on = ["aws_internet_gateway.gw"]
}

resource "aws_internet_gateway" "gw" {
  vpc_id = "${aws_vpc.azelb.id}"

  tag {
    Name = "main"
  }
}
`
