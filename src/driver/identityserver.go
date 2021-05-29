package driver

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *S3csi) GetPluginInfo(ctx context.Context, request *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	ctx = FollowRequest(ctx)

	if s.name == "" {
		return nil, status.Error(codes.Unavailable, "Driver name not configured")
	}

	if s.version == "" {
		return nil, status.Error(codes.Unavailable, "Driver is missing version")
	}

	return &csi.GetPluginInfoResponse{
		Name:          s.name,
		VendorVersion: s.version,
	}, nil
}

func (s *S3csi) GetPluginCapabilities(ctx context.Context, request *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	ctx = FollowRequest(ctx)

	var caps []*csi.PluginCapability

	for _, cap := range s.identitySupportCapabilities {
		c := &csi.PluginCapability{
			Type: &csi.PluginCapability_Service_{
				Service: &csi.PluginCapability_Service{
					Type: cap,
				},
			},
		}
		caps = append(caps, c)
	}

	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: caps,
	}, nil
}

func (s *S3csi) Probe(ctx context.Context, request *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	ctx = FollowRequest(ctx)

	return &csi.ProbeResponse{}, nil
}
