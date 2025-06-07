package system

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "显示磁盘信息",
	Long:  color.Success.Render("\r\n显示磁盘使用情况、I/O统计、文件系统信息等"),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		showIO, _ := cmd.Flags().GetBool("io")
		showInodes, _ := cmd.Flags().GetBool("inodes")
		
		// 显示标题
		showDiskHeader()
		
		// 显示磁盘使用情况
		showDiskUsageDetailed(verbose)
		
		// 显示文件系统信息
		if verbose {
			showFilesystemInfo()
		}
		
		// 显示inode使用情况
		if showInodes {
			showInodeUsage()
		}
		
		// 显示磁盘I/O统计
		if showIO {
			showDiskIOStats(verbose)
		}
		
		// 显示挂载点信息
		if verbose {
			showMountPoints()
		}
	},
}

func init() {
	diskCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
	diskCmd.Flags().Bool("io", false, "显示磁盘I/O统计")
	diskCmd.Flags().Bool("inodes", false, "显示inode使用情况")
}

// DiskInfo 磁盘信息结构
type DiskInfo struct {
	Filesystem string
	Size       string
	Used       string
	Avail      string
	UsePercent string
	MountPoint string
	Type       string
}

// showDiskHeader 显示磁盘信息标题
func showDiskHeader() {
	color.Blue.Println("===========================================")
	color.Blue.Println("         💿 磁盘信息                  ")
	color.Blue.Println("===========================================")
	fmt.Println()
}

// showDiskUsageDetailed 显示详细的磁盘使用情况
func showDiskUsageDetailed(verbose bool) {
	color.Blue.Println("📊 磁盘使用情况")
	
	// 获取磁盘使用情况
	dfOutput := getCommandOutput("df", "-h")
	if dfOutput == "" {
		color.Error.Println("❌ 无法获取磁盘使用信息")
		return
	}
	
	lines := strings.Split(dfOutput, "\n")
	disks := []DiskInfo{}
	
	for i, line := range lines {
		if i == 0 {
			// 显示表头
			color.Yellow.Printf("%-20s %-8s %-8s %-8s %-8s %s\n", "文件系统", "大小", "已用", "可用", "使用率", "挂载点")
			color.Yellow.Println(strings.Repeat("-", 70))
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			disk := DiskInfo{
				Filesystem: fields[0],
				Size:       fields[1],
				Used:       fields[2],
				Avail:      fields[3],
				UsePercent: fields[4],
				MountPoint: fields[5],
			}
			
			// 过滤显示重要的挂载点
			if shouldShowDisk(disk.MountPoint, verbose) {
				disks = append(disks, disk)
				showDiskInfo(disk)
			}
		}
	}
	
	fmt.Println()
	
	// 显示磁盘使用总结
	showDiskSummary(disks)
}

// shouldShowDisk 判断是否应该显示该磁盘
func shouldShowDisk(mountPoint string, verbose bool) bool {
	if verbose {
		return true
	}
	
	// 只显示重要的挂载点
	importantMounts := []string{"/", "/home", "/var", "/tmp", "/usr"}
	for _, mount := range importantMounts {
		if mountPoint == mount {
			return true
		}
	}
	
	// 显示 /Volumes 下的挂载点（外部设备）
	if strings.HasPrefix(mountPoint, "/Volumes/") {
		return true
	}
	
	// 跳过系统内部挂载点
	skipPrefixes := []string{"/dev", "/sys", "/proc", "/run", "/snap"}
	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(mountPoint, prefix) {
			return false
		}
	}
	
	return false
}

// showDiskInfo 显示单个磁盘信息
func showDiskInfo(disk DiskInfo) {
	// 根据使用率着色
	usagePercent := parseUsagePercent(disk.UsePercent)
	var usageColor func(a ...interface{}) string
	
	if usagePercent > 90 {
		usageColor = color.Red.Sprint
	} else if usagePercent > 80 {
		usageColor = color.Yellow.Sprint
	} else {
		usageColor = color.Green.Sprint
	}
	
	fmt.Printf("%-20s %-8s %-8s %-8s %s %s\n",
		truncateString(disk.Filesystem, 20),
		disk.Size,
		disk.Used,
		disk.Avail,
		usageColor(fmt.Sprintf("%-8s", disk.UsePercent)),
		disk.MountPoint)
	
	// 显示使用率进度条
	if usagePercent >= 0 {
		showUsageBar(usagePercent)
	}
}

// showDiskSummary 显示磁盘使用总结
func showDiskSummary(disks []DiskInfo) {
	color.Blue.Println("📈 磁盘使用总结")
	
	totalDisks := len(disks)
	highUsageDisks := 0
	
	for _, disk := range disks {
		usagePercent := parseUsagePercent(disk.UsePercent)
		if usagePercent > 80 {
			highUsageDisks++
		}
	}
	
	color.Info.Printf("总磁盘数: %d\n", totalDisks)
	color.Info.Printf("高使用率磁盘 (>80%%): %d\n", highUsageDisks)
	
	if highUsageDisks > 0 {
		color.Yellow.Println("⚠️  警告: 发现高使用率磁盘，建议清理空间")
	}
	
	fmt.Println()
}

