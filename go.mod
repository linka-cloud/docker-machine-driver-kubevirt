module go.linka.cloud/docker-machine-driver-kubevirt

go 1.18

require (
	github.com/docker/machine v0.16.2
	github.com/stretchr/testify v1.8.1
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v12.0.0+incompatible
	kubevirt.io/api v0.58.0
	kubevirt.io/client-go v0.58.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/coreos/prometheus-operator v0.38.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/mock v1.5.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v0.0.0-20191119172530-79f836b90111 // indirect
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/openshift/api v0.0.0-20210105115604-44119421ec6b // indirect
	github.com/openshift/client-go v0.0.0-20210112165513-ebc401615f47 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samalba/dockerclient v0.0.0-20160531175551-a30362618471 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.3.0 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/term v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.23.5 // indirect
	k8s.io/klog/v2 v2.40.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220124234850-424119656bbf // indirect
	k8s.io/utils v0.0.0-20211116205334-6203023598ed // indirect
	kubevirt.io/containerized-data-importer-api v1.50.0 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20170925200725-a06ad2792ab9
	github.com/docker/machine => gitlab.com/gitlab-org/ci-cd/docker-machine v0.16.2-gitlab.18
	github.com/openshift/api => github.com/openshift/api v0.0.0-20210105115604-44119421ec6b

	k8s.io/api => k8s.io/api v0.23.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.23.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.23.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.23.5
	k8s.io/client-go => k8s.io/client-go v0.23.5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.23.5
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.23.5
	k8s.io/cri-api => k8s.io/cri-api v0.23.5
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.23.5
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.23.5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.23.5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20210113233702-8566a335510f
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.23.5
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.23.5
	k8s.io/kubectl => k8s.io/kubectl v0.23.5
	k8s.io/kubelet => k8s.io/kubelet v0.23.5
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.23.5
	k8s.io/metrics => k8s.io/metrics v0.23.5
	k8s.io/node-api => k8s.io/node-api v0.23.5
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.23.5
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.23.5
	k8s.io/sample-controller => k8s.io/sample-controller v0.23.5
)
