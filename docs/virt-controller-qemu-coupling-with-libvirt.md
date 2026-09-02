# `virt-controller` coupling to QEMU when Libvirt is retained

## 1. Scope and conclusion

This document asks a narrower question than the general virtualization-stack
analysis: what must change in `virt-controller` if KubeVirt keeps Libvirt and
the existing handler-launcher control model, but permits Libvirt to run a VMM
other than QEMU, for example Cloud Hypervisor?

The important distinction is between three choices that present-day KubeVirt
mostly treats as one:

1. the accelerator or execution mode, such as KVM, software emulation, or MSHV;
2. the Libvirt driver and VMM process, currently the QEMU driver and QEMU;
3. the KubeVirt launcher image and process topology.

`virt-controller` does not directly open a Libvirt connection. It renders the
`virt-launcher` pod. Nevertheless, that pod is a QEMU-driver pod: it starts a
QEMU-oriented launcher image and monitor, requests resources calculated from a
QEMU process model, mounts QEMU-specific runtime paths, and schedules against
capabilities discovered from the current QEMU/Libvirt stack.

Keeping Libvirt therefore preserves much of the controller's Kubernetes
orchestration, but does **not** make VMM replacement transparent. Beyond memory
overhead, the controller needs driver-aware rendering for CPU accounting,
launcher runtime arguments, private paths, capability labels, helper processes,
security requirements, feature compatibility, and migration support.

The repository contains no Cloud Hypervisor implementation or conformance tests.
This document identifies the KubeVirt assumptions that an alternative Libvirt
driver must satisfy; it does not assert that Cloud Hypervisor supports every
listed KubeVirt feature.

## 2. Current selection does not select a Libvirt driver

