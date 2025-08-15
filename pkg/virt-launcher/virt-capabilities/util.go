package virt_capabilities

const (
	DefaultMinCPUModel = "Penryn"
	RequirePolicy      = "require"
	// TODO Check if KVM Path is required here
	KVMPath    = "/dev/kvm"
	VmxFeature = "vmx"
)

// TODO Check which ones of these are required here
