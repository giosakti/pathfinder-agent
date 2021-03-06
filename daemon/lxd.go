package daemon

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	client "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/pathfinder-cm/pathfinder-agent/config"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
)

type apiContainers []api.Container

type LXD struct {
	hostname  string
	localSrv  client.ContainerServer
	targetSrv client.ContainerServer
}

func (a apiContainers) toContainerList() *pfmodel.ContainerList {
	containerList := make(pfmodel.ContainerList, len(a))
	for i, c := range a {
		containerList[i] = pfmodel.Container{
			Hostname: c.Name,
		}
	}
	return &containerList
}

func NewLXD(hostname string, socketPath string) (*LXD, error) {
	localSrv, err := client.ConnectLXDUnix(socketPath, nil)
	if err != nil {
		return nil, err
	}

	// If in clustered mode, specifically target the node
	var targetSrv client.ContainerServer
	if localSrv.IsClustered() {
		targetSrv = localSrv.UseTarget(hostname)
	} else {
		targetSrv = localSrv
	}

	return &LXD{
		hostname:  hostname,
		localSrv:  localSrv,
		targetSrv: targetSrv,
	}, nil
}

func (l *LXD) ListContainers() (*pfmodel.ContainerList, error) {
	res, err := l.targetSrv.GetContainers()
	if err != nil {
		return nil, err
	}

	containerList := apiContainers(res).toContainerList()

	return containerList, nil
}

func (l *LXD) CreateContainer(container pfmodel.Container) (bool, string, error) {
	var certificate string
	if container.Source.Remote.AuthType == "tls" {
		certificate = container.Source.Remote.Certificate
	}

	// Container creation request
	req := api.ContainersPost{
		Name: container.Hostname,
		Source: api.ContainerSource{
			Type:        container.Source.Type,
			Server:      container.Source.Remote.Server,
			Protocol:    container.Source.Remote.Protocol,
			Alias:       container.Source.Alias,
			Mode:        container.Source.Mode,
			Certificate: certificate,
		},
	}

	// Get LXD to create the container (background operation)
	op, err := l.targetSrv.CreateContainer(req)
	if err != nil {
		return false, "", err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, "", err
	}

	// Get LXD to start the container (background operation)
	startReq := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err = l.targetSrv.UpdateContainerState(container.Hostname, startReq, "")
	if err != nil {
		return false, "", err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, "", err
	}

	// Wait for ipaddress to be available
	ipaddress := ""
	found := false
	timeLimit := time.Now().Add(60 * time.Second)

	for !found && time.Now().Before(timeLimit) {
		state, _, err := l.targetSrv.GetContainerState(container.Hostname)
		if err != nil {
			return false, "", err
		}

		addresses := state.Network["eth0"].Addresses
		for _, address := range addresses {
			if address.Family == "inet" {
				ipaddress = address.Address
			}
		}

		if ipaddress != "" {
			found = true
			break
		}

		time.Sleep(time.Duration(3) * time.Second)
	}

	if !found {
		return false, "", errors.New("Missing ip address")
	}

	return true, ipaddress, nil
}

