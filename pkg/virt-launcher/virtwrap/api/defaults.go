package api

const (
	DefaultProtocol   = "TCP"
	DefaultVMCIDR     = "10.0.2.0/24"
	DefaultVMIpv6CIDR = "fd10:0:2::/120"
	DefaultBridgeName = "k6t-eth0"
)

func NewDefaulter(arch string) *Defaulter {
	return &Defaulter{Architecture: arch}
}

type Defaulter struct {
	Architecture string
}

func (d *Defaulter) IsPPC64() bool {
	if d.Architecture == "ppc64le" {
		return true
	}
	return false
}

func (d *Defaulter) IsARM64() bool {
	if d.Architecture == "arm64" {
		return true
	}
	return false
}

func (d *Defaulter) SetDefaults_Devices(devices *Devices) {

}

func (d *Defaulter) SetDefaults_OSType(ostype *OSType) {
	ostype.OS = "hvm"

	if ostype.Arch == "" {
		if d.IsPPC64() {
			ostype.Arch = "ppc64le"
		} else if d.IsARM64() {
			ostype.Arch = "aarch64"
		} else {
			ostype.Arch = "x86_64"
		}
	}

	// q35 is an alias of the newest q35 machine type.
	// TODO: we probably want to select concrete type in the future for "future-backwards" compatibility.
	if ostype.Machine == "" {
		if d.IsPPC64() {
			ostype.Machine = "pseries"
		} else if d.IsARM64() {
			ostype.Machine = "virt"
		} else {
			ostype.Machine = "q35"
		}
	}
}

func (d *Defaulter) SetDefaults_DomainSpec(spec *DomainSpec) {
	spec.XmlNS = "http://libvirt.org/schemas/domain/qemu/1.0"
	if spec.Type == "" {
		spec.Type = "kvm" // TODO Hermes. Needs to be generalized for MSHV
	}
}

func (d *Defaulter) SetDefaults_SysInfo(sysinfo *SysInfo) {
	if sysinfo.Type == "" {
		sysinfo.Type = "smbios"
	}
}

func (d *Defaulter) SetObjectDefaults_Domain(in *Domain) {
	d.SetDefaults_DomainSpec(&in.Spec)
	d.SetDefaults_OSType(&in.Spec.OS.Type)
	if in.Spec.SysInfo != nil {
		d.SetDefaults_SysInfo(in.Spec.SysInfo)
	}
	d.SetDefaults_Devices(&in.Spec.Devices)
}
