package shell

import (
	"errors"
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"log"
	"time"
)

// The unique ID for this builder.
const BuilderId = "packer.shell"

type config struct {
	common.PackerConfig `mapstructure:",squash"`

	ServerName           string `mapstructure:"server_name"`
	Hypervisor 		     string `mapstructure:"hypervisor"`
	InstanceCpu          int    `mapstructure:"instance_cpu"`
	InstanceMemory       int64  `mapstructure:"instance_memory"`
	InstanceDiskSize	 int    `mapstructure:"instance_disk_size"`
	OsName				 string `mapstructure:"os_name"`
	Network				 string `mapstructure:"network"`
	SshPort              int64  `mapstructure:"ssh_port"`
	SshUserName          string `mapstructure:"ssh_username"`
	SshPrivateKeyFile    string `mapstructure:"ssh_private_key_file"`

	RawSshTimeout   string `mapstructure:"ssh_timeout"`
	RawStateTimeout string `mapstructure:"instance_state_timeout"`

	SshTimeout   time.Duration
	StateTimeout time.Duration

	tpl *packer.ConfigTemplate
}

// Builder represents a Packer Builder.
type Builder struct {
	config config
	runner multistep.Runner
}

// Prepare processes the build configuration parameters.
func (self *Builder) Prepare(raws ...interface{}) (parms []string, retErr error) {
	metadata, err := common.DecodeConfig(&self.config, raws...)
	if err != nil {
		return nil, err
	}

	// Check that there aren't any unknown configuration keys defined
	errs := common.CheckUnusedConfig(metadata)
	if errs == nil {
		errs = &packer.MultiError{}
	}

	self.config.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return nil, err
	}
	self.config.tpl.UserVars = self.config.PackerUserVars

	if self.config.Hypervisor == "" {
		self.config.Hypervisor = "esxi01"
	}

	if self.config.ServerName == "" {
		self.config.ServerName = fmt.Sprintf("packer-shell-%s", time.Now().Unix())
	}

	if self.config.InstanceCpu == 0 {
		self.config.InstanceCpu = 1
	}

	if self.config.InstanceMemory == 0 {
		self.config.InstanceMemory = 1024
	}

	if self.config.InstanceDiskSize == 0 {
		self.config.InstanceDiskSize = 20
	}

	if self.config.OsName == "" {
		self.config.OsName = "RHEL"
	}

	if self.config.Network == "" {
		self.config.Network = "packer"
	}

	if self.config.SshPort == 0 {
		self.config.SshPort = 22
	}

	if self.config.SshUserName == "" {
		self.config.SshUserName = "root"
	}

	if self.config.RawSshTimeout == "" {
		self.config.RawSshTimeout = "5m"
	}

	if self.config.RawStateTimeout == "" {
		self.config.RawStateTimeout = "10m"
	}

	templates := map[string]*string{
		"ssh_timeout":            &self.config.RawSshTimeout,
		"instance_state_timeout": &self.config.RawStateTimeout,
		"ssh_username":           &self.config.SshUserName,
		"ssh_private_key_file":   &self.config.SshPrivateKeyFile,
	}

	for n, ptr := range templates {
		var err error
		*ptr, err = self.config.tpl.Process(*ptr, nil)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("Error processing %s: %s", n, err))
		}
	}

	// Translate date configuration data from string to time format
	sshTimeout, err := time.ParseDuration(self.config.RawSshTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing ssh_timeout: %s", err))
	}
	self.config.SshTimeout = sshTimeout

	stateTimeout, err := time.ParseDuration(self.config.RawStateTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing state_timeout: %s", err))
	}
	self.config.StateTimeout = stateTimeout

	if len(errs.Errors) > 0 {
		retErr = errors.New(errs.Error())
	}

	return nil, retErr
}

// fork script with args
func (self *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {

	// Set up the state which is used to share state between the steps
	state := new(multistep.BasicStateBag)
	state.Put("config", self.config)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Build the steps
	steps := []multistep.Step{
		new(stepCreateInstance),
		new(stepWaitforInstance),
		&common.StepConnectSSH{
			SSHAddress:     sshAddress,
			SSHConfig:      sshConfig,
			SSHWaitTimeout: self.config.SshTimeout,
		},
		new(common.StepProvision),
	}

	// Create the runner which will run the steps we just build
	self.runner = &multistep.BasicRunner{Steps: steps}
	self.runner.Run(state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	if _, ok := state.GetOk("image_id"); !ok {
		log.Println("Failed to find image_id in state. Bug?")
		return nil, nil
	}

	// Create an artifact and return it
	artifact := &Artifact{
	uuid: state.Get("uuid").(string),
	}

	return artifact, nil
}

// Cancel.
func (self *Builder) Cancel() {
	if self.runner != nil {
		log.Println("Cancelling the step runner...")
		self.runner.Cancel()
	}
	fmt.Println("Canceling the builder")
}
