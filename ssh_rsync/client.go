package ssh_rsync

import (
	"fmt"
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

type SSHClient interface {
	Execute(cmd string) string
	DoesPathExist(path string, dir bool) bool
	ExecuteLocalRsync(rsyncArgs ...string) string
	RsyncUpload(localPath, remotePath string, flags ...string) string
}

type sshClient struct {
	logger         Logger
	user           string
	host           string
	port           int
	privateKeyPath string
}

func (s *sshClient) getKeyFile() (key ssh.Signer) {
	buf, err := ioutil.ReadFile(s.privateKeyPath)
	CheckError(err)
	key, err = ssh.ParsePrivateKey(buf)
	CheckError(err)
	return
}

func (s *sshClient) connectSession() *ssh.Session {
	config := &ssh.ClientConfig{
		User: s.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(s.getKeyFile()),
		},
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.logger.Debug("Connecting on address: %s", addr)
	client, err := ssh.Dial("tcp", addr, config)
	CheckError(err)

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	CheckError(err)
	return session
}

func (s *sshClient) Execute(cmd string) string {
	session := s.connectSession()
	defer session.Close()

	timestamp := time.Now().UnixNano()

	s.logger.Debug("Executing (timestamp %d) cmd: %s", timestamp, cmd)
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		panic("ERROR, output: " + string(output) + ", error: " + err.Error())
	}

	s.logger.Debug("Output of execution (timestamp %d): %s", timestamp, output)
	return string(output)
}

func (s *sshClient) DoesPathExist(path string, dir bool) bool {
	var modeStr string
	if dir {
		modeStr = "-d"
	} else {
		modeStr = "-f"
	}
	output := s.Execute(fmt.Sprintf(`if [ %s "%s" ]; then echo yes; else echo no; fi`, modeStr, path))
	return strings.EqualFold("yes", strings.TrimSpace(output))
}

func (s *sshClient) ExecuteLocalRsync(rsyncArgs ...string) string {
	allArgs := []string{
		"-e",
		fmt.Sprintf("ssh -i %s -p %d", s.privateKeyPath, s.port),
	}
	allArgs = append(allArgs, rsyncArgs...)

	timestamp := time.Now().UnixNano()

	s.logger.Debug("Executing rsync (timestamp %d) with args: %q", timestamp, allArgs)
	output, err := exec.Command("rsync", allArgs...).CombinedOutput()
	if err != nil {
		panic("ERROR, output (local rsync): " + string(output) + ", error: " + err.Error())
	}
	s.logger.Debug("Output of rsync execution (timestamp %d): %s", timestamp, output)
	return string(output)
}

func (s *sshClient) RsyncUpload(localPath, remotePath string, flags ...string) string {
	allArgs := []string{}
	allArgs = append(allArgs, flags...)
	allArgs = append(allArgs, localPath)
	allArgs = append(allArgs, fmt.Sprintf("%s@%s:%s", s.user, s.host, remotePath))
	return s.ExecuteLocalRsync(allArgs...)
}
