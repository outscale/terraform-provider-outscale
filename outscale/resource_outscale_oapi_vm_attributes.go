package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIVMAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIVMAttributesCreate,
		Read:   resourceOAPIVMAttributesRead,
		Update: resourceOAPIVMAttributesUpdate,
		Delete: resourceOAPIVMAttributesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// ModifyInstanceAttribute schema

			"attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"block_device_mapping": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"bsu": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_on_vm_deletion": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"iops": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"no_device": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"virtual_device_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"bsu_optimized": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"firewall_rules_set_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"shutdown_automatic_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"activated_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// DescribeInstanceAttribute schema
			// same as above, but with attr and instance id required
			"group_set": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"firewall_rules_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"sriov_net_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kernel": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_codes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			// DescribeInstanceStatus schema
			// Computed
			"instance_status_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"events_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"not_after": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"not_before": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_state": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"instance_status": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"details": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"system_status": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"details": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceOAPIVMAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Get("vm_id").(string)
	var attr *string

	if v, ok := d.GetOk("attribute"); ok {
		attr = aws.String(v.(string))
	}

	if v, ok := d.GetOk("deletion_protection"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Attribute:  attr,
			DisableApiTermination: &fcu.AttributeBooleanValue{
				Value: aws.Bool(v.(bool)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "deletion_protection"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("firewall_rules_set_id"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Groups:     v.([]*string),
		}
		if err := modifyInstanceAttr(conn, opts, "firewall_rules_set_id"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("shutdown_automatic_behavior"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Attribute:  attr,
			InstanceInitiatedShutdownBehavior: &fcu.AttributeValue{
				Value: aws.String(v.(string)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "shutdown_automatic_behavior"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("activated_check"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Attribute:  attr,
			SourceDestCheck: &fcu.AttributeBooleanValue{
				Value: aws.Bool(v.(bool)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "activated_check"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("type"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Attribute:  attr,
			InstanceType: &fcu.AttributeValue{
				Value: aws.String(v.(string)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "type"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Attribute:  attr,
			UserData: &fcu.BlobAttributeValue{
				Value: v.([]byte),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "user_data"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("bsu_optimized"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Attribute:  attr,
			EbsOptimized: &fcu.AttributeBooleanValue{
				Value: v.(*bool),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "bsu_optimized"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("delete_on_vm_deletion"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Attribute:  attr,
			DeleteOnTermination: &fcu.AttributeBooleanValue{
				Value: v.(*bool),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "delete_on_vm_deletion"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("block_device_mapping"); ok {
		maps := v.(*schema.Set).List()
		mappings := []*fcu.BlockDeviceMapping{}

		for _, m := range maps {
			f := m.(map[string]interface{})
			mapping := &fcu.BlockDeviceMapping{
				DeviceName:  aws.String(f["device_name"].(string)),
				NoDevice:    aws.String(f["no_device"].(string)),
				VirtualName: aws.String(f["virtual_device_name"].(string)),
			}

			e := f["bsu"].(map[string]interface{})

			bsu := &fcu.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(e["delete_on_vm_deletion"].(bool)),
				Iops:                aws.Int64(int64(e["iops"].(int))),
				SnapshotId:          aws.String(e["snapshot_id"].(string)),
				VolumeSize:          aws.Int64(int64(e["volume_size"].(int))),
				VolumeType:          aws.String((e["type"].(string))),
			}

			mapping.Ebs = bsu

			mappings = append(mappings, mapping)
		}

		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId:          aws.String(id),
			BlockDeviceMappings: mappings,
		}
		if err := modifyInstanceAttr(conn, opts, "deletion_protection"); err != nil {
			return err
		}
	}

	d.SetId(resource.UniqueId())

	return resourceOAPIVMAttributesRead(d, meta)
}

func resourceOAPIVMAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	if err := readDescribeOAPIVMAttr(d, conn); err != nil {
		return err
	}

	return readDescribeOAPIVMStatus(d, conn)
}

func resourceOAPIVMAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	d.Partial(true)

	id := d.Get("vm_id").(string)
	// attr := aws.String(d.Get("attribute").(string))

	log.Printf("[DEBUG] updating the instance %s", id)

	if d.HasChange("type") && !d.IsNewResource() {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			InstanceType: &fcu.AttributeValue{
				Value: aws.String(d.Get("type").(string)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "type"); err != nil {
			return err
		}
	}

	if d.HasChange("user_data") && !d.IsNewResource() {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			UserData: &fcu.BlobAttributeValue{
				Value: d.Get("user_data").([]byte),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "user_data"); err != nil {
			return err
		}
	}

	if d.HasChange("bsu_optimized") && !d.IsNewResource() {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			EbsOptimized: &fcu.AttributeBooleanValue{
				Value: aws.Bool(d.Get("bsu_optimized").(bool)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "bsu_optimized"); err != nil {
			return err
		}
	}

	if d.HasChange("delete_on_vm_deletion") && !d.IsNewResource() {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			DeleteOnTermination: &fcu.AttributeBooleanValue{
				Value: d.Get("delete_on_vm_deletion").(*bool),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "delete_on_vm_deletion"); err != nil {
			return err
		}
	}

	if d.HasChange("deletion_protection") {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			DisableApiTermination: &fcu.AttributeBooleanValue{
				Value: aws.Bool(d.Get("deletion_protection").(bool)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "deletion_protection"); err != nil {
			return err
		}
	}

	if d.HasChange("shutdown_automatic_behavior") {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			InstanceInitiatedShutdownBehavior: &fcu.AttributeValue{
				Value: aws.String(d.Get("shutdown_automatic_behavior").(string)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "shutdown_automatic_behavior"); err != nil {
			return err
		}
	}

	if d.HasChange("group_set") {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Groups:     d.Get("group_set").([]*string),
		}
		if err := modifyInstanceAttr(conn, opts, "deletion_protection"); err != nil {
			return err
		}
	}

	if d.HasChange("activated_check") {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			SourceDestCheck: &fcu.AttributeBooleanValue{
				Value: aws.Bool(d.Get("activated_check").(bool)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "activated_check"); err != nil {
			return err
		}
	}

	if d.HasChange("block_device_mapping") {
		maps := d.Get("block_device_mapping").(*schema.Set).List()
		mappings := []*fcu.BlockDeviceMapping{}

		for _, m := range maps {
			f := m.(map[string]interface{})
			mapping := &fcu.BlockDeviceMapping{
				DeviceName:  aws.String(f["device_name"].(string)),
				NoDevice:    aws.String(f["no_device"].(string)),
				VirtualName: aws.String(f["virtual_device_name"].(string)),
			}

			e := f["bsu"].(map[string]interface{})

			bsu := &fcu.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(e["delete_on_vm_deletion"].(bool)),
				Iops:                aws.Int64(int64(e["iops"].(int))),
				SnapshotId:          aws.String(e["snapshot_id"].(string)),
				VolumeSize:          aws.Int64(int64(e["volume_size"].(int))),
				VolumeType:          aws.String((e["type"].(string))),
			}

			mapping.Ebs = bsu

			mappings = append(mappings, mapping)
		}

		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId:          aws.String(id),
			BlockDeviceMappings: mappings,
		}
		if err := modifyInstanceAttr(conn, opts, "deletion_protection"); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceOAPIVMAttributesRead(d, meta)
}

func resourceOAPIVMAttributesDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}

func readDescribeOAPIVMAttr(d *schema.ResourceData, conn *fcu.Client) error {
	input := &fcu.DescribeInstanceAttributeInput{
		Attribute:  aws.String(d.Get("attribute").(string)),
		InstanceId: aws.String(d.Get("vm_id").(string)),
	}

	var resp *fcu.DescribeInstanceAttributeOutput
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeInstanceAttribute(input)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error reading the DescribeInstanceAttribute %s", err)
	}

	d.Set("vm_id", resp.InstanceId)

	d.Set("block_device_mapping", getBlockDeviceMapping(resp.BlockDeviceMappings))

	d.Set("deletion_protection", resp.DisableApiTermination)

	d.Set("bsu_optimized", resp.EbsOptimized)

	err = d.Set("group_set", getGroupSet(resp.Groups))
	if err != nil {
		fmt.Println(getGroupSet(resp.Groups))
	}

	d.Set("shutdown_automatic_behavior", resp.InstanceInitiatedShutdownBehavior)

	d.Set("type", resp.InstanceType)

	d.Set("kernel", resp.KernelId)

	d.Set("product_codes", getProductCodes(resp.ProductCodes))

	d.Set("ramdisk", resp.RamdiskId)

	d.Set("root_device_name", resp.RootDeviceName)

	d.Set("activated_check", resp.SourceDestCheck)

	d.Set("sriov_net_support", resp.SriovNetSupport)

	d.Set("user_data", resp.UserData)

	return nil
}

func readDescribeOAPIVMStatus(d *schema.ResourceData, conn *fcu.Client) error {
	input := &fcu.DescribeInstanceStatusInput{
		InstanceIds: []*string{aws.String(d.Get("vm_id").(string))},
	}

	var resp *fcu.DescribeInstanceStatusOutput
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeInstanceStatus(input)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error reading the DescribeInstanceStatus %s", err)
	}

	if len(resp.InstanceStatuses) > 0 {
		instances := make([]map[string]interface{}, len(resp.InstanceStatuses))

		for k, v := range resp.InstanceStatuses {
			instance := make(map[string]interface{})

			if v.AvailabilityZone != nil {
				instance["availability_zone"] = *v.AvailabilityZone
			}
			if v.Events != nil {
				events := make([]map[string]interface{}, len(v.Events))
				for i, e := range v.Events {
					event := make(map[string]interface{})
					if e.Code != nil {
						event["code"] = *e.Code
					}
					if e.Description != nil {
						event["description"] = *e.Description
					}
					if e.NotAfter != nil {
						event["not_after"] = *e.NotAfter
					}
					if e.NotBefore != nil {
						event["not_before"] = *e.NotBefore
					}
					events[i] = event
				}
				instance["events"] = events
			}
			if v.InstanceId != nil {
				instance["vm_id"] = *v.InstanceId
			}
			if v.InstanceState != nil {
				state := make(map[string]interface{})

				if v.InstanceState.Code != nil {
					state["code"] = fmt.Sprint(*v.InstanceState.Code)
				}
				if v.InstanceState.Name != nil {
					state["name"] = *v.InstanceState.Name
				}
				instance["instance_state"] = state
			}
			if v.InstanceStatus != nil {
				state := make(map[string]interface{})

				if v.InstanceStatus.Details != nil {
					details := make([]map[string]interface{}, len(v.InstanceStatus.Details))
					for j, d := range v.InstanceStatus.Details {
						detail := make(map[string]interface{})
						if d.Name != nil {
							detail["name"] = *d.Name
						}
						if d.Status != nil {
							detail["status"] = *d.Status
						}
						details[j] = detail
					}
					state["details"] = details
				}
				if v.InstanceStatus.Status != nil {
					state["status"] = *v.InstanceStatus.Status
				}
				instance["instance_status"] = state
			}
			if v.SystemStatus != nil {
				state := make(map[string]interface{})

				if v.SystemStatus.Details != nil {
					details := make([]map[string]interface{}, len(v.SystemStatus.Details))
					for j, d := range v.SystemStatus.Details {
						detail := make(map[string]interface{})
						if d.Name != nil {
							detail["name"] = *d.Name
						}
						if d.Status != nil {
							detail["status"] = *d.Status
						}
						details[j] = detail
					}
					state["details"] = details
				}
				if v.SystemStatus.Status != nil {
					state["status"] = *v.SystemStatus.Status
				}
				instance["system_status"] = state
			}

			instances[k] = instance
		}

		fmt.Printf("\n\n[DEBUG] instance_status_set %s", instances)

		if err := d.Set("instance_status_set", instances); err != nil {
			return err
		}
	}
	return nil
}
