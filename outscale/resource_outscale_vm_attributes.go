package outscale

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleVMAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceVMAttributesCreate,
		Read:   resourceVMAttributesRead,
		Update: resourceVMAttributesUpdate,
		Delete: resourceVMAttributesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Argument
			"attribute": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disable_api_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// Attributes schema
			"block_device_mapping": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"ebs": {
							Type:     schema.TypeMap,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_on_termination": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"volume_id": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"ebs_optimized": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"group_set": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
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
			"instance_initiated_shutdown_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ramdisk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_dest_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"sriov_net_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
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
						"instance_id": {
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

func resourceVMAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	i, idOk := d.GetOk("instance_id")

	if !idOk {
		return fmt.Errorf("Please provide an instance_id")
	}

	fmt.Printf("\n\n[DEBUG] INSTANCE TO MODIFY (%s)", i)

	id := i.(string)

	if v, ok := d.GetOk("disable_api_termination"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			DisableApiTermination: &fcu.AttributeBooleanValue{
				Value: aws.Bool(v.(bool)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "disable_api_termination"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("group_id"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Groups:     v.([]*string),
		}
		if err := modifyInstanceAttr(conn, opts, "group_id"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("instance_initiated_shutdown_behavior"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			InstanceInitiatedShutdownBehavior: &fcu.AttributeValue{
				Value: aws.String(v.(string)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "instance_initiated_shutdown_behavior"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("source_dest_check"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			SourceDestCheck: &fcu.AttributeBooleanValue{
				Value: aws.Bool(v.(bool)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "source_dest_check"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("instance_type"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			InstanceType: &fcu.AttributeValue{
				Value: aws.String(v.(string)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "instance_type"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			UserData: &fcu.BlobAttributeValue{
				Value: v.([]byte),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "user_data"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("ebs_optimized"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			EbsOptimized: &fcu.AttributeBooleanValue{
				Value: aws.Bool(v.(bool)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "ebs_optimized"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("delete_on_termination"); ok {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			DeleteOnTermination: &fcu.AttributeBooleanValue{
				Value: aws.Bool(v.(bool)),
			},
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := modifyInstanceAttr(conn, opts, "delete_on_termination"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("block_device_mapping"); ok {
		if err := setBlockDevice(v, conn, id); err != nil {
			return err
		}
	}

	d.SetId(id)

	return resourceVMAttributesRead(d, meta)
}

func resourceVMAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	if err := readDescribeVMAttr(d, conn); err != nil {
		return err
	}

	// return readDescribeVMStatus(d, conn)

	return nil
}

func resourceVMAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	d.Partial(true)

	log.Printf("[DEBUG] updating the instance %s", d.Id())

	if d.HasChange("instance_type") && !d.IsNewResource() {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(d.Id()),
			InstanceType: &fcu.AttributeValue{
				Value: aws.String(d.Get("instance_type").(string)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "instance_type"); err != nil {
			return err
		}
	}

	if d.HasChange("user_data") && !d.IsNewResource() {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(d.Id()),
			UserData: &fcu.BlobAttributeValue{
				Value: d.Get("user_data").([]byte),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "user_data"); err != nil {
			return err
		}
	}

	if d.HasChange("ebs_optimized") && !d.IsNewResource() {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(d.Id()),
			EbsOptimized: &fcu.AttributeBooleanValue{
				Value: aws.Bool(d.Get("ebs_optimized").(bool)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "ebs_optimized"); err != nil {
			return err
		}
	}

	if d.HasChange("delete_on_termination") && !d.IsNewResource() {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(d.Id()),
			DeleteOnTermination: &fcu.AttributeBooleanValue{
				Value: aws.Bool(d.Get("delete_on_termination").(bool)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "delete_on_termination"); err != nil {
			return err
		}
	}

	if d.HasChange("disable_api_termination") {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(d.Id()),
			DisableApiTermination: &fcu.AttributeBooleanValue{
				Value: aws.Bool(d.Get("disable_api_termination").(bool)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "disable_api_termination"); err != nil {
			return err
		}
	}

	if d.HasChange("instance_initiated_shutdown_behavior") {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(d.Id()),
			InstanceInitiatedShutdownBehavior: &fcu.AttributeValue{
				Value: aws.String(d.Get("instance_initiated_shutdown_behavior").(string)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "instance_initiated_shutdown_behavior"); err != nil {
			return err
		}
	}

	if d.HasChange("group_set") {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(d.Id()),
			Groups:     d.Get("group_set").([]*string),
		}
		if err := modifyInstanceAttr(conn, opts, "disable_api_termination"); err != nil {
			return err
		}
	}

	if d.HasChange("source_dest_check") {
		opts := &fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(d.Id()),
			SourceDestCheck: &fcu.AttributeBooleanValue{
				Value: aws.Bool(d.Get("source_dest_check").(bool)),
			},
		}
		if err := modifyInstanceAttr(conn, opts, "source_dest_check"); err != nil {
			return err
		}
	}

	if d.HasChange("block_device_mapping") {
		if err := setBlockDevice(d.Get("block_device_mapping"), conn, d.Id()); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceVMAttributesRead(d, meta)
}

func resourceVMAttributesDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}

func readDescribeVMAttr(d *schema.ResourceData, conn *fcu.Client) error {
	input := &fcu.DescribeInstanceAttributeInput{
		Attribute:  aws.String(d.Get("attribute").(string)),
		InstanceId: aws.String(d.Id()),
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

	pretty, err := json.MarshalIndent(resp, "", "  ")

	fmt.Print("\n\n[DEBUG] RESPONSE ", string(pretty))

	d.Set("instance_id", *resp.InstanceId)

	d.Set("block_device_mapping", getBlockDeviceMapping(resp.BlockDeviceMappings))

	d.Set("product_codes", getProductCodes(resp.ProductCodes))

	if resp.DisableApiTermination != nil {
		d.Set("disable_api_termination", *resp.DisableApiTermination.Value)
	}

	if resp.EbsOptimized != nil {
		d.Set("ebs_optimized", *resp.EbsOptimized.Value)
	}

	if resp.Groups != nil {
		err = d.Set("group_set", getGroupSet(resp.Groups))
		if err != nil {
			fmt.Println(getGroupSet(resp.Groups))
		}
	}

	if resp.InstanceInitiatedShutdownBehavior != nil {
		d.Set("instance_initiated_shutdown_behavior", *resp.InstanceInitiatedShutdownBehavior.Value)
	} else {
		d.Set("instance_initiated_shutdown_behavior", "")
	}

	if resp.InstanceType != nil {
		d.Set("instance_type", *resp.InstanceType.Value)
	} else {
		d.Set("instance_type", "")
	}

	if resp.KernelId != nil {
		d.Set("kernel", *resp.KernelId.Value)
	} else {
		d.Set("kernel", "")
	}

	if resp.RamdiskId != nil {
		d.Set("ramdisk", *resp.RamdiskId.Value)
	} else {
		d.Set("ramdisk", "")
	}

	if resp.RootDeviceName != nil {
		d.Set("root_device_name", *resp.RootDeviceName.Value)
	} else {
		d.Set("root_device_name", "")
	}

	if resp.SourceDestCheck != nil {
		d.Set("source_dest_check", *resp.SourceDestCheck.Value)
	} else {
		d.Set("source_dest_check", "")
	}

	if resp.SriovNetSupport != nil {
		d.Set("sriov_net_support", *resp.SriovNetSupport.Value)
	} else {
		d.Set("sriov_net_support", "")
	}

	if resp.UserData != nil {
		d.Set("user_data", *resp.UserData.Value)
	} else {
		d.Set("user_data", "")
	}

	d.Set("request_id", resp.RequestId)

	return nil
}

func readDescribeVMStatus(d *schema.ResourceData, conn *fcu.Client) error {
	input := &fcu.DescribeInstanceStatusInput{
		InstanceIds: []*string{aws.String(d.Get("instance_id").(string))},
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

			fmt.Printf("\n\n[DEBUG] RESPONSEINSTANCE %+v", v)

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
				instance["instance_id"] = *v.InstanceId
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

func setBlockDevice(v interface{}, conn *fcu.Client, id string) error {
	maps := v.([]interface{})
	mappings := []*fcu.BlockDeviceMapping{}

	for _, m := range maps {
		f := m.(map[string]interface{})
		mapping := &fcu.BlockDeviceMapping{
			DeviceName: aws.String(f["device_name"].(string)),
		}

		e := f["ebs"].(map[string]interface{})
		var del bool
		if e["delete_on_termination"].(string) == "0" {
			del = false
		} else {
			del = true
		}

		ebs := &fcu.EbsBlockDevice{
			DeleteOnTermination: aws.Bool(del),
		}

		mapping.Ebs = ebs

		mappings = append(mappings, mapping)
	}

	opts := &fcu.ModifyInstanceAttributeInput{
		InstanceId:          aws.String(id),
		BlockDeviceMappings: mappings,
	}
	if err := modifyInstanceAttr(conn, opts, "block_device_mapping"); err != nil {
		return err
	}
	return nil
}
