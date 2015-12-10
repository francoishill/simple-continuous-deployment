package main

import (
	"github.com/francoishill/golang-common-ddd/Implementations/Logger/DefaultLogger"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/francoishill/simple-continuous-deployment/pipeline"
	"github.com/francoishill/simple-continuous-deployment/ssh_rsync"
)

type setup struct {
	BaseLocalDir string `json:"base_local_dir"`

	SshPrivateKeyPath string `json:"ssh_private_key_path"`
	RemoteUser        string `json:"remote_user"`
	RemoteHost        string `json:"remote_host"`
	RemotePort        int    `json:"remote_port"`
	BaseRemoteDir     string `json:"base_remote_dir"`

	StaticDirs  []string `json:"static_dirs"`
	StaticFiles []string `json:"static_files"`
	ExecFiles   []string `json:"exec_files"`

	StopCmd  string `json:"stop_cmd"`
	StartCmd string `json:"start_cmd"`
}

func (s *setup) Validate() {
	if strings.TrimSpace(s.BaseRemoteDir) == "" {
		panic("Please specify BaseRemoteDir")
	}
}

func loadSetup() *setup {
	if len(os.Args) < 2 {
		log.Fatal("Command-line argument missing for yaml file")
	}

	ymlSetupFile := os.Args[1]

	yamlBytes, err := ioutil.ReadFile(ymlSetupFile)
	CheckError(err)

	s := &setup{}
	err = yaml.Unmarshal(yamlBytes, s)
	CheckError(err)
	s.Validate()
	return s
}

func main() {
	logger := DefaultLogger.New("rolling.log", "[SCD]", true)
	defer func() {
		if r := recover(); r != nil {
			logger.Emergency("Application startup failed: %+v", r)
		}
	}()

	s := loadSetup()

	client := ssh_rsync.NewBuilder(logger).
		PrivateKeyPath(s.SshPrivateKeyPath).
		User(s.RemoteUser).Host(s.RemoteHost).Port(s.RemotePort).
		Build()

	pendingCtx := &pipeline.PipelineContext{
		LocalBaseDir:  s.BaseLocalDir,
		RemoteBaseDir: s.BaseRemoteDir + "_pending",
	}

	backupCtx := &pipeline.PipelineContext{
		LocalBaseDir:  s.BaseLocalDir,
		RemoteBaseDir: s.BaseRemoteDir + "_backup",
	}

	finalCtx := &pipeline.PipelineContext{
		LocalBaseDir:  s.BaseLocalDir,
		RemoteBaseDir: s.BaseRemoteDir,
	}

	dirFlags := "-vrz"
	fileFlags := "-vz"

	steps := []pipeline.PipelineStep{}

	steps = append(steps, pipeline.NewForceRemoveRemoteDirectoryStep(pendingCtx, true, "/"))

	for _, d := range s.StaticDirs {
		steps = append(steps, pipeline.NewRsyncUploadStep(pendingCtx, true, d, d, dirFlags))
	}
	for _, f := range s.StaticFiles {
		steps = append(steps, pipeline.NewRsyncUploadStep(pendingCtx, false, f, f, fileFlags))
	}

	steps = append(steps, pipeline.NewSetRemotePermissionsStep(pendingCtx, true, "/", "755"))
	steps = append(steps, pipeline.NewSetRemotePermissionsStep(pendingCtx, false, "/", "645"))

	for _, ef := range s.ExecFiles {
		steps = append(steps, pipeline.NewSetRemotePermissionsStep(pendingCtx, false, ef, "755"))
	}

	steps = append(steps, pipeline.NewForceRemoveRemoteDirectoryStep(backupCtx, true, "/"))

	steps = append(steps, pipeline.NewRemoteSSHCommandStep(s.StopCmd))
	steps = append(steps, pipeline.NewMoveRemoteDirectoryStep(finalCtx, backupCtx, "/"))
	steps = append(steps, pipeline.NewMoveRemoteDirectoryStep(pendingCtx, finalCtx, "/"))
	steps = append(steps, pipeline.NewRemoteSSHCommandStep(s.StartCmd))

	pipeline := pipeline.NewPipeline(steps)
	pipeline.ExecuteAll(logger, client)
}
