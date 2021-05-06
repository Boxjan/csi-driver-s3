module github.com/boxjan/csi-driver-s3

go 1.16

require (
	github.com/container-storage-interface/spec v1.4.0
	github.com/digitalocean/godo v1.60.0
	github.com/kubernetes-csi/csi-lib-utils v0.9.1
	github.com/minio/minio-go/v7 v7.0.10
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/grpc v1.36.1
	k8s.io/apimachinery v0.19.0
	k8s.io/klog/v2 v2.2.0
)
