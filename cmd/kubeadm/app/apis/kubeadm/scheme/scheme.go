package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var Scheme = runtime.NewScheme()

var Codecs = serializer.NewCodecFactory(Scheme)

func init() {
	//metav1.AddToGroupVersion(Scheme, scheme.GroupVersion{Version:})
}

func AddToScheme(schme *runtime.Scheme) {

}
