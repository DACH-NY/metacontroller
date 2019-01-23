/*
Copyright The Kubernetes Authors.

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

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "metacontroller.app/apis/metacontroller/v1alpha1"
	scheme "metacontroller.app/client/generated/clientset/internalclientset/scheme"
)

// CompositeControllersGetter has a method to return a CompositeControllerInterface.
// A group's client should implement this interface.
type CompositeControllersGetter interface {
	CompositeControllers(namespace string) CompositeControllerInterface
}

// CompositeControllerInterface has methods to work with CompositeController resources.
type CompositeControllerInterface interface {
	Create(*v1alpha1.CompositeController) (*v1alpha1.CompositeController, error)
	Update(*v1alpha1.CompositeController) (*v1alpha1.CompositeController, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.CompositeController, error)
	List(opts v1.ListOptions) (*v1alpha1.CompositeControllerList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CompositeController, err error)
	CompositeControllerExpansion
}

// compositeControllers implements CompositeControllerInterface
type compositeControllers struct {
	client rest.Interface
	ns     string
}

// newCompositeControllers returns a CompositeControllers
func newCompositeControllers(c *MetacontrollerV1alpha1Client, namespace string) *compositeControllers {
	return &compositeControllers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the compositeController, and returns the corresponding compositeController object, and an error if there is any.
func (c *compositeControllers) Get(name string, options v1.GetOptions) (result *v1alpha1.CompositeController, err error) {
	result = &v1alpha1.CompositeController{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("compositecontrollers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CompositeControllers that match those selectors.
func (c *compositeControllers) List(opts v1.ListOptions) (result *v1alpha1.CompositeControllerList, err error) {
	result = &v1alpha1.CompositeControllerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("compositecontrollers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested compositeControllers.
func (c *compositeControllers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("compositecontrollers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a compositeController and creates it.  Returns the server's representation of the compositeController, and an error, if there is any.
func (c *compositeControllers) Create(compositeController *v1alpha1.CompositeController) (result *v1alpha1.CompositeController, err error) {
	result = &v1alpha1.CompositeController{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("compositecontrollers").
		Body(compositeController).
		Do().
		Into(result)
	return
}

// Update takes the representation of a compositeController and updates it. Returns the server's representation of the compositeController, and an error, if there is any.
func (c *compositeControllers) Update(compositeController *v1alpha1.CompositeController) (result *v1alpha1.CompositeController, err error) {
	result = &v1alpha1.CompositeController{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("compositecontrollers").
		Name(compositeController.Name).
		Body(compositeController).
		Do().
		Into(result)
	return
}

// Delete takes name of the compositeController and deletes it. Returns an error if one occurs.
func (c *compositeControllers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("compositecontrollers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *compositeControllers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("compositecontrollers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched compositeController.
func (c *compositeControllers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CompositeController, err error) {
	result = &v1alpha1.CompositeController{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("compositecontrollers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
