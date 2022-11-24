# Docker Machine driver for KubeVirt

This is a **Docker Machine** driver for *KubeVirt*. It allows you to use *Docker Machine* to create Docker hosts on top of *KubeVirt*.

*Status*: **Alpha**

## Overview

The driver is designed to be used with *gitlab-runner* and run inside a *Kubernetes* cluster.
No state is persisted between restart, so the driver is not suitable for long-running Docker hosts.

## Requirements

* A Kubernetes cluster with KubeVirt installed
* Access to the Kubernetes API server
* Access to the Kubernetes Pod network

## Example Usage

Deploy the driver as a sleeping Kubernetes Deployment in the **default** namespace:

```bash
# Deploy docker-machine
kubectl apply -k ./deploy/default
```

Exec into the driver pod:

```bash
kubectl exec -i -t docker-machine-5b595cf65d-q76bn -- sh

/ $ docker-machine create --driver kubevirt --help
Usage: docker-machine create [OPTIONS] [arg...]

Create a machine

Description:
   Run 'docker-machine create --driver name --help' to include the create flags for that driver in the help text.

Options:

   --driver, -d "virtualbox"                                                                            Driver to create machine with. [$MACHINE_DRIVER]
   --engine-env [--engine-env option --engine-env option]                                               Specify environment variables to set in the engine
   --engine-insecure-registry [--engine-insecure-registry option --engine-insecure-registry option]     Specify insecure registries to allow with the created engine
   --engine-install-url "https://get.docker.com"                                                        Custom URL to use for engine installation [$MACHINE_DOCKER_INSTALL_URL]
   --engine-label [--engine-label option --engine-label option]                                         Specify labels for the created engine
   --engine-opt [--engine-opt option --engine-opt option]                                               Specify arbitrary flags to include with the created engine in the form flag=value
   --engine-registry-mirror [--engine-registry-mirror option --engine-registry-mirror option]           Specify registry mirrors to use [$ENGINE_REGISTRY_MIRROR]
   --engine-storage-driver                                                                              Specify a storage driver to use with the engine
   --kubevirt-cpu-count "1"                                                                             Number of CPUs
   --kubevirt-image "linkacloud/d2vm-docker-machine:alpine"                                             Container Disk Image to use for the VM, it should provide docker, sshd and qemu-guest-agent
   --kubevirt-kubeconfig                                                                                Path to the Kubernetes config file, if not specified, the default config path or in-cluster config will be used [$KUBECONFIG]
   --kubevirt-memory "1024"                                                                             Size of memory for host in MB
   --kubevirt-namespace "default"                                                                       Namespace to use for the VM, if not specified, the default namespace will be used
   --swarm                                                                                              Configure Machine to join a Swarm cluster
   --swarm-addr                                                                                         addr to advertise for Swarm (default: detect and use the machine IP)
   --swarm-discovery                                                                                    Discovery service to use with Swarm
   --swarm-experimental                                                                                 Enable Swarm experimental features
   --swarm-host "tcp://0.0.0.0:3376"                                                                    ip/socket to listen on for Swarm master
   --swarm-image "swarm:latest"                                                                         Specify Docker image to use for Swarm [$MACHINE_SWARM_IMAGE]
   --swarm-join-opt [--swarm-join-opt option --swarm-join-opt option]                                   Define arbitrary flags for Swarm join
   --swarm-master                                                                                       Configure Machine to be a Swarm master
   --swarm-opt [--swarm-opt option --swarm-opt option]                                                  Define arbitrary flags for Swarm master
   --swarm-strategy "spread"                                                                            Define a default scheduling strategy for Swarm
   --tls-san [--tls-san option --tls-san option]                                                        Support extra SANs for TLS certs
   
# Create a new Docker host with the driver
/ $ docker-machine create --driver kubevirt --kubevirt-cpu-count 2 --kubevirt-memory 2048 docker-0
Creating CA: /root/.docker/machine/certs/ca.pem
Creating client certificate: /root/.docker/machine/certs/cert.pem
Running pre-create checks...
Creating machine...
(docker-0) Creating SSH key...
(docker-0) Creating VM docker-0 in namespace default
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with alpine...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env docker-0

/ $ docker-machine ls
NAME       ACTIVE   DRIVER     STATE     URL                      SWARM   DOCKER      ERRORS
docker-0   -        kubevirt   Running   tcp://10.42.0.152:2376           v20.10.20

# Connect to the Docker host
/ $ eval $(docker-machine env docker-0)

# Run a container on the Docker host
/ $ docker run -d --name whoami -p 80:80 traefik/whoami
Unable to find image 'traefik/whoami:latest' locally
latest: Pulling from traefik/whoami
029cd1bf7e7c: Pull complete
e73b694ead4f: Pull complete
99df6e9e9886: Pull complete
Digest: sha256:24829edb0dbaea072dabd7d902769168542403a8c78a6f743676af431166d7f0
Status: Downloaded newer image for traefik/whoami:latest
e3cb2a13140635bbaefab25034bbcc9b17abe547d66d99e72168cd4703db6a44

/ $ docker ps
CONTAINER ID   IMAGE            COMMAND     CREATED          STATUS          PORTS                               NAMES
e3cb2a131406   traefik/whoami   "/whoami"   24 seconds ago   Up 21 seconds   0.0.0.0:80->80/tcp, :::80->80/tcp   whoami

# Connect to the container
/ $ curl $(docker-machine ip docker-0)
Hostname: e3cb2a131406
IP: 127.0.0.1
IP: 172.17.0.2
RemoteAddr: 10.42.0.162:60974
GET / HTTP/1.1
Host: 10.42.0.163
User-Agent: curl/7.86.0
Accept: */*

# Destroy the Docker host
/ $ docker-machine rm -f -y docker-0
About to remove docker-0
WARNING: This action will delete both local reference and remote instance.
(docker-0) Getting IP address for VM docker-0 in namespace default
(docker-0) Killing vm docker-0 in namespace default
Successfully removed docker-0

```

## Docker images

The **KubeVirt** container disk images are build using [d2vm](https://github.com/linka-cloud/d2vm) and
available in two flavors:
- ubuntu-22.04
- alpine 3.16

The definitions are located in the [images directory](./images).

In order to be able to provision Alpine based docker host, a fork of the [gitlab's docker-machine](https://gitlab.com/linka-cloud/docker-machine)
is used.

The *KubeVirt* images must have the **qemu-guest-agent** installed and running in order to provision the ssh public key.

### TODO:
- [ ] Use recommended resources labels
- [ ] Customize DNS configuration
- [ ] Run as non-root user
- [ ] *Gitlab Runner* example
