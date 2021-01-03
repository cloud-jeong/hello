package cmd

import (
	kubeadmapi "github.com/cloud-jeong/sandbox/cmd/kubeadm/app/apis/kubeadm"
	kubeadmapiv1beta2 "github.com/cloud-jeong/sandbox/cmd/kubeadm/app/apis/kubeadm/v1beta2"
	"github.com/cloud-jeong/sandbox/cmd/kubeadm/app/cmd/options"
	cmdutil "github.com/cloud-jeong/sandbox/cmd/kubeadm/app/cmd/util"
	kubeadmconstants "github.com/cloud-jeong/sandbox/cmd/kubeadm/app/constants"
	phases "github.com/cloud-jeong/sandbox/cmd/kubeadm/app/phases/init"
	"github.com/lithammer/dedent"
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	"path/filepath"
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
	bto                     *options.BootstrapTokenOptions
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
	outputWriter            io.Writer
	uploadCerts             bool
	skipCertificateKeyPrint bool
	patchesDir              string
}

func (d *initData) UploadCerts() bool {
	return d.uploadCerts
}

func (d *initData) CertificateKey() string {
	return d.cfg.CertificateKey
}

func (d *initData) SetCertificateKey(key string) {
	d.cfg.CertificateKey = key
}

func (d *initData) SkipCertificateKeyPrint() bool {
	return d.skipCertificateKeyPrint
}

func (d *initData) Cfg() *kubeadmapi.InitConfiguration {
	return d.cfg
}

func (d *initData) DryRun() bool {
	return d.dryRun
}

func (d *initData) SkipTokenPrint() bool {
	return d.skipTokenPrint
}

func (d *initData) IgnorePreflightErrors() sets.String {
	return d.ignorePreflightErrors
}

func (d *initData) CertificateWriteDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return d.certificatesDir
}

func (d *initData) CertificateDir() string {
	return d.certificatesDir
}

func (d *initData) KubeConfigDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return d.kubeconfigDir
}

func (d *initData) KubeConfigPath() string {
	if d.dryRun {
		d.kubeconfigPath = filepath.Join(d.dryRunDir, kubeadmconstants.AdminKubeConfigFileName)
	}
	return d.kubeconfigPath
}

func (d *initData) ManifestDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return kubeadmconstants.GetStaticPodDirectory()
}

func (d *initData) KubeletDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return kubeadmconstants.KubeletRunDirectory
}

func (d *initData) ExternalCA() bool {
	return d.externalCA
}

func (d *initData) OutputWriter() io.Writer {
	return d.outputWriter
}

func (d *initData) Client() (clientset.Interface, error) {

}

func (d *initData) Tokens() []string {
	tokens := []string{}
	for _, bt := range d.cfg.BootstrapTokens {
		tokens = append(tokens, bt.Token.String())
	}
	return tokens
}

func (d *initData) PatchesDir() string {
	return d.patchesDir
}

func printJoinCommand(out io.Writer, adminKubeConfigPath, token string, i *initData) error {
	joinControlPlaneCommand, err := cmdutil.GetJoinControlPlaneCommand(adminKubeConfigPath, token, i.CertificateKey(), i.skipTokenPrint, i.skipCertificateKeyPrint)
	if err != nil {
		return err
	}

	joinWorkerCommand, err := cmdutil.GetJoinWorkerCommand(adminKubeConfigPath, token, i.skipTokenPrint)
	if err != nil {
		return err
	}

	ctx := map[string]interface{}{
		"KubeConfigPath":          adminKubeConfigPath,
		"ControlPlaneEndpoint":    i.Cfg().ControlPlaneEndpoint,
		"UploadCerts":             i.uploadCerts,
		"joinControlPlaneCommand": joinControlPlaneCommand,
		"joinWorkerCommand":       joinWorkerCommand,
	}

	return initDoneTempl.Execute(out, ctx)
}

func showJoinCommand(i *initData, out io.Writer) error {
	adminKubeConfigPath := i.KubeConfigPath()

	for _, token := range i.tokens() {
		if err := printJoinCommand(out, adminKubeConfigPath, token, i); err != nil {
			return error.Wrap(err, "failed to print join command")
		}
	}

	return nil
}
