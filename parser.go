package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// MappedFolder holds a host/container folders path pair
type MappedFolder struct {
	Host      string
	Container string
}

// AddDropCapabilities holds the capabilities to add/drop to a container
type AddDropCapabilities struct {
	Add  []string
	Drop []string
}

// ImageConfiguration config parameter for a given docker image.
// Fields must start with an uppercase letter
type ImageConfiguration struct {
	Name         string
	Image        string
	Tag          string
	SaveTag      string `yaml:"save_tag"`
	Runtime      string
	Network      string
	Shell        string
	ExtraFlags   string `yaml:"extra_flags"`
	Autosave     bool
	Privileged   bool
	Gui          bool
	Folders      []MappedFolder
	Capabilities AddDropCapabilities
}

// Print outputs all the fields to the standard output
func (configuration *ImageConfiguration) Print() {
	fmt.Println("Name:", configuration.Name)
	fmt.Println("\tImage:", configuration.Image)
	fmt.Println("\tTag:", configuration.Tag)
	fmt.Println("\tSaveTag:", configuration.SaveTag)
	fmt.Println("\tRuntime:", configuration.Runtime)
	fmt.Println("\tNetwork:", configuration.Network)
	fmt.Println("\tShell:", configuration.Shell)
	fmt.Println("\tAutosave:", configuration.Autosave)
	fmt.Println("\tPrivileged:", configuration.Privileged)
	fmt.Println("\tGui:", configuration.Gui)
	fmt.Println("\tFolders:", configuration.Folders)
	fmt.Println("\tCapabilities:", configuration.Capabilities)
	fmt.Println("\tExtraFlags:", configuration.ExtraFlags)
}

// GenerateStartCommand generates the appropriate command to start the container with the given configuration
func (configuration *ImageConfiguration) GenerateStartCommand() []string {
	var command []string
	push := func(arg string) {
		command = append(command, arg)
	}
	push("docker")
	push("run")
	push("-ti")
	push("-d")

	if configuration.Runtime != "none" {
		push("--runtime=" + configuration.Runtime)
	}

	push("--network=" + configuration.Network)

	if configuration.Privileged {
		push("--privileged")
	}

	if configuration.Gui {
		push("--env")
		push("DISPLAY")
		push("--env")
		push("QT_X11_NO_MITSHM=1")
		push("--volume=/dev/video0:/dev/video0")
		push("--volume=/tmp/.X11-unix:/tmp/.X11-unix:ro")
	}

	for idx := range configuration.Folders {
		push("--volume=" + configuration.Folders[idx].Host + ":" + configuration.Folders[idx].Container)
	}

	for idx := range configuration.Capabilities.Add {
		push("--cap-add=" + configuration.Capabilities.Add[idx])
	}

	for idx := range configuration.Capabilities.Drop {
		push("--cap-drop=" + configuration.Capabilities.Drop[idx])
	}

    if len(configuration.ExtraFlags) > 0 {
        tokens := strings.Split(configuration.ExtraFlags, " ")
        for tokenIdx := range tokens {
            push(tokens[tokenIdx])
        }
    }

	push(configuration.GetImageWithTag())

	return command
}

// GetIDPath provides the full path to the ID file location
func (configuration *ImageConfiguration) GetIDPath() string {
	return "/tmp/" + configuration.Name + "_master_id"
}

// GetImageWithTag provides the full image expression (e.g image:tag)
func (configuration *ImageConfiguration) GetImageWithTag() string {
	return configuration.Image + ":" + configuration.Tag
}

// GetImageWithSaveTag provides the full image expression using the save tag (e.g image:save_tag)
func (configuration *ImageConfiguration) GetImageWithSaveTag() string {
	return configuration.Image + ":" + configuration.SaveTag
}

// ReadDockerConfiguration parses a YAML configuration file and populates the configurations array
func ReadDockerConfiguration(file string, configurations *[]ImageConfiguration) {
	source, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, configurations)
	if err != nil {
		panic(err)
	}
	for i := range *configurations {
		var conf = &(*configurations)[i]
		if conf.Name == "" {
			panic("The 'name' field is mandatory and must be a string")
		}
		// Remove all white spaces
		conf.Name = strings.Join(strings.Fields(conf.Name), "")
		if conf.Image == "" {
			panic("The 'image' field is mandatory and must be a string")
		}
		if conf.Tag == "" {
			conf.Tag = "latest"
		}
		if conf.SaveTag == "" {
			conf.SaveTag = conf.Tag
		}
		if conf.Runtime == "" {
			conf.Runtime = "none"
		}
		if conf.Network == "" {
			conf.Network = "bridge"
		}
		if conf.Shell == "" {
			conf.Shell = "bash"
		}
	}
}
