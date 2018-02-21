package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"syscall"
)

func main() {
	var configurations []ImageConfiguration
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	ReadDockerConfiguration(currentUser.HomeDir+"/.docker-manager.yaml", &configurations)
	var selectedConfiguration = -1
	if len(os.Args) > 1 {
		for idx := range configurations {
			if configurations[idx].Name == os.Args[1] {
				selectedConfiguration = idx
				break
			}
		}
	}

	if selectedConfiguration < 0 {
		printImageSelectionMenu(configurations)
		selectedConfiguration = getUserChoice(1, len(configurations)) - 1
	}

	var selectedAction = -1
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "start":
			selectedAction = startContainerAction
		case "connect":
			selectedAction = connectToContainerAction
		case "save":
			selectedAction = saveContainerAction
		case "stop":
			selectedAction = stopContainerAction
		case "show":
			selectedAction = showContainerConfiguration
		}
	}
	if selectedAction < 0 {
		printActionSelectionMenu()
		selectedAction = getUserChoice(1, 5) - 1
	}

	switch selectedAction {
	case startContainerAction:
		startContainer(&configurations[selectedConfiguration])
	case connectToContainerAction:
		connectToContainer(&configurations[selectedConfiguration])
	case saveContainerAction:
		saveContainer(&configurations[selectedConfiguration])
	case stopContainerAction:
		stopContainer(&configurations[selectedConfiguration])
	case showContainerConfiguration:
		configurations[selectedConfiguration].Print()
	}
}

func printImageSelectionMenu(configurations []ImageConfiguration) {
	fmt.Println("Please select an image:")
	for idx := range configurations {
		fmt.Print("\t")
		fmt.Print(idx + 1)
		fmt.Println(") " + configurations[idx].Name)
	}
}

func printActionSelectionMenu() {
	fmt.Println("Please select an action for this container:")
	fmt.Println("\t1) Start")
	fmt.Println("\t2) Connect")
	fmt.Println("\t3) Save")
	fmt.Println("\t4) Stop")
	fmt.Println("\t5) Show")
}

func getUserChoice(min, max int) int {
	var inputOk = false
	var selection int
	for !inputOk {
		fmt.Print(">>> ")
		fmt.Scanf("%d", &selection)
		if selection >= min && selection <= max {
			inputOk = true
		}
	}
	return selection
}

const (
	startContainerAction       = iota
	connectToContainerAction   = iota
	saveContainerAction        = iota
	stopContainerAction        = iota
	showContainerConfiguration = iota
)

func readContainerID(conf *ImageConfiguration) []byte {
	var id []byte
	if _, err := os.Stat(conf.GetIDPath()); !os.IsNotExist(err) {
		var readError error
		id, readError = ioutil.ReadFile(conf.GetIDPath())
		if readError != nil {
			log.Fatal("Cannot read file ", conf.GetIDPath(), ": ", readError)
		}
		if len(id) > 0 {
			id = id[:len(id)-1]
		}
	}
	return id
}

func startContainer(conf *ImageConfiguration) {
	if len(readContainerID(conf)) > 0 {
		fmt.Println("The container is already started")
		return
	}
	if conf.Gui {
		xhostError := exec.Command("xhost", []string{"+"}...).Run()
		if xhostError != nil {
			log.Fatal("Cannot disable xhost acces control: ", xhostError)
		}
	}
	var params = conf.GenerateStartCommand()
	cmd := exec.Command(params[0], params[1:]...)
	fmt.Println("Starting the container " + conf.Image + ":" + conf.Tag)
	runError := cmd.Run()
	if runError != nil {
		log.Fatal("Unable to start the container: ", runError)
	}
	out, psError := exec.Command("docker", []string{"ps", "-q", "-l"}...).Output()
	if psError != nil {
		log.Fatal("Unable to get the container ID: ", psError)
	}
	writeError := ioutil.WriteFile(conf.GetIDPath(), out, 0644)
	if writeError != nil {
		log.Fatal("Unable to save the container ID: ", writeError)
	}
}

func connectToContainer(conf *ImageConfiguration) {
	id := readContainerID(conf)
	if len(id) == 0 {
		startContainer(conf)
		id = readContainerID(conf)
	}

	dockerPath, lookError := exec.LookPath("docker")
	if lookError != nil {
		log.Fatal("The docker executable cannot be found in PATH")
	}
	runError := syscall.Exec(dockerPath, []string{"docker", "exec", "-ti", string(id), conf.Shell}, os.Environ())
	if runError != nil {
		log.Fatal("Unable to connect to the container: ", runError)
	}
}

func saveContainer(conf *ImageConfiguration) {
	id := readContainerID(conf)
	if len(id) == 0 {
		fmt.Println("The container is not running")
		return
	}

	commitError := exec.Command("docker", []string{"commit", string(id), conf.GetImageWithSaveTag()}...).Run()
	if commitError != nil {
		log.Fatal("Unable to save the container state: ", commitError)
	}
}

func stopContainer(conf *ImageConfiguration) {
	id := readContainerID(conf)
	if len(id) == 0 {
		fmt.Println("The container is not running")
		return
	}

	if conf.Autosave {
		saveContainer(conf)
	}

	stopError := exec.Command("docker", []string{"kill", string(id)}...).Run()
	if stopError != nil {
		log.Fatal("Unable to stop the container: ", stopError)
	}

	writeError := ioutil.WriteFile(conf.GetIDPath(), []byte{}, 0644)
	if writeError != nil {
		log.Fatal("Unable to reset the container ID: ", writeError)
	}

	if conf.Gui {
		xhostError := exec.Command("xhost", []string{"-"}...).Run()
		if xhostError != nil {
			log.Fatal("Cannot enable xhost acces control: ", xhostError)
		}
	}
}
