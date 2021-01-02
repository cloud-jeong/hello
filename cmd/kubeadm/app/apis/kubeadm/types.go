package kubeadm

import (
	"crypto/x509"
	"github.com/cloud-jeong/sandbox/cmd/kubeadm/app/features"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InitConfiguration struct {
	metav1.TypeMeta
}

type ClusterConfiguration struct {
	ImageRepository   string
	CIImageRepository string
	FeatureGates      map[string]bool
}

type ControlPlaneComponent struct {
}

type APIServer struct {
}

type DNSAddOnType string

const (
	// CoreDNS add-on type
	CoreDNS DNSAddOnType = "CoreDNS"

	// KubeDNS add-on type
	KubeDNS DNSAddOnType = "kube-dns"
)

type DNS struct {
}

type ImageMeta struct {
}

type ClusterStatus struct {
}

type APIEndpoint struct {
}

type NodeRegistrationOptions struct {
}

type Networking struct {
}

type BootstrapToken struct {
}

type Etcd struct {
}

type ExternalEtcd struct {
}

type JoinConfiguration struct {
}

type Discovery struct {
}

type BootstrapTokenDiscovery struct {
}

type FileDiscovery struct {
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
