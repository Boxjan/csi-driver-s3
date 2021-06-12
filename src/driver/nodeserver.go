package driver

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

func (s *S3csi) NodePublishVolume(ctx context.Context, request *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Infof("Got NodePublishVolume Request with the req: %v, with context: %v", CleanNodePublishVolumeRequest(request), ctx)

	panic("implement me")
}

func (s *S3csi) NodeUnpublishVolume(ctx context.Context, request *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	ctx = FollowRequest(ctx)

	panic("implement me")
}

func (s *S3csi) NodeStageVolume(ctx context.Context, request *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Infof("Got NodeStageVolumeRequest Request with req: %v, with context: %v", CleanNodeStageVolumeRequest(request), ctx)

	panic("implement me")
}

func (s *S3csi) NodeUnstageVolume(ctx context.Context, request *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	ctx = FollowRequest(ctx)

	panic("implement me")
}

func (s *S3csi) NodeGetVolumeStats(ctx context.Context, request *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Infof("Got NodeGetVolumeStats Request with req: %v, with context: %v", request, ctx)

	return nil, status.Error(codes.Unimplemented, "")
}

func (s *S3csi) NodeGetCapabilities(ctx context.Context, request *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Infof("Got NodeGetCapabilities Request with the req: %v, with context: %v", request, ctx)

	var caps []*csi.NodeServiceCapability

	for _, cap := range s.nodeSupportCapabilities {
		c := &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: cap,
				},
			},
		}
		caps = append(caps, c)
	}

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: caps,
	}, nil
}

func (s *S3csi) NodeExpandVolume(ctx context.Context, request *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Infof("Got NodeExpandVolume Request with req: %v, with context: %v", request, ctx)
	request.VolumeId

	return nil, status.Error(codes.Unimplemented, "")
}

func (s *S3csi) NodeGetInfo(ctx context.Context, request *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Infof("Got NodeGetInfo Request with the req: %v, with context: %v", request, ctx)

	return &csi.NodeGetInfoResponse{
		NodeId:             s.nodeId,
		AccessibleTopology: nil,
	}, nil
}
