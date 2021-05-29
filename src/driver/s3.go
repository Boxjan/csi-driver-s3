package driver

import (
	"context"
	csicommon "github.com/boxjan/csi-driver-s3/src/csi-common"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/google/uuid"
	"k8s.io/klog/v2"
	"time"
)

type S3csi struct {
	name     string
	endpoint string
	nodeId   string
	version  string

	bucketPrefix string

	Serv csicommon.NonBlockingGRPCServer

	identitySupportCapabilities   []csi.PluginCapability_Service_Type
	controllerSupportCapabilities []csi.ControllerServiceCapability_RPC_Type
	nodeSupportCapabilities       []csi.NodeServiceCapability_RPC_Type
}

type StorageClassConfig struct {
	Endpoint string
	Region   string
	UseSSL   bool
	Mounter  string

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

const (
	TokenGenInterval = 1 * time.Millisecond // 1000 QPS here
)

var (
	requestTokenPool chan<- string
)

func init() {

	go runTokenGen()
}

func runTokenGen() {
	poolSize := int64(1*time.Second/TokenGenInterval) + 4
	requestTokenPoolFull := make(chan string, poolSize)
	requestTokenPool = requestTokenPoolFull

	for _ = range time.NewTicker(TokenGenInterval).C {
		requestTokenPoolFull <- uuid.New().String()
	}
}

func FollowRequest(ctx context.Context) context.Context {
	// now ignore all context cancel request

	//if ctx != nil {
	//	return context.WithValue(ctx, "RequestId", requestTokenPool)
	//}

	return context.WithValue(context.Background(), "RequestId", requestTokenPool)
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
