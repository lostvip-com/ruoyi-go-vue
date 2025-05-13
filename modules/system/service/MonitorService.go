package service

import (
	"github.com/lostvip-com/lv_framework/utils/lv_file"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"system/vo"
	"time"
)

type MonitorService struct {
	BaseService
}

func (s *MonitorService) GetCpu() vo.Cpu {
	cpuCount, _ := cpu.Counts(true) //cpu逻辑数量
	return vo.Cpu{
		CpuNum: runtime.NumCPU(),
		Total:  cpuCount,
		Sys:    "",
		Used:   s.GetCpuPercent(),
		Wait:   "",
		Free:   "",
	}
}
func (s *MonitorService) GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func (s *MonitorService) GetGoInfo() vo.GoInfo {

	var gomem runtime.MemStats
	runtime.ReadMemStats(&gomem)
	goUsed := gomem.Sys / 1024 / 1024
	return vo.GoInfo{
		Total:   goUsed,
		Used:    goUsed,
		Version: "",
		Home:    "",
	}
}

func (s *MonitorService) GetMem() vo.Mem {
	info, _ := mem.VirtualMemory()
	return vo.Mem{
		Total: info.Total / 1024 / 1024 / 1024,
		Used:  info.Used / 1024 / 1024 / 1024,
		Free:  (info.Total - info.Used) / 1024 / 1024 / 1024,
	}
}

func (s *MonitorService) GetSys(ip string) vo.Sys {
	info, _ := host.Info()
	return vo.Sys{
		ComputerName: info.Hostname,
		ComputerIp:   ip,
		UserDir:      s.GetAppPath(),
		OsName:       info.OS,
		OsArch:       runtime.GOARCH,
	}
}

// GetAppPath 获取运行时的相对目录
func (s *MonitorService) GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func (s *MonitorService) GetSysFile() []vo.SysFile {
	var sysFiles []vo.SysFile
	infos, _ := disk.Partitions(true) //所有分区
	for i := 0; i < len(infos); i++ {
		var info = infos[i]
		info2, _ := disk.Usage(info.Mountpoint) //指定某路径的硬盘使用情况
		var sysFile = vo.SysFile{
			DirName:     info2.Path,
			SysTypeName: info.Fstype,
			TypeName:    "",
			Total:       lv_file.FormatFileSize(info2.Total),
			Free:        info2.Free,
			Used:        info2.Used,
			Usage:       info2.UsedPercent,
		}
		sysFiles = append(sysFiles, sysFile)
	}
	return sysFiles
}
