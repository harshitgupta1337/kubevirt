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
)

// LauncherRenderMode identifies the lifecycle operation for which a launcher
// Pod is being rendered.
type LauncherRenderMode string

const (
	LauncherRenderModeLaunch          LauncherRenderMode = "launch"
	LauncherRenderModeMigrationTarget LauncherRenderMode = "migration-target"
	LauncherRenderModeProvisioning    LauncherRenderMode = "provisioning"
)

// LauncherManifestRendererContext contains request metadata shared by launcher
// renderers. Configuration must contain the effective, defaulted KubeVirt
// configuration.
type LauncherManifestRendererContext struct {
	Configuration *virtv1.KubeVirtConfiguration
	Mode          LauncherRenderMode
}

// LauncherManifestRenderRequest is the complete renderer API boundary. Pod
// contains the stack-neutral base launcher Pod.
type LauncherManifestRenderRequest struct {
	VMI     *virtv1.VirtualMachineInstance
	Pod     *k8sv1.Pod
	Context *LauncherManifestRendererContext
}

// LauncherManifestRenderer completes a stack-neutral virt-launcher Pod for a
// virtualization stack. Implementations update request.Pod in place.
type LauncherManifestRenderer interface {
	Render(ctx context.Context, request *LauncherManifestRenderRequest) error
}
