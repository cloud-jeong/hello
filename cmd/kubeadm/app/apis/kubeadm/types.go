package kubeadm

import (
	"crypto/x509"
	"github.com/cloud-jeong/sandbox/cmd/kubeadm/app/features"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type InitConfiguration struct {
	metav1.TypeMeta
	ClusterConfiguration `json:"-"`
	BootstrapTokens      []BootstrapToken
	NodeRegistration     NodeRegistrationOptions
	LocalAPIEndpoint     APIEndpoint
	CertificateKey       string
}

type ClusterConfiguration struct {
	metav1.TypeMeta
	ComponentConfigs     ComponentConfigMap
	Etcd                 Etcd
	Networking           Networking
	KubernetesVersion    string
	ControlPlaneEndpoint string
	APIServer            APIServer
	ControllerManager    ControlPlaneComponent
	Scheduler            ControlPlaneComponent
	DNS                  DNS
	CertificatesDir      string
	ImageRepository      string
	CIImageRepository    string
	UseHyperKubeImage    bool
	FeatureGates         map[string]bool
	ClusterName          string
}

type ControlPlaneComponent struct {
	ExtraArgs    map[string]string
	ExtraVolumes []HostPathMount
}

type APIServer struct {
	ControlPlaneComponent
	CertSANs               []string
	TimeoutForControlPlane *metav1.Duration
}

type DNSAddOnType string

const (
	// CoreDNS add-on type
	CoreDNS DNSAddOnType = "CoreDNS"

	// KubeDNS add-on type
	KubeDNS DNSAddOnType = "kube-dns"
)

type DNS struct {
	Type      DNSAddOnType
	ImageMeta `json:",inline"`
}

type ImageMeta struct {
	ImageRepository string
	ImageTag        string
}

type ClusterStatus struct {
	metav1.TypeMeta
	APIEndpoints map[string]APIEndpoint
}

type APIEndpoint struct {
	AdvertiseAddress string
	BindPort         int32
}

type NodeRegistrationOptions struct {
	Name                  string
	CRISocket             string
	Taints                []v1.Taint
	KubeletExtraArgs      map[string]string
	IgnorePreflightErrors []string
}

type Networking struct {
	ServiceSubnet string
	PodSubnet     string
	DNSDomain     string
}

type BootstrapToken struct {
	Token       *BootstrapTokenString
	Description string
	TTL         *metav1.Duration
	Expires     *metav1.Time
	Updates     []string
	Groups      []string
}

type Etcd struct {
	Local    *LocalEtcd
	External *ExternalEtcd
}

type LocalEtcd struct {
	ImageMeta       `json:",inline"`
	DataDir         string
	ExtraArgs       map[string]string
	ServiceCertSANs []string
	PeerCertSANs    []string
}

type ExternalEtcd struct {
	Endpoints []string
	CAFile    string
	CertFile  string
	KeyFile   string
}

type JoinConfiguration struct {
	metav1.TypeMeta
	NodeRegistration NodeRegistrationOptions
	CACertPath       string
	Discovery        Discovery
	ControlPlane     *JoinControlPlane
}

type JoinControlPlane struct {
	LocalAPIEndpoint APIEndpoint
	CertificateKey   string
}

type Discovery struct {
	BootstrapToken    *BootstrapTokenDiscovery
	File              *FileDiscovery
	TLSBootstrapToken string
	Timeout           *metav1.Duration
}

type BootstrapTokenDiscovery struct {
	Token                   string
	APIServerEndpoint       string
	CACertHashes            []string
	UnsafeSkipCAVerfication bool
}

type FileDiscovery struct {
	KubeConfigPath string
}

func (cfg *ClusterConfiguration) GetControlPlaneImageRepository() string {
	if cfg.CIImageRepository != "" {
		return cfg.CIImageRepository
	}

	return cfg.ImageRepository
}

func (cfg *ClusterConfiguration) PublicKeyAlgorithm() x509.PublicKeyAlgorithm {
	if features.Enabled(cfg.FeatureGates, features.PublicKeysECDSA) {
		return x509.ECDSA
	}

	return x509.RSA
}

type HostPathMount struct {
	Name      string
	HostPath  string
	MountPath string
	ReadOnly  bool
	PathType  v1.HostPathType
}

type DocumentMap map[schema.GroupVersionKind][]byte

type ComponentConfig interface {
	DeepCopy() ComponentConfig
	Marshal() ([]byte, error)
	Unmarshal(docmap DocumentMap) error
	Default(cfg *ClusterConfiguration, localAPIEndpoint *APIEndpoint, nodeRegOpts *NodeRegistrationOptions)
	IsUserSupplied() bool
	SetUserSupplied(userSupplied bool)
}

type ComponentConfigMap map[string]ComponentConfig
