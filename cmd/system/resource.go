package system

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "显示系统资源占用情况",
	Long:  color.Success.Render("\r\n显示系统资源占用情况，包括CPU、内存、磁盘使用率等"),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		watch, _ := cmd.Flags().GetBool("watch")
		
		if watch {
			color.Info.Println("监控模式 (按 Ctrl+C 退出)")
			fmt.Println()
			// 这里可以实现循环监控，暂时只显示一次
		}
		
		// 显示标题
		showResourceHeader()
		
		// 显示CPU使用情况
		showCPUUsage(verbose)
		
		// 显示内存使用情况
		showMemoryUsage(verbose)
		
		// 显示磁盘使用情况
		showDiskUsage(verbose)
		
		// 显示负载信息
		showLoadAverage(verbose)
	},
}

func init() {
	resourceCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
	resourceCmd.Flags().BoolP("watch", "w", false, "监控模式")
}

// showResourceHeader 显示资源信息标题
func showResourceHeader() {
	color.Blue.Println("===========================================")
	color.Blue.Println("         📊 系统资源占用情况          ")
	color.Blue.Println("===========================================")
	fmt.Println()
}

// showCPUUsage 显示CPU使用情况
func showCPUUsage(verbose bool) {
	color.Blue.Println("🔥 CPU 使用情况")
	
	color.Info.Printf("CPU 核心数: %d\n", runtime.NumCPU())
	
	if runtime.GOOS == "darwin" {
		// 获取CPU使用率
		if cpuUsage := getCPUUsage(); cpuUsage >= 0 {
			color.Info.Printf("CPU 使用率: %.1f%%\n", cpuUsage)
			showUsageBar(cpuUsage)
		}
		
		if verbose {
			// 显示每个核心的使用情况
			if topOutput := getCommandOutput("top", "-l", "1", "-n", "0"); topOutput != "" {
				lines := strings.Split(topOutput, "\n")
				for _, line := range lines {
					if strings.Contains(line, "CPU usage:") {
						color.Gray.Printf("详细信息: %s\n", strings.TrimSpace(line))
						break
					}
				}
			}
		}
	}
	
	fmt.Println()
}

// showMemoryUsage 显示内存使用情况
func showMemoryUsage(verbose bool) {
	color.Blue.Println("💾 内存使用情况")
	
	if runtime.GOOS == "darwin" {
		memInfo := getMemoryInfo()
		if memInfo != nil {
			color.Info.Printf("总内存: %.1f GB\n", float64(memInfo.Total)/1024/1024/1024)
			color.Info.Printf("已使用: %.1f GB\n", float64(memInfo.Used)/1024/1024/1024)
			color.Info.Printf("可用内存: %.1f GB\n", float64(memInfo.Free)/1024/1024/1024)
			
			usagePercent := float64(memInfo.Used) / float64(memInfo.Total) * 100
			color.Info.Printf("使用率: %.1f%%\n", usagePercent)
			showUsageBar(usagePercent)
			
			if verbose {
				color.Gray.Printf("缓存: %.1f GB\n", float64(memInfo.Cached)/1024/1024/1024)
				color.Gray.Printf("缓冲区: %.1f GB\n", float64(memInfo.Buffer)/1024/1024/1024)
			}
		}
	}
	
	fmt.Println()
}

// showDiskUsage 显示磁盘使用情况
func showDiskUsage(verbose bool) {
	color.Blue.Println("💿 磁盘使用情况")
	
	// 获取磁盘使用情况
	if dfOutput := getCommandOutput("df", "-h"); dfOutput != "" {
		lines := strings.Split(dfOutput, "\n")
		for i, line := range lines {
			if i == 0 {
				// 跳过标题行
				continue
			}
			
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				filesystem := fields[0]
				size := fields[1]
				used := fields[2]
				avail := fields[3]
				usageStr := fields[4]
				mountPoint := fields[5]
				
				// 只显示主要的挂载点
				if mountPoint == "/" || strings.HasPrefix(mountPoint, "/Volumes") || verbose {
					color.Info.Printf("挂载点: %s\n", mountPoint)
					color.Info.Printf("  文件系统: %s\n", filesystem)
					color.Info.Printf("  总大小: %s\n", size)
					color.Info.Printf("  已使用: %s\n", used)
					color.Info.Printf("  可用空间: %s\n", avail)
					color.Info.Printf("  使用率: %s\n", usageStr)
					
					// 解析使用率百分比
					if usagePercent := parseUsagePercent(usageStr); usagePercent >= 0 {
						showUsageBar(usagePercent)
					}
					fmt.Println()
				}
			}
		}
	}
}

