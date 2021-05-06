package main

import (
	"flag"
	"fmt"
	"github.com/boxjan/csi-driver-s3/src/driver"
	"os"
	"path"
)

var (
	endpoint    = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName  = flag.String("name", "driver.csi.k8s.io", "name of the driver")
	nodeID      = flag.String("nodeid", "", "node id")
	showVersion = flag.Bool("version", false, "Show version.")

	// Set by the build process
	version = ""
)

func init() {
	flag.Parse()
}

func main() {
	if *showVersion {
		baseName := path.Base(os.Args[0])
		fmt.Println(baseName, version)
		return
	}

	// run
	driver.NewS3Csi(*driverName, *endpoint, *nodeID, version).Run()
}
