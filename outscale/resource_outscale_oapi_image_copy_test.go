package outscale

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPIImageCopy(t *testing.T) {
	t.Skip()
	var amiID string

	snapshots := []string{}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIImageCopyConfig,
				Check: func(state *terraform.State) error {
					rs, ok := state.RootModule().Resources["outscale_image_copy.test"]
					if !ok {
						return fmt.Errorf("Image resource not found")
					}

					amiID = rs.Primary.ID

					if amiID == "" {
						return fmt.Errorf("Image id is not set")
					}

					conn := testAccProvider.Meta().(*OutscaleClient).FCU
					req := &fcu.DescribeImagesInput{
						ImageIds: []*string{aws.String(amiID)},
					}

					var describe *fcu.DescribeImagesOutput

					err := resource.Retry(5*time.Minute, func() *resource.RetryError {
						var err error
						describe, err = conn.VM.DescribeImages(req)
						if err != nil {
							if strings.Contains(err.Error(), "RequestLimitExceeded:") {
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						return nil
					})

					if err != nil {

						return err
					}

					if len(describe.Images) != 1 ||
						*describe.Images[0].ImageId != rs.Primary.ID {
						return fmt.Errorf("Image not found")
					}

					image := describe.Images[0]
					if expected := "available"; *image.State != expected {
						return fmt.Errorf("invalid image state; expected %v, got %v", expected, image.State)
					}
					if expected := "machine"; *image.ImageType != expected {
						return fmt.Errorf("wrong image type; expected %v, got %v", expected, image.ImageType)
					}
					// if expected := "terraform-acc-ami-copy"; *image.Name != expected {
					// 	return fmt.Errorf("wrong name; expected %v, got %v", expected, image.Name)
					// }

					for _, bdm := range image.BlockDeviceMappings {
						// The snapshot ID might not be set,
						// even for a block device that is an
						// EBS volume.
						if bdm.Ebs != nil && bdm.Ebs.SnapshotId != nil {
							snapshots = append(snapshots, *bdm.Ebs.SnapshotId)
						}
					}

					if expected := 1; len(snapshots) != expected {
						return fmt.Errorf("wrong number of snapshots; expected %v, got %v", expected, len(snapshots))
					}

					return nil
				},
			},
		},
		CheckDestroy: func(state *terraform.State) error {
			conn := testAccProvider.Meta().(*OutscaleClient).FCU
			diReq := &fcu.DescribeImagesInput{
				ImageIds: []*string{aws.String(amiID)},
			}

			var diRes *fcu.DescribeImagesOutput

			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				var err error
				diRes, err = conn.VM.DescribeImages(diReq)
				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidAMIID.NotFound") {
					return nil
				}
				return err
			}

			if len(diRes.Images) > 0 {
				state := diRes.Images[0].State
				return fmt.Errorf("Image %v remains in state %v", amiID, state)
			}

			stillExist := make([]string, 0, len(snapshots))
			checkErrors := make(map[string]error)
			for _, snapshotID := range snapshots {
				dsReq := &fcu.DescribeSnapshotsInput{
					SnapshotIds: []*string{aws.String(snapshotID)},
				}

				var err error
				err = resource.Retry(5*time.Minute, func() *resource.RetryError {

					_, err = conn.VM.DescribeSnapshots(dsReq)
					if err != nil {
						if strings.Contains(err.Error(), "RequestLimitExceeded:") {
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})

				if err == nil {
					stillExist = append(stillExist, snapshotID)
					continue
				}

				awsErr, ok := err.(awserr.Error)
				if !ok {
					checkErrors[snapshotID] = err
					continue
				}

				if awsErr.Code() != "InvalidSnapshot.NotFound" {
					checkErrors[snapshotID] = err
					continue
				}
			}

			if len(stillExist) > 0 || len(checkErrors) > 0 {
				errParts := []string{
					"Expected all snapshots to be gone, but:",
				}
				for _, snapshotID := range stillExist {
					errParts = append(
						errParts,
						fmt.Sprintf("- %v still exists", snapshotID),
					)
				}
				for snapshotID, err := range checkErrors {
					errParts = append(
						errParts,
						fmt.Sprintf("- checking %v gave error: %v", snapshotID, err),
					)
				}
				return errors.New(strings.Join(errParts, "\n"))
			}

			return nil
		},
	})
}

var testAccOutscaleOAPIImageCopyConfig = `
	resource "outscale_vm" "outscale_vm" {
		count    = 1
		image_id = "ami-880caa66"
		type     = "c4.large"
	}

	resource "outscale_image" "outscale_image" {
		name  = "image_${outscale_vm.outscale_vm.id}"
		vm_id = "${outscale_vm.outscale_vm.id}"
	}

	resource "outscale_image_copy" "test" {
		count = 1

		source_image_id    = "${outscale_image.outscale_image.image_id}"
		source_region_name = "eu-west-2"
	}
`
