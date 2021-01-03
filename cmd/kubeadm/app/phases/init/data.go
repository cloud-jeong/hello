package init

import (
	kubeadmapi "github.com/cloud-jeong/sandbox/cmd/kubeadm/app/apis/kubeadm"
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
)

type InitData interface {
	UploadCerts() bool
	CertificateKey() string
	SetCertificateKey(key string)
	SkipCertificateKeyPrint() bool
	Cfg() *kubeadmapi.InitConfiguration
	DryRun() bool
	SkipTokenPrint() bool
	IgnorePreflightErrors() sets.String
	CertificateWriteDir() string
	CertificateDir() string
	KubeConfigDir() string
	KubeConfigPath() string
	ManifestDir() string
	KubeletDir() string
	ExternalCA() bool
	OutputWriter() io.Writer
	Client() (clientset.Interface, error)
	Tokens() []string
	PatchesDir() string
}
