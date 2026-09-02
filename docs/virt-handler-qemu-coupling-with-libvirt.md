# `virt-handler` coupling to QEMU when Libvirt is retained

## 1. Scope and conclusion

This document identifies what must change in `virt-handler` if KubeVirt keeps
Libvirt and the existing handler-launcher interaction, but Libvirt may control a
VMM other than QEMU, for example Cloud Hypervisor.

Two boundary facts prevent overestimating the required handler rewrite:

- `virt-handler` normally calls the versioned `virt-launcher` command socket; it
  does not define domains through Libvirt itself.
- lifecycle reconciliation is driven by the shared `api.Domain` state and
  metadata. That model can remain if the launcher translates the selected
  Libvirt driver's XML, state, and events into the existing contract.

The direct QEMU coupling is concentrated in privileged node-side mechanics:
process and thread discovery, CPU placement, memory-lock adjustment, migration
proxy endpoints, QEMU migration policy, passt socket repair, capability labels,
and several QEMU-named optional RPCs and status checks.

Replacing only the vCPU regex and VMM executable name is insufficient. The
handler also assumes a `virtqemud` parent/child lifecycle, QEMU emulator and PIT
threads, QEMU-driver migration sockets and ports, QEMU guest-agent channel
identity, QEMU domain type for software emulation, and capabilities generated
from the current QEMU-driver stack.

The repository contains no Cloud Hypervisor handler implementation or tests.
This document therefore states the contract an alternative Libvirt driver must
meet and flags features that require capability negotiation; it does not claim
that Cloud Hypervisor implements all of them.

## 2. Current architecture: handler-to-launcher is reusable, launcher-to-Libvirt is not

