package oapi

import (
	"context"
	"fmt"
	"slices"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLBUAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourceLBUAttachmentCreate,
		Read:   ResourceLBUAttachmentRead,
		Update: ResourceLBUAttachmentUpdate,
		Delete: ResourceLBUAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"backend_vm_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"backend_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceLBUAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	lbuName := d.Get("load_balancer_name").(string)
	vmIds := utils.SetToStringSlice(d.Get("backend_vm_ids").(*schema.Set))
	vmIps := d.Get("backend_ips").(*schema.Set)
	if len(vmIds) == 0 && vmIps.Len() == 0 {
		return fmt.Errorf("error: the 'backend_vm_ids' and 'backend_ips' parameters cannot both be empty")
	}
	if vmIps.Len() > 0 {
		vm_ids, err := getVmIdsThroughVmIps(conn, vmIps)
		if err != nil {
			return err
		}
		vmIds = append(vmIds, vm_ids...)
	}
	req := oscgo.RegisterVmsInLoadBalancerRequest{
		LoadBalancerName: lbuName,
		BackendVmIds:     vmIds,
	}

	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.LoadBalancerApi.RegisterVmsInLoadBalancer(
			context.Background()).RegisterVmsInLoadBalancerRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf(" failure linking loadbalancer backend_vm_ids/backend_ips with lbu: %w", err)
	}
	d.SetId(lbuName)
	return ResourceLBUAttachmentRead(d, meta)
}

func ResourceLBUAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	lbuName := d.Get("load_balancer_name").(string)
	lbu, _, err := readResourceLb(conn, lbuName)
	if err != nil {
		return err
	}
	if lbu == nil {
		utils.LogManuallyDeleted("LoadBalancerVms", d.Id())
		d.SetId("")
		return nil
	}
	if len(lbu.GetBackendVmIds()) == 0 {
		utils.LogManuallyDeleted("LoadBalancerVms", d.Id())
		d.SetId("")
		return nil
	}

	expectedVmIds := d.Get("backend_vm_ids").(*schema.Set)
	expectedIps := d.Get("backend_ips").(*schema.Set)
	copyTypeSet := schema.CopySet(expectedVmIds)
	all_backendVms := schema.NewSet(copyTypeSet.F, []interface{}{})
	all_backendIps := schema.NewSet(copyTypeSet.F, []interface{}{})

	for _, vmId := range lbu.GetBackendVmIds() {
		all_backendVms.Add(vmId)
	}
	publicIps, err := getVmIpsThroughVmIds(conn, all_backendVms)
	if err != nil {
		return err
	}
	for _, vmIp := range publicIps {
		all_backendIps.Add(vmIp)
	}

	managedVmIds := all_backendVms.Intersection(expectedVmIds)
	managedIps := all_backendIps.Intersection(expectedIps)

	if managedIps.Len() > 0 {
		vmIdsLinkedByIps, err := getVmIdsThroughVmIps(conn, managedIps)
		if err != nil {
			return err
		}
		for _, vmId := range vmIdsLinkedByIps {
			all_backendVms.Remove(vmId)
		}
	}

	if manVmIdsLink := all_backendVms.Difference(expectedVmIds); manVmIdsLink.Len() > 0 {
		for _, vmId := range manVmIdsLink.List() {
			managedVmIds.Add(vmId)
		}
	}

	if managedVmIds.Len() == 0 && managedIps.Len() == 0 {
		d.SetId("")
		return nil
	}
	if err := d.Set("backend_vm_ids", managedVmIds); err != nil {
		return err
	}
	if err := d.Set("backend_ips", managedIps); err != nil {
		return err
	}
	if err := d.Set("load_balancer_name", lbu.GetLoadBalancerName()); err != nil {
		return err
	}
	d.SetId(lbu.GetLoadBalancerName())
	return nil
}

func ResourceLBUAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	lbuName := d.Get("load_balancer_name").(string)
	linkReq, unLinkReq, err := buildUpdateBackendsRequest(d, conn, lbuName)
	if err != nil {
		return err
	}

	if unLinkReq.HasBackendVmIds() {
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			_, httpResp, err := conn.LoadBalancerApi.
				UnlinkLoadBalancerBackendMachines(context.Background()).
				UnlinkLoadBalancerBackendMachinesRequest(*unLinkReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failure unlinking backends from lbu: %s", err)
		}
	}

	if linkReq.HasBackendVmIds() {
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			_, httpResp, err := conn.LoadBalancerApi.
				LinkLoadBalancerBackendMachines(context.Background()).
				LinkLoadBalancerBackendMachinesRequest(*linkReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failure linking backends to lbu: %s", err)
		}
	}
	return ResourceLBUAttachmentRead(d, meta)
}

