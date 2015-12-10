package pipeline

import (
	"strings"
)

type PipelineContext struct {
	LocalBaseDir       string
	RemoteBaseDir      string
	ServiceManagerPath string
}

func (p *PipelineContext) getFullLocalPath(relativeLocalPath string) string {
	return strings.TrimSpace(strings.TrimSuffix(p.LocalBaseDir, "/") + "/" + strings.TrimPrefix(relativeLocalPath, "/"))
}

func (p *PipelineContext) getFullRemotePath(relativeRemotePath string) string {
	return strings.TrimSpace(strings.TrimSuffix(p.RemoteBaseDir, "/") + "/" + strings.TrimPrefix(relativeRemotePath, "/"))
}
