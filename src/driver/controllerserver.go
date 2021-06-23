package driver

import (
	"context"
	"github.com/boxjan/csi-driver-s3/src/s3"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	"path"
)

func (s *S3csi) CreateVolume(ctx context.Context, request *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Info("Got CreateVolume Request with the req: %v, with context: %v", CleanCreateVolumeRequestSecret(request), ctx)

	params := request.GetParameters()
	// because params will not appear in other request, so it will be encode into volumeId

	secret := request.GetSecrets()
	volumeId := s.sanitizeVolumeName(request.GetName())
	volumeSize := request.GetCapacityRange().GetRequiredBytes()
	_ = request.GetCapacityRange().GetLimitBytes() // driver hard to limit volume size. ignore now

	if v, ok := params["bucketName"]; ok && v != "" {
		volumeId = path.Join(v, volumeId)
	}

	klog.V(5).Info("get s3 client from params and secret")
	client, err := s3.NewS3Client(&params, &secret)
	if err != nil {
		klog.Error("new s3 Client failed with error:", err)
		return nil, err
	}

	klog.V(5).Infof("check volume: %s is exists", volumeId)
	exists, err := client.VolumeExists(ctx, volumeId)
	if err != nil {
		klog.Error("check volume failed with error:", err)
		return nil, err
	}

	if exists {
		klog.V(5).Info("volume: ", volumeId, "exists, will not create it")
	} else {
		klog.V(5).Info("volume: ", volumeId, "not exists, will create it")
		err := client.CreateVolume(ctx, volumeId)
		if err != nil {
			klog.Error("create volume failed with error:", err)
			return nil, err
		}
	}

	err = s.PutMete(volumeId, params)
	if err != nil {
		return nil, err
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			CapacityBytes: volumeSize,
			VolumeId:      volumeId,
			VolumeContext: request.GetParameters(),
		},
	}, nil

}

func (s *S3csi) DeleteVolume(ctx context.Context, request *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Info("Got CreateVolume Request with the req: %v, with context: %v", CleanDeleteVolumeRequestSecret(request), ctx)

	volumeId := request.GetVolumeId()
	params, err := s.GetMeta(volumeId)
	if err != nil {
		klog.Errorf("get storage: %s meta info failed with err: %+v", volumeId, err)
		klog.Error("we have lost ", volumeId, "control, will return a success to k8s, need to delete file by manual")
		return &csi.DeleteVolumeResponse{}, nil
	}

	secret := request.GetSecrets()
	if v, ok := params["bucketName"]; ok && v != "" {
		volumeId = path.Join(v, volumeId)
	}

	client, err := s3.NewS3Client(&params, &secret)
	if err != nil {
		klog.Error("new s3 Client failed with error:", err)
		return nil, err
	}

	if err := client.DeleteVolume(ctx, volumeId); err != nil {
		// if have any err when delete volume, will put meta to ensure not out of ctl.
		_ = s.PutMete(volumeId, params)
		return nil, err
	}
	return &csi.DeleteVolumeResponse{}, nil
}

func (s *S3csi) ControllerPublishVolume(ctx context.Context, request *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	ctx = FollowRequest(ctx)

	return nil, status.Error(codes.Unimplemented, "not support")
}

func (s *S3csi) ControllerUnpublishVolume(ctx context.Context, request *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	ctx = FollowRequest(ctx)
	//
	return nil, status.Error(codes.Unimplemented, "not support")
}

func (s *S3csi) ValidateVolumeCapabilities(ctx context.Context, request *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	ctx = FollowRequest(ctx)

	panic("implement me")
}

func (s *S3csi) ListVolumes(ctx context.Context, request *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	ctx = FollowRequest(ctx)

	panic("implement me")
}

func (s *S3csi) GetCapacity(ctx context.Context, request *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	ctx = FollowRequest(ctx)

	panic("implement me")
}

func (s *S3csi) ControllerGetCapabilities(ctx context.Context, request *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	ctx = FollowRequest(ctx)

	var caps []*csi.ControllerServiceCapability

	for _, cap := range s.controllerSupportCapabilities {
		c := &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: cap,
				},
			},
		}
		caps = append(caps, c)
	}

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: caps,
	}, nil
}

func (s *S3csi) CreateSnapshot(ctx context.Context, request *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	ctx = FollowRequest(ctx)
	// s3 snapshot will not a atom operate.
	return nil, status.Error(codes.Unimplemented, "not support")
}

func (s *S3csi) DeleteSnapshot(ctx context.Context, request *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	ctx = FollowRequest(ctx)
	// s3 snapshot will not a atom operate.
	return nil, status.Error(codes.Unimplemented, "not support")
}

func (s *S3csi) ListSnapshots(ctx context.Context, request *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	ctx = FollowRequest(ctx)
	// s3 snapshot will not a atom operate.
	return nil, status.Error(codes.Unimplemented, "not support")
}

func (s *S3csi) ControllerExpandVolume(ctx context.Context, request *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Info("Got CreateVolume Request with the req: %v, with context: %v", CleanControllerExpandVolume(request), ctx)

	// will success every time
	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         request.CapacityRange.GetRequiredBytes(),
		NodeExpansionRequired: false,
	}, nil
}

func (s *S3csi) ControllerGetVolume(ctx context.Context, request *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	ctx = FollowRequest(ctx)
	klog.V(5).Info("Got ControllerGetVolume Request with the req: %v, with context: %v", request, ctx)

	return &csi.ControllerGetVolumeResponse{
		Volume: &csi.Volume{
			CapacityBytes:      0,
			VolumeId:           "",
			VolumeContext:      nil,
			ContentSource:      nil,
			AccessibleTopology: nil,
		},
		Status: &csi.ControllerGetVolumeResponse_VolumeStatus{
			PublishedNodeIds: nil,
			VolumeCondition:  nil,
		},
	}, nil
}
