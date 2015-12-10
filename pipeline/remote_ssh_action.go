package pipeline

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	"github.com/francoishill/simple-continuous-deployment/ssh_rsync"
)

type remoteSSHCommand struct {
	cmd string
}

func (r *remoteSSHCommand) Execute(logger Logger, client ssh_rsync.SSHClient) {
	client.Execute(r.cmd)
}

func NewRemoteSSHCommandStep(cmd string) PipelineStep {
	return &remoteSSHCommand{
		cmd,
	}
}
