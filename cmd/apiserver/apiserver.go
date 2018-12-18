package main

import (
	"net"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	registryHost string
	shadowHost   string
)

func init() {
	viper.SetDefault("REGISTRY_HOST", "device-registry:8080")
	viper.SetDefault("SHADOW_HOST", "shadow-api:8096")
	viper.AutomaticEnv()

	registryHost = viper.GetString("REGISTRY_HOST")
	shadowHost = viper.GetString("SHADOW_HOST")
}

func main() {
	srv := grpc.NewServer()

	registryConn, err := grpc.Dial(registryHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	devicesClient := registrypb.NewDevicesClient(registryConn)

	shadowConn, err := grpc.Dial(shadowHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	shadowClient := shadowpb.NewShadowsClient(shadowConn)

	apipb.RegisterDevicesServer(srv, &deviceAPI{client: devicesClient})
	apipb.RegisterShadowsServer(srv, &shadowAPI{client: shadowClient})
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	reflection.Register(srv)

	err = srv.Serve(listener)
	if err != nil {
		panic(err)
	}

}
