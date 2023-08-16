package main

import (
	"context"
	"fmt"

	"github.com/hazelcast/hazelcast-go-client"

	gettingstarted "getting-started"
)

func main() {
	// Connection details for cluster
	config, err := gettingstarted.ClientConfig()
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()
	// create the client and connect to the cluster
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Welcome to your Hazelcast Cluster!")
	if err = client.Shutdown(ctx); err != nil {
		panic(err)
	}
}
