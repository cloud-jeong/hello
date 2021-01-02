package cmd

import (
	kubeadmapiv1beta2 "github.com/cloud-jeong/sandbox/cmd/kubeadm/app/apis/kubeadm/v1beta2"
	"github.com/lithammer/dedent"
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	"text/template"
)

var (
	initDoneTempl = template.Must(template.New("init").Parse(dedent.Dedent(`
		Your Kubernetes control-plane has initialized successfully!

		To start using your cluster, you need to run the following as a regular user:

		  mkdir -p $HOME/.kube
		  sudo cp -i {{.KubeConfigPath}} $HOME/.kube/config
		  sudo chown $(id -u):$(id -g) $HOME/.kube/config

		Alternatively, if you are the root user, you can run:

		  export KUBECONFIG=/etc/kubernetes/admin.conf

		You should now deploy a pod network to the cluster.
		Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
		  https://kubernetes.io/docs/concepts/cluster-administration/addons/

		{{if .ControlPlaneEndpoint -}}
		{{if .UploadCerts -}}
		You can now join any number of the control-plane node running the following command on each as root:

		  {{.joinControlPlaneCommand}}

		Please note that the certificate-key gives access to cluster sensitive data, keep it secret!
		As a safeguard, uploaded-certs will be deleted in two hours; If necessary, you can use
		"kubeadm init phase upload-certs --upload-certs" to reload certs afterward.

		{{else -}}
		You can now join any number of control-plane nodes by copying certificate authorities
		and service account keys on each node and then running the following as root:

		  {{.joinControlPlaneCommand}}

		{{end}}{{end}}Then you can join any number of worker nodes by running the following on each as root:

		{{.joinWorkerCommand}}
		`)))
)

type initOptions struct {
	cfgPath                 string
	skipTokenPrint          bool
	dryRun                  bool
	kubeconfigDir           string
	kubeconfigPath          string
	featureGatesString      string
	ignorePreflightErrors   []string
	bto                     *optons.BootstrapTokenOptions
	externalInitCfg         *kubeadmapiv1beta2.InitConfiguration
	externalClusterCfg      *kubeadmapiv1beta2.ClusterConfiguration
	uploadCerts             bool
	skipCertificateKeyPrint bool
	patchesDir              string
}

var _ phases.InitData = &initData{}

type initData struct {
	cfg                     *kubeadmapi.InitConfiguration
	skipTokenPrint          bool
	dryRun                  bool
	kubeconfigDir           string
	kubeconfigPath          string
	ignorePreflightErrors   sets.String
	certificatesDir         string
	dryRunDir               string
	externalCA              bool
	client                  clientset.Interface
	ouputWriter             io.Writer
	uploadCerts             bool
	skipCertificateKeyPrint bool
	patchesDir              string
}
