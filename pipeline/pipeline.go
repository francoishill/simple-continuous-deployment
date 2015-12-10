package pipeline

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	"github.com/francoishill/simple-continuous-deployment/ssh_rsync"
)

type Pipeline interface {
	ExecuteAll(logger Logger, client ssh_rsync.SSHClient)
}

type pipeline struct {
	steps []PipelineStep
}

func (p *pipeline) ExecuteAll(logger Logger, client ssh_rsync.SSHClient) {
	for _, s := range p.steps {
		s.Execute(logger, client)
	}
}

func NewPipeline(steps []PipelineStep) Pipeline {
	return &pipeline{
		steps,
	}
}
