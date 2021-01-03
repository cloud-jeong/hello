package util

import (
	"bytes"
	"crypto/x509"
	kubeconfigutil "github.com/cloud-jeong/sandbox/cmd/kubeadm/app/util/kubeconfig"
	"github.com/cloud-jeong/sandbox/cmd/kubeadm/app/util/pubkeypin"
	"github.com/pkg/errors"
	"html/template"
	"k8s.io/client-go/tools/clientcmd"
	clientcertutil "k8s.io/client-go/util/cert"
	"strings"
)

var joinCommandTemplate = template.Must(template.New("join").Parse(`` +
	`kubeadm join {{.ControlPlaneHostPort}} -- token {{.Token}} \
	{{range $h := .CAPubKeyPins}}--discovery-token-ca-cert-hash {{$h}} {{end}}{if .ControlPlane}}\
	--control-plane {{if .CertificateKey}}--certificate-key {{.CertificateKey}}{{end}}{{end}}`,
))

func GetJoinWorkerCommand(kubeConfigFile, token string, skipTokenPrint bool) (string, error) {
	return getJoinCommand(kubeConfigFile, token, "", false, skipTokenPrint, false)
}

func GetJoinControlPlaneCommand(kubeConfigFile, token, key string, skipTokenPrint, skipCertificateKeyPrint bool) (string, error) {
	return getJoinCommand(kubeConfigFile, token, key, true, skipTokenPrint, skipCertificateKeyPrint)
}

func getJoinCommand(kubeConfigFile, token, key string, controlPlane, skipTokenPrint, skipCertificateKeyPrint bool) (string, error) {
	config, err := clientcmd.LoadFromFile(kubeConfigFile)
	if err != nil {
		return "", errors.Wrap(err, "failed to load kubeconfig")
	}

	clusterConfig := kubeconfigutil.GetClusterFromKubeConfig(config)
	if clusterConfig == nil {
		return "", errors.New("failed to get default cluster config")
	}

	var caCerts []*x509.Certificate
	if clusterConfig.CertificateAuthorityData != nil {
		caCerts, err = clientcertutil.ParseCertsPEM(clusterConfig.CertificateAuthorityData)
		if err != nil {
			return "", errors.Wrap(err, "failed to parse CA certificate from kubeconfig")
		}
	} else if clusterConfig.CertificateAuthority != "" {
		caCerts, err = clientcertutil.CertsFromFile(clusterConfig.CertificateAuthority)
		if err != nil {
			return "", errors.Wrap(err, "failed to load CA certificate referenced by kubeconfig")
		}
	} else {
		return "", errors.Wrap(err, "no CA certificates found in kubeconfig")
	}

	publicKeyPins := make([]string, 0, len(caCerts))
	for _, caCert := range caCerts {
		publicKeyPins = append(publicKeyPins, pubkeypin.Hash(caCert))
	}

	ctx := map[string]interface{}{
		"Token":                token,
		"CAPubKeyPins":         publicKeyPins,
		"ControlPlaneHostPort": strings.Replace(clusterConfig.Server, "https://", "", -1),
		"CertificateKey":       key,
		"ControlPlane":         controlPlane,
	}

	if skipTokenPrint {
		ctx["Token"] = template.HTML("<value withheld>")
	}
	if skipCertificateKeyPrint {
		ctx["CertificateKey"] = template.HTML("<value withheld>")
	}

	var out bytes.Buffer
	err = joinCommandTemplate.Execute(&out, ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to render join command template")
	}
	return out.String(), nil
}
