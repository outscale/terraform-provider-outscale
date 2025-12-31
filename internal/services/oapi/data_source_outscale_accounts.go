package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceAccounts() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAccountsRead,
		Schema: map[string]*schema.Schema{
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceAccountsRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.ReadAccountsRequest{}

	var resp oscgo.ReadAccountsResponse
	var err error

	err = retry.Retry(30*time.Second, func() *retry.RetryError {
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
		return fmt.Errorf("Unable to find Accounts")
	}

	if err := d.Set("accounts", flattenAccounts(resp.GetAccounts())); err != nil {
		return err
	}
	d.SetId(id.UniqueId())

	return nil
}

func flattenAccounts(accounts []oscgo.Account) []map[string]interface{} {
	accountsMap := make([]map[string]interface{}, len(accounts))

	for i, account := range accounts {
		accountsMap[i] = map[string]interface{}{
			"account_id":        account.GetAccountId(),
			"additional_emails": utils.StringSlicePtrToInterfaceSlice(account.AdditionalEmails),
			"city":              account.GetCity(),
			"company_name":      account.GetCompanyName(),
			"country":           account.GetCountry(),
			"customer_id":       account.GetCustomerId(),
			"email":             account.GetEmail(),
			"first_name":        account.GetFirstName(),
			"job_title":         account.GetJobTitle(),
			"last_name":         account.GetLastName(),
			"mobile_number":     account.GetMobileNumber(),
			"phone_number":      account.GetPhoneNumber(),
			"state_province":    account.GetStateProvince(),
			"vat_number":        account.GetVatNumber(),
			"zip_code":          account.GetZipCode(),
		}
	}
	return accountsMap
}
