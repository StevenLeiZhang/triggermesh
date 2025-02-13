/*
Copyright 2022 TriggerMesh Inc.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/triggermesh/triggermesh/pkg/apis/flow/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeJQTransformations implements JQTransformationInterface
type FakeJQTransformations struct {
	Fake *FakeFlowV1alpha1
	ns   string
}

var jqtransformationsResource = schema.GroupVersionResource{Group: "flow.triggermesh.io", Version: "v1alpha1", Resource: "jqtransformations"}

var jqtransformationsKind = schema.GroupVersionKind{Group: "flow.triggermesh.io", Version: "v1alpha1", Kind: "JQTransformation"}

// Get takes name of the jQTransformation, and returns the corresponding jQTransformation object, and an error if there is any.
func (c *FakeJQTransformations) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.JQTransformation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(jqtransformationsResource, c.ns, name), &v1alpha1.JQTransformation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JQTransformation), err
}

// List takes label and field selectors, and returns the list of JQTransformations that match those selectors.
func (c *FakeJQTransformations) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.JQTransformationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(jqtransformationsResource, jqtransformationsKind, c.ns, opts), &v1alpha1.JQTransformationList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.JQTransformationList{ListMeta: obj.(*v1alpha1.JQTransformationList).ListMeta}
	for _, item := range obj.(*v1alpha1.JQTransformationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested jQTransformations.
func (c *FakeJQTransformations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(jqtransformationsResource, c.ns, opts))

}

// Create takes the representation of a jQTransformation and creates it.  Returns the server's representation of the jQTransformation, and an error, if there is any.
func (c *FakeJQTransformations) Create(ctx context.Context, jQTransformation *v1alpha1.JQTransformation, opts v1.CreateOptions) (result *v1alpha1.JQTransformation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(jqtransformationsResource, c.ns, jQTransformation), &v1alpha1.JQTransformation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JQTransformation), err
}

// Update takes the representation of a jQTransformation and updates it. Returns the server's representation of the jQTransformation, and an error, if there is any.
func (c *FakeJQTransformations) Update(ctx context.Context, jQTransformation *v1alpha1.JQTransformation, opts v1.UpdateOptions) (result *v1alpha1.JQTransformation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(jqtransformationsResource, c.ns, jQTransformation), &v1alpha1.JQTransformation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JQTransformation), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeJQTransformations) UpdateStatus(ctx context.Context, jQTransformation *v1alpha1.JQTransformation, opts v1.UpdateOptions) (*v1alpha1.JQTransformation, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(jqtransformationsResource, "status", c.ns, jQTransformation), &v1alpha1.JQTransformation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JQTransformation), err
}

// Delete takes name of the jQTransformation and deletes it. Returns an error if one occurs.
func (c *FakeJQTransformations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(jqtransformationsResource, c.ns, name), &v1alpha1.JQTransformation{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeJQTransformations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(jqtransformationsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.JQTransformationList{})
	return err
}

// Patch applies the patch and returns the patched jQTransformation.
func (c *FakeJQTransformations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.JQTransformation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(jqtransformationsResource, c.ns, name, pt, data, subresources...), &v1alpha1.JQTransformation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JQTransformation), err
}
