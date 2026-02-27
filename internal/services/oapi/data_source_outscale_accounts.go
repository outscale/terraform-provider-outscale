package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceAccountsRead,
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

func DataSourceAccountsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadAccountsRequest{}

	resp, err := client.ReadAccounts(ctx, req, options.WithRetryTimeout(30*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Accounts == nil || len(*resp.Accounts) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if err := d.Set("accounts", flattenAccounts(*resp.Accounts)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id.UniqueId())

	return nil
}

func flattenAccounts(accounts []osc.Account) []map[string]interface{} {
	accountsMap := make([]map[string]interface{}, len(accounts))

	for i, account := range accounts {
		accountsMap[i] = map[string]interface{}{
			"account_id":        account.AccountId,
			"additional_emails": utils.StringSlicePtrToInterfaceSlice(account.AdditionalEmails),
			"city":              account.City,
			"company_name":      account.CompanyName,
			"country":           account.Country,
			"customer_id":       account.CustomerId,
			"email":             account.Email,
			"first_name":        account.FirstName,
			"job_title":         account.JobTitle,
			"last_name":         account.LastName,
			"mobile_number":     account.MobileNumber,
			"phone_number":      account.PhoneNumber,
			"state_province":    account.StateProvince,
			"vat_number":        account.VatNumber,
			"zip_code":          account.ZipCode,
		}
	}
	return accountsMap
}
