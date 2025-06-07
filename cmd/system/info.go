package system

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "显示系统基础信息",
	Long:  color.Success.Render("\r\n显示操作系统基础信息，包括版本、内核、硬件等"),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		
		// 显示标题
		showSystemInfoHeader()
		
		// 显示系统基础信息
		showBasicSystemInfo(verbose)
		
		// 显示硬件信息
		showHardwareInfo(verbose)
		
		// 显示环境信息
		showEnvironmentInfo(verbose)
	},
}

func init() {
	infoCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
}

// showSystemInfoHeader 显示系统信息标题
func showSystemInfoHeader() {
	color.Blue.Println("===========================================")
	color.Blue.Println("         💻 系统基础信息              ")
	color.Blue.Println("===========================================")
	fmt.Println()
}

// showBasicSystemInfo 显示基础系统信息
func showBasicSystemInfo(verbose bool) {
	color.Blue.Println("🖥️  操作系统信息")
	
	// Go runtime 信息
	color.Info.Printf("操作系统: %s\n", runtime.GOOS)
	color.Info.Printf("系统架构: %s\n", runtime.GOARCH)
	color.Info.Printf("CPU 核心数: %d\n", runtime.NumCPU())
	
	// 主机名
	if hostname, err := os.Hostname(); err == nil {
		color.Info.Printf("主机名: %s\n", hostname)
	}
	
	// 当前用户
	if user := os.Getenv("USER"); user != "" {
		color.Info.Printf("当前用户: %s\n", user)
	}
	
	// 工作目录
	if cwd, err := os.Getwd(); err == nil {
		color.Info.Printf("工作目录: %s\n", cwd)
	}
	
	// macOS 特定信息
	if runtime.GOOS == "darwin" {
		showMacOSInfo(verbose)
	}
	
	fmt.Println()
}

// showMacOSInfo 显示 macOS 特定信息
func showMacOSInfo(verbose bool) {
	// 系统版本
	if version := getCommandOutput("sw_vers", "-productVersion"); version != "" {
		color.Info.Printf("macOS 版本: %s\n", version)
	}
	
	// 构建版本
	if build := getCommandOutput("sw_vers", "-buildVersion"); build != "" {
		color.Info.Printf("构建版本: %s\n", build)
	}
	
	// 内核版本
	if kernel := getCommandOutput("uname", "-r"); kernel != "" {
		color.Info.Printf("内核版本: %s\n", kernel)
	}
	
	if verbose {
		// 系统启动时间
		if uptime := getCommandOutput("uptime"); uptime != "" {
			color.Info.Printf("系统运行时间: %s\n", strings.TrimSpace(uptime))
		}
	}
}

// showHardwareInfo 显示硬件信息
func showHardwareInfo(verbose bool) {
	color.Blue.Println("🔧 硬件信息")
	
	if runtime.GOOS == "darwin" {
		// CPU 信息
		if cpuBrand := getCommandOutput("sysctl", "-n", "machdep.cpu.brand_string"); cpuBrand != "" {
			color.Info.Printf("处理器: %s\n", cpuBrand)
		}
		
		// 内存信息
		if memSize := getCommandOutput("sysctl", "-n", "hw.memsize"); memSize != "" {
			// 转换字节为 GB
			if size := parseMemorySize(memSize); size > 0 {
				color.Info.Printf("内存大小: %.1f GB\n", float64(size)/1024/1024/1024)
			}
		}
		
		if verbose {
			// CPU 频率
			if cpuFreq := getCommandOutput("sysctl", "-n", "hw.cpufrequency_max"); cpuFreq != "" {
				if freq := parseMemorySize(cpuFreq); freq > 0 {
					color.Info.Printf("CPU 最大频率: %.2f GHz\n", float64(freq)/1000000000)
				}
			}
			
			// 缓存信息
			if l1Cache := getCommandOutput("sysctl", "-n", "hw.l1icachesize"); l1Cache != "" {
				color.Info.Printf("L1 指令缓存: %s bytes\n", l1Cache)
			}
			if l2Cache := getCommandOutput("sysctl", "-n", "hw.l2cachesize"); l2Cache != "" {
				color.Info.Printf("L2 缓存: %s bytes\n", l2Cache)
			}
			if l3Cache := getCommandOutput("sysctl", "-n", "hw.l3cachesize"); l3Cache != "" {
				color.Info.Printf("L3 缓存: %s bytes\n", l3Cache)
			}
		}
	}
	
	fmt.Println()
}

// showEnvironmentInfo 显示环境信息
func showEnvironmentInfo(verbose bool) {
	color.Blue.Println("🌍 环境信息")
	
	// Shell 信息
	if shell := os.Getenv("SHELL"); shell != "" {
		color.Info.Printf("默认 Shell: %s\n", shell)
	}
	
	// 终端信息
	if term := os.Getenv("TERM"); term != "" {
		color.Info.Printf("终端类型: %s\n", term)
	}
	
	// 语言环境
	if lang := os.Getenv("LANG"); lang != "" {
		color.Info.Printf("语言环境: %s\n", lang)
	}
	
	if verbose {
		// PATH 环境变量
		if path := os.Getenv("PATH"); path != "" {
			color.Info.Println("PATH 环境变量:")
			paths := strings.Split(path, ":")
			for i, p := range paths {
				if i < 10 { // 只显示前10个路径
					color.Gray.Printf("  %s\n", p)
				} else if i == 10 {
					color.Gray.Printf("  ... 还有 %d 个路径\n", len(paths)-10)
					break
				}
			}
		}
	}
	
	fmt.Println()
}

// getCommandOutput 执行命令并返回输出
func getCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// parseMemorySize 解析内存大小字符串为整数
func parseMemorySize(sizeStr string) int64 {
	var size int64
	fmt.Sscanf(sizeStr, "%d", &size)
	return size
}