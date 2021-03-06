// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/christopherhein/node-operator/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// AuthorizedKeys returns a AuthorizedKeyInformer.
	AuthorizedKeys() AuthorizedKeyInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// AuthorizedKeys returns a AuthorizedKeyInformer.
func (v *version) AuthorizedKeys() AuthorizedKeyInformer {
	return &authorizedKeyInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
