package shell

type InstanceType struct {
	ServerName           string
	Hypervisor           string
	Cpus                 int
	Memory               int64
    DiskSize             int
	OsName               string
	Network 			 string
}
