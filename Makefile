CMDS=s3plugin
all: build

include release-tools/build.make

REGISTRY_NAME=boxjan/csi-s3

main.version=dev