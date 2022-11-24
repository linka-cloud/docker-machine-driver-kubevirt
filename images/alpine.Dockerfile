FROM alpine:3.16

RUN apk add --no-cache curl \
    openrc \
    qemu-guest-agent \
    htop \
    openssh-server \
    docker && \
      rc-update add sshd default && \
      rc-update add qemu-guest-agent default && \
      rc-update add docker default && \
      echo "PermitRootLogin prohibit-password" >> /etc/ssh/sshd_config
