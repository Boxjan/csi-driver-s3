package driver

import (
	"context"
	csicommon "github.com/boxjan/csi-driver-s3/src/csi-common"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"k8s.io/klog/v2"
	"net/url"
	"os"
	"path"
	"time"
)

type S3csi struct {
	name     string
	endpoint string
	nodeId   string
	version  string

	// directory for storing state information across driver restarts
	stateDir    string
	nodeStatDir string

	// the s3 bucket for storing state information across driver restarts
	controllerBucketClient *minio.Client

	// controller server will not init node server
	isController bool

	// bucket
	bucketPrefix string

	emptyDriver *csicommon.CSIDriver
	Serv        csicommon.NonBlockingGRPCServer

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
	Name                    string
	NodeId                  string
	ListenEndpoint          string
	Version                 string
	StateDir                string
	ControllerStorageBucket *S3ConnInfo
	IamController           bool
}

type S3ConnInfo struct {
	Endpoint    string
	BucketName  string
	AccessToken string
	SecretKey   string
	Region      string
	SubPath     string
}

type MountInfo struct {
	S3ConnInfo
	Mounter       string
	MountEndpoint string
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

func NewS3Csi(cfg *S3csiDriverConfig) *S3csi {
	emptyDriver := csicommon.NewCSIDriver(cfg.Name, cfg.Version, cfg.NodeId)
	if emptyDriver == nil {
		klog.Fatalf("driver init failed")
	}

	if err := os.MkdirAll(cfg.StateDir, 0750); err != nil {
		klog.Fatalf("create data dir: %s failed with err: %v", cfg.StateDir, err)
	}

	s3 := &S3csi{
		name:         cfg.Name,
		endpoint:     cfg.ListenEndpoint,
		nodeId:       cfg.NodeId,
		version:      cfg.Version,
		stateDir:     cfg.StateDir,
		emptyDriver:  emptyDriver,
		isController: cfg.IamController,
	}

	if s3.isController {
		// only controller need init controller storage bucket client
		endpointU, err := url.Parse(cfg.ControllerStorageBucket.Endpoint)
		if err != nil {
			klog.Fatalf("parse controller storage bucket endpoint failed with err: %v", err)
		}

		secure := endpointU.Scheme != "http"

		controllerBucketClient, err := minio.New(cfg.ControllerStorageBucket.Endpoint,
			&minio.Options{
				Creds: credentials.NewStaticV4(cfg.ControllerStorageBucket.AccessToken,
					cfg.ControllerStorageBucket.SecretKey, ""),
				Secure:       secure,
				Region:       cfg.ControllerStorageBucket.Region,
				BucketLookup: 0,
			})
		if err != nil {
			klog.Fatalf("init ControllerStorageBucket client failed with err: %v", err)
		}
		exists, err := controllerBucketClient.BucketExists(context.Background(),
			cfg.ControllerStorageBucket.BucketName)
		if err != nil {
			klog.Fatalf("use ControllerStorageBucket Client get bucket: %v with err: %v",
				cfg.ControllerStorageBucket.BucketName, err)
		} else if !exists {
			klog.Warningf("controller-storage-bucket-client not found bucket in endpoint: %s"+
				", region: %v; will try to create it",
				cfg.ControllerStorageBucket.Endpoint, cfg.ControllerStorageBucket.Region)
			if err := controllerBucketClient.MakeBucket(context.Background(),
				cfg.ControllerStorageBucket.BucketName, minio.MakeBucketOptions{}); err != nil {
				klog.Fatalf("controller-storage-bucket-client create bucket failed with err: %v", err)
			}
		}
		s3.controllerBucketClient = controllerBucketClient
	} else {
		// run as node server
		s3.nodeStatDir = path.Join(s3.stateDir, "node", cfg.StateDir)

		if err := os.MkdirAll(s3.nodeStatDir, 0750); err != nil {
			klog.Fatalf("create data dir: %s failed with err: %v", s3.nodeStatDir, err)
		}
	}

	return s3
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

	if s.isController {
		// controller will only
		klog.Warningf("server will now run under ")
		s.Serv.Start(s.endpoint, s, csicommon.NewDefaultControllerServer(s.emptyDriver), s)
	} else {
		s.Serv.Start(s.endpoint, s, s, csicommon.NewDefaultNodeServer(s.emptyDriver))
	}

	s.Serv.Wait()
}
