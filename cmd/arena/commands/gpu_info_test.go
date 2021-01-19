package commands

import (
	"context"
	"fmt"
	"github.com/kubeflow/arena/pkg/operators/et-operator/api/v1alpha1"
	scheme2 "github.com/kubeflow/arena/pkg/operators/et-operator/client/clientset/versioned/scheme"
	tfv1 "github.com/kubeflow/arena/pkg/operators/tf-operator/apis/tensorflow/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
	"testing"

	"github.com/kubeflow/arena/pkg/util"

	"sigs.k8s.io/controller-runtime/pkg/cache"
)

func TestQueryMetricByPrometheus(t *testing.T) {
	clientset := util.GetClientSetForTest(t)
	if clientset == nil {
		t.Skip("kubeclient not setup")
	}
	gpuMetrics, _ := QueryMetricByPrometheus(clientset, "prometheus-svc", KUBEFLOW_NAMESPACE, fmt.Sprintf(`{__name__=~"%s"}`, strings.Join(GPU_METRIC_LIST, "|")))

	for _, m := range gpuMetrics {
		t.Logf("metric %++v", m)
		t.Logf("metric name %s, value: %s", m.MetricName, m.Value)
	}
}


func TestCacheClient(t *testing.T) {

	// init config
	func () {
		loadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
		overrides := clientcmd.ConfigOverrides{}
		clientConfig = clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &overrides, os.Stdin)
		restConfig, _ = clientConfig.ClientConfig()
	}()


	// registry crd scheme
	func () {
		scheme2.AddToScheme(scheme.Scheme)
		tfv1.AddToScheme(scheme.Scheme)
	}()

	ctx := context.Background()

	cacheClient, err := cache.New(restConfig, cache.Options{Resync: nil, Namespace: ""})

	go cacheClient.Start(ctx)

	cacheClient.WaitForCacheSync(ctx)

	trainingJobs := &v1alpha1.TrainingJobList{}

	err = cacheClient.List(ctx, trainingJobs)
	if err != nil {
		t.Errorf("failed to list trainingJobs")
	}else {
		t.Log(len(trainingJobs.Items))
	}

	tfJobs := &tfv1.TFJobList{}

	err = cacheClient.List(ctx, tfJobs)
	if err != nil {
		t.Errorf("failed to list tfjobs")
	}else {
		t.Log(len(tfJobs.Items))
	}


	pods := &v1.PodList{}

	err = cacheClient.List(ctx, pods)
	if err != nil {
		t.Errorf("failed to list pods")
	}else {
		t.Log(len(pods.Items))
	}
}
