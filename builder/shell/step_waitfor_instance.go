package shell

import (
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"os/exec"
	"strings"
	"bytes"
)

type stepWaitforInstance struct{}

func (self *stepWaitforInstance) Run(state multistep.StateBag) multistep.StepAction {
	uuid := state.Get("server_uuid").(string)
	//config := state.Get("config").(config)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Waiting for IP of server...")

	// FIXME - loop ontil timeout
	cmd := exec.Command("tr", "a-z", "A-Z")
	cmd.Stdin = strings.NewReader(uuid)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}
	state.Put("IpAddress", out)
	return multistep.ActionContinue
}

func (self *stepWaitforInstance) Cleanup(state multistep.StateBag) {
}
