package shell

import (
	"fmt"
	"log"
)
// Artifact represents a Softlayer image as the result of a Packer build.
type Artifact struct {
	uuid string
}

// BuilderId returns the builder Id.
func (*Artifact) BuilderId() string {
	return BuilderId
}

// Destroy destroys the Softlayer image represented by the artifact.
func (self *Artifact) Destroy() error {
	log.Printf("Destroying image: %s", self.String())
	return nil
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
