# IMPORTANT

### These images are only for testing purposes.

## OpenVMM COM2 socket

The OpenVMM manager always passes `--com2 listen=/var/run/kubevirt/sockets/openvmm-com2.sock`.
This UNIX domain socket lives beside the virt-launcher Command API socket
(`launcher-sock`); a custom launcher shared directory changes the socket directory
accordingly. COM2 remains available even when serial console attachment or
logging is disabled. The existing COM1 console and logging are unchanged.

Guest command execution connects directly to this socket using QEMU GA's
`guest-exec` and `guest-exec-status` commands, returning base64-decoded stdout
and the guest exit code. The guest must run QEMU GA on COM2. Each Exec sends
`guest-exec` directly without a synchronization handshake and keeps the
connection open while awaiting replies; no `socat` dependency is required.
There is no synchronization-based protection against stale serial replies
from previous connections. Exec calls are serialized, and the positive
timeout in seconds covers waiting for COM2, connecting, and receiving results.
A timeout stops waiting; it does not terminate the command in the guest.

An exec probe succeeds only after a reply containing a positive PID, followed
by a status reply with `exited: true` and an explicit `exitcode: 0`. Writing to
the socket alone is not success. Missing or malformed replies, missing exit
codes, and agent errors fail the probe and are reported on its stderr.
This checks the QEMU GA protocol, not the identity of the socket peer.

Launcher log verbosity 4 records socket connection,
request submission, received PIDs, completion, and timeouts. Verbosity 5 also
prints guest-agent JSON replies, polling requests, and decoded stdout. Enable
verbosity 5 only when needed: replies and stdout may contain sensitive data.
Command arguments are not logged.

## OpenVMM VMI validation

`openvmm-validation-policy.yaml` implements the final verdicts in section 0 of
`validation-constraints.md`, with the subsequently approved compatibility
allowances below. It contains a `ValidatingAdmissionPolicy` and a
cluster-wide `Deny` binding for VMI **CREATE** requests.
`openvmm-validation-plugin.yaml` registers its name through the KubeVirt Plugin
API. It does not replace the EFI mutating policy.

**Defaulting compatibility:** the VMI in the repository-root `vmi-example` is
accepted. The following allowances supersede the corresponding original verdicts:

- `firmware.serial`, `cpu.maxSockets`, and `devices.rng` are allowed without extra
  value restrictions.
- `machine.type` may be absent, empty, or `q35`.
- `features.acpi` may be absent or enabled/default; explicitly disabling ACPI is
  rejected.

Admission cannot distinguish injected values from identical explicit requests.
These allowances do not imply implemented guest serial-number configuration,
CPU hotplug, RNG attachment, or configurable QEMU machine/ACPI semantics.
`memory.maxGuest` and the other unapproved features remain rejected, including
if a different deployment defaults them. Other example manifests may request
unsupported features beyond this compatibility set.

Apply after reviewing the cluster-wide scope:

```bash
kubectl apply -f examples/openvmm/openvmm-validation-policy.yaml
kubectl apply -f examples/openvmm/openvmm-validation-plugin.yaml
```

The API server enforces the policy through its binding; the Plugin references
the policy rather than directly calling CEL. The Plugins feature and Plugin CRD
must already be available. The admission policy APIs are built-in Kubernetes
APIs, not CRDs.

This policy implements the final verdict list, not every candidate constraint in
the inventory. It allows valid direct kernel boot without requiring EFI, and it
does not reject DataVolumes or inspect referenced PVCs. Empty strings/lists and
disabled boolean options are not active feature requests, except that `bios`
and `initrdPath` are rejected whenever their fields are present.
Sound/TPM objects are rejected even if a TPM explicitly sets `enabled: false`.

Parent VM submission, VMI updates/status, runtime operations, backing PVC mode,
and artifact readiness are outside this CREATE-only policy. Ordinary KubeVirt
validation still applies, including general source/target one-of rules and disk
bus validation. This policy adds the specifically approved VirtIO CD-ROM check,
not a new disk-bus allowlist.

## OpenVMM VirtualMachine validation

`openvmm-vm-validation-policy.yaml` applies the same 32 validation expressions
and compatibility allowances to `spec.template.spec` inside a VirtualMachine's
`spec.template`. It includes a separate cluster-wide `Deny` binding for VM
**CREATE** requests. `openvmm-vm-validation-plugin.yaml` references this policy
through a separate Plugin, so both VM and VMI policies can coexist.

```bash
kubectl apply -f examples/openvmm/openvmm-vm-validation-policy.yaml
kubectl apply -f examples/openvmm/openvmm-vm-validation-plugin.yaml
```

The VM policy examines the template as it exists at VM admission, before the
controller creates a VMI. A mutating policy matching only VMIs cannot fix a
rejected VM template; for example, configure EFI or direct kernel boot in the
template rather than relying on VMI-only mutation to remove explicit BIOS.
Fields supplied only at VMI creation, including instancetype/preference
expansion and VMI defaults, are not resolved by this policy. The existing VMI
policy remains useful for checking the resulting instance. VM updates and
status requests are not matched.
