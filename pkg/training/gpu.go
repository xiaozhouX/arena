// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package training

import (
	"strconv"

	"github.com/kubeflow/arena/pkg/apis/types"
	"github.com/kubeflow/arena/pkg/apis/utils"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// filter out the pods with GPU
func gpuPods(pods []v1.Pod) (podsWithGPU []v1.Pod) {
	for _, pod := range pods {
		if gpuInPod(pod) > 0 {
			podsWithGPU = append(podsWithGPU, pod)
		}
	}
	return podsWithGPU
}

// The way to get total GPU Count of Node: nvidia.com/gpu
func totalGpuInNode(node v1.Node) int64 {
	val, ok := node.Status.Capacity[NVIDIAGPUResourceName]

	if !ok {
		return gpuInNodeDeprecated(node)
	}

	return val.Value()
}

// The way to get allocatble GPU Count of Node: nvidia.com/gpu
func allocatableGpuInNode(node v1.Node) int64 {
	val, ok := node.Status.Allocatable[NVIDIAGPUResourceName]

	if !ok {
		return gpuInNodeDeprecated(node)
	}

	return val.Value()
}

// The way to get GPU Count of Node: alpha.kubernetes.io/nvidia-gpu
func gpuInNodeDeprecated(node v1.Node) int64 {
	val, ok := node.Status.Allocatable[DeprecatedNVIDIAGPUResourceName]

	if !ok {
		return 0
	}

	return val.Value()
}

func gpuInPod(pod v1.Pod) (gpuCount int64) {
	return gpuInPodSpec(pod.Spec)
}

func gpuInPodSpec(spec v1.PodSpec) (gpuCount int64) {
	containers := spec.Containers
	for _, container := range containers {
		gpuCount += gpuInContainer(container)
	}

	return gpuCount
}

func getRequestGPUsOfJobFromPodAnnotation(pods []*v1.Pod) int64 {
	for _, pod := range pods {
		val := pod.Annotations[types.RequestGPUsOfJobAnnoKey]
		if val == "" {
			continue
		}
		gpuCount, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Debugf("failed to get requests gpus of job from pod annotations,%v", err)
			continue
		}
		return gpuCount

	}
	return 0
}

// Get gpu number from the active pod
func gpuInActivePod(pod v1.Pod) (gpuCount int64) {
	if pod.Status.StartTime == nil {
		return 0
	}

	if utils.IsCompletedPod(&pod) {
		return 0
	}

	containers := pod.Spec.Containers
	for _, container := range containers {
		gpuCount += gpuInContainer(container)
	}

	return gpuCount
}

func gpuInContainer(container v1.Container) int64 {
	val, ok := container.Resources.Limits[NVIDIAGPUResourceName]

	if ok {
		return val.Value()
	}

	val, ok = container.Resources.Limits[AliyunGPUResourceName]

	if ok {
		return val.Value()
	}

	return gpuInContainerDeprecated(container)
}

func gpuInContainerDeprecated(container v1.Container) int64 {
	val, ok := container.Resources.Limits[DeprecatedNVIDIAGPUResourceName]

	if !ok {
		return 0
	}

	return val.Value()
}