func (l *LXD) DeleteContainer(hostname string) (bool, error) {
	// Get LXD to stop the container (background operation)
	stopReq := api.ContainerStatePut{
		Action:  "stop",
		Timeout: 60,
	}

	// Stop the container
	// We don't care if container already stopped
	op, _ := l.targetSrv.UpdateContainerState(hostname, stopReq, "")
	op.Wait()

	// Get LXD to delete the container (background operation)
	op, err := l.targetSrv.DeleteContainer(hostname)
	if err != nil {
		return false, err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (l *LXD) CreateContainerBootstrapScript(container pfmodel.Container) (bool, error) {
	contentType := "file"
	writeMode := "overwrite"

	for _, bs := range container.Bootstrappers {
		content, mode, err := util.GenerateBootstrapScriptContent(bs)
		if err != nil {
			return false, err
		}

		bootstrapContent := strings.NewReader(content)

		fileArgs := client.ContainerFileArgs{
			Content:   bootstrapContent,
			UID:       0,
			GID:       0,
			Mode:      mode,
			Type:      contentType,
			WriteMode: writeMode,
		}

		err = l.targetSrv.CreateContainerFile(container.Hostname, config.AbsoluteBootstrapScriptPath, fileArgs)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (l *LXD) bootstrapContainer(container pfmodel.Container) (bool, error) {
	commands := []string{
		"bash",
	}
	commands = append(commands, config.AbsoluteBootstrapScriptPath)

	req := api.ContainerExecPost{
		Command:     commands,
		WaitForWS:   true,
		Interactive: false,
	}

	args := client.ContainerExecArgs{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	op, err := l.targetSrv.ExecContainer(container.Hostname, req, &args)
	if err != nil {
		return false, err
	}

	err = op.Wait()
	if err != nil {
		return false, err
	}

	opAPI := op.Get()
	retVal, ok := opAPI.Metadata["return"].(float64)
	if !ok {
		return false, fmt.Errorf("Error Status: Unable to parse exit status")
	}
	if retVal != 0 {
		return false, fmt.Errorf("Error Status: %v when executing bootstrap command", retVal)
	}

	return true, nil
}

func (l *LXD) ValidateAndBootstrapContainer(container pfmodel.Container) (bool, error) {
	var ok bool
	var err error

	maxTry := config.BootstrapContainerMaxRetry + 1

	for maxTry > 0 {
		ok, err = l.bootstrapContainer(container)
		if ok {
			break
		}

		// Add delay between processing
		delay := 5 + util.RandomIntRange(1, 5)
		time.Sleep(time.Duration(delay) * time.Second)

		maxTry--
	}
	return ok, err
}

func (l *LXD) MigrateContainer(container pfmodel.Container) (success bool, ipaddress string, err error) {
	// create params
	migrateReq := api.ContainerPost{
		Name:          container.Hostname,
		Migration:     true,
		Live:          false,
		ContainerOnly: false,
		Target: &api.ContainerPostTarget{
			Certificate: container.Source.Remote.Certificate,
		},
	}

	isRunning, err := l.isContainerRunning(container.Hostname)
	if err != nil {
		return false, "", err
	}

	// stop container first
	if isRunning {
		ok, err := l.stopContainer(container.Hostname)
		if !ok {
			return false, "", err
		}
	}

	// migrate the container
	// don't know if we should stop the container first
	op, err := l.targetSrv.MigrateContainer(container.Hostname, migrateReq)
	if err != nil {
		return false, "", err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, "", err
	}

	// start container
	ok, err := l.startContainer(container.Hostname)
	if !ok {
		return false, "", err
	}

	// Wait for ipaddress to be available
	found := false
	timeLimit := time.Now().Add(time.Duration(config.ContainerRelocationTimeoutInMinute) * time.Minute)

	for !found && time.Now().Before(timeLimit) {
		state, _, err := l.targetSrv.GetContainerState(container.Hostname)
		if err != nil {
			return false, "", err
		}

		addresses := state.Network["eth0"].Addresses
		for _, address := range addresses {
			if address.Family == "inet" {
				ipaddress = address.Address
			}
		}

		if ipaddress != "" {
			found = true
			break
		}

		time.Sleep(time.Duration(5) * time.Second)
	}

	return true, ipaddress, nil
}

func (l *LXD) stopContainer(hostname string) (bool, error) {
	req := api.ContainerStatePut{
		Action:  "stop",
		Timeout: -1,
	}

	op, err := l.targetSrv.UpdateContainerState(hostname, req, "")

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, err
	}

	return true, err
}

func (l *LXD) startContainer(hostname string) (bool, error) {
	req := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err := l.targetSrv.UpdateContainerState(hostname, req, "")

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, err
	}

	return true, err
}

func (l *LXD) isContainerStopped(hostname string) (bool, error) {
	c, _, err := l.targetSrv.GetContainer(hostname)
	if err != nil {
		return false, err
	}

	return c.StatusCode == api.Stopped, nil
}

func (l *LXD) isContainerRunning(hostname string) (bool, error) {
	c, _, err := l.targetSrv.GetContainer(hostname)
	if err != nil {
		return false, err
	}

	return c.StatusCode == api.Running, nil
}
