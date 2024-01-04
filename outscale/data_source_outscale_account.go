package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccountRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"additional_emails": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Computed: true,
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
			"state_province": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vat_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zip_code": {
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

func dataSourceAccountRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadAccountsRequest{}

	var resp oscgo.ReadAccountsResponse
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.AccountApi.ReadAccounts(context.Background()).ReadAccountsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetAccounts()) == 0 {
		return fmt.Errorf("Unable to find Account")
	}

	if len(resp.GetAccounts()) > 1 {
		return fmt.Errorf("multiple results returned, please use a more specific criteria in your query")
	}

	account := resp.GetAccounts()[0]

	d.SetId(resource.UniqueId())

	if err := d.Set("account_id", account.GetAccountId()); err != nil {
		return err
	}
	if err := d.Set("additional_emails", utils.StringSlicePtrToInterfaceSlice(account.AdditionalEmails)); err != nil {
		return err
	}
	if err := d.Set("city", account.GetCity()); err != nil {
		return err
	}
	if err := d.Set("company_name", account.GetCompanyName()); err != nil {
		return err
	}
	if err := d.Set("country", account.GetCountry()); err != nil {
		return err
	}
	if err := d.Set("customer_id", account.GetCustomerId()); err != nil {
		return err
	}
	if err := d.Set("email", account.GetEmail()); err != nil {
		return err
	}
	if err := d.Set("first_name", account.GetFirstName()); err != nil {
		return err
	}
	if err := d.Set("job_title", account.GetJobTitle()); err != nil {
		return err
	}
	if err := d.Set("last_name", account.GetLastName()); err != nil {
		return err
	}
	if err := d.Set("mobile_number", account.GetMobileNumber()); err != nil {
		return err
	}
	if err := d.Set("phone_number", account.GetPhoneNumber()); err != nil {
		return err
	}
	if err := d.Set("state_province", account.GetStateProvince()); err != nil {
		return err
	}
	if err := d.Set("vat_number", account.GetVatNumber()); err != nil {
		return err
	}
	if err := d.Set("zip_code", account.GetZipCode()); err != nil {
		return err
	}
	return nil
}
