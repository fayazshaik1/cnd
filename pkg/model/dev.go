package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	// CNDLabel is the label added to a dev deployment in k8
	CNDLabel = "cnd.okteto.com/deployment"

	// CNDDeploymentAnnotation is the original deployment manifest
	CNDDeploymentAnnotation = "cnd.okteto.com/manifest"

	cndInitSyncContainerTemplate = "cnd-init-%s"
	cndSyncContainerTemplate     = "cnd-%s"
	cndSyncVolumeTemplate        = "cnd-data-%s"
	cndSyncMountTemplate         = "/var/cnd-sync/%s"
)

//Dev represents a cloud native development environment
type Dev struct {
	Swap    Swap              `yaml:"swap"`
	Mount   Mount             `yaml:"mount"`
	Scripts map[string]string `yaml:"scripts"`
}

//Swap represents the metadata for the container to be swapped
type Swap struct {
	Deployment Deployment `yaml:"deployment"`
}

//Deployment represents the container to be swapped
type Deployment struct {
	Name      string   `yaml:"name"`
	Container string   `yaml:"container,omitempty"`
	Image     string   `yaml:"image"`
	Command   []string `yaml:"command,omitempty"`
	Args      []string `yaml:"args,omitempty"`
}

//Mount represents how the local filesystem is mounted
type Mount struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

//NewDev returns a new instance of dev with default values
func NewDev() *Dev {
	return &Dev{
		Swap: Swap{
			Deployment: Deployment{},
		},
		Mount: Mount{
			Source: ".",
			Target: "/app",
		},
		Scripts: make(map[string]string),
	}
}

func (dev *Dev) validate() error {
	file, err := os.Stat(dev.Mount.Source)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("Source mount folder %s does not exists", dev.Mount.Source)
	}
	if !file.Mode().IsDir() {
		return fmt.Errorf("Source mount folder is not a directory")
	}

	if dev.Swap.Deployment.Name == "" {
		return fmt.Errorf("Swap deployment name cannot be empty")
	}

	return nil
}

//ReadDev returns a Dev object from a given file
func ReadDev(devPath string) (*Dev, error) {
	b, err := ioutil.ReadFile(devPath)
	if err != nil {
		return nil, err
	}

	d, err := loadDev(b)
	if err != nil {
		return nil, err
	}

	if err := d.validate(); err != nil {
		return nil, err
	}

	d.fixPath(devPath)
	return d, nil
}

func loadDev(b []byte) (*Dev, error) {
	dev := Dev{
		Mount: Mount{
			Source: ".",
			Target: "/src",
		},
	}

	err := yaml.Unmarshal(b, &dev)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(dev.Mount.Source, "~/") {
		home := os.Getenv("HOME")
		dev.Mount.Source = filepath.Join(home, dev.Mount.Source[2:])
	}

	return &dev, nil
}

func (dev *Dev) fixPath(originalPath string) {
	wd, _ := os.Getwd()

	if !filepath.IsAbs(dev.Mount.Source) {
		if filepath.IsAbs(originalPath) {
			dev.Mount.Source = path.Join(path.Dir(originalPath), dev.Mount.Source)
		} else {

			dev.Mount.Source = path.Join(wd, path.Dir(originalPath), dev.Mount.Source)
		}
	}
}

// GetCNDInitSyncContainer returns the CND init sync container name for a given container
func (dev *Dev) GetCNDInitSyncContainer() string {
	return fmt.Sprintf(cndInitSyncContainerTemplate, dev.Swap.Deployment.Container)
}

// GetCNDSyncContainer returns the CND sync container name for a given container
func (dev *Dev) GetCNDSyncContainer() string {
	return fmt.Sprintf(cndSyncContainerTemplate, dev.Swap.Deployment.Container)
}

// GetCNDSyncVolume returns the CND sync volume name for a given container
func (dev *Dev) GetCNDSyncVolume() string {
	return fmt.Sprintf(cndSyncVolumeTemplate, dev.Swap.Deployment.Container)
}

// GetCNDSyncMount returns the CND sync mount for a given container
func (dev *Dev) GetCNDSyncMount() string {
	return fmt.Sprintf(cndSyncMountTemplate, dev.Swap.Deployment.Container)
}
