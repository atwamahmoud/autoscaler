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
	core_utils "k8s.io/autoscaler/cluster-autoscaler/core/utils"
	caerrors "k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	klog "k8s.io/klog/v2"
)

type filterOutExpandable struct {
}

// NewFilterOutExpandablePodListProcessor creates a PodListProcessor filtering out expendable pods
func NewFilterOutExpandablePodListProcessor() *filterOutExpandable {
	return &filterOutExpandable{}
}

// Process filters out pods which are expendable and adds pods which is waiting for lower priority pods preemption to the cluster snapshot
func (p *filterOutExpandable) Process(context *context.AutoscalingContext, pods []*apiv1.Pod) ([]*apiv1.Pod, error) {
	klog.V(4).Infof("Filtering out expandable pods")
	nodes, err := context.AllNodeLister().List()
	if err != nil {
		return nil, err
	}
	expendablePodsPriorityCutoff := context.AutoscalingOptions.ExpendablePodsPriorityCutoff

	unschedulablePods, waitingForLowerPriorityPreemption := core_utils.FilterOutExpendableAndSplit(pods, nodes, expendablePodsPriorityCutoff)
	if err = p.addPreemptiblePodsToSnapshot(waitingForLowerPriorityPreemption, context); err != nil {
		return nil, err
	}

	return unschedulablePods, nil
}

// addPreemptiblePodsToSnapshot modifies the snapshot simulating scheduling of pods waiting for preemption.
// this is not strictly correct as we are not simulating preemption itself but it matches
// CA logic from before migration to scheduler framework. So let's keep it for now
func (p *filterOutExpandable) addPreemptiblePodsToSnapshot(pods []*apiv1.Pod, ctx *context.AutoscalingContext) error {
	for _, p := range pods {
		if err := ctx.ClusterSnapshot.AddPod(p, p.Status.NominatedNodeName); err != nil {
			klog.Errorf("Failed to update snapshot with pod %s waiting for preemption", err)
			return caerrors.ToAutoscalerError(caerrors.InternalError, err)
		}
	}
	return nil
}

func (p *filterOutExpandable) CleanUp() {
}
