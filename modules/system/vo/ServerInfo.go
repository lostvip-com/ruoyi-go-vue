package vo

import (
	"fmt"
	"github.com/lostvip-com/lv_framework/utils/lv_file"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ServerInfo struct {
	Cpu     Cpu       `json:"cpu"`
	GoInfo  GoInfo    `json:"jvm"`
	Mem     Mem       `json:"mem"`
	Sys     Sys       `json:"sys"`
	SysFile []SysFile `json:"sysFile"`
}

type Cpu struct {
	CpuNum int     `json:"cpuNum"`
	Total  int     `json:"total"`
	Sys    string  `json:"sys"`
	Used   float64 `json:"used"`
	Wait   string  `json:"wait"`
	Free   string  `json:"free"`
}

type GoInfo struct {
	Total   uint64 `json:"total"`
	Version string `json:"version"`
	Home    string `json:"home"`
	Used    uint64 `json:"used"`
}

type Mem struct {
	Total uint64  `json:"total"`
	Used  uint64  `json:"used"`
	Free  uint64  `json:"free"`
	Usage float64 `json:"usage"`
}

type Sys struct {
	ComputerName string `json:"computerName"`
	ComputerIp   string `json:"computerIp"`
	UserDir      string `json:"userDir"`
	OsName       string `json:"osName"`
	OsArch       string `json:"osArch"`
}

type SysFile struct {
	DirName     string  `json:"dirName"`
	SysTypeName string  `json:"sysTypeName"`
	TypeName    string  `json:"typeName"`
	Total       string  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	Usage       float64 `json:"usage"`
}

func (s *ServerInfo) getCpu() Cpu {
	cpuCount, _ := cpu.Counts(true) //cpu逻辑数量
	return Cpu{
		CpuNum: runtime.NumCPU(),
		Total:  cpuCount,
		Sys:    "",
		Used:   s.getCpuPercent(),
		Wait:   "",
		Free:   "",
	}
}
func (s *ServerInfo) getCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func (s *ServerInfo) getJvm() GoInfo {
	var gomem runtime.MemStats
	runtime.ReadMemStats(&gomem)
	goUsed := gomem.Sys / 1024 / 1024

	return GoInfo{
		Total:   goUsed,
		Used:    goUsed,
		Version: runtime.Version(), //版本
		Home:    "",
	}
}

func (s *ServerInfo) getMem() Mem {
	info, _ := mem.VirtualMemory()
	memUsage, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", info.UsedPercent), 64)
	return Mem{
		Total: info.Total,
		Used:  info.Used,
		Free:  info.Total - info.Used,
		Usage: memUsage,
	}
}

func (s *ServerInfo) getSys(ip string) Sys {
	info, _ := host.Info()
	return Sys{
		ComputerName: info.Hostname,
		ComputerIp:   ip,
		UserDir:      s.GetAppPath(),
		OsName:       info.OS,
		OsArch:       runtime.GOARCH,
	}
}

// GetAppPath 获取运行时的相对目录
func (s *ServerInfo) GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func (s *ServerInfo) getSysFile() []SysFile {
	var sysFiles []SysFile
	infos, _ := disk.Partitions(true) //所有分区
	for i := 0; i < len(infos); i++ {
		var info = infos[i]
		info2, _ := disk.Usage(info.Mountpoint) //指定某路径的硬盘使用情况
		var sysFile = SysFile{
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
