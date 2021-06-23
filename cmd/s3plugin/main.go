package main

import (
	"flag"
	"fmt"
	"github.com/boxjan/csi-driver-s3/src/driver"
	"os"
	"path"
)

var (
	showVersion bool
	cfg driver.S3csiDriverConfig
	// Set by the build process
	version = ""
)

func init() {

	flag.StringVar(&cfg.ListenEndpoint, "endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	flag.StringVar(&cfg.Name, "name", "driver.csi.k8s.io", "name of the driver")
	flag.StringVar(&cfg.NodeId, "node-id", "", "node id")

	flag.StringVar(&cfg.StateDir, "state-dir", "/csi-s3-state",
		"directory for storing state information across driver restarts")

	// another bucket for storing state information across different controller in same
	flag.StringVar(&cfg.ControllerStorageBucket.Endpoint, "s3-endpoint", "",
		"the Endpoint of the s3 bucket for storing state information")
	flag.StringVar(&cfg.ControllerStorageBucket.BucketName, "s3-bucket-name", "",
		"the BucketName of the s3 bucket for storing state information")

	// if is empty will try to get from env by key: CTL-S3-ACCESS-KEY
	flag.StringVar(&cfg.ControllerStorageBucket.AccessToken, "s3-access-token", "",
		"the AccessToken of the s3 bucket for storing state information")
	// if is empty will try to get from env by key: CTL-S3-SECRET-KEY
	flag.StringVar(&cfg.ControllerStorageBucket.SecretKey, "s3-secret-key", "",
		"the SecretKey of the s3 bucket for storing state information")

	flag.StringVar(&cfg.ControllerStorageBucket.Region, "s3-region", "",
		"the Region of the s3 bucket for storing state information")
	flag.StringVar(&cfg.ControllerStorageBucket.SubPath, "s3-sub-dir", "",
		"the SubPath directory of the s3 bucket for storing state information")
}

func main() {

	flag.Parse()

	if len(cfg.ControllerStorageBucket.AccessToken) == 0 {
		cfg.ControllerStorageBucket.AccessToken = os.Getenv("CTL-S3-ACCESS-KEY")
	}

	if len(cfg.ControllerStorageBucket.SecretKey) == 0 {
		cfg.ControllerStorageBucket.SecretKey = os.Getenv("CTL-S3-SECRET-KEY")
	}

	cfg.Version = version

	if showVersion {
		baseName := path.Base(os.Args[0])
		fmt.Println(baseName, version)
		return
	}

	// run
	driver.NewS3Csi(&cfg).Run()
}
