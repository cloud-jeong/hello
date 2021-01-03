package kubeadm

import (
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
)

type BootstrapTokenString struct {
	ID     string
	Secret string
}

func (bts *BootstrapTokenString) String() string {
	if len(bts.ID) > 0 && len(bts.Secret) > 0 {
		return bootstraputil.TokenFromIDAndSecret(bts.ID, bts.Secret)
	}
	return ""
}
