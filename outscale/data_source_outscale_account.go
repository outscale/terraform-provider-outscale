package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleAccountRead,

		Schema: map[string]*schema.Schema{
			"account_pid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"city": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"company_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"country": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"customer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mobile_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phone_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vat_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zipcode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleAccountRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).ICU
	request := &icu.ReadAccountInput{}

	var resp *icu.ReadAccountOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.GetAccount(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NotFound") {
			d.SetId("")
			return nil
		}
		fmt.Printf("[ERROR] Error getting account info: %s", err)
		return err
	}

	account := resp.Account

	d.Set("account_pid", aws.StringValue(account.AccountPid))
	d.Set("city", aws.StringValue(account.City))
	d.Set("company_name", aws.StringValue(account.CompanyName))
	d.Set("country", aws.StringValue(account.Country))
	d.Set("customer_id", aws.StringValue(account.CustomerId))
	d.Set("email", aws.StringValue(account.Email))
	d.Set("first_name", aws.StringValue(account.FirstName))
	d.Set("job_title", aws.StringValue(account.JobTitle))
	d.Set("last_name", aws.StringValue(account.LastName))
	d.Set("mobile_number", aws.StringValue(account.MobileNumber))
	d.Set("phone_number", aws.StringValue(account.PhoneNumber))
	d.Set("state", aws.StringValue(account.State))
	d.Set("vat_number", aws.StringValue(account.VatNumber))
	d.Set("zipcode", aws.StringValue(account.Zipcode))

	d.SetId(aws.StringValue(account.AccountPid))

	return d.Set("request_id", resp.ResponseMetadata.RequestID)

}
