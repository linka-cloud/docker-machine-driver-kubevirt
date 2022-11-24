FROM ubuntu

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
        qemu-guest-agent \
        openssh-server \
        docker.io \
        sudo \
        net-tools \
        && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    echo "PermitRootLogin prohibit-password" >> /etc/ssh/sshd_config
