package types

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"time"
)

var timeFormat = time.RFC3339

func toTags(in []types.Tag) map[string]string {
	out := make(map[string]string)
	for _, v := range in {
		out[aws.ToString(v.Key)] = aws.ToString(v.Value)
	}
	return out
}
