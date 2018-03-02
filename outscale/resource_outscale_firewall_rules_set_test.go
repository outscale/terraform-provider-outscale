package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestIpPermissionIDHash(t *testing.T) {
	simple := &fcu.IpPermission{
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(int64(80)),
		ToPort:     aws.Int64(int64(8000)),
		IpRanges: []*fcu.IpRange{
			{
				CidrIp: aws.String("10.0.0.0/8"),
			},
		},
	}

	egress := &fcu.IpPermission{
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(int64(80)),
		ToPort:     aws.Int64(int64(8000)),
		IpRanges: []*fcu.IpRange{
			{
				CidrIp: aws.String("10.0.0.0/8"),
			},
		},
	}

	egress_all := &fcu.IpPermission{
		IpProtocol: aws.String("-1"),
		IpRanges: []*fcu.IpRange{
			{
				CidrIp: aws.String("10.0.0.0/8"),
			},
		},
	}

	vpc_security_group_source := &fcu.IpPermission{
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(int64(80)),
		ToPort:     aws.Int64(int64(8000)),
		UserIdGroupPairs: []*fcu.UserIdGroupPair{
			{
				UserId:  aws.String("987654321"),
				GroupId: aws.String("sg-12345678"),
			},
			{
				UserId:  aws.String("123456789"),
				GroupId: aws.String("sg-987654321"),
			},
			{
				UserId:  aws.String("123456789"),
				GroupId: aws.String("sg-12345678"),
			},
		},
	}

	security_group_source := &fcu.IpPermission{
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(int64(80)),
		ToPort:     aws.Int64(int64(8000)),
		UserIdGroupPairs: []*fcu.UserIdGroupPair{
			{
				UserId:    aws.String("987654321"),
				GroupName: aws.String("my-security-group"),
			},
			{
				UserId:    aws.String("123456789"),
				GroupName: aws.String("my-security-group"),
			},
			{
				UserId:    aws.String("123456789"),
				GroupName: aws.String("my-other-security-group"),
			},
		},
	}

	// hardcoded hashes, to detect future change
	cases := []struct {
		Input  *fcu.IpPermission
		Type   string
		Output string
	}{
		{simple, "ingress", "sgrule-3403497314"},
		{egress, "egress", "sgrule-1173186295"},
		{egress_all, "egress", "sgrule-766323498"},
		{vpc_security_group_source, "egress", "sgrule-351225364"},
		{security_group_source, "egress", "sgrule-2198807188"},
	}

	for _, tc := range cases {
		actual := ipPermissionIDHash("sg-12345", tc.Type, tc.Input)
		if actual != tc.Output {
			t.Errorf("input: %s - %s\noutput: %s", tc.Type, tc.Input, actual)
		}
	}
}

func TestAccOutscaleSecurityGroupRule_Ingress_VPC(t *testing.T) {
	var group fcu.SecurityGroup
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleIngressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSecurityGroupRuleExists("outscale_firewall_rules_set.web", &group),
					resource.TestCheckResourceAttr(
						"outscale_firewall_rules_set.web", "group_name", fmt.Sprintf("terraform_test_%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckOutscaleSecurityGroupRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_firewall_rules_set" {
			continue
		}

		// Retrieve our group
		req := &fcu.DescribeSecurityGroupsInput{
			GroupIds: []*string{aws.String(rs.Primary.ID)},
		}

		var resp *fcu.DescribeSecurityGroupsOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeSecurityGroups(req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err == nil {
			if len(resp.SecurityGroups) > 0 && *resp.SecurityGroups[0].GroupId == rs.Primary.ID {
				return fmt.Errorf("Security Group (%s) still exists.", rs.Primary.ID)
			}

			return nil
		}

		if resp == nil {
			return nil
		}

		if strings.Contains(err.Error(), "InvalidGroup.NotFound") {
			return err
		}

		return err
	}

	return nil
}

func testAccCheckOutscaleSecurityGroupRuleExists(n string, group *fcu.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		req := &fcu.DescribeSecurityGroupsInput{
			GroupIds: []*string{aws.String(rs.Primary.ID)},
		}

		var resp *fcu.DescribeSecurityGroupsOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeSecurityGroups(req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		if len(resp.SecurityGroups) > 0 && *resp.SecurityGroups[0].GroupId == rs.Primary.ID {
			*group = *resp.SecurityGroups[0]
			return nil
		}

		return fmt.Errorf("Security Group not found")
	}
}

func testAccOutscaleSecurityGroupRuleIngressConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "web" {
		group_name = "terraform_test_%d"
		group_description = "Used in the terraform acceptance tests"
		tags {
						Name = "tf-acc-test"
		}
		vpc_id = "vpc-e9d09d63"
	}`, rInt)
}
