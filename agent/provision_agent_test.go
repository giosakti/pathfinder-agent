package agent

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pathfinder-cm/pathfinder-agent/mock"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
)

func TestProvisionProcess(t *testing.T) {
	node := "test-01"

	scs := make(pfmodel.ContainerList, 4)
	scs[0] = pfmodel.Container{Hostname: "test-c-01", Status: "SCHEDULED", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Mode: "pull",
		Remote: pfmodel.Remote{
			Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
		},
	},
	}
	scs[1] = pfmodel.Container{Hostname: "test-c-02", Status: "SCHEDULED", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Mode: "pull",
		Remote: pfmodel.Remote{
			Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
		},
	},
	}
	scs[2] = pfmodel.Container{Hostname: "test-c-03", Status: "SCHEDULED", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Mode: "pull",
		Remote: pfmodel.Remote{
			Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
		},
	},
	}
	scs[3] = pfmodel.Container{Hostname: "test-c-04", Status: "SCHEDULE_DELETION", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Mode: "pull",
		Remote: pfmodel.Remote{
			Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
		},
	},
	}

	lcs := make(pfmodel.ContainerList, 3)
	lcs[0] = pfmodel.Container{Hostname: "test-c-01", Source: pfmodel.Source{Alias: "16.04"}}
	lcs[1] = pfmodel.Container{Hostname: "test-c-02", Source: pfmodel.Source{Alias: "16.04"}}
	lcs[2] = pfmodel.Container{Hostname: "test-c-04", Source: pfmodel.Source{Alias: "16.04"}}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockContainerDaemon := mock.NewMockContainerDaemon(mockCtrl)
	mockContainerDaemon.EXPECT().ListContainers().Return(&lcs, nil).AnyTimes()
	mockContainerDaemon.EXPECT().CreateContainer(scs[2]).Return(true, "127.0.0.3", nil)
	mockContainerDaemon.EXPECT().DeleteContainer(scs[3].Hostname).Return(true, nil)

	mockPfClient := mock.NewMockPfclient(mockCtrl)
	mockPfClient.EXPECT().FetchScheduledContainersFromServer(node).Return(&scs, nil)
	mockPfClient.EXPECT().UpdateIpaddress(node, "test-c-03", "127.0.0.3").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-01").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-02").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-03").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsDeleted(node, "test-c-04").Return(true, nil)

	provisionAgent := NewProvisionAgent(node, mockContainerDaemon, mockPfClient)
	ok := provisionAgent.Process()
	if !ok {
		t.Errorf("Agent does not process properly")
	}
}

func TestRelocationProcess(t *testing.T) {
	node := "test-01"

	scs := make(pfmodel.ContainerList, 4)
	scs[0] = pfmodel.Container{Hostname: "test-c-01", Status: "SCHEDULE_RELOCATION", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Mode: "pull",
		Remote: pfmodel.Remote{
			Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
		},
	},
	}
	scs[1] = pfmodel.Container{Hostname: "test-c-04", Status: "SCHEDULE_RELOCATION", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Mode: "pull",
		Remote: pfmodel.Remote{
			Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
		},
	},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockContainerDaemon := mock.NewMockContainerDaemon(mockCtrl)
	mockContainerDaemon.EXPECT().ListContainers().Return(&pfmodel.ContainerList{}, nil).Times(len(scs))
	mockContainerDaemon.EXPECT().MigrateContainer(scs[0]).Return(true, "127.0.0.10", nil)
	mockContainerDaemon.EXPECT().MigrateContainer(scs[1]).Return(true, "127.0.0.11", nil)

	mockPfClient := mock.NewMockPfclient(mockCtrl)
	mockPfClient.EXPECT().FetchScheduledContainersFromServer(node).Return(&scs, nil)
	mockPfClient.EXPECT().MarkContainerAsRelocateStarted(node, "test-c-01").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsRelocateStarted(node, "test-c-04").Return(true, nil)
	mockPfClient.EXPECT().UpdateIpaddress(node, "test-c-01", "127.0.0.10").Return(true, nil)
	mockPfClient.EXPECT().UpdateIpaddress(node, "test-c-04", "127.0.0.11").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-01").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsProvisioned(node, "test-c-04").Return(true, nil)

	provisionAgent := NewProvisionAgent(node, mockContainerDaemon, mockPfClient)
	ok := provisionAgent.Process()
	if !ok {
		t.Errorf("Agent does not process properly")
	}
}

func TestRelocationProcess_error(t *testing.T) {
	// should mark the container as relocate_error
	node := "test-01"

	scs := make(pfmodel.ContainerList, 4)
	scs[0] = pfmodel.Container{Hostname: "test-c-01", Status: "SCHEDULE_RELOCATION", Source: pfmodel.Source{
		Type: "image", Alias: "16.04", Mode: "pull",
		Remote: pfmodel.Remote{
			Server: "https://cloud-images.ubuntu.com/releases", Certificate: "random", Protocol: "simplestreams", AuthType: "none",
		},
	},
	}

	lcs := make(pfmodel.ContainerList, 3)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockContainerDaemon := mock.NewMockContainerDaemon(mockCtrl)
	mockContainerDaemon.EXPECT().ListContainers().Return(&lcs, nil).Times(len(scs))
	mockContainerDaemon.EXPECT().MigrateContainer(scs[0]).Return(false, "", errors.New("unexpected_error"))

	mockPfClient := mock.NewMockPfclient(mockCtrl)
	mockPfClient.EXPECT().FetchScheduledContainersFromServer(node).Return(&scs, nil)
	mockPfClient.EXPECT().MarkContainerAsRelocateStarted(node, "test-c-01").Return(true, nil)
	mockPfClient.EXPECT().MarkContainerAsRelocateError(node, "test-c-01")

	provisionAgent := NewProvisionAgent(node, mockContainerDaemon, mockPfClient)
	ok := provisionAgent.Process()
	if !ok {
		t.Errorf("Agent does not process properly")
	}
}
