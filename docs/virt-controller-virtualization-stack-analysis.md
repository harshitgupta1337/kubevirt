# `virt-controller` Virtualization-Stack Dependency Analysis

## 1. Executive summary

### Verified facts

The proposed plugin boundary is too narrow. `virt-controller` calculates launcher
overhead, but it also renders a launcher pod with a fixed process entrypoint,
launcher image, filesystem and socket layout, device-plugin resources,
capability set, probe mechanism, migration transport marker, firmware path, and
node-label vocabulary.

Per-VMI stack selection is not currently supported:

- The API restricts `KubeVirtConfiguration.Hypervisors` to one item
  ([staging/src/kubevirt.io/api/core/v1/types.go](../staging/src/kubevirt.io/api/core/v1/types.go#L3183-L3187)).
- `GetHypervisorFromKvConfig` explicitly returns only the first cluster
  configuration and defaults to KVM
  ([pkg/virt-config/virt-config.go](../pkg/virt-config/virt-config.go#L520-L539)).
- One launcher image enters `virt-controller` through `--launcher-image`
  ([pkg/virt-controller/watch/application.go](../pkg/virt-controller/watch/application.go#L1049-L1056)).
- One `TemplateService` caches that image and one
  `LauncherHypervisorResources` implementation at controller initialization
  ([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L1406-L1441)).
- Unknown hypervisor names silently select KVM resources
  ([pkg/hypervisor/launcherhypervisorresources.go](../pkg/hypervisor/launcherhypervisorresources.go#L40-L48)).

A second backend can be introduced without changing core pod ownership, storage
resolution, or scheduling policy, but launcher rendering must first be
decomposed into validated, typed contributions. Full migration support also
requires backend separation inside `virt-launcher` and `virt-handler`: the
current migration implementation constructs Libvirt domain XML and Libvirt
migration proxies.

### Architectural assessment

The initial plugin should not return an unrestricted pod mutation. It should
return typed backend requirements that core validates and merges into a pod.
Core should retain final authority over security, ownership, orchestration,
resource accounting, and scheduling policy.

This means that the plugins should expose specific pieces of information needed to construct to `virt-launcher` pod.

## 2. Launcher pod construction call graph

```text
VMI controller sync
  pkg/virt-controller/watch/vmi/lifecycle.go:125-164
    -> TemplateService.RenderLaunchManifest
       pkg/virt-controller/services/template.go:355-365
       -> CalculateMemoryOverhead
          template.go:1687-1705
          -> LauncherHypervisorResources.GetMemoryOverhead
       -> TemplateService.renderLaunchManifest
          template.go:395-780
          -> newResourceRenderer
          -> sidecar creators
          -> fixed virt-launcher-monitor command and arguments
          -> newVolumeRenderer
          -> newContainerSpecRenderer
          -> newNodeSelectorRenderer
          -> PodSpec construction
    -> Kubernetes Pod create

Migration controller createTargetPod
  pkg/virt-controller/watch/migration/migration.go:940-965
    -> TemplateService.RenderMigrationManifest
       template.go:324-351
       -> same renderLaunchManifest path
       -> source image-digest preservation and target annotations
```

The controller-side entry is
[`Controller.sync`](../pkg/virt-controller/watch/vmi/lifecycle.go#L125-L164).
The common controlling renderer is
[`TemplateService.renderLaunchManifest`](../pkg/virt-controller/services/template.go#L395-L780).

## 3. Findings table

| ID | Classification | Severity | Exact file and symbol | Assumption | Recommended disposition |
| --- | --- | --- | --- | --- | --- |
| F1 | Explicit current-stack dependency | Blocking | [`GetHypervisorFromKvConfig`](../pkg/virt-config/virt-config.go#L520-L539), [`NewTemplateService`](../pkg/virt-controller/services/template.go#L1406-L1441) | One stack and launcher image apply to the controller and cluster. | Generalize behind a core registry keyed by an immutable per-VMI stack ID. |
| F2 | Explicit current-stack dependency | Blocking | [`TemplateService.renderLaunchManifest`](../pkg/virt-controller/services/template.go#L452-L500) | Every compute container runs `virt-launcher-monitor` with QEMU, Libvirt, OVMF, and hypervisor flags. | Plugin returns structured set of launcher-container requirements; core validates them and constructs the final pod. |
| F3 | Explicit current-stack dependency | Blocking | [`VolumeRenderer.Mounts` and `Volumes`](../pkg/virt-controller/services/rendervolumes.go#L69-L96) | Every launcher has Libvirt runtime and KubeVirt socket directories. | Keep core shared volumes; Plugin returns stack-specific volumes and mounts. |
| F4 | Extension point for future backends | Blocking | [`LauncherHypervisorResources`](../pkg/hypervisor/launcherhypervisorresources.go#L30-L48) | The backend abstraction exposes only device name and memory overhead; unknown names become KVM. | Replace defaulting with explicit lookup and broaden the typed contribution contract. |
| F5 | Explicit current-stack dependency | Significant | [`KvmHypervisorBackend.GetMemoryOverhead`](../pkg/hypervisor/kvm/hypervisorbackend.go#L62-L160) | The formula models Libvirt/QEMU processes and QEMU-specific memory behavior. | Delegate stack overhead calculation and return an explainable breakdown. |
| F6 | Implicit stack-specific behavior | Significant | [`getRequiredResources`](../pkg/virt-controller/services/renderresources.go#L467-L486), [`VMIResourcePredicates`](../pkg/virt-controller/services/template.go#L1641-L1685) | Core assumes TUN, vhost-net, vhost-vsock, KVM/MSHV, SEV, TDX, and IOMMUFD resources. | VMI Resource Predicates are just launcher spec modification rules that execute one after the other. We can query plugin for specific devices, etc. and construct list of renderers. ASSUMPTION: No arbitrary mod of launcher spec. |
| F7 | Implicit stack-specific behavior | Significant | [`newNodeSelectorRenderer`](../pkg/virt-controller/services/template.go#L795-L859), [`NodeSelectorRenderer.Render`](../pkg/virt-controller/services/nodeselectorrenderer.go#L54-L96) | Core knows fixed CPU model, machine type, Hyper-V/KVM, SEV, TDX, and TSC labels. | Plugin should return a set of stack-specific labels that need to be added to launcher pod. |
| F8 | Implicit stack-specific behavior | Significant | [`computePodSecurityContext`](../pkg/virt-controller/services/template.go#L374-L393), [`requiredCapabilities`](../pkg/virt-controller/services/rendercontainer.go#L292-L303) | Launcher UID 107, capabilities, seccomp, SELinux, and runtime class are globally shaped. | Plugin declares requirements; core validates and applies final policy. |
| F9 | Explicit current-stack dependency | Significant | [`updateReadinessProbe` and `updateLivenessProbe`](../pkg/virt-controller/services/rendercontainer.go#L245-L286) | Probes use `virt-probe`, a `virtwrap/api` domain key, and a fixed Libvirt startup delay. | Define a stack-neutral launcher health RPC call (In Command API?) or delegate a bounded probe adapter. |
| F10 | Explicit current-stack dependency | Blocking for migration | [`generatePodAnnotations`](../pkg/virt-controller/services/template.go#L1548-L1564), [`prepareMigrationTarget`](../pkg/virt-launcher/virtwrap/live-migration-target.go#L160-L265) | Every pod advertises Unix migration while launcher migration converts to Libvirt `api.Domain` and uses Libvirt ports and sockets. | Keep migration orchestration in core; delegate backend migration mechanics and capabilities. |
| F11 | Explicit current-stack dependency | Significant | [`withBackendStorage`](../pkg/virt-controller/services/rendervolumes.go#L387-L445), [`withHugepages`](../pkg/virt-controller/services/rendervolumes.go#L500-L530) | Persistent TPM, NVRAM, and hugepages use Libvirt/QEMU directory layouts. | Plugin provides stack-specific runtime paths and persistent-state layout. IMPORTANT: Plugin does not allow arbitrary pod mods. It just provides values to virt-controller for certain fixed fields.|
| F12 | Extension point for future backends | Minor | [`SidecarCreatorFunc`](../pkg/virt-controller/services/sidecar.go#L26-L35), renderer options | Existing options compose pod fragments but are in-process, startup-registered, and not a stack-selection contract. | Reuse their internal composition patterns, not the APIs as an unrestricted plugin boundary. |

### Observed MSHV pod: backend selection inside a Libvirt/QEMU launcher

The following excerpts come from a running `virt-launcher` pod configured with
the `ConfigurableHypervisor` feature and `hyperv-direct`. They are runtime
evidence supplied with this analysis, rather than source-controlled test data.
They show that selecting MSHV changes the hypervisor device and launcher
argument, but does not replace the surrounding Libvirt/QEMU launcher contract.

The pod selects MSHV explicitly and receives its device-plugin resource:

```yaml
args:
- --hypervisor
- hyperv-direct
resources:
  limits:
    devices.kubevirt.io/mshv: "1"
  requests:
    devices.kubevirt.io/mshv: "1"
```

This is the variation modeled by `LauncherHypervisorResources`: the MSHV
backend supplies `mshv` as its device name
([pkg/hypervisor/mshv/hypervisorbackend.go](../pkg/hypervisor/mshv/hypervisorbackend.go#L45-L59)),
and core turns it into `devices.kubevirt.io/mshv`
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L1024-L1042)).

Despite that selection, the compute container still starts the standard
launcher monitor with QEMU- and Libvirt-specific controls:

```yaml
command:
- /usr/bin/virt-launcher-monitor
args:
- --qemu-timeout
- 345s
- --ovmf-path
- /usr/share/edk2/ovmf
- --hypervisor
- hyperv-direct
- --libvirt-hook-server-and-client
image: harshitg/virt-launcher:vmbus
```

The MSHV choice is therefore a mode within the existing launcher image and
process topology, not selection of an independent stack-specific launcher.
The source renderer fixes the monitor command and QEMU timeout before appending
the cluster hypervisor
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L452-L500)).
The additional `--libvirt-hook-server-and-client` flag is direct runtime
evidence that this particular MSHV launcher still expects Libvirt integration;
the exact flag is not emitted by the upstream renderer lines analyzed here and
may come from the supplied build or configuration.

The generated pod also retains an unconditional Libvirt runtime volume:

```yaml
volumeMounts:
- mountPath: /var/run/libvirt
  name: libvirt-runtime
- mountPath: /var/run/kubevirt/sockets
  name: sockets
volumes:
- emptyDir: {}
  name: libvirt-runtime
- emptyDir: {}
  name: sockets
```

This exactly matches the default mounts and volumes in `VolumeRenderer`
([pkg/virt-controller/services/rendervolumes.go](../pkg/virt-controller/services/rendervolumes.go#L69-L96)).
An MSHV backend that did not use Libvirt would still receive these directories.

Launcher-owned auxiliary behavior is coupled to the same image and binary set:

```yaml
initContainers:
- name: guest-console-log
  image: harshitg/virt-launcher:vmbus
  command:
  - /usr/bin/virt-tail
  args:
  - --logfile
  - /var/run/kubevirt-private/dfd95d56-513c-4f67-9164-c2a1e6491119/virt-serial0-log
```

The compute container and console logger use the same launcher image. This
illustrates why image selection alone is insufficient: a stack contract may
also need to describe, replace, or explicitly support launcher helper binaries
and init containers. The upstream renderer passes `t.launcherImage` to the
serial-console container generator
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L630-L636)).

The security identity remains the QEMU-oriented KubeVirt identity even under
MSHV:

```yaml
securityContext:
  fsGroup: 107
  runAsUser: 0
initContainers:
- name: guest-console-log
  securityContext:
    runAsNonRoot: true
    runAsUser: 107
```

UID 107 is named `qemu` in KubeVirt utilities
([pkg/util/util.go](../pkg/util/util.go#L23-L40)), and pod security always sets
that value as `fsGroup`
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L374-L393)).
The numeric identity may remain a useful compatibility convention, but it is
not stack-neutral by origin and should not become an undocumented plugin
requirement.

Finally, the pod advertises Libvirt-era migration and machine-model concepts
even though its selected hypervisor resource is MSHV:

```yaml
metadata:
  annotations:
    kubevirt.io/migrationTransportUnix: "true"
spec:
  nodeSelector:
    machine-type.node.kubevirt.io/q35: "true"
```

Core currently adds the migration annotation to every launcher pod
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L1548-L1564))
and translates the requested machine type directly into a fixed node-labeler
label
([pkg/virt-controller/services/nodeselectorrenderer.go](../pkg/virt-controller/services/nodeselectorrenderer.go#L159-L170)).
The annotation does not prove that this VMI was migrating; it proves that the
rendered launcher contract preselects the Unix migration mechanism. Likewise,
`q35` may be valid for the current MSHV implementation, but a future VMM may
use a different machine-model vocabulary or none at all.

The pod reports `kubevirt.io/memory-overhead-bytes: "281018368"` while using
the MSHV device. The manifest alone cannot identify the formula that produced
the value. Source evidence supplies that missing link: the MSHV backend uses
the same named `virtlogd`, `virtqemud`, and QEMU process overhead components as
the KVM backend
([pkg/hypervisor/mshv/hypervisorbackend.go](../pkg/hypervisor/mshv/hypervisorbackend.go#L35-L83)).
Together, the runtime value and implementation show that overhead is selected
for MSHV, but still calculated from the shared Libvirt/QEMU process topology.

## 4. Detailed evidence for each finding

### F1: Singleton stack and launcher image

**Call path:** `VirtControllerApp.initCommon` -> `NewTemplateService` -> cached
backend -> every `RenderLaunchManifest`.

`initCommon` constructs one service from `vca.launcherImage`
([pkg/virt-controller/watch/application.go](../pkg/virt-controller/watch/application.go#L689-L722)).
The service stores one image and initializes `launcherHypervisorResources` from
the cluster hypervisor
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L1423-L1441)).
Rendering also reads the cluster hypervisor for the `--hypervisor` argument
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L460-L474)).

The API permits only `kvm` and `hyperv-direct` and caps the list at one
([staging/src/kubevirt.io/api/core/v1/types.go](../staging/src/kubevirt.io/api/core/v1/types.go#L3232-L3247)).
Tests enforce global defaulting and selection
([pkg/virt-config/configuration_test.go](../pkg/virt-config/configuration_test.go#L1023-L1060)).
Launcher version tracking compares every pod to the one service image
([pkg/virt-controller/watch/vmi/lifecycle.go](../pkg/virt-controller/watch/vmi/lifecycle.go#L846-L858)).

**Concrete assumption:** one stack configuration and one launcher image apply to
all VMIs observed by the controller.

**Why another stack differs:** stacks need distinct images, runtime contracts,
and upgrade streams while coexisting in one cluster or node.

**Recommended disposition:** add an immutable per-VMI stack assignment and use a
core-owned resolver to select the plugin for each render. Selection policy
belongs in core admission; rendering details belong in the plugin.

**Confidence:** High.

### F2: Fixed compute-container runtime

The compute container is always generated from `t.launcherImage`
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L996-L1015))
and invokes `virt-launcher-monitor` with `--qemu-timeout`, OVMF, disk
verification, hook count, and cluster hypervisor arguments
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L452-L500)).
Libvirt debug annotations become launcher flags, and verbosity becomes
Libvirt-specific environment variables
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L493-L496),
[pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L530-L548)).

Tests assert the complete command and argument vector
([pkg/virt-controller/services/template_test.go](../pkg/virt-controller/services/template_test.go#L588-L606)).

**Concrete assumption:** the launcher image contains KubeVirt's monitor, Libvirt
integration, QEMU controls, firmware layout, and diagnostic binaries.

**Why another stack differs:** Firecracker, Cloud Hypervisor, and OpenVMM have
different monitors, firmware handling, startup lifecycles, and configuration
mechanisms.

**Recommended disposition:** delegate a typed image, entrypoint, arguments,
environment requirements in the Plugin API. Core should reject unknown
or unsafe fields.

**Confidence:** High.

### F3 and F11: Fixed filesystem and runtime-state layout

`VolumeRenderer` unconditionally creates and mounts `libvirt-runtime` at
`/var/run/libvirt`, plus fixed private, public, and socket directories
([pkg/virt-controller/services/rendervolumes.go](../pkg/virt-controller/services/rendervolumes.go#L69-L96)).
Tests explicitly require these volumes and mounts
([pkg/virt-controller/services/rendervolumes_test.go](../pkg/virt-controller/services/rendervolumes_test.go#L44-L50)).

Persistent TPM setup creates `/var/run/kubevirt-private/libvirt/qemu`
([pkg/virt-controller/services/rendervolumes.go](../pkg/virt-controller/services/rendervolumes.go#L407-L442));
hugepages are mounted under `libvirt/qemu`
([pkg/virt-controller/services/rendervolumes.go](../pkg/virt-controller/services/rendervolumes.go#L510-L529)).
Shared helpers encode Libvirt/QEMU locations for swtpm and NVRAM
([pkg/util/util.go](../pkg/util/util.go#L193-L216)).

**Concrete assumption:** Libvirt and QEMU own launcher runtime state, firmware
state, TPM state, hugepage files, and sockets at fixed paths.

**Why another stack differs:** another VMM may not use Libvirt, may use one
process, and may expose different API sockets and persistent state.

**Recommended disposition:** keep VMI storage and PVC resolution in core;
delegate backend-private runtime volumes, mounts, socket names, hugepage mount
targets, and persistent-state subpaths through typed fields.

**Confidence:** High.

### F4 and F5: Narrow backend interface and QEMU overhead

`LauncherHypervisorResources` exposes only a device name and overhead
calculation
([pkg/hypervisor/launcherhypervisorresources.go](../pkg/hypervisor/launcherhypervisorresources.go#L30-L48)).
Its default branch is KVM, including for unknown names; this behavior is tested
([pkg/hypervisor/launcherhypervisorresources_test.go](../pkg/hypervisor/launcherhypervisorresources_test.go#L30-L37)).

KVM overhead explicitly adds `virt-launcher-monitor`, `virt-launcher`,
`virtlogd`, `virtqemud`, and QEMU RSS, followed by QEMU page tables, vCPU and
IOThread memory, video RAM, ARM pflash, VFIO MMIO, SEV, and swtpm costs
([pkg/hypervisor/kvm/hypervisorbackend.go](../pkg/hypervisor/kvm/hypervisorbackend.go#L35-L160)).
The MSHV implementation currently repeats the same Libvirt/QEMU process formula
([pkg/hypervisor/mshv/hypervisorbackend.go](../pkg/hypervisor/mshv/hypervisorbackend.go#L35-L83)).
Tests fix the static value at `243Mi` and exercise the components
([pkg/hypervisor/kvm/hypervisorbackend_test.go](../pkg/hypervisor/kvm/hypervisorbackend_test.go#L41-L96)).

**Concrete assumption:** backend variation can be represented by one device name
and a memory formula, while the surrounding launcher topology remains fixed.

**Why another stack differs:** process count, page-table cost, firmware,
I/O-thread model, video device, and locked-memory behavior vary by VMM.

**Recommended disposition:** delegate overhead components and return an
explainable breakdown rather than only one aggregate quantity. Core should
combine guest memory, plugin overhead, network-plugin overhead, sidecar
overhead, hugepages, and Kubernetes QoS rules.

**Confidence:** High. Exact formulas for future stacks are intentionally left to
their implementations.

### F6: Device and resource assumptions

`getRequiredResources` chooses TUN, vhost-net, stack device, and vhost-vsock
resources; its comment explicitly relies on Libvirt falling back to QEMU
user-mode networking
([pkg/virt-controller/services/renderresources.go](../pkg/virt-controller/services/renderresources.go#L467-L486)).
The stack device becomes `devices.kubevirt.io/<backend-device>`
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L1024-L1042)).
Additional rules request fixed SEV, TDX, persistent-reservation, and IOMMUFD
resources
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L1641-L1685)).

Tests require KVM unless global emulation is enabled
([pkg/virt-controller/services/template_test.go](../pkg/virt-controller/services/template_test.go#L269-L303)).

**Concrete assumption:** all launchers consume the current KubeVirt device
plugins and Libvirt/QEMU networking fallback behavior.

**Why another stack differs:** backends may use userspace networking, different
accelerator devices, DRA resources, or no `/dev/kvm`-style device at all.

**Recommended disposition:** core owns resources directly requested by the VMI
and generic Kubernetes/DRA semantics. Plugins map abstract VMI requirements to
stack-specific resource names and quantities.

**Confidence:** High.

### F7: Fixed node-label vocabulary

The controller converts VMI CPU model and features, machine type, TSC, realtime,
SEV/SNP/ES, Secure Execution, TDX, and cross-architecture state into node
selectors
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L795-L859)).
`NodeSelectorRenderer` writes fixed KubeVirt labels directly
([pkg/virt-controller/services/nodeselectorrenderer.go](../pkg/virt-controller/services/nodeselectorrenderer.go#L54-L96)).
Hyper-V rules themselves contain QEMU/KVM capability assumptions
([pkg/virt-controller/services/nodeselectorrenderer.go](../pkg/virt-controller/services/nodeselectorrenderer.go#L228-L241)).
Tests enforce examples such as Intel selection for Hyper-V
([pkg/virt-controller/services/nodeselectorrenderer_test.go](../pkg/virt-controller/services/nodeselectorrenderer_test.go#L50-L67)).

**Concrete assumption:** core knows how each guest requirement maps to fixed
host and hypervisor labels.

**Why another stack differs:** capability discovery and the relationship between
guest CPU features and VMM support differ by stack and version.

**Recommended disposition:** core retains VMI and user selectors, affinity
merging, anti-affinity, architecture policy, and safety validation. A plugin
returns named capability requirements; a trusted core mapping converts those
into labels advertised by the selected stack's node-labeler. Do not allow
unrestricted selector injection.

**Confidence:** High.

### F8 and F9: Security and probe ABI

Pod security chooses root or UID/GID 107 and always sets FSGroup 107
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L374-L393));
UID 107 is named `qemu` in shared utilities
([pkg/util/util.go](../pkg/util/util.go#L23-L40)).
Root launchers receive `NET_BIND_SERVICE` and `SYS_NICE`; non-root launchers
retain `NET_BIND_SERVICE`
([pkg/virt-controller/services/rendercontainer.go](../pkg/virt-controller/services/rendercontainer.go#L292-L303)).
Tests enforce these capabilities
([pkg/virt-controller/services/rendercontainer_test.go](../pkg/virt-controller/services/rendercontainer_test.go#L51-L100)).
Seccomp, SELinux type, and runtime class are applied from cluster-wide
configuration during pod construction
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L680-L756)).

Readiness and liveness probes are wrapped with `virt-probe`, use a
`virtwrap/api` domain key, and add ten seconds because Libvirt is expected to
start
([pkg/virt-controller/services/rendercontainer.go](../pkg/virt-controller/services/rendercontainer.go#L245-L286),
[pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L108-L110)).

**Concrete assumption:** the launcher runs as the QEMU account, requires the
current capability set, follows cluster-wide runtime security settings, and
exposes Libvirt-compatible probe behavior.

**Why another stack differs:** a launcher may be rootless, need different
devices or capabilities, start at a different rate, and have no Libvirt domain
identity.

**Recommended disposition:** plugins declare minimal security and health
requirements. Core owns the allowlist, Pod Security compatibility,
seccomp/SELinux policy, service account, privilege-escalation rules, and final
security context. Prefer a stack-neutral launcher health command over
backend-defined arbitrary probes.

**Confidence:** High.

### F10: Migration is a separate backend boundary

Every rendered launcher receives the Unix migration annotation
([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L1548-L1564));
tests require it even for ordinary launch pods
([pkg/virt-controller/services/template_test.go](../pkg/virt-controller/services/template_test.go#L1117-L1126)).
Migration target rendering reuses the same launcher renderer
([pkg/virt-controller/watch/migration/migration.go](../pkg/virt-controller/watch/migration/migration.go#L940-L965)).

Inside the launcher, target preparation creates a Libvirt converter context,
converts the VMI to `api.Domain`, invokes domain hooks, and creates Libvirt
migration socket or TCP proxies
([pkg/virt-launcher/virtwrap/live-migration-target.go](../pkg/virt-launcher/virtwrap/live-migration-target.go#L160-L265)).
The nominal `DomainManager` interface returns `api.DomainSpec` and exposes
Libvirt-shaped lifecycle and migration operations
([pkg/virt-launcher/virtwrap/manager.go](../pkg/virt-launcher/virtwrap/manager.go#L169-L190)).

**Concrete assumption:** the launcher implements Libvirt domain conversion and
Libvirt migration semantics behind a common pod annotation and socket layout.

**Why another stack differs:** migration support, transport, compatibility,
state format, pre-copy/post-copy support, and endpoint count vary by backend.

**Recommended disposition:** core continues migration orchestration,
source/target placement, status, timeouts, and compatibility checks. Backend
implementations own migration preparation, transport endpoints,
capability/version negotiation, state transfer, and finalization.

**Confidence:** High.

### F12: Existing extension machinery is not a backend contract

Reusable composition mechanisms include:

- `ResourceRendererOption` and predicate rules
  ([pkg/virt-controller/services/renderresources.go](../pkg/virt-controller/services/renderresources.go#L22-L59)).
- `VolumeRendererOption`
  ([pkg/virt-controller/services/rendervolumes.go](../pkg/virt-controller/services/rendervolumes.go#L30-L52)).
- `NodeSelectorRendererOption`
  ([pkg/virt-controller/services/nodeselectorrenderer.go](../pkg/virt-controller/services/nodeselectorrenderer.go#L31-L42)).
- Per-VMI sidecar creators
  ([pkg/virt-controller/services/sidecar.go](../pkg/virt-controller/services/sidecar.go#L26-L35)).
- Network memory and annotation generators
  ([pkg/virt-controller/services/template.go](../pkg/virt-controller/services/template.go#L1758-L1781)).

These are useful internal implementation units, but they are registered
in-process at controller startup and can mutate broad pod fragments. Hook
sidecars also assume the existing launcher and socket ecosystem; they are not a
complete VMM backend boundary.

**Recommended disposition:** reuse the composition model when applying
validated plugin output, but do not expose these function options directly as a
remote or unrestricted plugin API.

**Confidence:** High.

## 5. Existing abstractions that can be reused

1. Preserve `ResourceRenderer`, `VolumeRenderer`, `ContainerSpecRenderer`, and
   `NodeSelectorRenderer` as core assemblers, but feed them validated backend
   contributions.
2. Evolve `LauncherHypervisorResources` into explicit backend lookup rather
   than expanding its defaulting switch indefinitely.
3. Reuse sidecar and annotation-generator patterns for bounded optional
   integrations.
4. Preserve the migration controller's source/target orchestration and
   image-digest handling.
5. Preserve the launcher-side `DomainManager` concept, but remove
   `api.DomainSpec` from a stack-neutral interface or make the existing
   interface explicitly Libvirt-specific.

## 6. Gaps in the proposed plugin boundary

A contract returning only memory overhead and related pod properties is
insufficient.

### Inputs core KubeVirt should provide

- Immutable selected stack ID and requested stack API/ABI version.
- VMI spec plus the status fields relevant to rendering.
- Effective architecture and target node architecture, when known.
- Launch versus migration-target operation and source stack/version metadata.
- Resolved cluster policy and feature gates, preferably as a bounded capability
  view rather than an unrestricted configuration object.
- Resolved storage and network facts, image digests, sidecar inventory, and
  runtime security policy.
- Node-labeler capability schema and version, not raw access to all node labels.

### Outputs a stack plugin may need to provide

- Resource overhead breakdown and backend device resources.
- Launcher image identity and supported launcher ABI.
- Entrypoint, bounded arguments, and bounded environment requirements.
- Backend-private volumes, mounts, volume devices, sockets, and persistent-state
  paths.
- Required Linux devices, capabilities, UID/GID constraints, seccomp features,
  and SELinux needs.
- Named node capabilities and migration compatibility constraints.
- Firmware and runtime artifacts plus backend-specific init or auxiliary
  containers.
- Migration transport capabilities and endpoint requirements.
- A health-check capability compatible with a core-defined probe ABI.

### Pod properties core must retain

Core should retain final control of:

- Owner references and VMI identity labels.
- Service accounts and image-pull credentials.
- Restart and termination policy.
- Readiness gates and core lifecycle annotations.
- DNS configuration.
- User-requested affinity, tolerations, and topology constraints.
- PVC, DataVolume, and network resolution.
- Final resource accounting and Kubernetes requests/limits.
- Final security context and Pod Security enforcement.
- Migration orchestration and status.

### Properties that must not be unrestricted plugin mutations

A plugin should not return an arbitrary `Pod` or `PodSpec`. In particular, it
must not have unrestricted control over:

- Owner references or service accounts.
- Host namespaces or arbitrary host paths.
- Privileged mode, privilege escalation, or arbitrary capabilities.
- Arbitrary node selectors, affinity, scheduler name, or tolerations.
- Image-pull secrets.
- Admission-sensitive annotations.
- Arbitrary init or sidecar containers.

The plugin should return typed intent. Core should apply allowlists, validate
cross-field invariants, and construct the final pod.

### Logic that belongs elsewhere

- **Admission:** stack selection, authorization, defaulting, immutability, and
  basic stack/VMI compatibility.
- **`virt-handler`:** node-local backend discovery, lifecycle routing, and
  migration proxy integration.
- **Node-labeler plugin:** versioned capability advertisement for the stack on a
  node.
- **Stack-specific `virt-launcher`:** VMM configuration, process supervision,
  device realization, runtime files, and guest lifecycle.
- **Libvirt backend:** Libvirt domain XML and Libvirt-specific hooks.

## 7. Recommended ownership split between core and plugins

| Responsibility | Recommended owner |
| --- | --- |
| Stack selection policy, defaulting, authorization, and immutable VMI assignment | Core admission and API |
| Plugin discovery, version negotiation, timeout, caching, and output validation | `virt-controller` core |
| Pod lifecycle, ownership, VMI status, PVC resolution, and network resolution | `virt-controller` core |
| Resource accounting and final Kubernetes requests and limits | Core, using a plugin contribution |
| Image, entrypoint, backend runtime paths, and backend device requirements | Stack plugin |
| Final security context and scheduling policy | Core, validating typed plugin requirements |
| Capability-to-node-label mapping | Core plus a versioned node-labeler contract |
| VMM process, configuration format, domain XML, and runtime sockets | Stack-specific launcher |
| Migration orchestration and policy | Core |
| Migration mechanics and transport support | Stack launcher and handler integration |
| Node device discovery and lifecycle operations | `virt-handler` stack integration |

## 8. Sidecar versus cluster-wide service assessment

| Concern | Sidecar co-located with each `virt-controller` | Cluster-wide service |
| --- | --- | --- |
| Availability | Shares the controller pod lifecycle; no network dependency beyond local IPC. | Requires an independently highly available deployment and service discovery. |
| Latency | Lowest and predictable. | Adds a network hop to uncached launcher rendering. |
| Scaling | Scales with controller replicas, including idle replicas. | Scales independently with rendering volume. |
| Version compatibility | Strong controller/plugin pairing through co-deployment. | Must negotiate multiple controller and plugin versions concurrently. |
| Caching | Each controller has a cache that may briefly differ from peers. | A central cache is consistent but can become a bottleneck. |
| Failure behavior | Failure disables one co-located controller replica; leader failover can recover. | A service outage can block all launcher pod creation. |
| Upgrade ordering | Controller and plugin can roll together as one compatibility unit. | Requires explicit old/new version skew and careful upgrade ordering. |
| Security | Smaller network surface, though plugin output still affects privileged pod fields. | Requires mTLS, authorization, tenant isolation, and network policy. |
| Effect on pod creation | Local timeout is on the critical path; another controller replica may take over. | A central dependency is on the critical path for every uncached pod render. |

### Recommendation

Start with a version-pinned sidecar per `virt-controller`, using local
authenticated IPC, strict timeouts, deterministic responses, and no fail-open
behavior. A cluster-wide service is attractive only when plugins are expensive,
centrally licensed, or independently operated. In that model it needs HA,
version-skew support, request idempotency, and bounded controller-side caching.

In either model, plugin absence, timeout, invalid output, and version mismatch
should produce a clear VMI condition and prevent pod creation. They must never
silently fall back to KVM.

## 9. Suggested focused follow-up VEP topics

1. VMI stack-selection API, immutability, defaults, and migration
   compatibility.
2. Typed launcher-pod contribution model and core security validation.
3. Versioned stack capability vocabulary shared with node labelers.
4. Launcher image provenance, signing, upgrade skew, and rollback.
5. Resource-overhead accounting and an observable overhead breakdown.
6. Stack-neutral launcher and handler lifecycle ABI.
7. Migration capability negotiation and transport ownership.
8. Backend persistent-state, hotplug, firmware, and snapshot contracts.
9. Plugin availability semantics, timeout behavior, caching, and VMI
   conditions.
10. Conformance tests that run the same core orchestration against two stacks.

## 10. Open questions requiring maintainer input

- Should stack selection be explicit in every VMI, inherited from VM or
  preference objects, selected by namespace policy, or assigned by admission
  defaulting?
- Must live migration across different stack implementations ever be supported,
  or only migration between compatible versions of the same stack?
- Is a stack allowed to supply its own complete launcher image, or must it
  implement a KubeVirt-owned launcher ABI?
- Which node-label names are public API and which are internal node-labeler
  implementation details?
- Can a plugin request host paths or elevated privileges, and which authority
  approves them?
- Should launcher image updates mark only VMIs assigned to that stack as
  outdated?
- Must plugin evaluation be deterministic and side-effect free so responses can
  be safely cached?
- How should plugin absence, timeout, invalid output, and version
  incompatibility be represented in VMI conditions?
- Is MSHV intended as the evolutionary basis for multi-stack support despite
  currently sharing the Libvirt/QEMU overhead and launcher topology?

## Prioritized first refactorings

1. Add an immutable per-VMI stack ID and replace the one-item/global hypervisor
   lookup with explicit backend resolution. Remove silent unknown-to-KVM
   fallback.
2. Introduce a typed, validated launcher contribution object covering image and
   runtime, overhead, resources, private volumes, capabilities, and named
   scheduling requirements.
3. Split `renderLaunchManifest` into core pod assembly plus a Libvirt/QEMU
   backend implementation that exactly reproduces the current pod.
4. Move Libvirt paths, QEMU flags, process overhead, probe delay, device
   resources, and migration annotation out of generic renderers.
5. Add golden pod tests comparing the new Libvirt backend output to current
   output, then add a minimal second-backend test proving that two VMIs can
   select different images, resources, paths, and node capabilities
   concurrently.
6. Separately define stack-neutral `virt-handler` and launcher lifecycle and
   migration interfaces. Pod rendering alone will not provide a functional
   second stack.

---

This report is based on read-only static analysis of the repository as of
2026-08-26. Architectural recommendations are identified separately from
verified implementation facts.
