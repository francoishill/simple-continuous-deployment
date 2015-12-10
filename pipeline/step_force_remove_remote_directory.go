package pipeline

import (
	"fmt"
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	"github.com/francoishill/simple-continuous-deployment/ssh_rsync"
)

type forceRemoveRemoteDirectory struct {
	ctx                *PipelineContext
	recursive          bool
	relativeRemotePath string
}

func (r *forceRemoveRemoteDirectory) Execute(logger Logger, client ssh_rsync.SSHClient) {
	remoteFull := r.ctx.getFullRemotePath(r.relativeRemotePath)

	//Just a security check
	if remoteFull == "" || remoteFull == "./" || remoteFull == "/" {
		panic("Invalid path to remove forcfully, security check failed: " + remoteFull)
	}

	if r.recursive {
		client.Execute(fmt.Sprintf(`rm -rf "%s"`, remoteFull))
	} else {
		client.Execute(fmt.Sprintf(`rm "%s"`, remoteFull))
	}
}

func NewForceRemoveRemoteDirectoryStep(ctx *PipelineContext, recursive bool, relativeRemotePath string) PipelineStep {
	return &forceRemoveRemoteDirectory{
		ctx,
		recursive,
		relativeRemotePath,
	}
}
