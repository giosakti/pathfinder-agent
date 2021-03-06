package metrics

import (
	"github.com/mistifyio/go-zfs"
	"github.com/pathfinder-cm/pathfinder-agent/config"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

func Collect() *pfmodel.Metrics {
	return &pfmodel.Metrics{
		Memory:   collectMemory(),
		Load:     collectLoad(),
		RootDisk: collectRootDisk(),
		ZFSDisk:  collectZFSDisk(),
	}
}

func collectMemory() *pfmodel.Memory {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return &pfmodel.Memory{
		Used:  util.BToMb(vmStat.Used),
		Free:  util.BToMb(vmStat.Free),
		Total: util.BToMb(vmStat.Total),
	}
}

func collectLoad() *pfmodel.Load {
	nCpus, err := cpu.Counts(true)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	avgStat, err := load.Avg()
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return &pfmodel.Load{
		Capacity:   nCpus,
		LoadAvg1M:  avgStat.Load1,
		LoadAvg5M:  avgStat.Load5,
		LoadAvg15M: avgStat.Load15,
	}
}

func collectRootDisk() *pfmodel.Disk {
	usageStat, err := disk.Usage("/")
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return &pfmodel.Disk{
		Total: util.BToMb(usageStat.Total),
		Used:  util.BToMb(usageStat.Used),
	}
}

func collectZFSDisk() *pfmodel.Disk {
	if config.MetricsZpoolName == "" {
		return nil
	}

	zpoolStat, err := zfs.GetZpool(config.MetricsZpoolName)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return &pfmodel.Disk{
		Total: util.BToMb(zpoolStat.Size),
		Used:  util.BToMb(zpoolStat.Allocated),
	}
}
