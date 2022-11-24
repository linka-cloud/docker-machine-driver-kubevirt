package driver

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "kubevirt.io/api/core/v1"
)

func sshSecretName(name string) string {
	return fmt.Sprintf("%s-ssh", name)
}

func makeSSHSecret(namespace, name, pubkey string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sshSecretName(name),
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"authorized_keys": []byte(pubkey),
		},
	}
}

func makeVM(namespace, name, image string, cpus uint32, mem resource.Quantity) *v1.VirtualMachine {
	return &v1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.VirtualMachineSpec{
			RunStrategy: P(v1.RunStrategyRerunOnFailure),
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"kubevirt.io/domain": name,
					},
				},
				Spec: v1.VirtualMachineInstanceSpec{
					Hostname: name,
					AccessCredentials: []v1.AccessCredential{{
						SSHPublicKey: &v1.SSHPublicKeyAccessCredential{
							Source: v1.SSHPublicKeyAccessCredentialSource{
								Secret: &v1.AccessCredentialSecretSource{
									SecretName: sshSecretName(name),
								},
							},
							PropagationMethod: v1.SSHPublicKeyAccessCredentialPropagationMethod{
								QemuGuestAgent: &v1.QemuGuestAgentSSHPublicKeyAccessCredentialPropagation{
									Users: []string{defaultUser},
								},
							},
						},
					}},
					Domain: v1.DomainSpec{
						CPU: &v1.CPU{
							Cores: cpus,
							Model: "host-passthrough",
						},
						Resources: v1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: mem,
							},
						},
						Devices: v1.Devices{
							AutoattachGraphicsDevice: P(false),
							Rng:                      &v1.Rng{},
							Disks: []v1.Disk{
								{
									Name: defaultUser,
									DiskDevice: v1.DiskDevice{
										Disk: &v1.DiskTarget{Bus: "virtio"},
									},
								},
							},
							Interfaces: []v1.Interface{
								{
									Name: "default",
									InterfaceBindingMethod: v1.InterfaceBindingMethod{
										Masquerade: &v1.InterfaceMasquerade{},
									},
								},
							},
						},
					},
					Networks: []v1.Network{
						{
							Name: "default",
							NetworkSource: v1.NetworkSource{
								Pod: &v1.PodNetwork{},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "root",
							VolumeSource: v1.VolumeSource{
								ContainerDisk: &v1.ContainerDiskSource{
									Image: image,
								},
							},
						},
					},
				},
			},
		},
	}
}

func P[T any](v T) *T {
	return &v
}
