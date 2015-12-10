package ssh_rsync

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	"github.com/francoishill/golang-web-dry/osutils"
	"strings"
)

type Builder interface {
	User(user string) Builder
	Host(h string) Builder
	Port(p int) Builder
	PrivateKeyPath(p string) Builder
	Build() SSHClient
}

type builder struct {
	sshClient
}

func (b *builder) User(u string) Builder {
	b.user = u
	return b
}

func (b *builder) Host(h string) Builder {
	b.host = h
	return b
}

func (b *builder) Port(p int) Builder {
	b.port = p
	return b
}

func (b *builder) PrivateKeyPath(p string) Builder {
	b.privateKeyPath = p
	return b
}

func (b *builder) Build() SSHClient {
	if strings.TrimSpace(b.user) == "" {
		panic("Please specify a user")
	}
	if strings.TrimSpace(b.host) == "" {
		panic("Please specify a host")
	}
	if strings.TrimSpace(b.privateKeyPath) != "" {
		if !osutils.FileExists(b.privateKeyPath) {
			panic("Private key file does not exist: " + b.privateKeyPath)
		}
	}

	if b.port == 0 {
		b.port = 22 //default
	}

	return &sshClient{
		b.logger,
		b.user,
		b.host,
		b.port,
		b.privateKeyPath,
	}
}

func NewBuilder(logger Logger) Builder {
	b := &builder{}
	b.logger = logger
	return b
}
