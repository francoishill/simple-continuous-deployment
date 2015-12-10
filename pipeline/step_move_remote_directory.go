package pipeline

import (
	"fmt"
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	"github.com/francoishill/simple-continuous-deployment/ssh_rsync"
)

type moveRemoteDirectory struct {
	srcCtx      *PipelineContext
	destCtx     *PipelineContext
	relativeDir string
}

func (m *moveRemoteDirectory) Execute(logger Logger, client ssh_rsync.SSHClient) {
	srcFullPath := m.srcCtx.getFullRemotePath(m.relativeDir)
	destFullPath := m.destCtx.getFullRemotePath(m.relativeDir)

	defer func() {
		if r := recover(); r != nil {
			if !client.DoesPathExist(srcFullPath, true) {
				logger.Info("SSH `MV` command failed but the source dir does not exist so assumed that is the error, here is the actual error: %q", r)
			} else {
				//If the path exists and we caught the error, it is probably another error, so re-throw
				panic(r)
			}
		}
	}()
	client.Execute(fmt.Sprintf(`mv "%s" "%s"`, srcFullPath, destFullPath))
}

func NewMoveRemoteDirectoryStep(srcCtx, destCtx *PipelineContext, relativeDir string) PipelineStep {
	return &moveRemoteDirectory{
		srcCtx,
		destCtx,
		relativeDir,
	}
}
