package pipeline

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	"github.com/francoishill/simple-continuous-deployment/ssh_rsync"
)

type PipelineStep interface {
	Execute(logger Logger, client ssh_rsync.SSHClient)
}
