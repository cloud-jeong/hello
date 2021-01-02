package componentconfigs

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog"
	"sort"
	"strings"
)

type UnsupportedConfigVersionError struct {
	OldVersion     schema.GroupVersion
	CurrentVersion schema.GroupVersion
	Document       []byte
}

func (err *UnsupportedConfigVersionError) Error() string {
	return fmt.Sprintf("unsupported apiVersion %q, you may have to do manual conversion to %q and run kubeadm again", err.OldVersion, err.CurrentVersion)
}

type UnsupportedConfigVersionsErrorMap map[string]*UnsupportedConfigVersionError

func (errs UnsupportedConfigVersionsErrorMap) Error() string {
	groups := make([]string, 0, len(errs))
	for group := range errs {
		groups = append(groups, group)
	}
	sort.Strings(groups)

	msgs := make([]string, 1, 1+len(errs))
	msgs[0] = "multiple unsupported config version errors encountered:"
	for _, group := range groups {
		msgs = append(msgs, errs[group].Error())
	}

	return strings.Join(msgs, "\n\t- ")
}

func warnDefaultComponentConfigValue(componentConfigKind, paramName string, defaultValue, userValue interface{}) {
	klog.Warningf("The recommanded value for %q in %q is: %v; the provided value is %v",
		paramName, componentConfigKind, defaultValue, userValue)
}
