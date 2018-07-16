package outscale

import (
	"fmt"
	"strings"
	"time"

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

	if account.AccountPid != nil {
		d.Set("account_pid", *account.AccountPid)
	}
	if account.City != nil {
		d.Set("city", *account.City)
	}
	if account.CompanyName != nil {
		d.Set("company_name", *account.CompanyName)
	}
	if account.Country != nil {
		d.Set("country", *account.Country)
	}
	if account.CustomerId != nil {
		d.Set("customer_id", *account.CustomerId)
	}
	if account.Email != nil {
		d.Set("email", *account.Email)
	}
	if account.FirstName != nil {
		d.Set("first_name", *account.FirstName)
	}
	if account.JobTitle != nil {
		d.Set("job_title", *account.JobTitle)
	}
	if account.LastName != nil {
		d.Set("last_name", *account.LastName)
	}
	if account.MobileNumber != nil {
		d.Set("mobile_number", *account.MobileNumber)
	}
	if account.PhoneNumber != nil {
		d.Set("phone_number", *account.PhoneNumber)
	}
	if account.State != nil {
		d.Set("state", *account.State)
	}
	if account.VatNumber != nil {
		d.Set("vat_number", *account.VatNumber)
	}
	if account.Zipcode != nil {
		d.Set("zipcode", *account.Zipcode)
	}

	d.SetId(resource.UniqueId())
	d.Set("request_id", resp.ResponseMetadata.RequestID)

	return nil
}
