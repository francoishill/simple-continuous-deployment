package pipeline

import (
	"fmt"
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	"github.com/francoishill/simple-continuous-deployment/ssh_rsync"
)

type setRemotePermissions struct {
	ctx                *PipelineContext
	dirMode            bool
	relativeRemotePath string
	permissionString   string
}

func (s *setRemotePermissions) Execute(logger Logger, client ssh_rsync.SSHClient) {
	var fileOrFolderFlag string
	if s.dirMode {
		fileOrFolderFlag = "d"
	} else {
		fileOrFolderFlag = "f"
	}

	remoteFull := s.ctx.getFullRemotePath(s.relativeRemotePath)

	cmd := fmt.Sprintf(`find "%s" -type %s -exec chmod %s {} \;`, remoteFull, fileOrFolderFlag, s.permissionString)
	client.Execute(cmd)
}

func NewSetRemotePermissionsStep(ctx *PipelineContext, dirMode bool, relativeRemotePath, permissionString string) PipelineStep {
	return &setRemotePermissions{
		ctx,
		dirMode,
		relativeRemotePath,
		permissionString,
	}
}
