package s3

import (
	csicommon "github.com/boxjan/csi-driver-s3/src/csi-common"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

type S3csi struct {
	name     string
	endpoint string
	nodeId   string
	version  string

	Serv csicommon.NonBlockingGRPCServer

	identitySupportCapabilities   []csi.PluginCapability_Service_Type
	controllerSupportCapabilities []csi.ControllerServiceCapability_RPC_Type
	nodeSupportCapabilities       []csi.NodeServiceCapability_RPC_Type
}

type StorageClassConfig struct {
	S3Endpoint string
	Region     string
	UseSSL     bool
	Mounter    string

	// if use exists bucket
	BucketName string

	// digitalocean S3
	DigitaloceanUseEdgeCdn bool
}

type S3csiDriverConfig struct {
	Name           string
	NodeId         string
	ListenEndpoint string
	Version        string
}

type AllNeedStruct struct {
	BucketName    string
	AccessToken   string
	SecretKey     string
	MountEndpoint string
	ApiEndpoint   string
	Region        string
	UseSSL        string
}

func NewS3Csi(name, endpoint, nodeId, version string) *S3csi {
	if csicommon.NewCSIDriver(name, nodeId, version) == nil {
		klog.Fatalf("driver init failed")
	}
	return &S3csi{
		name:     name,
		endpoint: endpoint,
		nodeId:   nodeId,
		version:  version,
	}
}

func (s *S3csi) Run() {

	s.identitySupportCapabilities = append(s.identitySupportCapabilities,
		csi.PluginCapability_Service_CONTROLLER_SERVICE,
	)

	s.controllerSupportCapabilities = append(s.controllerSupportCapabilities,
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	)

	s.nodeSupportCapabilities = append(s.nodeSupportCapabilities)

	s.Serv = csicommon.NewNonBlockingGRPCServer()
	s.Serv.Start(s.endpoint, s, s, s)
	s.Serv.Wait()
}
