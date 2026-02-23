package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

var (
	_ datasource.DataSource              = &dataSourceQuota{}
	_ datasource.DataSourceWithConfigure = &dataSourceQuota{}
)

func NewDataSourceQuota() datasource.DataSource {
	return &dataSourceQuota{}
}

func (d *dataSourceQuota) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client := req.ProviderData.(client.OutscaleClient)
	d.Client = client.OSC
}

// ExampleDataSource defines the data source implementation.
type dataSourceQuota struct {
	Client *osc.Client
}

// ExampleDataSourceModel describes the data source data model.
type quotaModel struct {
	// ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
	Id               types.String `tfsdk:"id"`
	Filter           types.Set    `tfsdk:"filter"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	MaxValue         types.Int64  `tfsdk:"max_value"`
	UsedValue        types.Int64  `tfsdk:"used_value"`
	QuotaType        types.String `tfsdk:"quota_type"`
	QuotaCollection  types.String `tfsdk:"quota_collection"`
	ShortDescription types.String `tfsdk:"short_description"`
	AccountId        types.String `tfsdk:"account_id"`
	RequestId        types.String `tfsdk:"request_id"`
}

func FwDataSourceFiltersSchema() *schema.SetNestedBlock {
	return &schema.SetNestedBlock{
		Validators: []validator.Set{
			setvalidator.IsRequired(),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Required: true,
				},
				"values": schema.SetAttribute{
					ElementType: types.StringType,
					Required:    true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
			},
		},
	}
}

func (d *dataSourceQuota) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quota"
}

func (d *dataSourceQuota) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"filter": FwDataSourceFiltersSchema(),
		},
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"max_value": schema.Int64Attribute{
				Computed: true,
			},
			"used_value": schema.Int64Attribute{
				Computed: true,
			},
			"quota_type": schema.StringAttribute{
				Computed: true,
			},
			"quota_collection": schema.StringAttribute{
				Computed: true,
			},
			"short_description": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"account_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *dataSourceQuota) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	reqApi := osc.ReadQuotasRequest{}
	mapTftypes := map[string]tftypes.Value{}
	var respApi *osc.ReadQuotasResponse
	var quota osc.Quota
	var quotaType osc.QuotaTypes
	var dataState quotaModel
	var filters *osc.FiltersQuota
	var listFilters []tftypes.Value
	var diags diag.Diagnostics
	var flatenFilters basetypes.SetValue

	err := req.Config.Raw.As(&mapTftypes)
	if err != nil {
		goto CHECK_ERR
	}
	err = mapTftypes["filter"].As(&listFilters)
	if err != nil {
		goto CHECK_ERR
	}
	filters, err = buildOutscaleQuotaDataSourceFrameworkFilters(ctx, listFilters)
	if err != nil {
		goto CHECK_ERR
	}
	reqApi.Filters = filters

	respApi, err = d.Client.ReadQuotas(ctx, reqApi, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		goto CHECK_ERR
	}
	if respApi.QuotaTypes == nil || len(*respApi.QuotaTypes) == 0 {
		err = ErrNoResults
		goto CHECK_ERR
	}
	if len(*respApi.QuotaTypes) > 1 {
		err = ErrMultipleResults
		goto CHECK_ERR
	}
	quotaType = (*respApi.QuotaTypes)[0]
	if quotaType.Quotas == nil || len(*quotaType.Quotas) == 0 {
		err = ErrNoResults
		goto CHECK_ERR
	}

	if len(*quotaType.Quotas) > 1 {
		err = ErrMultipleResults
		goto CHECK_ERR
	}

	flatenFilters, diags = flatenQuotaDataSourceFilters(listFilters)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	quota = (*quotaType.Quotas)[0]
	dataState.QuotaType = to.String(quotaType.QuotaType)
	dataState.Name = to.String(quota.Name)
	dataState.Description = to.String(quota.Description)
	dataState.MaxValue = to.Int64(int64(ptr.From(quota.MaxValue)))
	dataState.UsedValue = to.Int64(int64(ptr.From(quota.UsedValue)))
	dataState.QuotaCollection = to.String(quota.QuotaCollection)
	dataState.ShortDescription = to.String(quota.ShortDescription)
	dataState.AccountId = to.String(quota.AccountId)
	dataState.Filter = flatenFilters
	dataState.Id = to.String(id.UniqueId())
	dataState.RequestId = to.String(respApi.ResponseContext.RequestId)
	diags = resp.State.Set(ctx, &dataState)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

CHECK_ERR:
	if err != nil { // resp.Diagnostics.Error != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Outscale Quotas",
			"If the error is not clear, please contact the provider developers.\n\n"+
				"Outscale Client Error: "+err.Error(),
		)

		// resp.Diagnostics.Append(err...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

func buildOutscaleQuotaDataSourceFrameworkFilters(ctx context.Context, listFilters []tftypes.Value) (*osc.FiltersQuota, error) {
	var filters osc.FiltersQuota

	for _, val := range listFilters {
		var mapFilters map[string]tftypes.Value
		val.As(&mapFilters)
		var name string
		mapFilters["name"].As(&name)
		var listValues []tftypes.Value
		mapFilters["values"].As(&listValues)
		var filterValues []string
		for _, val := range listValues {
			var value string
			val.As(&value)
			filterValues = append(filterValues, value)
		}
		switch name {
		case "quota_types":
			filters.QuotaTypes = &filterValues
		case "quota_names":
			filters.QuotaNames = &filterValues
		case "collections":
			filters.Collections = &filterValues
		case "short_descriptions":
			filters.ShortDescriptions = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}

func flatenQuotaDataSourceFilters(listFilters []tftypes.Value) (basetypes.SetValue, diag.Diagnostics) {
	var setfil []attr.Value
	var setValue basetypes.SetValue
	var diags diag.Diagnostics

	filtersValuesType := basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":   basetypes.StringType{},
			"values": basetypes.SetType{ElemType: basetypes.StringType{}},
		},
	}
	mapObjectType := make(map[string]attr.Type)
	mapObjectType["values"] = basetypes.SetType{ElemType: basetypes.StringType{}}
	mapObjectType["name"] = basetypes.StringType{}

	for _, val := range listFilters {
		var mapFilters map[string]tftypes.Value
		val.As(&mapFilters)
		mapObject := make(map[string]attr.Value)
		var name string
		mapFilters["name"].As(&name)
		mapObject["name"] = to.String(name)

		var listValues []tftypes.Value
		mapFilters["values"].As(&listValues)
		var nSet []attr.Value
		for _, val := range listValues {
			var value string
			val.As(&value)

			nSet = append(nSet, to.String(value))
		}
		filtersNameType := basetypes.StringType{}
		obt, diag := types.SetValue(filtersNameType, nSet)
		if diag != nil {
			return setValue, diag
		}
		mapObject["values"] = obt
		retOK, diag := types.ObjectValue(mapObjectType, mapObject)
		if diag != nil {
			return setValue, diag
		}
		setfil = append(setfil, retOK)
	}
	setValue, diags = types.SetValue(filtersValuesType, setfil)
	if diags != nil {
		return setValue, diags
	}
	return setValue, nil
}
