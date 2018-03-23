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

package fake

import (
	node_io_v1 "github.com/christopherhein/node-operator/pkg/apis/node.io/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAuthorizedKeys implements AuthorizedKeyInterface
type FakeAuthorizedKeys struct {
	Fake *FakeNodeV1
	ns   string
}

var authorizedkeysResource = schema.GroupVersionResource{Group: "node.io", Version: "v1", Resource: "authorizedkeys"}

var authorizedkeysKind = schema.GroupVersionKind{Group: "node.io", Version: "v1", Kind: "AuthorizedKey"}

// Get takes name of the authorizedKey, and returns the corresponding authorizedKey object, and an error if there is any.
func (c *FakeAuthorizedKeys) Get(name string, options v1.GetOptions) (result *node_io_v1.AuthorizedKey, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(authorizedkeysResource, c.ns, name), &node_io_v1.AuthorizedKey{})

	if obj == nil {
		return nil, err
	}
	return obj.(*node_io_v1.AuthorizedKey), err
}

// List takes label and field selectors, and returns the list of AuthorizedKeys that match those selectors.
func (c *FakeAuthorizedKeys) List(opts v1.ListOptions) (result *node_io_v1.AuthorizedKeyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(authorizedkeysResource, authorizedkeysKind, c.ns, opts), &node_io_v1.AuthorizedKeyList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &node_io_v1.AuthorizedKeyList{}
	for _, item := range obj.(*node_io_v1.AuthorizedKeyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested authorizedKeys.
func (c *FakeAuthorizedKeys) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(authorizedkeysResource, c.ns, opts))

}

// Create takes the representation of a authorizedKey and creates it.  Returns the server's representation of the authorizedKey, and an error, if there is any.
func (c *FakeAuthorizedKeys) Create(authorizedKey *node_io_v1.AuthorizedKey) (result *node_io_v1.AuthorizedKey, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(authorizedkeysResource, c.ns, authorizedKey), &node_io_v1.AuthorizedKey{})

	if obj == nil {
		return nil, err
	}
	return obj.(*node_io_v1.AuthorizedKey), err
}

// Update takes the representation of a authorizedKey and updates it. Returns the server's representation of the authorizedKey, and an error, if there is any.
func (c *FakeAuthorizedKeys) Update(authorizedKey *node_io_v1.AuthorizedKey) (result *node_io_v1.AuthorizedKey, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(authorizedkeysResource, c.ns, authorizedKey), &node_io_v1.AuthorizedKey{})

	if obj == nil {
		return nil, err
	}
	return obj.(*node_io_v1.AuthorizedKey), err
}

// Delete takes name of the authorizedKey and deletes it. Returns an error if one occurs.
func (c *FakeAuthorizedKeys) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(authorizedkeysResource, c.ns, name), &node_io_v1.AuthorizedKey{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAuthorizedKeys) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(authorizedkeysResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &node_io_v1.AuthorizedKeyList{})
	return err
}

// Patch applies the patch and returns the patched authorizedKey.
func (c *FakeAuthorizedKeys) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *node_io_v1.AuthorizedKey, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(authorizedkeysResource, c.ns, name, data, subresources...), &node_io_v1.AuthorizedKey{})

	if obj == nil {
		return nil, err
	}
	return obj.(*node_io_v1.AuthorizedKey), err
}
