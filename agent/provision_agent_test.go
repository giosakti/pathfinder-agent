package agent

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pathfinder-cm/pathfinder-agent/mock"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
)

func TestProvisionProcess(t *testing.T) {
	node := "test-01"

	scs := make(pfmodel.ContainerList, 4)
	scs[0] = pfmodel.Container{Hostname: "test-c-01", Type: "image", Alias: "16.04", Certificate: "random", Mode: "pull", Server: "https://cloud-images.ubuntu.com/releases", Protocol: "simplestreams", Status: "SCHEDULED"}
	scs[1] = pfmodel.Container{Hostname: "test-c-02", Type: "image", Alias: "16.04", Certificate: "random", Mode: "pull", Server: "https://cloud-images.ubuntu.com/releases", Protocol: "simplestreams", Status: "SCHEDULED"}
	scs[2] = pfmodel.Container{Hostname: "test-c-03", Type: "image", Alias: "16.04", Certificate: "random", Mode: "pull", Server: "https://cloud-images.ubuntu.com/releases", Protocol: "simplestreams", Status: "SCHEDULED"}
	scs[3] = pfmodel.Container{Hostname: "test-c-04", Type: "image", Alias: "16.04", Certificate: "random", Mode: "pull", Server: "https://cloud-images.ubuntu.com/releases", Protocol: "simplestreams", Status: "SCHEDULE_DELETION"}

	lcs := make(pfmodel.ContainerList, 3)
	lcs[0] = pfmodel.Container{Hostname: "test-c-01", Alias: "16.04"}
	lcs[1] = pfmodel.Container{Hostname: "test-c-02", Alias: "16.04"}
	lcs[2] = pfmodel.Container{Hostname: "test-c-04", Alias: "16.04"}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockContainerDaemon := mock.NewMockContainerDaemon(mockCtrl)
	mockContainerDaemon.EXPECT().ListContainers().Return(&lcs, nil).AnyTimes()
	mockContainerDaemon.EXPECT().CreateContainer("test-c-03", "image", "16.04", "random", "pull", "https://cloud-images.ubuntu.com/releases", "simplestreams").Return(true, "127.0.0.1", nil)
	mockContainerDaemon.EXPECT().DeleteContainer("test-c-04").Return(true, nil)

	mockPfClient := mock.NewMockPfclient(mockCtrl)
	mockPfClient.EXPECT().FetchContainersFromServer(node).Return(&scs, nil)
	mockPfClient.EXPECT().UpdateIpaddress(node, "test-c-03", "127.0.0.1").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-01").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-02").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-03").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsDeleted(node, "test-c-04").Return(true, nil)

	provisionAgent := NewProvisionAgent(node, mockContainerDaemon, mockPfClient)
	ok := provisionAgent.Process()
	if ok != true {
		t.Errorf("Agent does not process properly")
	}
}
