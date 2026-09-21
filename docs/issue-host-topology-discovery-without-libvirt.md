# Decouple host NUMA topology discovery from Libvirt/QEMU

## Summary

`virt-handler` currently obtains host NUMA topology indirectly from Libvirt. A
privileged init container runs `virsh capabilities`, writes Libvirt capabilities
XML to a host-mounted directory, and `virt-handler` parses that file before it
starts its controllers. The topology is then converted to KubeVirt's command API
types and sent back into each `virt-launcher`.

Host NUMA topology is a property of the node's Linux kernel and hardware, not of
Libvirt or QEMU. Since `virt-handler` is already a privileged, node-local
component with access to host sysfs and procfs, it should discover this data
directly. Doing so removes an unnecessary dependency on the default
Libvirt/QEMU stack and is an incremental step toward the pluggable virtualization
stack described by [VEP 82](https://github.com/harshitgupta1337/kubevirt-enhancements/blob/plugin-virtstack-new/veps/sig-compute/82-plugin-virtstack/vep.md).

## Current flow

1. The `virt-handler` DaemonSet starts a privileged init container from the
   `virt-launcher` image and executes `node-labeller.sh`.
2. The script starts `virtqemud` and runs `virsh capabilities`, writing
   `/var/lib/kubevirt-node-labeller/capabilities.xml`.
3. At startup, `cmd/virt-handler/virt-handler.go` reads that file and unmarshals
   it into `libvirtxml.Caps`.
4. The same object has several independent uses:
   - `capabilities.Host.NUMA.Cells` supplies host NUMA topology.
   - `capabilities.Host.CPU.Counter` supplies TSC frequency/scalability labels.
   - guest machine entries supply supported machine-type labels and metrics.
5. `pkg/virt-handler/options.go` converts the NUMA cells into the
   `cmdv1.Topology` command API type. The payload contains NUMA cell IDs, CPUs,
   SMT sibling sets, NUMA distances, per-cell memory, and page-size information.
6. The topology is attached to `VirtualMachineOptions` for normal sync,
   migration-target sync, CPU and memory hotplug, and target migration topology
   reporting.
7. `virt-launcher` consumes the topology primarily to:
   - group an allocated pod CPU set by physical core and NUMA cell for dedicated
     CPU placement; and
   - construct guest NUMA and `numatune` mappings, including NUMA distances, for
     NUMA guest-mapping passthrough.

The data path is therefore:

```text
Linux host topology
  -> Libvirt/QEMU capability probing in a virt-launcher init container
  -> capabilities.xml on a hostPath
  -> libvirtxml.Caps in virt-handler
  -> cmdv1.Topology in VirtualMachineOptions
  -> virt-launcher CPU pinning and NUMA conversion
```

## Why this is undesirable

- Host topology discovery cannot complete unless the Libvirt/QEMU tooling starts
  successfully, even though the topology itself is virtualization-stack
  independent.
- Core `virt-handler` controller state carries a `*libvirtxml.Caps` solely so
  host topology can later be extracted while preparing command options.
- An alternative `virt-launcher` image must otherwise contain Libvirt/QEMU
  tooling for a node-local operation that `virt-handler` can perform itself.
- The extra init-container process, XML file, hostPath handoff, and XML schema
  conversion obscure ownership of host information and enlarge the startup
  failure surface.
- Topology is discovered once during `virt-handler` startup and treated as an
  immutable Libvirt capability snapshot.

This coupling is unnecessary. Linux exposes the required topology through
stable sysfs/procfs interfaces, and `virt-handler` already reads neighboring
host hardware information from those interfaces.

## Important scope qualification

This issue should separate host topology from virtualization capabilities. The
node-labeller init container also obtains CPU models and features, machine types,
TSC information, cross-architecture emulation support, and launch-security
capabilities from Libvirt/QEMU. Those values are stack-specific or require
separate discovery work.

Moving host topology discovery does not, by itself, allow removal of the init
container or all files under `/var/lib/kubevirt-node-labeller`. It does remove
`libvirtxml.Caps` from the topology path and makes the remaining dependencies
explicit for follow-up work.

## Proposed design

Introduce a small host-topology provider owned by `virt-handler`:

```go
type HostTopologyProvider interface {
    Topology() (*cmdv1.Topology, error)
}
```

The production implementation should read Linux host interfaces directly and
return the existing command API topology type, avoiding a wire-format change:

| Topology field | Linux source |
| --- | --- |
| online NUMA nodes | `/sys/devices/system/node/online` or `node[0-9]*` |
| CPUs in a NUMA node | `/sys/devices/system/node/nodeN/cpulist` |
| online CPUs | `/sys/devices/system/cpu/online` |
| SMT siblings | `/sys/devices/system/cpu/cpuN/topology/thread_siblings_list` |
| NUMA distances | `/sys/devices/system/node/nodeN/distance` |
| node memory | `/sys/devices/system/node/nodeN/meminfo` |
| huge-page sizes/counts | `/sys/devices/system/node/nodeN/hugepages/hugepages-*kB/nr_hugepages` |

Existing CPU-list parsing helpers in `pkg/util/hardware` should be reused. The
implementation must handle sparse CPU and NUMA IDs, offline CPUs, systems with a
single node and no explicit NUMA configuration, missing optional files, and
architectures with different topology exposure.

Only NUMA cell membership, CPU IDs, SMT siblings, and distances are used by the
current placement paths. Memory and page metadata are present in the command API
but do not currently drive that placement. Before implementation, we should
either preserve their Libvirt-compatible semantics from sysfs/procfs or formally
make them optional and add tests proving that omission is safe.

Discover topology once during `virt-handler` startup, fail startup when required
CPU/NUMA relationships are inconsistent, and inject the resulting
`*cmdv1.Topology` into `VirtualMachineController` and
`MigrationTargetController`. `virtualMachineOptions` should accept topology
directly instead of a `*libvirtxml.Caps` and should clone it before use if future
callers may mutate the value.

## Implementation plan

1. **Define topology semantics and fixtures**
   - Document which fields are required and which are optional.
   - Add sysfs/procfs fixtures for one-node, multi-node, SMT, sparse-ID,
     offline-CPU, and missing-optional-data hosts.
   - Capture representative `virsh capabilities` output for parity tests.

2. **Add direct host discovery**
   - Add a provider under `pkg/virt-handler` or `pkg/virt-handler/node-labeller`
     with injectable sysfs/procfs roots.
   - Reuse `hardware.ParseCPUSetLine` rather than adding another CPU-list parser.
   - Return deterministic ordering for cells, CPUs, siblings, and distances.

3. **Prove compatibility**
   - Unit test the provider against the fixtures.
   - Convert matching Libvirt XML fixtures with the current converter and assert
     equality with direct discovery for fields with equivalent semantics.
   - Add validation for duplicate CPUs, CPUs assigned to multiple cells,
     sibling references to unavailable CPUs, and malformed distance matrices.

4. **Change controller wiring**
   - Discover topology in `cmd/virt-handler/virt-handler.go`.
   - Pass `*cmdv1.Topology` to the VM and migration-target controllers.
   - Change `virtualMachineOptions` to consume that topology directly.
   - Keep reading Libvirt capabilities temporarily for TSC and supported machine
     data, but stop retaining or passing the full object in these controllers.

5. **Roll out with a parity window**
   - For one development cycle, optionally compare direct discovery with the
     Libvirt-derived topology and report mismatches through logs or a metric.
   - Use the direct provider as the source of truth once architecture CI and
     NUMA end-to-end tests show parity.

6. **Follow up on the remaining capability file users**
   - Move TSC discovery to a host-owned implementation where kernel interfaces
     provide equivalent semantics.
   - Keep machine types and other VMM capabilities behind the default
     virtualization-stack implementation or its future runtime plugin.
   - Remove `virsh capabilities` and `capabilities.xml` only after those users
     have migrated.

## Testing and acceptance criteria

- Dedicated CPU placement produces the same pinning on SMT and non-SMT hosts.
- NUMA guest-mapping passthrough produces the same host-node mapping and distance
  information.
- Source and target topology reporting during migration remains compatible with
  existing and version-skewed launchers.
- CPU and memory hotplug continue to receive topology in
  `VirtualMachineOptions`.
- Unit tests cover single-node systems, multiple NUMA nodes, sparse IDs, offline
  CPUs, partial SMT sibling allocation, and malformed or unavailable sysfs data.
- Architecture CI covers at least amd64, arm64, and s390x behavior.
- `VirtualMachineController`, `MigrationTargetController`, and
  `virtualMachineOptions` no longer depend on `libvirtxml.Caps` for host topology.
- No user-visible API or command API wire change is required.

## Alternatives considered

### Keep using Libvirt capabilities XML

This preserves current behavior but requires every virtualization stack to carry
or emulate a Libvirt/QEMU discovery path. It retains the coupling this issue is
intended to remove.

### Query topology in each virt-launcher

This duplicates node discovery in every VMI pod and may expose only the pod's
restricted CPU set rather than the full host topology. It also places
stack-independent host discovery on the stack-specific side of the proposed
plugin boundary.

### Add topology to the virt-runtime plugin API

Host topology is common to all virtualization stacks on a node. Asking each
plugin to report it creates duplicate implementations and potentially
inconsistent answers. Core `virt-handler` is the natural owner; plugins should
report only stack-specific capabilities and labels.

## Related code

- `cmd/virt-launcher/node-labeller/node-labeller.sh`
- `pkg/virt-operator/resource/generate/components/daemonsets.go`
- `cmd/virt-handler/virt-handler.go`
- `pkg/virt-handler/options.go`
- `pkg/virt-handler/vm.go`
- `pkg/virt-handler/migration-target.go`
- `pkg/handler-launcher-com/cmd/v1/cmd.proto`
- `pkg/virt-launcher/virtwrap/converter/vcpu/vcpu.go`
- `pkg/util/hardware/hw_utils.go`
