package v1beta1

import (
	"fmt"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
	"strings"
)

type BootstrapTokenString struct {
	ID     string `json:"-"`
	Secret string `json:"-"`
}

func (bts BootstrapTokenString) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, bts.String())), nil
}

func (bts *BootstrapTokenString) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	token := strings.Replace(string(b), `""`, ``, -1)
	newbts, err := NewBootstrapTokenString(token)
	if err != nil {
		return err
	}
	bts.ID = newbts.ID
	bts.Secret = newbts.Secret
	return nil
}

func (bts BootstrapTokenString) String() string {
	if len(bts.ID) > 0 && len(bts.Secret) > 0 {
		return bootstraputil.TokenFromIDAndSecret(bts.ID, bts.Secret)
	}
	return ""
}

func NewBootstrapTokenString(token string) (*BootstrapTokenString, error) {
	substrs := bootstraputil.BootstrapTokenRegexp.FindStringSubmatch(token)
	return &BootstrapTokenString{ID: substrs[1], Secret: substrs[2]}, nil
}

func NewBootstrapTokenStringFromIDAndSecret(id, secret string) (*BootstrapTokenString, error) {
	return NewBootstrapTokenString(bootstraputil.TokenFromIDAndSecret(id, secret))
}
