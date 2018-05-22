package outscale

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func resourceOutscaleAppCookieStickinessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleAppCookieStickinessPolicyCreate,
		Read:   resourceOutscaleAppCookieStickinessPolicyRead,
		Delete: resourceOutscaleAppCookieStickinessPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !regexp.MustCompile(`^[0-9A-Za-z-]+$`).MatchString(value) {
						es = append(es, fmt.Errorf(
							"only alphanumeric characters and hyphens allowed in %q", k))
					}
					return
				},
			},

			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"lb_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"cookie_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOutscaleAppCookieStickinessPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	l, lok := d.GetOk("load_balancer_name")
	p, pok := d.GetOk("policy_name")

	if !lok && !pok {
		return fmt.Errorf("please provide the required attributes load_balancer_name and policy_name")
	}

	acspOpts := &lbu.CreateAppCookieStickinessPolicyInput{
		LoadBalancerName: aws.String(l.(string)),
		PolicyName:       aws.String(p.(string)),
	}

	if v, ok := d.GetOk("cookie_name"); ok {
		acspOpts.CookieName = aws.String(v.(string))
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.CreateAppCookieStickinessPolicy(acspOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating AppCookieStickinessPolicy: %s", err)
	}

	setLoadBalancerOpts := &lbu.SetLoadBalancerPoliciesOfListenerInput{
		LoadBalancerName: aws.String(d.Get("load_balancer_name").(string)),
		LoadBalancerPort: aws.Int64(int64(d.Get("lb_port").(int))),
		PolicyNames:      []*string{aws.String(d.Get("policy_name").(string))},
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.SetLoadBalancerPoliciesOfListener(setLoadBalancerOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error setting AppCookieStickinessPolicy: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%d:%s",
		*acspOpts.LoadBalancerName,
		*setLoadBalancerOpts.LoadBalancerPort,
		*acspOpts.PolicyName))

	return resourceOutscaleAppCookieStickinessPolicyRead(d, meta)
}

func resourceOutscaleAppCookieStickinessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*OutscaleClient).LBU

	lbName, lbPort, policyName := resourceOutscaleAppCookieStickinessPolicyParseID(d.Id())

	request := &lbu.DescribeLoadBalancerPoliciesInput{
		LoadBalancerName: aws.String(lbName),
		PolicyNames:      []*string{aws.String(policyName)},
	}

	var getResp *lbu.DescribeLoadBalancerPoliciesOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = elbconn.API.DescribeLoadBalancerPolicies(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "PolicyNotFound") || strings.Contains(fmt.Sprint(err), "LoadBalancerNotFound") {
			log.Printf("[WARN] Load Balancer / Load Balancer Policy (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving policy: %s", err)
	}
	if len(getResp.PolicyDescriptions) != 1 {
		return fmt.Errorf("Unable to find policy %#v", getResp.PolicyDescriptions)
	}

	// we know the policy exists now, but we have to check if it's assigned to a listener
	assigned, err := resourceAwsELBSticknessPolicyAssigned(policyName, lbName, lbPort, elbconn)
	if err != nil {
		return err
	}
	if !assigned {
		// policy exists, but isn't assigned to a listener
		log.Printf("[DEBUG] policy '%s' exists, but isn't assigned to a listener", policyName)
		d.SetId("")
		return nil
	}

	policyDesc := getResp.PolicyDescriptions[0]
	cookieAttr := policyDesc.PolicyAttributeDescriptions[0]
	if *cookieAttr.AttributeName != "CookieName" {
		return fmt.Errorf("Unable to find cookie Name")
	}

	d.Set("cookie_name", cookieAttr.AttributeValue)
	d.Set("policy_name", policyName)
	d.Set("load_balancer_name", lbName)

	return nil
}

// Determine if a particular policy is assigned to an ELB listener
func resourceAwsELBSticknessPolicyAssigned(policyName, lbName, lbPort string, elbconn *lbu.Client) (bool, error) {
	describeElbOpts := &lbu.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String(lbName)},
	}

	var describeResp *lbu.DescribeLoadBalancersOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		describeResp, err = elbconn.API.DescribeLoadBalancers(describeElbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "LoadBalancerNotFound") {
			return false, nil
		}
		return false, fmt.Errorf("Error retrieving ELB description: %s", err)
	}

	if len(describeResp.LoadBalancerDescriptions) != 1 {
		return false, fmt.Errorf("Unable to find ELB: %#v", describeResp.LoadBalancerDescriptions)
	}

	lb := describeResp.LoadBalancerDescriptions[0]
	assigned := false
	for _, listener := range lb.ListenerDescriptions {
		if lbPort != strconv.Itoa(int(*listener.Listener.LoadBalancerPort)) {
			continue
		}

		for _, name := range listener.PolicyNames {
			if policyName == *name {
				assigned = true
				break
			}
		}
	}

	return assigned, nil
}

func resourceOutscaleAppCookieStickinessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*OutscaleClient).LBU

	lbName, _, policyName := resourceOutscaleAppCookieStickinessPolicyParseID(d.Id())

	setLoadBalancerOpts := &lbu.SetLoadBalancerPoliciesOfListenerInput{
		LoadBalancerName: aws.String(d.Get("load_balancer_name").(string)),
		PolicyNames:      []*string{},
	}

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = elbconn.API.SetLoadBalancerPoliciesOfListener(setLoadBalancerOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error removing AppCookieStickinessPolicy: %s", err)
	}

	request := &lbu.DeleteLoadBalancerPolicyInput{
		LoadBalancerName: aws.String(lbName),
		PolicyName:       aws.String(policyName),
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = elbconn.API.DeleteLoadBalancerPolicy(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting App stickiness policy %s: %s", d.Id(), err)
	}
	return nil
}

func resourceOutscaleAppCookieStickinessPolicyParseID(id string) (string, string, string) {
	parts := strings.SplitN(id, ":", 3)
	return parts[0], parts[1], parts[2]
}
