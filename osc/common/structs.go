package common

type Tag struct {
	_ struct{} `type:"structure"`

	Key *string `locationName:"key" type:"string"`

	Value *string `locationName:"value" type:"string"`
}