func buildUpdateBackendsRequest(d *schema.ResourceData, conn *oscgo.APIClient, lbuName string) (*oscgo.LinkLoadBalancerBackendMachinesRequest, *oscgo.UnlinkLoadBalancerBackendMachinesRequest, error) {
	linkReq := oscgo.NewLinkLoadBalancerBackendMachinesRequest(lbuName)
	unLinkReq := oscgo.NewUnlinkLoadBalancerBackendMachinesRequest(lbuName)
	linkVmIds := make([]string, 0)
	unlinkVmIds := make([]string, 0)
	if d.HasChange("backend_vm_ids") {
		oldBackends, newBackends := d.GetChange("backend_vm_ids")
		inter := oldBackends.(*schema.Set).Intersection(newBackends.(*schema.Set))
		created := newBackends.(*schema.Set).Difference(inter)
		removed := oldBackends.(*schema.Set).Difference(inter)

		if created.Len() > 0 {

			_, err := getVmIpsThroughVmIds(conn, created)
			if err != nil {
				return linkReq, unLinkReq, err
			}
			linkVmIds = append(linkVmIds, utils.SetToStringSlice(created)...)
		}
		if removed.Len() > 0 {
			_, err := getVmIpsThroughVmIds(conn, removed)
			if err != nil {
				return linkReq, unLinkReq, err
			}
			unlinkVmIds = append(unlinkVmIds, utils.SetToStringSlice(removed)...)
		}
	}

	if d.HasChange("backend_ips") {
		oldBackends, newBackends := d.GetChange("backend_ips")
		inter := oldBackends.(*schema.Set).Intersection(newBackends.(*schema.Set))
		created := newBackends.(*schema.Set).Difference(inter)
		removed := oldBackends.(*schema.Set).Difference(inter)

		if created.Len() > 0 {
			vmIdsToCreate, err := getVmIdsThroughVmIps(conn, created)
			if err != nil {
				return linkReq, unLinkReq, err
			}
			linkVmIds = vmIdsToCreate
		}

		if removed.Len() > 0 {
			vmIdsToRemove, err := getVmIdsThroughVmIps(conn, removed)
			if err != nil {
				return linkReq, unLinkReq, err
			}
			unlinkVmIds = vmIdsToRemove
		}
	}
	if len(linkVmIds) > 0 {
		linkReq.SetBackendVmIds(linkVmIds)
	}
	if len(unlinkVmIds) > 0 {
		unLinkReq.SetBackendVmIds(unlinkVmIds)
	}
	return linkReq, unLinkReq, nil
}

func ResourceLBUAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	lbuName := d.Get("load_balancer_name").(string)
	unlinkVmIds := utils.SetToStringSlicePtr(d.Get("backend_vm_ids").(*schema.Set))
	if ips := d.Get("backend_ips").(*schema.Set); ips.Len() > 0 {
		vmIps, err := getVmIdsThroughVmIps(conn, ips)
		if err != nil {
			return err
		}
		*unlinkVmIds = append(*unlinkVmIds, vmIps...)
	}
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.LoadBalancerApi.
			UnlinkLoadBalancerBackendMachines(context.Background()).
			UnlinkLoadBalancerBackendMachinesRequest(
				oscgo.UnlinkLoadBalancerBackendMachinesRequest{
					LoadBalancerName: lbuName,
					BackendVmIds:     unlinkVmIds,
				}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failure unlinking backend_ips from lbu: %s", err)
	}
	return nil
}

func getVmIdsThroughVmIps(conn *oscgo.APIClient, vmIps *schema.Set) ([]string, error) {
	filterIps := oscgo.NewFiltersVm()
	ipsList := utils.SetToStringSlice(vmIps)
	filterIps.SetPublicIps(ipsList)
	var resp oscgo.ReadVmsResponse
	err := retry.Retry(30*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
			Filters: filterIps,
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return nil, err
	}
	vms := resp.GetVms()
	if len(vms) == 0 {
		return nil, fmt.Errorf("cannot find vms with public_ips: %v", ipsList)
	}
	vmsIds := make([]string, 0, len(vms))
	vmsIpsList := make([]string, 0, len(vms))
	for _, vm := range vms {
		vmsIds = append(vmsIds, vm.GetVmId())
		vmsIpsList = append(vmsIpsList, vm.GetPublicIp())
	}
	slices.Sort(vmsIpsList)
	slices.Sort(ipsList)
	if slices.Compare(ipsList, vmsIpsList) != 0 {
		return nil, fmt.Errorf("some public_ips are not linked to any vm in this list: %v", ipsList)
	}
	return vmsIds, nil
}

func getVmIpsThroughVmIds(conn *oscgo.APIClient, vmIds *schema.Set) ([]string, error) {
	filters := oscgo.NewFiltersVm()
	vmIdsList := utils.SetToStringSlice(vmIds)
	filters.SetVmIds(vmIdsList)
	var resp oscgo.ReadVmsResponse
	err := retry.Retry(30*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
			Filters: filters,
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return nil, err
	}
	vms := resp.GetVms()
	if len(vms) == 0 {
		return nil, fmt.Errorf("cannot find vms with vm_ids %v", vmIdsList)
	}
	publicIps := make([]string, 0, len(vms))
	readVmIds := make([]string, 0, len(vms))
	for _, vm := range vms {
		publicIps = append(publicIps, vm.GetPublicIp())
		readVmIds = append(readVmIds, vm.GetVmId())
	}
	slices.Sort(readVmIds)
	slices.Sort(vmIdsList)
	if slices.Compare(vmIdsList, readVmIds) != 0 {
		return nil, fmt.Errorf("some vm_ids are not existed in this list: %v", vmIdsList)
	}
	return publicIps, nil
}