The handler command boundary is
[`LauncherClient`](../pkg/virt-handler/cmd-client/client.go#L82-L137). General
operations such as sync, pause, shutdown, delete, hotplug, statistics, and
migration are sent to `virt-launcher` over the launcher socket.

The launcher then starts `virtqemud` and connects to
`qemu+unix:///system` or `qemu+unix:///session`
([launcher connection](../cmd/virt-launcher/virt-launcher.go#L121-L132),
[`StartVirtqemud`](../pkg/virt-launcher/virtwrap/util/libvirt_helper.go#L256-L313)).
Thus, selecting another Libvirt driver is primarily a launcher implementation
change, but `virt-handler` must stop reaching through that boundary with
QEMU-specific process, thread, path, and migration knowledge.

The existing `hypervisor.VirtRuntime` is the natural node-side seam, but it has
only two methods
([`VirtRuntime`](../pkg/hypervisor/virt-runtime.go#L33-L46)). Unknown names
silently select KVM, and both current implementations retain a QEMU-shaped
process model. It should become an explicitly resolved Libvirt-VMM runtime
adapter, while the common launcher lifecycle API remains stable.

## 3. Handler control flow

```text
VMI or domain event
  -> VirtualMachineController.execute / sync
    -> read VMI and shared api.Domain
    -> perform node mounts, network, devices, cgroups, and isolation
    -> call LauncherClient lifecycle operations
    -> VirtRuntime.HandleHousekeeping / AdjustResources
    -> update VMI status and conditions from api.Domain

Migration source or target event
  -> prepare node resources
  -> find launcher PID and enter its mount namespace
  -> construct QEMU-driver sockets and TLS proxies
  -> pass QEMU-oriented migration options to LauncherClient
  -> observe api.Domain migration metadata and transfer ownership
```

The first flow is mostly reusable. The second reaches deeply into the current
QEMU-driver implementation.

## 4. Findings

| ID | Severity | Present QEMU assumption | Required behavior with another Libvirt VMM |
| --- | --- | --- | --- |
| H1 | Blocking | `AdjustResources` finds QEMU when running and `virtqemud` before domain creation, then changes that process's memlock limit. | The selected runtime identifies process roles and the correct pre-spawn daemon or VMM resource target. Core applies only validated rlimit policy. |
| H2 | Blocking for dedicated/realtime CPU | vCPUs are QEMU threads named `CPU <n>/KVM`; housekeeping identifies non-vCPU threads by QEMU names; PIT and emulator-thread handling assume QEMU/KVM internals. | The runtime enumerates typed vCPU, emulator/housekeeping, I/O, and timer thread roles, or rejects unsupported CPU-isolation/realtime features. |
| H3 | Significant | The runtime finds QEMU by executable prefixes and assumes its parent is the launcher process. | Discover the selected VMM by a typed process role or launcher-reported PID, not by arbitrary plugin regex alone. Validate that the PID belongs to the launcher pod cgroup/namespace. |
| H4 | Significant | Node labels and migration host CPU data are loaded from capabilities and domain-capabilities files generated for the current QEMU driver. | Query and publish capabilities for each available Libvirt driver/version; preserve Libvirt XML parsing where applicable but keep label sets driver-specific. |
| H5 | Blocking for migration | Target proxying always includes `virtqemud-sock`, ports 49152/49153, and QEMU-derived Unix socket keys. Source launcher migration uses `qemu+unix` and QEMU migration parameters. | Negotiate a driver migration capability and named endpoint plan. Reject migration if the selected driver/version lacks a compatible implementation. |
| H6 | Significant for migration policy | Handler exposes auto-converge, post-copy, parallel threads, dirty-rate stall detection, compression, and downtime tuning as common migration options. | Separate general policy goals from driver options and pass only negotiated capabilities. Required unsupported policy must fail explicitly. |
| H7 | Significant for passt migration | Passt repair is hard-coded beneath `/var/run/libvirt/qemu/run/passt`. | Resolve the network-backend repair endpoint from a validated driver/network profile or disable repair where the selected driver does not use that socket layout. |
| H8 | Significant | Guest-agent connection detection searches for `org.qemu.guest_agent.0`; QEMU agent operations back guest info, ping, exec, users, filesystems, and access credentials. | Preserve the APIs only if the selected stack exposes a compatible guest-agent channel; otherwise capability-gate guest-agent-dependent features and status. |
| H9 | Significant | The common launcher client exposes `GetQemuVersion`, QEMU dirty-rate statistics, SEV operations, and backend checkpoint redefinition without capability negotiation. | Move QEMU-named and optional functions behind versioned capabilities or return explicit unsupported status without breaking general lifecycle calls. |
| H10 | Minor | Cross-architecture software-emulation status is detected by `domain.Spec.Type == "qemu"`. | Use an explicit execution-mode/status field rather than a QEMU Libvirt domain type. |
| H11 | Driver-dependent | Machine type, memory, volume status, pause reasons, statistics, snapshots, backup, hotplug, and migration progress are read from `api.Domain` or launcher RPCs. | Keep handler logic if the launcher normalizes equivalent Libvirt-driver data; add capability checks or normalized fields where the driver lacks equivalent semantics. |

## 5. Changes beyond process name and vCPU regex

### 5.1 Process lifecycle and memlock targets

`KvmVirtRuntime.AdjustResources` chooses the QEMU process after a domain is
running, but chooses `virtqemud` before the VMM exists or while a migration
target awaits its domain
([resource adjustment](../pkg/hypervisor/kvm/runtime.go#L49-L95)). This depends
on a QEMU-driver behavior: the Libvirt daemon applies or inherits limits before
forking/execing QEMU.

A different Libvirt driver may use a different daemon, process ancestry, or
resource-application mechanism. Merely replacing `qemu-system` with another
executable leaves the pre-start path wrong. The runtime contract should identify
at least:

```text
process role: libvirt-driver-daemon | vmm | helper
lifecycle phase: pre-domain | running | migration-target-pre-domain
operation: set memlock | inspect threads | collect identity
```

The adapter should return PIDs or process handles; trusted handler code should
verify pod/cgroup membership and apply approved rlimits. It should not execute
an arbitrary plugin command as root.

Current executable-prefix and parent-process assumptions are in
[`GetQEMUProcess`](../pkg/hypervisor/kvm/runtime.go#L239-L251) and
[`FindVirtqemudProcess`](../pkg/hypervisor/common/qemu.go#L40-L59).

### 5.2 CPU pinning, housekeeping, and realtime scheduling

There are four distinct QEMU/KVM assumptions, not one regex:

1. vCPU threads match `^CPU (\d+)/KVM$`
   ([`VcpuRegex`](../pkg/hypervisor/kvm/realtime.go#L32-L39));
2. non-vCPU threads are moved into a housekeeping cgroup by excluding names
   containing both the `CPU` prefix and `KVM`
   ([`configureHousekeepingCgroup`](../pkg/hypervisor/kvm/runtime.go#L163-L218));
3. realtime vCPUs receive `SCHED_FIFO` after thread discovery
   ([`configureVCPUScheduler`](../pkg/hypervisor/kvm/realtime.go#L41-L72));
4. a kernel `kvm-pit/<nspid>` thread is found, optionally promoted to FIFO, and
   affined to vCPU 0
   ([`affinePitThread`](../pkg/hypervisor/kvm/runtime.go#L120-L161),
   [`KvmPitPid`](../pkg/hypervisor/kvm/runtime.go#L255-L276)).

The `api.Domain` CPU-tune intent can remain. The selected runtime must map that
intent to its actual thread model. A robust interface should return thread IDs
by typed role rather than accept unbounded regular expressions. If a VMM cannot
provide stable thread identity, admission must reject dedicated-CPU features
whose guarantees cannot be met.

This handler contract must match the controller's CPU accounting. Reserving an
"emulator thread CPU" in the pod without a corresponding alternative-VMM
thread placement mechanism would waste CPU and violate isolation expectations.

### 5.3 Node capability and machine-type discovery

At startup, `virt-handler` unconditionally reads `capabilities.xml`, derives
supported machines, and constructs one node labeller
([capability loading](../cmd/virt-handler/virt-handler.go#L345-L374)). The node
labeller then reads domain capabilities, usable CPU models, host model,
launch-security features, and machine types
([domain-capability loading](../pkg/virt-handler/node-labeller/cpu_plugin.go#L82-L139),
[node-labeller state](../pkg/virt-handler/node-labeller/node_labeller.go#L54-L95)).

Libvirt XML is not itself the problem; retaining Libvirt makes it a useful
source. The problem is that the files and resulting labels describe one active
QEMU-driver environment. Supporting another driver requires:

- discovery against each installed/available Libvirt driver;
- driver ID and version associated with every advertised capability;
- per-driver machine types and usable CPU models;
- explicit reporting of optional device, lifecycle, snapshot, and migration
  capabilities;
- scheduler labels that cannot accidentally match a VMI to another driver's
  capabilities.

Host facts such as NUMA topology and hardware CPU features can remain shared.
Driver realizability of a guest CPU model, machine, or device cannot.

### 5.4 Migration transport and control

The target handler constructs a path to `libvirt/virtqemud-sock`, then appends
sockets derived from ports 49152 and 49153
([target proxy setup](../pkg/virt-handler/migration-target.go#L739-L759),
[port constants](../pkg/virt-handler/migration-proxy/migration-proxy.go#L36-L91)).
The source handler creates corresponding Unix-to-TCP proxies
([source proxy setup](../pkg/virt-handler/migration-source.go#L462-L484)).

This is not generic simply because traffic eventually enters Libvirt. The
launcher uses `qemu+unix` as the destination URI and `MigrateToURI3` parameters
whose state and disk URIs are built from those proxy sockets
([migration parameters](../pkg/virt-launcher/virtwrap/live-migration-source.go#L1077-L1130),
[migration invocation](../pkg/virt-launcher/virtwrap/live-migration-source.go#L1340-L1372)).

A retained-Libvirt design should keep handler ownership of TLS, host listeners,
status publication, and source/target orchestration, but replace implicit QEMU
conventions with a bounded plan:

```text
MigrationPlan
  driver ID and compatibility version
  supported modes: none | cold | live | block
  named endpoints: control, state, optional disk or driver-defined roles
  endpoint direction and protocol
  supported policy options
```

The launcher adapter resolves endpoint names into its driver-local Libvirt
socket/URI. The handler validates endpoint count, protocol, and location beneath
approved launcher runtime roots. It must not accept arbitrary host paths from a
plugin.

If the alternative Libvirt driver does not support KubeVirt's live-migration
requirements, the correct first implementation is to report the VMI as not
live-migratable. Keeping Libvirt must not imply migration support.

### 5.5 Migration policy knobs and dirty-rate metrics

`MigrationSourceController.migrateVMI` always builds a common options object
containing bandwidth, downtime, auto-converge, post-copy, workload disruption,
stall detection, compression, downtime tuning, and parallel migration threads
([option construction](../pkg/virt-handler/migration-source.go#L520-L581)).
Several of these encode QEMU migration semantics.

Dirty-rate monitoring also calls the QEMU-oriented launcher operation
`GetDomainDirtyRateStats`
([client contract](../pkg/virt-handler/cmd-client/client.go#L123-L132)).
A different Libvirt driver might expose a subset, different semantics, or no
live migration. Handler policy should be expressed as goals and requirements;
the driver capability maps them to supported Libvirt operations. Metrics must
report unsupported/not available rather than silently producing QEMU-shaped
values.

### 5.6 Passt migration repair path

Handler migration resolves
`/var/run/libvirt/qemu/run/passt` inside the launcher root
([`passtSocketDirOnHost`](../pkg/virt-handler/unsafepath.go#L27-L42)). This path
is QEMU-driver-private even though namespace entry and safe path resolution are
generic.

The network binding and selected VMM jointly determine whether socket repair is
required and where the socket lives. Keep mount-namespace isolation and safe
path validation in core; make the repair endpoint a named, validated runtime
path supplied by the selected driver/network integration.

### 5.7 Guest agent and QEMU-specific launcher methods

`guestAgentConnected` recognizes only the Libvirt channel target
`org.qemu.guest_agent.0`
([guest-agent detection](../pkg/virt-handler/vm.go#L664-L674)). The launcher
client then exposes guest ping, exec, users, filesystems, and guest information,
and KubeVirt access credentials explicitly offer QEMU guest-agent propagation.

If another VMM presents the same QEMU Guest Agent protocol and Libvirt channel,
this handler behavior can remain as a compatibility capability. Otherwise the
handler needs a normalized `GuestAgentConnected` field or launcher capability,
and admission must reject QEMU-agent-only access credentials and probes.

The common `LauncherClient` also includes `GetQemuVersion`, SEV launch
measurement/secret operations, dirty-rate statistics, and checkpoint
redefinition
([optional methods](../pkg/virt-handler/cmd-client/client.go#L120-L133)). These
should not block a launcher driver that implements only general lifecycle.
Split them into negotiated optional interfaces or add a versioned capability
response with an explicit unsupported error.

`GetQemuVersion` is consumed by downward metrics rather than the main VM
reconciler
([downward-metrics scraper](../pkg/downwardmetrics/scraper/scraper.go#L90-L103));
a generic VMM identity/version field is needed if those metrics are intended to
work with another driver.

### 5.8 Status translation and lifecycle

The following handler behavior can remain if the launcher normalizes the
selected Libvirt driver's domain representation into `api.Domain`:

- existence and active/inactive domain state;
- shutdown grace-period metadata and deletion timestamps;
- pause/resume and normalized pause reason;
- current memory, machine type, volume status, and interface status;
- migration UID, timestamps, completion, failure, and abort metadata;
- backup, freeze, and changed-block-tracking state where supported.

The handler's shutdown policy is general: it sends graceful shutdown, observes
metadata/state, and kills after the deadline
([shutdown processing](../pkg/virt-handler/vm.go#L1549-L1652)). No QEMU process
inspection is needed there.

Driver differences should be handled in the launcher adapter first. Handler
changes are required only when current status fields have QEMU-only meaning or
a feature can be absent. One explicit QEMU leak is software-emulation status,
which tests `domain.Spec.Type == "qemu"`
([software-emulation condition](../pkg/virt-handler/vm.go#L752-L774)); this
should become a normalized execution-mode field.

## 6. What can remain unchanged

The following `virt-handler` responsibilities are not inherently QEMU-specific:

- VMI/domain informers, workqueues, retries, events, and status patching;
- launcher command-socket discovery and version negotiation;
- pod PID detection and safe access to mount/network namespaces;
- container-disk and hotplug-volume mounting;
- generic cgroup v1/v2 operations;
- network and device setup driven by validated requirements;
- shutdown deadlines and kill fallback;
- source/target migration orchestration, ownership transfer, and status;
- TLS certificate policy and host-side migration proxy lifecycle;
- `api.Domain` as the normalized desired/observed state contract.

These remain unchanged only when the selected launcher preserves the command
ABI and translates its Libvirt driver into the normalized state. Unsupported
features must be advertised and rejected; they must not be simulated by empty
or misleading `api.Domain` fields.

## 7. Recommended retained-Libvirt runtime boundary

Expand or replace `VirtRuntime` with an explicitly selected adapter whose
privileged outputs are typed:

```text
LibvirtVMMRuntime
  Identity() -> driver ID/version
  ProcessTargets(phase, operation) -> validated process roles
  ThreadTopology(vmi, domain) -> typed vCPU/housekeeping/I/O/timer thread IDs
  CapabilitySnapshot() -> CPU, machine, device, lifecycle, migration features
  RuntimeEndpoints(purpose) -> named paths beneath approved launcher roots
  MigrationCapabilities() -> modes, endpoints, policy-option support
```

General lifecycle operations should continue through `LauncherClient`. Native
Libvirt-driver mechanics should execute in the launcher. The runtime adapter
must not receive arbitrary host command execution; handler remains responsible
for validating PIDs and paths and applying privileged cgroup, affinity, rlimit,
and proxy operations.

Driver selection must be explicit. The current default branch in
[`GetVirtRuntime`](../pkg/hypervisor/virt-runtime.go#L38-L46) should return an
error for unknown drivers rather than select KVM.

## 8. Test changes

1. Preserve VM reconciliation tests based on `api.Domain` as driver-neutral
   conformance tests.
2. Convert QEMU executable, `virtqemud`, vCPU-regex, housekeeping, and PIT tests
   into QEMU-runtime adapter tests.
3. Add an alternative runtime test with different process ancestry and thread
   names/roles.
4. Verify PID and path validation rejects processes or sockets outside the
   launcher pod.
5. Convert migration tests for ports 49152/49153 and `virtqemud-sock` into QEMU
   migration-adapter compatibility tests.
6. Add a no-live-migration driver case and a fake named-endpoint migration case.
7. Test that unsupported QEMU guest agent, SEV, dirty-rate, checkpoint, and
   migration options fail explicitly.
8. Run the same shutdown, crash, delete, and normalized status tests against
   QEMU and a fake alternative Libvirt-driver launcher.

## 9. Prioritized implementation sequence

1. Add explicit Libvirt VMM-driver identity to VMI/launcher selection and remove
   unknown-to-KVM runtime fallback.
2. Replace QEMU/`virtqemud` process discovery with typed process roles,
   including the pre-domain memlock target.
3. Replace regex-only CPU handling with typed thread topology and align it with
   controller CPU accounting.
4. Publish driver-specific node capabilities and machine/CPU labels.
5. Split optional QEMU launcher methods from the general lifecycle contract.
6. Replace migration ports, socket paths, and options with a negotiated driver
   migration plan; initially allow an alternative driver to report migration as
   unsupported.
7. Parameterize passt and other driver-private runtime paths.
8. Add alternative-driver conformance tests around the unchanged `api.Domain`
   lifecycle state machine.

---

This report is based on read-only static analysis of the repository as of
2026-08-31. It deliberately retains Libvirt and separates verified KubeVirt
behavior from unverified capabilities of any particular alternative Libvirt VMM
driver.
