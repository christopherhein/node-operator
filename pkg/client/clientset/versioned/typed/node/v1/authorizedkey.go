/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	v1 "github.com/christopherhein/node-operator/pkg/apis/node.io/v1"
	scheme "github.com/christopherhein/node-operator/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AuthorizedKeysGetter has a method to return a AuthorizedKeyInterface.
// A group's client should implement this interface.
type AuthorizedKeysGetter interface {
	AuthorizedKeys(namespace string) AuthorizedKeyInterface
}

// AuthorizedKeyInterface has methods to work with AuthorizedKey resources.
type AuthorizedKeyInterface interface {
	Create(*v1.AuthorizedKey) (*v1.AuthorizedKey, error)
	Update(*v1.AuthorizedKey) (*v1.AuthorizedKey, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.AuthorizedKey, error)
	List(opts meta_v1.ListOptions) (*v1.AuthorizedKeyList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.AuthorizedKey, err error)
	AuthorizedKeyExpansion
}

// authorizedKeys implements AuthorizedKeyInterface
type authorizedKeys struct {
	client rest.Interface
	ns     string
}

// newAuthorizedKeys returns a AuthorizedKeys
func newAuthorizedKeys(c *NodeV1Client, namespace string) *authorizedKeys {
	return &authorizedKeys{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the authorizedKey, and returns the corresponding authorizedKey object, and an error if there is any.
func (c *authorizedKeys) Get(name string, options meta_v1.GetOptions) (result *v1.AuthorizedKey, err error) {
	result = &v1.AuthorizedKey{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("authorizedkeys").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AuthorizedKeys that match those selectors.
func (c *authorizedKeys) List(opts meta_v1.ListOptions) (result *v1.AuthorizedKeyList, err error) {
	result = &v1.AuthorizedKeyList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("authorizedkeys").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested authorizedKeys.
func (c *authorizedKeys) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("authorizedkeys").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a authorizedKey and creates it.  Returns the server's representation of the authorizedKey, and an error, if there is any.
func (c *authorizedKeys) Create(authorizedKey *v1.AuthorizedKey) (result *v1.AuthorizedKey, err error) {
	result = &v1.AuthorizedKey{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("authorizedkeys").
		Body(authorizedKey).
		Do().
		Into(result)
	return
}

// Update takes the representation of a authorizedKey and updates it. Returns the server's representation of the authorizedKey, and an error, if there is any.
func (c *authorizedKeys) Update(authorizedKey *v1.AuthorizedKey) (result *v1.AuthorizedKey, err error) {
	result = &v1.AuthorizedKey{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("authorizedkeys").
		Name(authorizedKey.Name).
		Body(authorizedKey).
		Do().
		Into(result)
	return
}

// Delete takes name of the authorizedKey and deletes it. Returns an error if one occurs.
func (c *authorizedKeys) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("authorizedkeys").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *authorizedKeys) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("authorizedkeys").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched authorizedKey.
func (c *authorizedKeys) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.AuthorizedKey, err error) {
	result = &v1.AuthorizedKey{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("authorizedkeys").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
