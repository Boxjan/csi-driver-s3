FROM debian
LABEL maintainers="Boxjan"
LABEL description="HostPath Driver"
ARG binary=./bin/s3plugin

# Add util-linux to get a new version of losetup.
RUN apt-get update && \
  apt-get install -y \
  s3fs curl unzip && \
  rm -rf /var/lib/apt/lists/*

# install rclone
ARG RCLONE_VERSION=v1.55.0
RUN cd /tmp \
  && curl -O https://downloads.rclone.org/${RCLONE_VERSION}/rclone-${RCLONE_VERSION}-linux-amd64.zip \
  && unzip /tmp/rclone-${RCLONE_VERSION}-linux-amd64.zip \
  && mv /tmp/rclone-*-linux-amd64/rclone /usr/bin \
  && rm -r /tmp/rclone*

COPY ${binary} /s3plugin
ENTRYPOINT ["/s3plugin"]