// showFilesystemInfo 显示文件系统信息
func showFilesystemInfo() {
	color.Blue.Println("🗂️  文件系统信息")
	
	if runtime.GOOS == "darwin" {
		// 使用 mount 命令获取挂载信息
		mountOutput := getCommandOutput("mount")
		if mountOutput != "" {
			lines := strings.Split(mountOutput, "\n")
			fsTypes := make(map[string]int)
			
			for _, line := range lines {
				if strings.Contains(line, " on ") && strings.Contains(line, " type ") {
					// 解析挂载信息
					parts := strings.Split(line, " type ")
					if len(parts) >= 2 {
						fsType := strings.Fields(parts[1])[0]
						fsTypes[fsType]++
					}
				}
			}
			
			color.Info.Println("文件系统类型统计:")
			for fsType, count := range fsTypes {
				color.Gray.Printf("  %s: %d 个挂载点\n", fsType, count)
			}
		}
	}
	
	fmt.Println()
}

// showInodeUsage 显示inode使用情况
func showInodeUsage() {
	color.Blue.Println("🔢 Inode 使用情况")
	
	// 使用 df -i 获取inode信息
	dfInodeOutput := getCommandOutput("df", "-i")
	if dfInodeOutput == "" {
		color.Error.Println("❌ 无法获取inode使用信息")
		return
	}
	
	lines := strings.Split(dfInodeOutput, "\n")
	for i, line := range lines {
		if i == 0 {
			// 显示表头
			color.Yellow.Printf("%-20s %-10s %-10s %-10s %-8s %s\n", "文件系统", "Inodes", "已用", "可用", "使用率", "挂载点")
			color.Yellow.Println(strings.Repeat("-", 75))
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			filesystem := fields[0]
			inodes := fields[1]
			used := fields[2]
			avail := fields[3]
			usePercent := fields[4]
			mountPoint := fields[5]
			
			// 只显示重要的挂载点
			if shouldShowDisk(mountPoint, false) {
				// 根据使用率着色
				usagePercent := parseUsagePercent(usePercent)
				var usageColor func(a ...interface{}) string
				
				if usagePercent > 90 {
					usageColor = color.Red.Sprint
				} else if usagePercent > 80 {
					usageColor = color.Yellow.Sprint
				} else {
					usageColor = color.Green.Sprint
				}
				
				fmt.Printf("%-20s %-10s %-10s %-10s %s %s\n",
					truncateString(filesystem, 20),
					inodes,
					used,
					avail,
					usageColor(fmt.Sprintf("%-8s", usePercent)),
					mountPoint)
			}
		}
	}
	
	fmt.Println()
}

// showDiskIOStats 显示磁盘I/O统计
func showDiskIOStats(verbose bool) {
	color.Blue.Println("⚡ 磁盘 I/O 统计")
	
	if runtime.GOOS == "darwin" {
		// 使用 iostat 命令获取I/O统计
		iostatOutput := getCommandOutput("iostat", "-d")
		if iostatOutput != "" {
			lines := strings.Split(iostatOutput, "\n")
			for i, line := range lines {
				if i < 3 { // 跳过前几行说明
					continue
				}
				
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					device := fields[0]
					if strings.HasPrefix(device, "disk") {
						color.Info.Printf("设备: %s\n", device)
						if len(fields) >= 6 {
							color.Gray.Printf("  读取 KB/t: %s\n", fields[1])
							color.Gray.Printf("  写入 KB/t: %s\n", fields[2])
							color.Gray.Printf("  读取次数: %s\n", fields[3])
							color.Gray.Printf("  写入次数: %s\n", fields[4])
							color.Gray.Printf("  读取 MB: %s\n", fields[5])
							if len(fields) >= 7 {
								color.Gray.Printf("  写入 MB: %s\n", fields[6])
							}
						}
						fmt.Println()
					}
				}
			}
		} else {
			color.Yellow.Println("⚠️  无法获取I/O统计信息 (可能需要安装 iostat)")
		}
		
		// 显示磁盘活动
		if verbose {
			showDiskActivity()
		}
	}
	
	fmt.Println()
}

// showDiskActivity 显示磁盘活动
func showDiskActivity() {
	color.Info.Println("💾 磁盘活动:")
	
	if runtime.GOOS == "darwin" {
		// 使用 iotop 或类似命令（如果可用）
		if _, err := exec.LookPath("iotop"); err == nil {
			iotopOutput := getCommandOutput("iotop", "-a", "-o", "-d", "1", "-n", "1")
			if iotopOutput != "" {
				lines := strings.Split(iotopOutput, "\n")
				for i, line := range lines {
					if i < 5 && strings.TrimSpace(line) != "" { // 只显示前几行
						color.Gray.Printf("  %s\n", line)
					}
				}
			}
		} else {
			color.Gray.Println("  iotop 命令不可用")
		}
	}
}

// showMountPoints 显示挂载点信息
func showMountPoints() {
	color.Blue.Println("🔗 挂载点信息")
	
	if runtime.GOOS == "darwin" {
		mountOutput := getCommandOutput("mount")
		if mountOutput != "" {
			lines := strings.Split(mountOutput, "\n")
			for _, line := range lines {
				if strings.Contains(line, " on ") && strings.Contains(line, " type ") {
					// 解析挂载信息
					parts := strings.Split(line, " on ")
					if len(parts) >= 2 {
						device := parts[0]
						mountInfo := parts[1]
						
						// 进一步解析挂载点和类型
						typeParts := strings.Split(mountInfo, " type ")
						if len(typeParts) >= 2 {
							mountPoint := typeParts[0]
							fsType := strings.Fields(typeParts[1])[0]
							
							// 只显示重要的挂载点
							if shouldShowDisk(mountPoint, false) {
								color.Info.Printf("设备: %s\n", device)
								color.Gray.Printf("  挂载点: %s\n", mountPoint)
								color.Gray.Printf("  文件系统: %s\n", fsType)
								fmt.Println()
							}
						}
					}
				}
			}
		}
	}
	
	fmt.Println()
}