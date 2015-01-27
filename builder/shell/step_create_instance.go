package shell

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
 	"os/exec"
 	"strings"
 	"bytes"
)

type stepCreateInstance struct {
	instanceId string
}

func (self *stepCreateInstance) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(config)
	ui := state.Get("ui").(packer.Ui)

	instanceDefinition := &InstanceType{
		ServerName:           config.ServerName,
		Hypervisor:           config.Hypervisor,
		Cpus:                 config.InstanceCpu,
		Memory:               config.InstanceMemory,
        DiskSize:             config.InstanceDiskSize,
		OsName:               config.OsName,
		Network:			  config.Network,
	}

	ui.Say("Creating an instance...")

	cmd := exec.Command("tr", "a-z", "A-Z")
	cmd.Stdin = strings.NewReader("1234uuid")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}
	ui.Say(fmt.Printf("in all caps: %q\n", out.String()))
	state.Put("uuid", out.String())
	return multistep.ActionContinue
}

func (self *stepCreateInstance) Cleanup(state multistep.StateBag) {
	config := state.Get("config").(config)
	ui := state.Get("ui").(packer.Ui)

	if self.instanceId == "" {
		return
	}

	ui.Say("Destroying instance...")
}
