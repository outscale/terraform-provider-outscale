package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func main() {

	ak := os.Getenv("OUTSCALE_ACCESSKEYID")
	sk := os.Getenv("OUTSCALE_SECRETKEYID")

	config := osc.Config{
		Credentials: &osc.Credentials{
			AccessKey: ak,
			SecretKey: sk,
			Region:    "eu-west-2",
		},
	}

	c, err := fcu.NewFCUClient(config)
	if err != nil {
		fmt.Println(err)
	}
	// keyname := "TestKey"
	// var maxC int64
	// imageID := "ami-8a6a0120"
	// maxC = 1
	// instanceType := "t2.micro"
	// input := fcu.RunInstancesInput{
	// 	ImageId:      &imageID,
	// 	MaxCount:     &maxC,
	// 	MinCount:     &maxC,
	// 	KeyName:      &keyname,
	// 	InstanceType: &instanceType,
	// }
	// output, err := c.VM.RunInstance(&input)
	// fmt.Println(err)
	// fmt.Println(output)

	// input2 := fcu.DescribeInstancesInput{
	// 	InstanceIds: []*string{output.Instances[0].InstanceId},
	// }

	// output2, err := c.VM.DescribeInstances(&input2)
	// fmt.Println(err)
	// fmt.Println(output2)
	// // id := "i-751ebdf6"
	// // input3 := fcu.GetPasswordDataInput{
	// // 	InstanceId: &id,
	// // }
	// //
	// time.Sleep(120 * time.Second)
	//

	// output3, err := c.VM.GetPasswordData(&input3)
	// fmt.Println(err)
	// fmt.Println(output3)

	// fmt.Printf("Key (%+v)\n", output3)
	// fmt.Printf("ID (%+v)\n", *output3.InstanceId)
	// fmt.Printf("Passw (%+v)\n", *output3.PasswordData)

	// output3, err := c.VM.StopInstances(&fcu.StopInstancesInput{
	// 	InstanceIds: []*string{output.Instances[0].InstanceId},
	// })

	// fmt.Println(output3)
	// fmt.Println(err)

	// output3, err := c.VM.StartInstances(&fcu.StartInstancesInput{
	// 	InstanceIds: []*string{aws.String("i-bab4810b")},
	// })

	// fmt.Println(output3)
	// fmt.Println(err)

	output4, err := c.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
		InstanceId: aws.String("i-bab4810b"),
		DisableApiTermination: &fcu.AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	})

	fmt.Println(output4)
	fmt.Println(err)

	//
	// input4 := fcu.ModifyInstanceKeyPairInput{
	// 	InstanceId: output.Instances[0].InstanceId,
	// 	KeyName:    &keyname,
	// }
	// err = c.VM.ModifyInstanceKeyPair(&input4)
	// fmt.Println(err)
	//
	// input6 := fcu.TerminateInstancesInput{
	// 	InstanceIds: []*string{&id},
	// }
	//
	// output6, err := c.VM.TerminateInstances(&input6)
	// fmt.Println(err)
	// fmt.Println(output6)

	// var runResp *ec2.Reservation
	// err = resource.Retry(30*time.Second, func() *resource.RetryError {
	// 	var err error
	// 	runResp, err = conn.RunInstances(runOpts)
	// 	// IAM instance profiles can take ~10 seconds to propagate in AWS:
	// 	// http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html#launch-instance-with-role-console
	// 	if isAWSErr(err, "InvalidParameterValue", "Invalid IAM Instance Profile") {
	// 		log.Print("[DEBUG] Invalid IAM Instance Profile referenced, retrying...")
	// 		return resource.RetryableError(err)
	// 	}
	// 	// IAM roles can also take time to propagate in AWS:
	// 	if isAWSErr(err, "InvalidParameterValue", " has no associated IAM Roles") {
	// 		log.Print("[DEBUG] IAM Instance Profile appears to have no IAM roles, retrying...")
	// 		return resource.RetryableError(err)
	// 	}
	// 	return resource.NonRetryableError(err)
	// })

	// Read the content
	// var bodyBytes []byte
	// if r.Body != nil {
	// 	bodyBytes, _ = ioutil.ReadAll(r.Body)
	// }
	// // Restore the io.ReadCloser to its original state
	// r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// // Use the content
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)
}
