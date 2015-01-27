package shell

import (
	"fmt"
	"log"
)
// Artifact represents a Softlayer image as the result of a Packer build.
type Artifact struct {
	uuid string
}

// uuid returns the uuid Id.
func (*Artifact) uuid() string {
	return self.uuid
}

// Files returns the files represented by the artifact.
func (*Artifact) Files() []string {
	return nil
}

func (self *Artifact) State(name string) interface{} {
	return nil
}

// String returns the string representation of the artifact.
func (self *Artifact) String() string {
	return fmt.Sprintf("%s", self.uuid)
}