// showLoadAverage 显示系统负载
func showLoadAverage(verbose bool) {
	color.Blue.Println("⚡ 系统负载")
	
	if runtime.GOOS == "darwin" {
		if uptime := getCommandOutput("uptime"); uptime != "" {
			// 解析 uptime 输出中的负载信息
			if strings.Contains(uptime, "load averages:") {
				parts := strings.Split(uptime, "load averages:")
				if len(parts) > 1 {
					loadInfo := strings.TrimSpace(parts[1])
					color.Info.Printf("负载平均值: %s\n", loadInfo)
					
					if verbose {
						color.Gray.Println("说明: 1分钟 5分钟 15分钟平均负载")
						color.Gray.Printf("CPU核心数: %d (负载超过此值表示系统繁忙)\n", runtime.NumCPU())
					}
				}
			}
		}
	}
	
	fmt.Println()
}

// MemoryInfo 内存信息结构
type MemoryInfo struct {
	Total  int64
	Used   int64
	Free   int64
	Cached int64
	Buffer int64
}

// getCPUUsage 获取CPU使用率
func getCPUUsage() float64 {
	if runtime.GOOS == "darwin" {
		// 使用 top 命令获取CPU使用率
		output := getCommandOutput("top", "-l", "1", "-n", "0")
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "CPU usage:") {
				// 解析类似 "CPU usage: 10.0% user, 5.0% sys, 85.0% idle" 的行
				parts := strings.Split(line, ",")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.Contains(part, "idle") {
						// 提取idle百分比
						fields := strings.Fields(part)
						if len(fields) > 0 {
							idleStr := strings.TrimSuffix(fields[0], "%")
							if idle, err := strconv.ParseFloat(idleStr, 64); err == nil {
								return 100.0 - idle // CPU使用率 = 100% - idle%
							}
						}
					}
				}
				break
			}
		}
	}
	return -1
}

// getMemoryInfo 获取内存信息
func getMemoryInfo() *MemoryInfo {
	if runtime.GOOS == "darwin" {
		// 使用 vm_stat 命令获取内存信息
		output := getCommandOutput("vm_stat")
		if output == "" {
			return nil
		}
		
		// 获取页面大小
		pageSize := int64(4096) // 默认4KB
		if pageSizeStr := getCommandOutput("sysctl", "-n", "hw.pagesize"); pageSizeStr != "" {
			if ps, err := strconv.ParseInt(pageSizeStr, 10, 64); err == nil {
				pageSize = ps
			}
		}
		
		memInfo := &MemoryInfo{}
		lines := strings.Split(output, "\n")
		
		for _, line := range lines {
			if strings.Contains(line, "Pages free:") {
				if pages := extractPages(line); pages > 0 {
					memInfo.Free = pages * pageSize
				}
			} else if strings.Contains(line, "Pages active:") {
				if pages := extractPages(line); pages > 0 {
					memInfo.Used += pages * pageSize
				}
			} else if strings.Contains(line, "Pages inactive:") {
				if pages := extractPages(line); pages > 0 {
					memInfo.Used += pages * pageSize
				}
			} else if strings.Contains(line, "Pages wired down:") {
				if pages := extractPages(line); pages > 0 {
					memInfo.Used += pages * pageSize
				}
			}
		}
		
		memInfo.Total = memInfo.Used + memInfo.Free
		return memInfo
	}
	return nil
}

// extractPages 从vm_stat输出行中提取页面数
func extractPages(line string) int64 {
	fields := strings.Fields(line)
	for _, field := range fields {
		field = strings.TrimSuffix(field, ".")
		if pages, err := strconv.ParseInt(field, 10, 64); err == nil && pages > 0 {
			return pages
		}
	}
	return 0
}

// parseUsagePercent 解析使用率百分比
func parseUsagePercent(usageStr string) float64 {
	usageStr = strings.TrimSuffix(usageStr, "%")
	if usage, err := strconv.ParseFloat(usageStr, 64); err == nil {
		return usage
	}
	return -1
}

// showUsageBar 显示使用率进度条
func showUsageBar(percent float64) {
	barLength := 30
	filledLength := int(percent / 100.0 * float64(barLength))
	
	bar := "["
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			if percent > 80 {
				bar += color.Red.Sprint("█")
			} else if percent > 60 {
				bar += color.Yellow.Sprint("█")
			} else {
				bar += color.Green.Sprint("█")
			}
		} else {
			bar += "░"
		}
	}
	bar += "]"
	
	fmt.Printf("  %s %.1f%%\n", bar, percent)
}