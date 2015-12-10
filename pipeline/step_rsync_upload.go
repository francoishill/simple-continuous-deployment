package pipeline

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	"github.com/francoishill/simple-continuous-deployment/ssh_rsync"
	"path/filepath"
)

type rsyncUpload struct {
	ctx                *PipelineContext
	isDir              bool
	relativeLocalPath  string
	relativeRemotePath string
	flags              []string
}

func (u *rsyncUpload) Execute(logger Logger, client ssh_rsync.SSHClient) {
	localFull := u.ctx.getFullLocalPath(u.relativeLocalPath)
	remoteFull := u.ctx.getFullRemotePath(u.relativeRemotePath)

	if u.isDir {
		client.Execute(`mkdir -p "` + remoteFull + `"`)
	} else {
		dirOfFile := filepath.Dir(remoteFull)
		client.Execute(`mkdir -p "` + dirOfFile + `"`)
	}

	client.RsyncUpload(localFull, remoteFull, u.flags...)
}

func NewRsyncUploadStep(ctx *PipelineContext, isDir bool, relativeLocalPath, relativeRemotePath string, flags ...string) PipelineStep {
	return &rsyncUpload{
		ctx,
		isDir,
		relativeLocalPath,
		relativeRemotePath,
		flags,
	}
}
