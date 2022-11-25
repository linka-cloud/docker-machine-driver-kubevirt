// Copyright 2022 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package driver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	v1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
)

const (
	isoFilename      = "boot2docker.iso"
	defaultUser      = "root"
	defaultImage     = "linkacloud/d2vm-docker-machine:alpine"
	defaultNamespace = "default"
)

var (
	_ drivers.Driver = (*Driver)(nil)

	ErrUnimplemented = errors.New("unimplemented")
)

func New(hostName, storePath string) *Driver {
	log.Debugf("Creating new driver for host %s (%s)", hostName, storePath)
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}

type Driver struct {
	*drivers.BaseDriver
	Memory     int
	CPU        int
	Namespace  string
	Image      string
	KubeConfig string
	c          kubecli.KubevirtClient
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.IntFlag{
			Name:  flagName("memory"),
			Usage: "Size of memory for host in MB",
			Value: 1024,
		},
		mcnflag.IntFlag{
			Name:  flagName("cpu-count"),
			Usage: "Number of CPUs",
			Value: 1,
		},
		mcnflag.StringFlag{
			Name:  flagName("image"),
			Usage: "Container Disk Image to use for the VM, it should provide docker, sshd and qemu-guest-agent",
			Value: defaultImage,
		},
		mcnflag.StringFlag{
			Name:  flagName("namespace"),
			Usage: "Namespace to use for the VM, if not specified, the default namespace will be used",
			Value: defaultNamespace,
		},
		mcnflag.StringFlag{
			Name:   flagName("kubeconfig"),
			Usage:  "Path to the Kubernetes config file, if not specified, the default config path or in-cluster config will be used",
			EnvVar: clientcmd.RecommendedConfigPathEnvVar,
		},
	}
}

func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	log.Debugf("Setting config from flags %+v", opts)
	d.Memory = opts.Int(flagName("memory"))
	d.CPU = opts.Int(flagName("cpu-count"))
	d.Image = opts.String(flagName("image"))
	d.Namespace = opts.String(flagName("namespace"))
	d.KubeConfig = opts.String(flagName("kubeconfig"))
	if _, err := d.client(); err != nil {
		return err
	}
	return nil
}

func (d *Driver) Create() error {
	log.Debugf("Checking if VM %s already exists in namespace %s", d.MachineName, d.Namespace)
	c, err := d.client()
	if err != nil {
		return err
	}
	if _, err := c.VirtualMachine(d.Namespace).Get(d.MachineName, &metav1.GetOptions{}); !apierrors.IsNotFound(err) {
		if err == nil {
			return fmt.Errorf("vm %s already exists in namespace %s", d.MachineName, d.Namespace)
		}
		return err
	}
	log.Infof("Creating SSH key...")
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}
	publicKey, err := os.ReadFile(d.GetSSHKeyPath() + ".pub")
	if err != nil {
		return err
	}
	mem, err := resource.ParseQuantity(fmt.Sprintf("%dMi", d.Memory))
	if err != nil {
		return fmt.Errorf("failed to parse memory quantity: %w", err)
	}
	if _, err := c.CoreV1().Secrets(d.Namespace).Create(context.TODO(), makeSSHSecret(d.Namespace, d.MachineName, string(publicKey)), metav1.CreateOptions{}); err != nil {
		return err
	}
	vm := makeVM(d.Namespace, d.MachineName, d.Image, uint32(d.CPU), mem)
	log.Infof("Creating VM %s in namespace %s", d.MachineName, d.Namespace)
	if _, err := c.VirtualMachine(d.Namespace).Create(vm); err != nil {
		if err := c.CoreV1().Secrets(d.Namespace).Delete(context.TODO(), d.MachineName, metav1.DeleteOptions{}); err != nil {
			log.Warnf("Failed to delete SSH secret: %v", err)
		}
		return fmt.Errorf("failed to create vm: %w", err)
	}
	return nil
}

