package features

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/component-base/featuregate"
	"sort"
	"strconv"
	"strings"
)

const (
	// IPv6DualStack is expected to be alpha in v1.16
	IPv6DualStack = "IPv6DualStack"
	// PublicKeysECDSA is expected to be alpha in v1.19
	PublicKeysECDSA = "PublicKeysECDSA"
)

var InitFeatureGates = FeatureList{
	IPv6DualStack:   {FeatureSpec: featuregate.FeatureSpec{Default: false, PreRelease: featuregate.Alpha}},
	PublicKeysECDSA: {FeatureSpec: featuregate.FeatureSpec{Default: false, PreRelease: featuregate.Alpha}},
}

type Feature struct {
	featuregate.FeatureSpec
	MinimumVersion     *version.Version
	HiddenInHelpText   bool
	DeprecationMessage string
}

type FeatureList map[string]Feature

func Enabled(featureList map[string]bool, featureName string) bool {
	if enabled, ok := featureList[string(featureName)]; ok {
		return enabled
	}
	return InitFeatureGates[string(featureName)].Default
}

func Supports(featureList FeatureList, featureName string) bool {
	for k, v := range featureList {
		if featureName == string(k) {
			return v.PreRelease != featuregate.Deprecated
		}
	}
	return false
}

func Keys(featureList FeatureList) []string {
	var list []string
	for k := range featureList {
		list = append(list, string(k))
	}
	return list
}

func KnownFeatures(f *FeatureList) []string {
	var known []string
	for k, v := range *f {
		if v.HiddenInHelpText {
			continue
		}

		pre := ""
		if v.PreRelease != featuregate.GA {
			pre = fmt.Sprintf("%s - ", v.PreRelease)
		}
		known = append(known, fmt.Sprintf("%s=true|false (%sdefault=%t)", k, pre, v.Default))
	}
	sort.Strings(known)
	return known
}

func NewFeatureGate(f *FeatureList, value string) (map[string]bool, error) {
	featureGate := map[string]bool{}
	for _, s := range strings.Split(value, ",") {
		if len(s) == 0 {
			continue
		}

		arr := strings.SplitN(s, "=", 2)
		if len(arr) != 2 {
			return nil, errors.Errorf("missing bool value for feature-gate key:%s", s)
		}

		k := strings.TrimSpace(arr[0])
		v := strings.TrimSpace(arr[1])

		featureSpec, ok := (*f)[k]
		if !ok {
			return nil, errors.Errorf("unrecognized feature-gate key: %s", k)
		}

		if featureSpec.PreRelease == featuregate.Deprecated {
			return nil, errors.Errorf("feature-gate key is deprecated: %s", k)
		}

		boolValue, err := strconv.ParseBool(v)
		if err != nil {
			return nil, errors.Errorf("invalid value %v for feature-gate key: %s, use true|false instead", v, k)
		}
		featureGate[k] = boolValue
	}

	return featureGate, nil
}

func CheckDeprecatedFlags(f *FeatureList, features map[string]bool) map[string]string {
	deprecatedMsg := map[string]string{}
	for k := range features {
		featureSpec, ok := (*f)[k]
		if !ok {
			// This case should never happen, it is implemented only as a sentinel
			// for removal of flags executed when flags are still in use (always before deprecate, then after one cycle remove)
			deprecatedMsg[k] = fmt.Sprintf("Unknown feature gate flag: %s", k)
		}

		if featureSpec.PreRelease == featuregate.Deprecated {
			if _, ok := deprecatedMsg[k]; !ok {
				deprecatedMsg[k] = featureSpec.DeprecationMessage
			}
		}
	}

	return deprecatedMsg
}
