package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceAccountRead,
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

func DataSourceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadAccountsRequest{}

	resp, err := client.ReadAccounts(ctx, req, options.WithRetryTimeout(30*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Accounts == nil || len(*resp.Accounts) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.Accounts) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	account := (*resp.Accounts)[0]

	d.SetId(id.UniqueId())

	if err := d.Set("account_id", ptr.From(account.AccountId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("additional_emails", utils.StringSlicePtrToInterfaceSlice(account.AdditionalEmails)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("city", ptr.From(account.City)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("company_name", ptr.From(account.CompanyName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("country", ptr.From(account.Country)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("customer_id", ptr.From(account.CustomerId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", ptr.From(account.Email)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("first_name", ptr.From(account.FirstName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("job_title", ptr.From(account.JobTitle)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_name", ptr.From(account.LastName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mobile_number", ptr.From(account.MobileNumber)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("phone_number", ptr.From(account.PhoneNumber)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state_province", ptr.From(account.StateProvince)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vat_number", ptr.From(account.VatNumber)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("zip_code", ptr.From(account.ZipCode)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
