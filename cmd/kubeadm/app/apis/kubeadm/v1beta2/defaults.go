package v1beta2

import (
	"k8s.io/apimachinery/pkg/runtime"
	"net/url"
	"time"
)

const (
	DefaultServiceDNSDomain  = "cluster.local"
	DefaultServiceSubnet     = "10.96.0.0/12"
	DefaultClusterDNSIP      = "10.96.0.10"
	DefaultKubernetesVersion = "stable-1"
	DefaultAPIBindPort       = 6443
	DefaultCertificatesDir   = "/etc/kubernetes/pki"
	DefaultImageRepository   = "k8s.gcr.io"
	DefaultManifestsDir      = "/etc/kubernetes/manifests"
	DefaultClusterName       = "kubernetes"

	DefaultEtcdDataDir        = "/var/lib/etcd"
	DefaultProxyBindAddressv4 = "0.0.0.0"
	DefaultProxyBindAddressv6 = "::"
	DefaultDiscoveryTimeout   = 5 * time.Minute
)

var (
	DefaultAuditPolicyLogMaxAge = int32(2)
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(schme)
}

// SetDefaults_InitConfiguration assigns default values for the InitConfiguration
func SetDefaults_InitConfiguration(obj *InitConfiguration) {
	SetDefaults_BootstrapTokens(obj)
	SetDefaults_APIEndpoint(&obj.LocalAPIEndpoint)
}

// SetDefaults_ClusterConfiguration assigns default values for the ClusterConfiguration
func SetDefaults_ClusterConfiguration(obj *ClusterConfiguration) {
	if obj.KubernetesVersion == "" {
		obj.KubernetesVersion = DefaultKubernetesVersion
	}

	if obj.Networking.ServiceSubnet == "" {
		obj.Networking.ServiceSubnet = DefaultServicesSubnet
	}

	if obj.Networking.DNSDomain == "" {
		obj.Networking.DNSDomain = DefaultServiceDNSDomain
	}

	if obj.CertificatesDir == "" {
		obj.CertificatesDir = DefaultCertificatesDir
	}

	if obj.ImageRepository == "" {
		obj.ImageRepository = DefaultImageRepository
	}

	if obj.ClusterName == "" {
		obj.ClusterName = DefaultClusterName
	}

	SetDefaults_DNS(obj)
	SetDefaults_Etcd(obj)
	SetDefaults_APIServer(&obj.APIServer)
}

// SetDefaults_APIServer assigns default values for the API Server
func SetDefaults_APIServer(obj *APIServer) {
	if obj.TimeoutForControlPlane == nil {
		obj.TimeoutForControlPlane = &metav1.Duration{
			Duration: constants.DefaultControlPlaneTimeout,
		}
	}
}

// SetDefaults_DNS assigns default values for the DNS component
func SetDefaults_DNS(obj *ClusterConfiguration) {
	if obj.DNS.Type == "" {
		obj.DNS.Type = CoreDNS
	}
}

// SetDefaults_Etcd assigns default values for the proxy
func SetDefaults_Etcd(obj *ClusterConfiguration) {
	if obj.Etcd.External == nil && obj.Etcd.Local == nil {
		obj.Etcd.Local = &LocalEtcd{}
	}
	if obj.Etcd.Local != nil {
		if obj.Etcd.Local.DataDir == "" {
			obj.Etcd.Local.DataDir = DefaultEtcdDataDir
		}
	}
}

// SetDefaults_JoinConfiguration assigns default values to a regular node
func SetDefaults_JoinConfiguration(obj *JoinConfiguration) {
	if obj.CACertPath == "" {
		obj.CACertPath = DefaultCACertPath
	}

	SetDefaults_JoinControlPlane(obj.ControlPlane)
	SetDefaults_Discovery(&obj.Discovery)
}

func SetDefaults_JoinControlPlane(obj *JoinControlPlane) {
	if obj != nil {
		SetDefaults_APIEndpoint(&obj.LocalAPIEndpoint)
	}
}

// SetDefaults_Discovery assigns default values for the discovery process
func SetDefaults_Discovery(obj *Discovery) {
	if len(obj.TLSBootstrapToken) == 0 && obj.BootstrapToken != nil {
		obj.TLSBootstrapToken = obj.BootstrapToken.Token
	}

	if obj.Timeout == nil {
		obj.Timeout = &metav1.Duration{
			Duration: DefaultDiscoveryTimeout,
		}
	}

	if obj.File != nil {
		SetDefaults_FileDiscovery(obj.File)
	}
}

// SetDefaults_FileDiscovery assigns default values for file based discovery
func SetDefaults_FileDiscovery(obj *FileDiscovery) {
	// Make sure file URL becomes path
	if len(obj.KubeConfigPath) != 0 {
		u, err := url.Parse(obj.KubeConfigPath)
		if err == nil && u.Scheme == "file" {
			obj.KubeConfigPath = u.Path
		}
	}
}

func SetDefaults_BootstrapTokens(obj *InitConfiguration) {

	if obj.BootstrapTokens == nil || len(obj.BootstrapTokens) == 0 {
		obj.BootstrapTokens = []BootstrapToken{{}}
	}

	for i := range obj.BootstrapTokens {
		SetDefaults_BootstrapToken(&obj.BootstrapTokens[i])
	}
}

// SetDefaults_BootstrapToken sets the defaults for an individual Bootstrap Token
func SetDefaults_BootstrapToken(bt *BootstrapToken) {
	if bt.TTL == nil {
		bt.TTL = &metav1.Duration{
			Duration: constants.DefaultTokenDuration,
		}
	}
	if len(bt.Usages) == 0 {
		bt.Usages = constants.DefaultTokenUsages
	}

	if len(bt.Groups) == 0 {
		bt.Groups = constants.DefaultTokenGroups
	}
}

// SetDefaults_APIEndpoint sets the defaults for the API server instance deployed on a node.
func SetDefaults_APIEndpoint(obj *APIEndpoint) {
	if obj.BindPort == 0 {
		obj.BindPort = DefaultAPIBindPort
	}
}
