/*
 Copyright 2022 The KubeSphere Authors.

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
package precheck

import (
	"github.com/kubesphere/kubekey/pkg/bootstrap/precheck"
	"github.com/kubesphere/kubekey/pkg/common"
	"github.com/kubesphere/kubekey/pkg/core/prepare"
	"github.com/kubesphere/kubekey/pkg/core/task"
)

type UprgadePreCheckModule struct {
	common.KubeModule
}

func (c *UprgadePreCheckModule) Init() {
	c.Name = "UprgadePreCheckModule"
	c.Desc = "Do pre-check on for upgrade phase"

	nodePreCheck := &task.RemoteTask{
		Name:     "NodePreCheck",
		Desc:     "A pre-check on nodes",
		Hosts:    c.Runtime.GetAllHosts(),
		Action:   new(precheck.NodePreCheck),
		Parallel: true,
	}

	getKubeConfig := &task.RemoteTask{
		Name:     "GetKubeConfig",
		Desc:     "Get KubeConfig file",
		Hosts:    c.Runtime.GetHostsByRole(common.Master),
		Prepare:  new(common.OnlyFirstMaster),
		Action:   new(precheck.GetKubeConfig),
		Parallel: true,
	}

	getAllNodesK8sVersion := &task.RemoteTask{
		Name:     "GetAllNodesK8sVersion",
		Desc:     "Get all nodes Kubernetes version",
		Hosts:    c.Runtime.GetHostsByRole(common.K8s),
		Action:   new(precheck.GetAllNodesK8sVersion),
		Parallel: true,
	}

	calculateMinK8sVersion := &task.LocalTask{
		Name:   "CalculateMinK8sVersion",
		Desc:   "Calculate min Kubernetes version",
		Action: new(precheck.CalculateMinK8sVersion),
	}

	calculateMaxK8sVersion := &task.LocalTask{
		Name:   "CalculateMaxK8sVersion",
		Desc:   "Calculate max Kubernetes version",
		Action: new(CalculateMaxK8sVersion),
	}

	checkDesiredK8sVersion := &task.LocalTask{
		Name:   "CheckDesiredK8sVersion",
		Desc:   "Check desired Kubernetes version",
		Action: new(precheck.CheckDesiredK8sVersion),
	}

	checkUpgradeK8sVersion := &task.LocalTask{
		Name:   "checkUpgradeK8sVersion",
		Desc:   "Check the Kubernetes version can correctly upgrade",
		Action: new(CheckUpgradeK8sVersion),
	}

	ksVersionCheck := &task.RemoteTask{
		Name:     "KsVersionCheck",
		Desc:     "Check KubeSphere version",
		Hosts:    c.Runtime.GetHostsByRole(common.Master),
		Prepare:  new(common.OnlyFirstMaster),
		Action:   new(precheck.KsVersionCheck),
		Parallel: true,
	}

	dependencyCheck := &task.RemoteTask{
		Name:  "DependencyCheck",
		Desc:  "Check dependency matrix for KubeSphere and Kubernetes",
		Hosts: c.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(precheck.KubeSphereExist),
		},
		Action:   new(precheck.DependencyCheck),
		Parallel: true,
	}

	getKubernetesNodesStatus := &task.RemoteTask{
		Name:     "GetKubernetesNodesStatus",
		Desc:     "Get kubernetes nodes status",
		Hosts:    c.Runtime.GetHostsByRole(common.Master),
		Prepare:  new(common.OnlyFirstMaster),
		Action:   new(precheck.GetKubernetesNodesStatus),
		Parallel: true,
	}

	c.Tasks = []task.Interface{
		nodePreCheck,
		getKubeConfig,
		getAllNodesK8sVersion,
		calculateMinK8sVersion,
		calculateMaxK8sVersion,
		checkDesiredK8sVersion,
		checkUpgradeK8sVersion,
		ksVersionCheck,
		dependencyCheck,
		getKubernetesNodesStatus,
	}
}
