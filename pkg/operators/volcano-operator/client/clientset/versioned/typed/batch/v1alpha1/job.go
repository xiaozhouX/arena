/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha1

import (
	"context"
	v1alpha1 "github.com/kubeflow/arena/pkg/operators/volcano-operator/apis/batch/v1alpha1"
	scheme "github.com/kubeflow/arena/pkg/operators/volcano-operator/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// JobsGetter has a method to return a JobInterface.
// A group's client should implement this interface.
type JobsGetter interface {
	Jobs(namespace string) JobInterface
}

// JobInterface has methods to work with Job resources.
type JobInterface interface {
	Create(*v1alpha1.Job) (*v1alpha1.Job, error)
	Update(*v1alpha1.Job) (*v1alpha1.Job, error)
	UpdateStatus(*v1alpha1.Job) (*v1alpha1.Job, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Job, error)
	List(opts v1.ListOptions) (*v1alpha1.JobList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Job, err error)
	JobExpansion
}

// jobs implements JobInterface
type jobs struct {
	client rest.Interface
	ns     string
}

// newJobs returns a Jobs
func newJobs(c *BatchV1alpha1Client, namespace string) *jobs {
	return &jobs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the job, and returns the corresponding job object, and an error if there is any.
func (c *jobs) Get(name string, options v1.GetOptions) (result *v1alpha1.Job, err error) {
	result = &v1alpha1.Job{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("jobs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(context.Background()).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Jobs that match those selectors.
func (c *jobs) List(opts v1.ListOptions) (result *v1alpha1.JobList, err error) {
	result = &v1alpha1.JobList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("jobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested jobs.
func (c *jobs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("jobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}

// Create takes the representation of a job and creates it.  Returns the server's representation of the job, and an error, if there is any.
func (c *jobs) Create(job *v1alpha1.Job) (result *v1alpha1.Job, err error) {
	result = &v1alpha1.Job{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("jobs").
		Body(job).
		Do(context.Background()).
		Into(result)
	return
}

// Update takes the representation of a job and updates it. Returns the server's representation of the job, and an error, if there is any.
func (c *jobs) Update(job *v1alpha1.Job) (result *v1alpha1.Job, err error) {
	result = &v1alpha1.Job{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("jobs").
		Name(job.Name).
		Body(job).
		Do(context.Background()).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *jobs) UpdateStatus(job *v1alpha1.Job) (result *v1alpha1.Job, err error) {
	result = &v1alpha1.Job{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("jobs").
		Name(job.Name).
		SubResource("status").
		Body(job).
		Do(context.Background()).
		Into(result)
	return
}

// Delete takes name of the job and deletes it. Returns an error if one occurs.
func (c *jobs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("jobs").
		Name(name).
		Body(options).
		Do(context.Background()).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *jobs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("jobs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do(context.Background()).
		Error()
}

// Patch applies the patch and returns the patched job.
func (c *jobs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Job, err error) {
	result = &v1alpha1.Job{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("jobs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do(context.Background()).
		Into(result)
	return
}
