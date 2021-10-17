package train

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
)

const (
	prefixInput  = "input/data/train/"
)

func mustLoadConfig() aws.Config {

	c, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Panicf("unable to load aws default config: %v", err)

	}
	return c
}
