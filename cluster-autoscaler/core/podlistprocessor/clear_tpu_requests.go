/*
Copyright 2023 The Kubernetes Authors.

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

package podlistprocessor

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"k8s.io/autoscaler/cluster-autoscaler/utils/tpu"
)

type clearTPURequests struct {
}

// NewFilterOutDaemonSetPodListProcessor creates a PodListProcessor filtering out daemon set pods
func NewClearTPURequestPodListProcessor() *clearTPURequests {
	return &clearTPURequests{}
}

// Process filters out pods which are daemon set pods.
func (p *clearTPURequests) Process(context *context.AutoscalingContext, unschedulablePods []*apiv1.Pod) ([]*apiv1.Pod, error) {
	return tpu.ClearTPURequests(unschedulablePods), nil
}

func (p *clearTPURequests) CleanUp() {
}