func (d *Driver) GetSSHHostname() (string, error) {
	c, err := d.client()
	if err != nil {
		return "", err
	}
	log.Debugf("Getting IP address for VM %s in namespace %s", d.MachineName, d.Namespace)
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}
	vm, err := c.VirtualMachineInstance(d.Namespace).Get(d.MachineName, &metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if len(vm.Status.Interfaces) == 0 {
		return "", fmt.Errorf("no interfaces found for vm %s", d.MachineName)
	}
	return vm.Status.Interfaces[0].IP, nil
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetSSHHostname()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

func (d *Driver) GetIP() (string, error) {
	return d.GetSSHHostname()
}

func (d *Driver) GetState() (state.State, error) {
	log.Debugf("Getting state for VM %s in namespace %s", d.MachineName, d.Namespace)
	c, err := d.client()
	if err != nil {
		return state.None, err
	}
	vm, err := c.VirtualMachine(d.Namespace).Get(d.MachineName, &metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return state.None, nil
		}
		return state.Error, fmt.Errorf("failed to get vm: %w", err)
	}
	switch vm.Status.PrintableStatus {
	case v1.VirtualMachineStatusProvisioning, v1.VirtualMachineStatusStarting, v1.VirtualMachineStatusWaitingForVolumeBinding:
		return state.Starting, nil
	case v1.VirtualMachineStatusRunning, v1.VirtualMachineStatusMigrating:
		return state.Running, nil
	case v1.VirtualMachineStatusPaused:
		return state.Paused, nil
	case v1.VirtualMachineStatusStopping, v1.VirtualMachineStatusTerminating:
		return state.Stopping, nil
	case v1.VirtualMachineStatusStopped:
		return state.Stopped, nil
	case v1.VirtualMachineStatusCrashLoopBackOff,
		v1.VirtualMachineStatusUnknown,
		v1.VirtualMachineStatusUnschedulable,
		v1.VirtualMachineStatusErrImagePull,
		v1.VirtualMachineStatusImagePullBackOff,
		v1.VirtualMachineStatusPvcNotFound,
		v1.VirtualMachineStatusDataVolumeError:
		return state.Error, nil
	default:
		return state.None, nil
	}
}

func (d *Driver) Kill() error {
	log.Infof("Killing vm %s in namespace %s", d.MachineName, d.Namespace)
	c, err := d.client()
	if err != nil {
		return err
	}
	return c.VirtualMachine(d.Namespace).Stop(d.MachineName, &v1.StopOptions{GracePeriod: P(int64(0))})
}

func (d *Driver) Remove() error {
	log.Infof("Getting IP address for VM %s in namespace %s", d.MachineName, d.Namespace)
	c, err := d.client()
	if err != nil {
		return err
	}
	s, err := d.GetState()
	if err != nil {
		return err
	}
	if s == state.Running {
		if err := d.Kill(); err != nil {
			return err
		}
	}
	if err := c.VirtualMachine(d.Namespace).Delete(d.MachineName, &metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	if err := c.CoreV1().Secrets(d.Namespace).Delete(context.TODO(), sshSecretName(d.MachineName), metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (d *Driver) Restart() error {
	log.Infof("Restarting vm %s in namespace %s", d.MachineName, d.Namespace)
	c, err := d.client()
	if err != nil {
		return err
	}
	return c.VirtualMachine(d.Namespace).Restart(d.MachineName, &v1.RestartOptions{})
}

func (d *Driver) Start() error {
	log.Infof("Starting vm %s in namespace %s", d.MachineName, d.Namespace)
	c, err := d.client()
	if err != nil {
		return err
	}
	return c.VirtualMachine(d.Namespace).Start(d.MachineName, &v1.StartOptions{})
}

func (d *Driver) Stop() error {
	log.Infof("Stopping vm %s in namespace %s", d.MachineName, d.Namespace)
	c, err := d.client()
	if err != nil {
		return err
	}
	return c.VirtualMachine(d.Namespace).Stop(d.MachineName, &v1.StopOptions{GracePeriod: P(int64(0))})
}

func (d *Driver) DriverName() string {
	return "kubevirt"
}

func (d *Driver) client() (kubecli.KubevirtClient, error) {
	if d.c != nil {
		return d.c, nil
	}
	config, err := loadKubeConfig(d.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubernetes config: %w", err)
	}
	d.c, err = kubecli.GetKubevirtClientFromRESTConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubevirt client: %w", err)
	}
	return d.c, nil
}

func flagName(name string) string {
	return "kubevirt-" + name
}

func loadKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		kubeconfigPath = path.Join(homedir.HomeDir(), clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName)
	}
	if _, err := os.Stat(kubeconfigPath); err == nil {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	return rest.InClusterConfig()
}
