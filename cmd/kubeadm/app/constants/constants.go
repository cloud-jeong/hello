package constants

import (
	v1 "k8s.io/api/core/v1"
	"time"
)

const (
	KubernetesDir               = "/etc/kubernetes"
	ManifestsSubDirName         = "manifests"
	TempDirForKubeadm           = "tmp"
	CertificateValidity         = time.Hour * 24 * 365
	CACertAndKeyBaseName        = "ca"
	CACertName                  = "ca.crt"
	CAKeyName                   = "ca.key"
	APIServerCertAndKeyBaseName = "apiserver"
	APIServerCertName           = "apiserver.crt"
	APIServerKeyName            = "apiserver.key"
	APIServerCertCommonName     = "kube-apiserver"

	LabelNodeRoleOldControlPlane = "node-role.kubernetes.io/master"
	LabelNodeRoleControlPlane    = "node-role.kubernetes.io/control-plane"
	AnnotationKubeadmCRISocket   = "kubeadm.alpha.kubernetes.io/cri-socket"
)

var (
	OldControlPlaneTaint = v1.Taint{
		Key:    LabelNodeRoleOldControlPlane,
		Effect: v1.TaintEffectNoSchedule,
	}
)
