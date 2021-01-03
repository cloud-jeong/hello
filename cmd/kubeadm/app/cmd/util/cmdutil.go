package util

import (
	"github.com/cloud-jeong/sandbox/cmd/kubeadm/app/cmd/options"
	kubeadmconstants "github.com/cloud-jeong/sandbox/cmd/kubeadm/app/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func SubCmdRunE(name string) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.Errorf("missing subcommand: %q is not meant to be run on its own", name)
		}

		return errors.Errorf("invalid subcommand: %q", args[0])
	}
}

func ValidateExactArgNumber(args []string, supportedArgs []string) error {
	lenSupported := len(supportedArgs)
	validArgs := 0

	for _, arg := range args {
		if len(arg) > 0 {
			validArgs++
		}
		if validArgs > lenSupported {
			return errors.Errorf("too many arguments. Required arguments: %v", supportedArgs)
		}
	}

	if validArgs < lenSupported {
		return errors.Errorf("missing one or more required arguments. Required arguments: %v", supportedArgs)
	}
	return nil
}

func GetkubeConfigPath(file string) string {
	if file != "" {
		return file
	}

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.Precedence = append(rules.Precedence, kubeadmconstants.GetAdminKubeConfigPath())
	file = rules.GetDefaultFilename()
	klog.V(1).Infof("Using kubeconfig file: %s", file)
	return file
}

func AddCRISocketFlag(flagSet *pflag.FlagSet, criSocket *string) {
	flagSet.StringVar(
		criSocket, options.NodeCRISocket, *criSocket,
		"path to the CRI socket to connect. If empty kubeadm will try to auto-detect this value; use this option only if you have more than one CRI installed or if you have non-standard CRI socket.",
	)
}
