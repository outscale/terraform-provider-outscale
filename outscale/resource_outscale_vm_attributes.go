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
			"attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"additional_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// Describe
			"filter": dataSourceFiltersSchema(),
			"include_all_instances": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

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

			// modify

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
						"ebs": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_on_termination": {
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
									"volume_type": {
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
						"virtual_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"disable_api_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ebs_optimized": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"source_dest_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Computed
			"sriov_net_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ramdisk_id": {
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
			"instances_set": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"current_state": {
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
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"previous_state": {
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
					},
				},
			},
		},
	}
}

func resourceVMAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := aws.String(d.Get("instance_id").(string))
	var attr *string

	if attribute, attributeOk := d.GetOk("attribute"); attributeOk {
		attr = aws.String(attribute.(string))
	}
	blockDevices, err := readBlockDeviceMappingsFromConfig(d, conn)
	if err != nil {
		return err
	}
	if len(blockDevices) > 0 {
		instanceAttrOpts := &fcu.ModifyInstanceAttributeInput{}
		instanceAttrOpts.InstanceId = id
		instanceAttrOpts.Attribute = attr
		instanceAttrOpts.BlockDeviceMappings = blockDevices
		if _, err = conn.VM.ModifyInstanceAttribute(instanceAttrOpts); err != nil {
			return err
		}
	}
	if t, tOk := d.GetOk("disable_api_termination"); tOk {
		instanceAttrOpts := &fcu.ModifyInstanceAttributeInput{}
		instanceAttrOpts.InstanceId = id
		instanceAttrOpts.Attribute = attr
		instanceAttrOpts.DisableApiTermination = &fcu.AttributeBooleanValue{Value: aws.Bool(t.(bool))}
		if _, err = conn.VM.ModifyInstanceAttribute(instanceAttrOpts); err != nil {
			return err
		}
	}

	if t, tOk := d.GetOk("group_id"); tOk {
		g := t.([]interface{})
		gr := make([]*string, len(g))
		for k, v := range g {
			gr[k] = aws.String(v.(string))
		}
		instanceAttrOpts := &fcu.ModifyInstanceAttributeInput{}
		instanceAttrOpts.InstanceId = id
		instanceAttrOpts.Attribute = attr
		instanceAttrOpts.Groups = gr
		if _, err = conn.VM.ModifyInstanceAttribute(instanceAttrOpts); err != nil {
			return err
		}
	}
	if instanceInit, instanceInitOk := d.GetOk("instance_initiated_shutdown_behavior"); instanceInitOk {
		instanceAttrOpts := &fcu.ModifyInstanceAttributeInput{}
		instanceAttrOpts.InstanceId = id
		instanceAttrOpts.Attribute = attr
		instanceAttrOpts.InstanceInitiatedShutdownBehavior = &fcu.AttributeValue{Value: aws.String(instanceInit.(string))}
		if _, err = conn.VM.ModifyInstanceAttribute(instanceAttrOpts); err != nil {
			return err
		}
	}
	if t, tOk := d.GetOk("source_dest_check"); tOk {
		instanceAttrOpts := &fcu.ModifyInstanceAttributeInput{}
		instanceAttrOpts.InstanceId = id
		instanceAttrOpts.Attribute = attr
		instanceAttrOpts.SourceDestCheck = &fcu.AttributeBooleanValue{Value: aws.Bool(t.(bool))}
		if _, err = conn.VM.ModifyInstanceAttribute(instanceAttrOpts); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("instance_type"); ok {
		instanceAttrOpts := &fcu.ModifyInstanceAttributeInput{}
		instanceAttrOpts.InstanceId = id
		instanceAttrOpts.Attribute = attr
		instanceAttrOpts.InstanceType = &fcu.AttributeValue{Value: aws.String(v.(string))}
		fmt.Printf("\n\n[INFO] Stopping Instance %q for instance_type change", instanceAttrOpts.InstanceId)
		fmt.Printf("\n\n[INFO] instnace_type (%s)", v.(string))

		_, err := conn.VM.StopInstances(&fcu.StopInstancesInput{
			InstanceIds: []*string{instanceAttrOpts.InstanceId},
		})
		if err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
			Target:     []string{"stopped"},
			Refresh:    InstanceStateRefreshFunc(conn, *instanceAttrOpts.InstanceId, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to stop: %s", *instanceAttrOpts.InstanceId, err)
		}

		fmt.Printf("[INFO] Modifying instance type %s", *instanceAttrOpts.InstanceId)
		_, err = conn.VM.ModifyInstanceAttribute(instanceAttrOpts)
		if err != nil {
			return err
		}

		log.Printf("[INFO] Starting Instance %q after instance_type change", instanceAttrOpts.InstanceId)
		_, err = conn.VM.StartInstances(&fcu.StartInstancesInput{
			InstanceIds: []*string{instanceAttrOpts.InstanceId},
		})

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"pending", "stopped"},
			Target:     []string{"running"},
			Refresh:    InstanceStateRefreshFunc(conn, *instanceAttrOpts.InstanceId, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s", *instanceAttrOpts.InstanceId, err)
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		instanceAttrOpts := &fcu.ModifyInstanceAttributeInput{}
		instanceAttrOpts.InstanceId = id
		instanceAttrOpts.Attribute = attr
		instanceAttrOpts.UserData = &fcu.BlobAttributeValue{
			Value: v.([]byte),
		}

		fmt.Printf("\n\n[INFO] Stopping Instance %q for instance_type change", instanceAttrOpts.InstanceId)
		_, err := conn.VM.StopInstances(&fcu.StopInstancesInput{
			InstanceIds: []*string{instanceAttrOpts.InstanceId},
		})
		if err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
			Target:     []string{"stopped"},
			Refresh:    InstanceStateRefreshFunc(conn, *instanceAttrOpts.InstanceId, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to stop: %s", *instanceAttrOpts.InstanceId, err)
		}

		log.Printf("[INFO] Modifying instance type %s", *instanceAttrOpts.InstanceId)
		_, err = conn.VM.ModifyInstanceAttribute(instanceAttrOpts)
		if err != nil {
			return err
		}

		log.Printf("[INFO] Starting Instance %q after user_data change", *instanceAttrOpts.InstanceId)
		_, err = conn.VM.StartInstances(&fcu.StartInstancesInput{
			InstanceIds: []*string{instanceAttrOpts.InstanceId},
		})
		if err != nil {
			return err
		}
		stateConf = &resource.StateChangeConf{
			Pending:    []string{"pending", "stopped"},
			Target:     []string{"running"},
			Refresh:    InstanceStateRefreshFunc(conn, *instanceAttrOpts.InstanceId, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s", *instanceAttrOpts.InstanceId, err)
		}
	}

	if v, ok := d.GetOk("ebs_optimized"); ok {
		instanceAttrOpts := &fcu.ModifyInstanceAttributeInput{}
		instanceAttrOpts.InstanceId = id
		instanceAttrOpts.Attribute = attr
		instanceAttrOpts.EbsOptimized = &fcu.AttributeBooleanValue{
			Value: aws.Bool(v.(bool)),
		}
		fmt.Printf("\n\n[INFO] Stopping Instance %q for ebs_optimized change", *instanceAttrOpts.InstanceId)
		fmt.Printf("\n\n[INFO] ebs_optimized (%v)", v.(bool))
		fmt.Printf("\n\n[INFO] ebs_optimized ok (%v)", ok)

		if instanceAttrOpts.Attribute == nil {
			return fmt.Errorf("the attribute argument must be set to be able, to modify the ebs_optimized attr")
		}

		_, err := conn.VM.StopInstances(&fcu.StopInstancesInput{
			InstanceIds: []*string{instanceAttrOpts.InstanceId},
		})

		if err != nil {
			fmt.Printf("[ERROR] DEBUG (%s)", err)
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
			Target:     []string{"stopped"},
			Refresh:    InstanceStateRefreshFunc(conn, *instanceAttrOpts.InstanceId, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to stop: %s", *instanceAttrOpts.InstanceId, err)
		}

		fmt.Printf("\n\n[INFO] Modifying instance type %s", *instanceAttrOpts.InstanceId)

		_, err = conn.VM.ModifyInstanceAttribute(instanceAttrOpts)
		if err != nil {
			fmt.Printf("[ERROR] DEBUG (%s)", err)
			return err
		}

		fmt.Printf("\n\n[INFO] Starting Instance %q after ebs_optimized change", *instanceAttrOpts.InstanceId)
		_, err = conn.VM.StartInstances(&fcu.StartInstancesInput{
			InstanceIds: []*string{instanceAttrOpts.InstanceId},
		})
		if err != nil {
			fmt.Printf("[ERROR] DEBUG (%s)", err)
			return err
		}
		stateConf = &resource.StateChangeConf{
			Pending:    []string{"pending", "stopped"},
			Target:     []string{"running"},
			Refresh:    InstanceStateRefreshFunc(conn, *instanceAttrOpts.InstanceId, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s", *instanceAttrOpts.InstanceId, err)
		}
	}

	return resourceVMAttributesRead(d, meta)
}

func resourceVMAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	if err := readDescribeVMAttr(d, conn); err != nil {
		return err
	}

	if err := readDescribeVMStatus(d, conn); err != nil {
		return err
	}

	return nil
}

func resourceVMAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	d.Partial(true)

	id := d.Get("instance_id").(string)

	log.Printf("[DEBUG] updating the instance %s", id)

	if d.HasChange("instance_type") && !d.IsNewResource() {
		log.Printf("[INFO] Stopping Instance %q for instance_type change", id)
		_, err := conn.VM.StopInstances(&fcu.StopInstancesInput{
			InstanceIds: []*string{aws.String(id)},
		})

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
			Target:     []string{"stopped"},
			Refresh:    InstanceStateRefreshFunc(conn, id, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to stop: %s", id, err)
		}

		log.Printf("[INFO] Modifying instance type %s", id)
		_, err = conn.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			InstanceType: &fcu.AttributeValue{
				Value: aws.String(d.Get("instance_type").(string)),
			},
		})
		if err != nil {
			return err
		}

		log.Printf("[INFO] Starting Instance %q after instance_type change", id)
		_, err = conn.VM.StartInstances(&fcu.StartInstancesInput{
			InstanceIds: []*string{aws.String(id)},
		})

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"pending", "stopped"},
			Target:     []string{"running"},
			Refresh:    InstanceStateRefreshFunc(conn, id, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				id, err)
		}
	}

	if d.HasChange("user_data") && !d.IsNewResource() {
		log.Printf("[INFO] Stopping Instance %q for instance_type change", id)
		_, err := conn.VM.StopInstances(&fcu.StopInstancesInput{
			InstanceIds: []*string{aws.String(id)},
		})

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
			Target:     []string{"stopped"},
			Refresh:    InstanceStateRefreshFunc(conn, id, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to stop: %s", id, err)
		}

		log.Printf("[INFO] Modifying instance type %s", id)
		_, err = conn.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			UserData: &fcu.BlobAttributeValue{
				Value: d.Get("user_data").([]byte),
			},
		})
		if err != nil {
			return err
		}

		log.Printf("[INFO] Starting Instance %q after user_data change", id)
		_, err = conn.VM.StartInstances(&fcu.StartInstancesInput{
			InstanceIds: []*string{aws.String(id)},
		})

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"pending", "stopped"},
			Target:     []string{"running"},
			Refresh:    InstanceStateRefreshFunc(conn, id, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				id, err)
		}
	}

	if d.HasChange("ebs_optimized") && !d.IsNewResource() {
		log.Printf("[INFO] Stopping Instance %q for instance_type change", id)
		_, err := conn.VM.StopInstances(&fcu.StopInstancesInput{
			InstanceIds: []*string{aws.String(id)},
		})

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
			Target:     []string{"stopped"},
			Refresh:    InstanceStateRefreshFunc(conn, id, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to stop: %s", id, err)
		}

		log.Printf("[INFO] Modifying instance type %s", id)
		_, err = conn.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			EbsOptimized: &fcu.AttributeBooleanValue{
				Value: d.Get("ebs_optimized").(*bool),
			},
		})
		if err != nil {
			return err
		}

		log.Printf("[INFO] Starting Instance %q after ebs_optimized change", id)
		_, err = conn.VM.StartInstances(&fcu.StartInstancesInput{
			InstanceIds: []*string{aws.String(id)},
		})

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"pending", "stopped"},
			Target:     []string{"running"},
			Refresh:    InstanceStateRefreshFunc(conn, id, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				id, err)
		}
	}

	if d.HasChange("delete_on_termination") && !d.IsNewResource() {
		log.Printf("[INFO] Stopping Instance %q for instance_type change", id)
		_, err := conn.VM.StopInstances(&fcu.StopInstancesInput{
			InstanceIds: []*string{aws.String(id)},
		})

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
			Target:     []string{"stopped"},
			Refresh:    InstanceStateRefreshFunc(conn, id, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to stop: %s", id, err)
		}

		log.Printf("[INFO] Modifying instance type %s", id)
		_, err = conn.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			DeleteOnTermination: &fcu.AttributeBooleanValue{
				Value: d.Get("delete_on_termination").(*bool),
			},
		})
		if err != nil {
			return err
		}

		log.Printf("[INFO] Starting Instance %q after delete_on_termination change", id)
		_, err = conn.VM.StartInstances(&fcu.StartInstancesInput{
			InstanceIds: []*string{aws.String(id)},
		})

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"pending", "stopped"},
			Target:     []string{"running"},
			Refresh:    InstanceStateRefreshFunc(conn, id, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				id, err)
		}
	}

	if d.HasChange("disable_api_termination") {
		_, err := conn.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			DisableApiTermination: &fcu.AttributeBooleanValue{
				Value: aws.Bool(d.Get("disable_api_termination").(bool)),
			},
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("instance_initiated_shutdown_behavior") {
		log.Printf("[INFO] Modifying instance %s", id)
		_, err := conn.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			InstanceInitiatedShutdownBehavior: &fcu.AttributeValue{
				Value: aws.String(d.Get("instance_initiated_shutdown_behavior").(string)),
			},
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("group_id") {
		log.Printf("[INFO] Modifying instance %s", id)
		_, err := conn.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			Groups:     d.Get("group_set").([]*string),
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("source_dest_check") {
		log.Printf("[INFO] Modifying instance %s", id)
		_, err := conn.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
			InstanceId: aws.String(id),
			SourceDestCheck: &fcu.AttributeBooleanValue{
				Value: aws.Bool(d.Get("source_dest_check").(bool)),
			},
		})

		if err != nil {
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
				VirtualName: aws.String(f["virtual_name"].(string)),
			}

			e := f["ebs"].(map[string]interface{})

			ebs := &fcu.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(e["delete_on_termination"].(bool)),
				Iops:                aws.Int64(int64(e["iops"].(int))),
				SnapshotId:          aws.String(e["snapshot_id"].(string)),
				VolumeSize:          aws.Int64(int64(e["volume_size"].(int))),
				VolumeType:          aws.String((e["volume_type"].(string))),
			}

			mapping.Ebs = ebs

			mappings = append(mappings, mapping)
		}

		log.Printf("[INFO] Modifying instance %s", id)
		_, err := conn.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
			InstanceId:          aws.String(id),
			BlockDeviceMappings: mappings,
		})

		if err != nil {
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
		InstanceId: aws.String(d.Get("instance_id").(string)),
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

	fmt.Printf("\n\n[DEBUG] RESPONSE %v", resp)

	d.Set("block_device_mapping", getBlockDeviceMapping(resp.BlockDeviceMappings))

	d.Set("disable_api_termination", resp.DisableApiTermination)

	d.Set("ebs_optimized", resp.EbsOptimized)

	err = d.Set("group_set", getGroupSet(resp.Groups))
	if err != nil {
		fmt.Println(getGroupSet(resp.Groups))
	}

	d.Set("instance_id", resp.InstanceId)
	d.SetId(*resp.InstanceId)

	d.Set("instance_initiated_shutdown_behavior", resp.InstanceInitiatedShutdownBehavior)

	d.Set("instance_type", resp.InstanceType)

	d.Set("kernel", resp.KernelId)

	d.Set("product_codes", getProductCodes(resp.ProductCodes))

	d.Set("ramdisk", resp.RamdiskId)

	d.Set("root_device_name", resp.RootDeviceName)

	d.Set("source_dest_check", resp.SourceDestCheck)

	d.Set("sriov_net_support", resp.SriovNetSupport)

	d.Set("user_data", resp.UserData)

	return nil
}

func readDescribeVMStatus(d *schema.ResourceData, conn *fcu.Client) error {
	input := &fcu.DescribeInstanceStatusInput{}

	filters, filtersOk := d.GetOk("filter")
	instancesIds, instancesIdsOk := d.GetOk("instance_ids")
	includeIds, includeIdsOk := d.GetOk("include_all_instances")

	if instancesIdsOk {
		var ids []*string

		for _, id := range instancesIds.([]interface{}) {
			ids = append(ids, aws.String(id.(string)))
		}

		input.InstanceIds = ids
	}

	if filtersOk {
		input.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	if includeIdsOk {
		input.IncludeAllInstances = aws.Bool(includeIds.(bool))
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

			fmt.Printf("\n\n[DEBUG] RESPONSEINSTANCE %v", v)

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
