package system

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "显示网络信息",
	Long:  color.Success.Render("\r\n显示网络接口、连接状态、流量统计等信息"),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		showConnections, _ := cmd.Flags().GetBool("connections")
		showStats, _ := cmd.Flags().GetBool("stats")
		
		// 显示标题
		showNetworkHeader()
		
		// 显示网络接口信息
		showNetworkInterfaces(verbose)
		
		// 显示网络连接
		if showConnections {
			showNetworkConnections(verbose)
		}
		
		// 显示网络统计
		if showStats {
			showNetworkStats(verbose)
		}
		
		// 显示路由信息
		if verbose {
			showRoutingInfo()
		}
	},
}

func init() {
	networkCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
	networkCmd.Flags().BoolP("connections", "c", false, "显示网络连接")
	networkCmd.Flags().BoolP("stats", "s", false, "显示网络统计")
}

// NetworkInterface 网络接口信息
type NetworkInterface struct {
	Name      string
	IPv4      []string
	IPv6      []string
	MAC       string
	MTU       int
	Flags     []string
	IsUp      bool
	IsLoopback bool
}

// showNetworkHeader 显示网络信息标题
func showNetworkHeader() {
	color.Blue.Println("===========================================")
	color.Blue.Println("         🌐 网络信息                  ")
	color.Blue.Println("===========================================")
	fmt.Println()
}

// showNetworkInterfaces 显示网络接口信息
func showNetworkInterfaces(verbose bool) {
	color.Blue.Println("🔌 网络接口")
	
	// 使用 Go 标准库获取网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		color.Error.Printf("❌ 获取网络接口失败: %v\n", err)
		return
	}
	
	for _, iface := range interfaces {
		// 获取接口地址
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		
		netIface := NetworkInterface{
			Name:       iface.Name,
			MAC:        iface.HardwareAddr.String(),
			MTU:        iface.MTU,
			IsUp:       iface.Flags&net.FlagUp != 0,
			IsLoopback: iface.Flags&net.FlagLoopback != 0,
		}
		
		// 解析地址
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					netIface.IPv4 = append(netIface.IPv4, ipnet.IP.String())
				} else {
					netIface.IPv6 = append(netIface.IPv6, ipnet.IP.String())
				}
			}
		}
		
		// 解析标志
		flags := []string{}
		if iface.Flags&net.FlagUp != 0 {
			flags = append(flags, "UP")
		}
		if iface.Flags&net.FlagBroadcast != 0 {
			flags = append(flags, "BROADCAST")
		}
		if iface.Flags&net.FlagLoopback != 0 {
			flags = append(flags, "LOOPBACK")
		}
		if iface.Flags&net.FlagPointToPoint != 0 {
			flags = append(flags, "POINTTOPOINT")
		}
		if iface.Flags&net.FlagMulticast != 0 {
			flags = append(flags, "MULTICAST")
		}
		netIface.Flags = flags
		
		// 显示接口信息
		showInterfaceInfo(netIface, verbose)
	}
	
	fmt.Println()
}

// showInterfaceInfo 显示单个接口信息
func showInterfaceInfo(iface NetworkInterface, verbose bool) {
	// 接口名称和状态
	statusColor := color.Red
	statusText := "DOWN"
	if iface.IsUp {
		statusColor = color.Green
		statusText = "UP"
	}
	
	color.Info.Printf("接口: %s [%s]\n", iface.Name, statusColor.Sprint(statusText))
	
	// IPv4 地址
	if len(iface.IPv4) > 0 {
		color.Info.Printf("  IPv4: %s\n", strings.Join(iface.IPv4, ", "))
	}
	
	// IPv6 地址
	if len(iface.IPv6) > 0 && verbose {
		color.Info.Printf("  IPv6: %s\n", strings.Join(iface.IPv6, ", "))
	}
	
	// MAC 地址
	if iface.MAC != "" {
		color.Info.Printf("  MAC: %s\n", iface.MAC)
	}
	
	// MTU
	if verbose {
		color.Info.Printf("  MTU: %d\n", iface.MTU)
		color.Info.Printf("  标志: %s\n", strings.Join(iface.Flags, ", "))
	}
	
	// 获取接口统计信息（macOS）
	if runtime.GOOS == "darwin" && verbose {
		showInterfaceStats(iface.Name)
	}
	
	fmt.Println()
}

