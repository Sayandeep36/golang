package common

import (
	"os"
	"time"

	"Golang/proto"
)

var (

	// topic and subscription-related variables
	TopicName           = "/data/AccountChangeEvent"
	ReplayPreset        = proto.ReplayPreset_EARLIEST
	ReplayId     []byte = nil
	Appetite     int32  = 5

	// gRPC server variables
	GRPCEndpoint    = "api.pubsub.salesforce.com:7443"
	GRPCDialTimeout = 5 * time.Second
	GRPCCallTimeout = 5 * time.Second

	// OAuth header variables
	GrantType    = os.Getenv("GrantType")
	ClientId     = os.Getenv("ClientId")
	ClientSecret = os.Getenv("ClientSecret")
	Username     = os.Getenv("Username")
	Password     = os.Getenv("Password")

	// OAuth server variables
	OAuthEndpoint    = os.Getenv("OAuthEndpoint")
	OAuthDialTimeout = 5 * time.Second
)
