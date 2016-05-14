package docker

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"time"
)

type ContainerConfig struct {
	Hostname        string
	Domainname      string
	Entrypoint      []string
	User            string
	Memory          int64
	MemorySwap      int64
	CpuShares       int64
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	PortSpecs       []string
	Tty             bool
	OpenStdin       bool
	StdinOnce       bool
	NetworkDisabled bool
	OnBuild         []string
	Env             []string
	Cmd             []string
	Dns             []string
	Image           string
	Volumes         map[string]struct{}
	VolumesFrom     string
	Labels          map[string]string
}

type Config struct {
	Hostname        string
	Domainname      string
	User            string
	Memory          int64
	MemorySwap      int64
	CpuShares       int64
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	PortSpecs       []string
	ExposedPorts    map[Port]struct{}
	OnBuild         []string
	Tty             bool
	OpenStdin       bool
	StdinOnce       bool
	Env             []string
	Cmd             []string
	Dns             []string // For Docker API v1.9 and below only
	Image           string
	Volumes         map[string]struct{}
	VolumesFrom     string
	WorkingDir      string
	Entrypoint      []string
	NetworkDisabled bool
	Labels          map[string]string
}

type LayerConfig struct {
	Id     string `json:"id"`               // Randomly generated, 256-bit, hexadecimal encoded. Uniquely identifies the image.
	Parent string `json:"parent,omitempty"` /* ID of the parent image. If there is no parent image then this field should be
	   omitted. A collection of images may share many of the same ancestor layers.
	   This organizational structure is strictly a tree with any one layer having
	   either no parent or a single parent and zero or more descendant layers. Cycles
	   are not allowed and implementations should be careful to avoid creating them or
	   iterating through a cycle indefinitely. */
	Comment           string           `json:"comment"`
	Created           time.Time        `json:"created"`                    //ISO-8601 formatted combined date and time at which the image was created.
	V1ContainerConfig *ContainerConfig `json:"ContainerConfig,omitempty"`  // Docker 1.0.0, 1.0.1
	V2ContainerConfig *ContainerConfig `json:"container_config,omitempty"` // All other versions
	Container         string           `json:"container"`
	Config            *Config          `json:"config,omitempty"`
	DockerVersion     string           `json:"docker_version"`
	Architecture      string           `json:"architecture"`
}

// ContainerConfig is a version-independent way to get the configuration from a layer.
func (l *LayerConfig) ContainerConfig() *ContainerConfig {
	if l.V2ContainerConfig != nil {
		return l.V2ContainerConfig
	}

	// If the exports use the 1.0.x json field name, convert it to the newer field
	// name which appears to work in all versions.
	if l.V1ContainerConfig != nil {
		l.V2ContainerConfig = l.V1ContainerConfig
		l.V1ContainerConfig = nil
		return l.V2ContainerConfig
	}

	l.V2ContainerConfig = &ContainerConfig{}

	return l.V2ContainerConfig
}

/*
Each layer is given an ID upon its creation. It is represented as a hexadecimal
encoding of 256 bits, e.g.,
a9561eb1b190625c9adb5a9513e72c4dedafc1cb2d4c5236c9a6957ec7dfd5a9. Image IDs
should be sufficiently random so as to be globally unique. 32 bytes read from
/dev/urandom is sufficient for all practical purposes. Alternatively, an image
ID may be derived as a cryptographic hash of image contents as the result is
considered indistinguishable from random. The choice is left up to implementors.
*/

func newID() (string, error) {
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return "", err
	}
	return hex.EncodeToString(id), nil
}
