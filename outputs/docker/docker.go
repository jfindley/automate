package docker

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/jfindley/testfs"
)

// NewLayer creates a new layer from a tarball
func NewLayer(tar []byte, parent, comment string) (layer []byte, err error) {
	layerFS := testfs.NewLocalTestFS()

	id, err := newID()
	if err != nil {
		return
	}

	err = layerFS.Mkdir(id, 0755)
	if err != nil {
		return
	}

	layerConfig := newLayerConfig(id, parent, comment)
	// layerConfig.ContainerConfig().Cmd = []string{"/bin/sh", "-c", fmt.Sprintf("#(squash) from %s", parent[:12])}
	entry := &exportedImage{
		Path:         id,
		JSONPath:     filepath.Join(id, "json"),
		VersionPath:  filepath.Join(id, "VERSION"),
		LayerTarPath: filepath.Join(id, "layer.tar"),
		LayerDirPath: filepath.Join(id, "layer"),
		LayerConfig:  layerConfig,
	}
	entry.LayerConfig.Created = time.Now().UTC()

	err = entry.writeJSON(layerFS)
	if err != nil {
		return
	}

	err = entry.writeVersion(layerFS)
	if err != nil {
		return
	}

	err = entry.writeData(layerFS, tar)
	if err != nil {
		return
	}

	return archive(layerFS)
}

// Port is a convient way of dealing with docker ports
type Port string

// Port returns the number of the port.
func (p Port) Port() string {
	return strings.Split(string(p), "/")[0]
}

// Proto returns the name of the protocol.
func (p Port) Proto() string {
	parts := strings.Split(string(p), "/")
	if len(parts) == 1 {
		return "tcp"
	}
	return parts[1]
}

type exportedImage struct {
	Path         string
	JSONPath     string
	VersionPath  string
	LayerTarPath string
	LayerDirPath string
	LayerConfig  *LayerConfig
}

func newLayerConfig(id, parent, comment string) *LayerConfig {
	return &LayerConfig{
		Id:            id,
		Parent:        parent,
		Comment:       comment,
		Created:       time.Now().UTC(),
		DockerVersion: "0.1.2",
		Architecture:  "x86_64",
	}
}

func (e *exportedImage) writeVersion(fs testfs.FileSystem) (err error) {
	fp := e.VersionPath
	f, err := fs.Create(fp)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.WriteString("1.0")
	if err != nil {
		return
	}

	return
}

func (e *exportedImage) writeJSON(fs testfs.FileSystem) (err error) {
	fp := e.JSONPath
	f, err := fs.Create(fp)
	if err != nil {
		return
	}
	defer f.Close()

	jb, err := json.Marshal(e.LayerConfig)
	if err != nil {
		return
	}

	_, err = f.WriteString(string(jb))
	if err != nil {
		return
	}

	return err
}

func (e *exportedImage) writeData(fs testfs.FileSystem, data []byte) (err error) {
	fp := e.LayerTarPath
	f, err := fs.Create(fp)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return
	}

	return err
}