The current `Hypervisors` configuration is limited to one configured item
([API configuration](../staging/src/kubevirt.io/api/core/v1/types.go#L3183-L3187)).
`GetHypervisorFromKvConfig` returns that first item and defaults to KVM
([configuration lookup](../pkg/virt-config/virt-config.go#L520-L539)). The
selected implementation exposes only a hypervisor device name, Libvirt domain
virtualization type, and memory-overhead calculation
([`LauncherHypervisorResources`](../pkg/hypervisor/launcherhypervisorresources.go#L30-L48)).
Unknown values silently select the KVM implementation.

This is not a Libvirt-driver registry. In the default path, `virt-launcher`
starts `virtqemud` and connects to a `qemu+unix` URI
([launcher connection](../cmd/virt-launcher/virt-launcher.go#L121-L132),
[daemon startup](../pkg/virt-launcher/virtwrap/util/libvirt_helper.go#L256-L313)).
The `--hypervisor` argument currently changes domain conversion details within
that QEMU-oriented launcher; it does not choose an arbitrary Libvirt driver.

A design that supports multiple VMMs must introduce an explicit driver identity,
for example a versioned `LibvirtVMMDriver` value, and must not overload the
current accelerator/hypervisor name. If QEMU and another VMM can coexist, the
selection must be immutable per VMI and used by admission, pod rendering,
node-capability matching, migration compatibility, and launcher image/version
tracking.

## 3. Controller call path

```text
VMI reconciliation
  -> TemplateService.RenderLaunchManifest
    -> resolve one cached LauncherHypervisorResources implementation
    -> CalculateMemoryOverhead
    -> renderLaunchManifest
      -> render resources and dedicated-CPU quantities
      -> choose fixed launcher image, monitor command, and arguments
      -> render fixed Libvirt/QEMU volumes and mounts
      -> render node selectors from QEMU-oriented node labels
      -> render security context, probes, and helper containers
      -> create final virt-launcher Pod
```

The central renderer is
[`TemplateService.renderLaunchManifest`](../pkg/virt-controller/services/template.go#L395-L780).
Migration target pods reuse the same rendering path
([migration target creation](../pkg/virt-controller/watch/migration/migration.go#L940-L965)).

## 4. Findings

| ID | Severity | Present QEMU assumption | Required behavior with another Libvirt VMM |
| --- | --- | --- | --- |
| C1 | Blocking | One cluster-wide hypervisor, one launcher image, and unknown-name fallback to KVM. | Resolve an explicit, versioned VMM driver for each VMI; reject unknown or unavailable drivers without fallback. |
| C2 | Blocking | The compute container always runs `virt-launcher-monitor` with `--qemu-timeout`, OVMF, Libvirt logging, and the current `--hypervisor` mode. | Select a launcher image/runtime contract that contains and supervises the chosen Libvirt driver daemon and VMM. Rename or capability-gate QEMU-only arguments. |
| C3 | Significant | CPU requests count vCPUs, QEMU IOThreads, and one or two isolated emulator-thread CPUs. | Ask the selected driver for its host-thread/resource model and reject unsupported `IOThreads` or `IsolateEmulatorThread` semantics. |
| C4 | Significant | Overhead models `virtqemud`, QEMU RSS, QEMU page tables, QEMU vCPU/IOThread memory, video RAM, pflash, VFIO MMIO, SEV, and `swtpm`. | Calculate an explainable driver-specific overhead while retaining generic KubeVirt, sidecar, probe, and network overhead in core. |
| C5 | Significant | `/var/run/libvirt` is shared with a QEMU-driver daemon; persistent state and hugepages use `libvirt/qemu` subpaths. | Keep a Libvirt runtime mount if the new driver needs it, but obtain driver-private daemon sockets, logs, hugepage paths, TPM state, NVRAM, and other persistent paths from a typed driver profile. |
| C6 | Significant | Scheduling uses CPU-model, CPU-feature, machine-type, Hyper-V, SEV, TDX, and architecture labels produced from the current stack's capabilities. | Advertise capabilities per Libvirt driver and map the VMI only to labels for the selected driver/version. Do not schedule from QEMU capabilities when another driver is selected. |
| C7 | Significant | Device requests assume the current KVM/QEMU networking and device realization, including `/dev/kvm`-style resources, TUN, vhost-net, vhost-vsock, SEV, TDX, and IOMMUFD. | Separate accelerator requirements from VMM requirements. Reuse a resource only when the selected driver advertises the corresponding feature; otherwise reject the VMI feature. |
| C8 | Significant | UID/GID 107, `SYS_NICE`, `NET_BIND_SERVICE`, seccomp/SELinux settings, and filesystem ownership fit the current launcher and QEMU-driver process tree. | Preserve core security policy, but validate a bounded driver security profile and derive process ownership/capabilities from it. |
| C9 | Significant | Probes run `virt-probe`, use the shared domain key, and add a fixed Libvirt startup delay. | The probe wrapper can remain if the launcher ABI remains stable; startup timing and supported guest-agent probing must be driver capabilities rather than QEMU-era constants. |
| C10 | Significant | Console logging, virtiofs, hook sidecars, access-credential propagation, and firmware assets are assumed to exist in the one launcher image. | Keep helpers only where the selected driver supports the same socket/device ABI; otherwise select compatible helpers or reject the feature. |
| C11 | Blocking for migration | Every launcher pod advertises Unix migration while handler/launcher migration uses QEMU-driver URIs, sockets, and tuning options. | Advertise migration only for a driver/version that supports the required mode and is compatible with the source. Driver-specific migration mechanics remain below controller orchestration. |

## 5. Changes beyond memory overhead

### 5.1 Launcher image, command, and daemon topology

The renderer fixes `/usr/bin/virt-launcher-monitor` and supplies
`--qemu-timeout`, `--ovmf-path`, and `--hypervisor`
([command construction](../pkg/virt-controller/services/template.go#L452-L500)).
It also uses the same launcher image for the serial-console logger and virtiofs
helpers
([helper construction](../pkg/virt-controller/services/template.go#L571-L636)).

The monitor binary may remain a KubeVirt-owned stable entrypoint, but its inputs
must describe the selected Libvirt driver rather than assume `virtqemud` and
QEMU. The driver profile needs at least:

- launcher image and supported handler-launcher ABI;
- Libvirt driver ID and connection mode;
- required daemon/VMM binaries and startup timeout;
- firmware and runtime artifacts;
- compatible helper binaries and feature set.

Core should still construct the final command and pod. A driver must not return
an unrestricted `PodSpec`.

### 5.2 CPU quantities are based on the QEMU thread model

Dedicated-CPU rendering adds one CPU per requested vCPU, adds supplemental
IOThread CPUs, and adds one or two CPUs for an isolated emulator thread
([`WithCPUPinning`](../pkg/virt-controller/services/renderresources.go#L316-L363)).
Non-dedicated CPU requests also include IOThreads
([`WithoutDedicatedCPU`](../pkg/virt-controller/services/renderresources.go#L128-L151)).

These quantities are not merely memory overhead. They reserve host CPUs for
QEMU's process/thread structure and must agree with `virt-handler`'s later
thread placement. An alternative VMM may expose different worker, I/O, or vCPU
threads, or may not implement KubeVirt's emulator-thread isolation semantics.
The selected driver should return a typed CPU-overhead plan, while core retains
final Kubernetes request/limit and QoS calculations.

### 5.3 Device-plugin resources are partly accelerator-specific, partly VMM-specific

`getRequiredResources` requests TUN, vhost-net, the selected hypervisor device,
and vhost-vsock. Its vhost-net behavior explicitly relies on Libvirt falling
back to QEMU userspace emulation when the device is absent
([resource selection](../pkg/virt-controller/services/renderresources.go#L467-L486)).
Additional rules request SEV, TDX, persistent-reservation, and IOMMUFD resources
([resource options](../pkg/virt-controller/services/renderresources.go#L400-L451)).

Another Libvirt VMM may still use KVM, TUN, vhost-net, or vhost-vsock; those
resources must not be renamed merely because QEMU is replaced. Instead, the
contract must distinguish:

- host accelerator requirements;
- networking backend requirements;
- optional guest-feature requirements;
- driver-specific process/device requirements.

Admission or rendering must fail clearly when a requested VMI feature is not in
the selected driver's advertised capability set.

### 5.4 Runtime and persistent paths

The generic renderer always mounts an empty directory at `/var/run/libvirt`
([default mounts](../pkg/virt-controller/services/rendervolumes.go#L69-L96)).
That path can remain useful because Libvirt is retained. The subpaths below it
are not generic:

- non-root TPM directories under `libvirt/qemu/swtpm`;
- NVRAM under `libvirt/qemu/nvram`;
- hugepages under `/dev/hugepages/libvirt/qemu`;
- QEMU-driver log and socket locations.

The hard-coded backend-storage and hugepage layouts are in
[`withBackendStorage`](../pkg/virt-controller/services/rendervolumes.go#L387-L445)
and [`withHugepages`](../pkg/virt-controller/services/rendervolumes.go#L500-L530).
Shared utility paths have the same assumption
([TPM and NVRAM paths](../pkg/util/util.go#L193-L216)).

A driver profile should return named, validated path purposes, not arbitrary
host paths. Core can then mount persistent TPM, NVRAM, changed-block-tracking,
and hugepage storage at approved locations.

### 5.5 Capability labels and machine vocabulary

The controller turns the VMI's CPU model/features and machine type directly
into node labels
([selector assembly](../pkg/virt-controller/services/template.go#L795-L859),
[label rendering](../pkg/virt-controller/services/nodeselectorrenderer.go#L54-L96)).
The Hyper-V label table even documents decisions derived from QEMU and KVM
capabilities
([Hyper-V mapping](../pkg/virt-controller/services/nodeselectorrenderer.go#L228-L255)).

Keeping Libvirt allows capability discovery to remain XML-based where the
selected driver exposes suitable capabilities and domain capabilities. It does
not make QEMU machine names, usable CPU models, or optional feature support
portable. Labels need a driver identity/version or a driver-neutral capability
mapping so that a VMI cannot match a node solely because QEMU supports a feature
there.

### 5.6 Security identity and probes

The pod always uses FSGroup 107 and runs either as root or UID/GID 107
([pod security context](../pkg/virt-controller/services/template.go#L374-L393)).
The compute container receives `NET_BIND_SERVICE`, plus `SYS_NICE` for root
launchers
([required capabilities](../pkg/virt-controller/services/rendercontainer.go#L292-L303)).
These may remain valid conventions for a KubeVirt-owned multi-driver launcher,
but they must be tested against the selected Libvirt daemon and VMM process.

Readiness and liveness probes are wrapped by `virt-probe` and receive a fixed
ten-second Libvirt startup allowance
([probe wrapping](../pkg/virt-controller/services/rendercontainer.go#L245-L286)).
The domain-key protocol is reusable if the launcher keeps the same ABI. The
startup delay and guest-agent support need driver-specific validation.

### 5.7 Optional features need capability gating

Many VMI APIs can still be represented in Libvirt XML, but an individual
Libvirt driver may implement only a subset. The selected-driver check must cover
at least:

- CPU modes, CPU features, machine types, and cross-architecture emulation;
- IOThreads, emulator-thread isolation, realtime scheduling, and hotplug;
- firmware, TPM/NVRAM, graphics, video, watchdog, and sound devices;
- host-device/VFIO, vhost-vsock, virtiofs, and persistent reservations;
- SEV/SEV-ES/SEV-SNP, TDX, and other launch-security modes;
- QEMU guest-agent based access credentials and guest probing;
- QEMU command-line extensions used by launcher conversion;
- snapshots, changed-block tracking, memory dump, and live migration.

The controller should reject unsupported combinations before creating a pod.
It should not render a best-effort pod and wait for Libvirt domain definition to
fail.

## 6. What can remain unchanged

The following responsibilities are independent of QEMU and should remain in
`virt-controller` core:

- VMI workqueues, retries, events, finalizers, and pod ownership;
- PVC, DataVolume, container-disk, cloud-init, Secret, and ConfigMap resolution;
- Multus attachment resolution and user-requested Kubernetes scheduling policy;
- final resource accounting, QoS rules, and Pod Security validation;
- service accounts, image-pull secrets, DNS, affinity, tolerations, and topology;
- migration object orchestration and target-pod lifecycle;
- deterministic assembly and validation of the final `Pod`.

The existing `ResourceRenderer`, `VolumeRenderer`, `ContainerSpecRenderer`, and
`NodeSelectorRenderer` can remain core implementation mechanisms. They need
validated driver-specific inputs rather than embedded QEMU defaults.

## 7. Recommended minimal design

Introduce a core-owned registry keyed by an immutable selected driver ID. A
driver implementation should provide typed contributions such as:

```text
LibvirtVMMDriverProfile
  identity and ABI: driver ID, version, launcher image, launcher ABI
  runtime: daemon role, startup timeout, bounded launcher arguments
  resources: accelerator/devices, CPU-thread overhead, memory breakdown
  filesystem: named runtime, log, firmware, TPM, NVRAM, hugepage paths
  capabilities: CPU/machine/features, devices, lifecycle, snapshots, migration
  security: UID/GID and bounded Linux capability requirements
  helpers: compatible probe, console, virtiofs, and hook contracts
```

Core validates and translates that profile into existing renderers. Plugin or
profile output must not control owner references, service accounts, host
namespaces, arbitrary host paths, arbitrary node selectors, privileged mode,
or arbitrary sidecars.

## 8. Test changes

1. Preserve a golden QEMU-driver launcher pod to prove no regression.
2. Add an alternative-driver pod golden test with a different image/runtime,
   overhead, private paths, and capability labels.
3. Run two VMI renders concurrently to prove selection is per VMI rather than a
   cached global choice.
4. Add negative tests for unsupported machine types, CPU/thread policies,
   launch security, devices, snapshots, and migration.
5. Split tests that currently require `--qemu-timeout`, QEMU paths, QEMU
   overhead, and QEMU machine labels into QEMU-driver compatibility tests.
6. Keep generic pod ownership, storage, network, QoS, and security-policy tests
   as core conformance tests.

## 9. Prioritized implementation sequence

1. Add explicit VMM-driver identity and remove unknown-to-KVM fallback.
2. Make launcher image/runtime selection driver-aware while preserving the
   existing QEMU output.
3. Extract memory and CPU-thread resource plans from the current KVM/QEMU
   implementation.
4. Parameterize QEMU private paths, firmware, helpers, and security needs.
5. Publish and consume per-driver capability labels; reject unsupported VMI
   features during admission/rendering.
6. Gate migration advertisement on same-driver/version compatibility and
   supported migration capabilities.
7. Add the first alternative Libvirt-driver profile and conformance tests.

---

This report is based on read-only static analysis of the repository as of
2026-08-31. It deliberately keeps Libvirt in the architecture and separates
verified KubeVirt behavior from unverified capabilities of any particular
alternative Libvirt VMM driver.
