/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright The KubeVirt Authors.
 *
 */

package services

import (
	"context"

	k8sv1 "k8s.io/api/core/v1"
	virtv1 "kubevirt.io/api/core/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type recordingLauncherManifestRenderer struct {
	request *LauncherManifestRenderRequest
}

var _ LauncherManifestRenderer = (*recordingLauncherManifestRenderer)(nil)

func (r *recordingLauncherManifestRenderer) Render(_ context.Context, request *LauncherManifestRenderRequest) error {
	r.request = request
	request.Pod.Spec.Containers = append(request.Pod.Spec.Containers, k8sv1.Container{
		Name:  "compute",
		Image: "example.com/launcher",
	})
	return nil
}

var _ = Describe("LauncherManifestRenderer", func() {
	It("receives the complete API boundary and updates the base Pod", func() {
		vmi := &virtv1.VirtualMachineInstance{}
		pod := &k8sv1.Pod{}
		configuration := &virtv1.KubeVirtConfiguration{}
		request := &LauncherManifestRenderRequest{
			VMI: vmi,
			Pod: pod,
			Context: &LauncherManifestRendererContext{
				Configuration: configuration,
				Mode:          LauncherRenderModeMigrationTarget,
			},
		}
		renderer := &recordingLauncherManifestRenderer{}

		err := renderer.Render(context.Background(), request)

		Expect(err).ToNot(HaveOccurred())
		Expect(renderer.request.VMI).To(BeIdenticalTo(vmi))
		Expect(renderer.request.Pod).To(BeIdenticalTo(pod))
		Expect(renderer.request.Context.Configuration).To(BeIdenticalTo(configuration))
		Expect(renderer.request.Context.Mode).To(Equal(LauncherRenderModeMigrationTarget))
		Expect(pod.Spec.Containers).To(ConsistOf(k8sv1.Container{
			Name:  "compute",
			Image: "example.com/launcher",
		}))
	})

	DescribeTable("uses stable render mode values",
		func(mode LauncherRenderMode, expected string) {
			Expect(string(mode)).To(Equal(expected))
		},
		Entry("for launch", LauncherRenderModeLaunch, "launch"),
		Entry("for migration targets", LauncherRenderModeMigrationTarget, "migration-target"),
		Entry("for temporary provisioning", LauncherRenderModeProvisioning, "provisioning"),
	)
})