// showInterfaceStats 显示接口统计信息
func showInterfaceStats(interfaceName string) {
	if runtime.GOOS == "darwin" {
		// 使用 netstat 获取接口统计
		output := getCommandOutput("netstat", "-i", "-b")
		if output != "" {
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) >= 10 && fields[0] == interfaceName {
					// 解析统计信息
					if rxPackets, err := strconv.ParseInt(fields[4], 10, 64); err == nil {
						color.Gray.Printf("  接收包数: %d\n", rxPackets)
					}
					if txPackets, err := strconv.ParseInt(fields[7], 10, 64); err == nil {
						color.Gray.Printf("  发送包数: %d\n", txPackets)
					}
					if rxBytes, err := strconv.ParseInt(fields[6], 10, 64); err == nil {
						color.Gray.Printf("  接收字节: %s\n", formatBytes(rxBytes))
					}
					if txBytes, err := strconv.ParseInt(fields[9], 10, 64); err == nil {
						color.Gray.Printf("  发送字节: %s\n", formatBytes(txBytes))
					}
					break
				}
			}
		}
	}
}

// showNetworkConnections 显示网络连接
func showNetworkConnections(verbose bool) {
	color.Blue.Println("🔗 网络连接")
	
	if runtime.GOOS == "darwin" {
		// 显示TCP连接
		color.Info.Println("TCP 连接:")
		tcpOutput := getCommandOutput("netstat", "-an", "-p", "tcp")
		if tcpOutput != "" {
			showConnectionsFromNetstat(tcpOutput, "tcp", verbose)
		}
		
		// 显示UDP连接
		if verbose {
			color.Info.Println("UDP 连接:")
			udpOutput := getCommandOutput("netstat", "-an", "-p", "udp")
			if udpOutput != "" {
				showConnectionsFromNetstat(udpOutput, "udp", verbose)
			}
		}
	}
	
	fmt.Println()
}

// showConnectionsFromNetstat 解析并显示netstat输出
func showConnectionsFromNetstat(output, protocol string, verbose bool) {
	lines := strings.Split(output, "\n")
	connectionCount := 0
	stateCount := make(map[string]int)
	
	for i, line := range lines {
		if i == 0 || i == 1 { // 跳过标题行
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) >= 6 && strings.ToLower(fields[0]) == protocol {
			connectionCount++
			
			localAddr := fields[3]
			foreignAddr := fields[4]
			state := fields[5]
			
			stateCount[state]++
			
			if verbose && connectionCount <= 10 { // 只显示前10个连接
				color.Gray.Printf("  %s -> %s [%s]\n", localAddr, foreignAddr, state)
			}
		}
	}
	
	color.Info.Printf("  总连接数: %d\n", connectionCount)
	
	// 显示状态统计
	if len(stateCount) > 0 {
		color.Info.Println("  状态统计:")
		for state, count := range stateCount {
			color.Gray.Printf("    %s: %d\n", state, count)
		}
	}
}

// showNetworkStats 显示网络统计
func showNetworkStats(verbose bool) {
	color.Blue.Println("📊 网络统计")
	
	if runtime.GOOS == "darwin" {
		// 显示协议统计
		protocolStats := getCommandOutput("netstat", "-s")
		if protocolStats != "" {
			parseProtocolStats(protocolStats, verbose)
		}
	}
	
	fmt.Println()
}

// parseProtocolStats 解析协议统计信息
func parseProtocolStats(output string, verbose bool) {
	lines := strings.Split(output, "\n")
	currentProtocol := ""
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 检查是否是协议标题
		if strings.HasSuffix(line, ":") && !strings.Contains(line, " ") {
			currentProtocol = strings.TrimSuffix(line, ":")
			if currentProtocol == "tcp" || currentProtocol == "udp" || currentProtocol == "ip" {
				color.Info.Printf("%s 统计:\n", strings.ToUpper(currentProtocol))
			}
			continue
		}
		
		// 显示重要的统计信息
		if currentProtocol == "tcp" || currentProtocol == "udp" || currentProtocol == "ip" {
			if verbose || isImportantStat(line) {
				color.Gray.Printf("  %s\n", line)
			}
		}
	}
}

// isImportantStat 判断是否是重要的统计信息
func isImportantStat(line string) bool {
	importantKeywords := []string{
		"packets sent",
		"packets received",
		"connections established",
		"connections failed",
		"packets dropped",
		"errors",
	}
	
	lineLower := strings.ToLower(line)
	for _, keyword := range importantKeywords {
		if strings.Contains(lineLower, keyword) {
			return true
		}
	}
	return false
}

// showRoutingInfo 显示路由信息
func showRoutingInfo() {
	color.Blue.Println("🛣️  路由信息")
	
	if runtime.GOOS == "darwin" {
		// 显示路由表
		routeOutput := getCommandOutput("netstat", "-rn")
		if routeOutput != "" {
			lines := strings.Split(routeOutput, "\n")
			for i, line := range lines {
				if i < 5 { // 只显示前几行重要路由
					if strings.TrimSpace(line) != "" {
						color.Gray.Printf("  %s\n", line)
					}
				} else {
					break
				}
			}
			color.Info.Println("  ... (使用 netstat -rn 查看完整路由表)")
		}
	}
	
	fmt.Println()
}

// formatBytes 格式化字节数
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}