package system

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "显示系统进程信息",
	Long:  color.Success.Render("\r\n显示系统进程信息，包括进程列表、资源占用等"),
	Run: func(cmd *cobra.Command, args []string) {
		top, _ := cmd.Flags().GetInt("top")
		sortBy, _ := cmd.Flags().GetString("sort")
		filter, _ := cmd.Flags().GetString("filter")
		verbose, _ := cmd.Flags().GetBool("verbose")
		
		// 显示标题
		showProcessHeader()
		
		// 显示进程统计
		showProcessStats()
		
		// 显示进程列表
		showProcessList(top, sortBy, filter, verbose)
	},
}

func init() {
	processCmd.Flags().IntP("top", "t", 10, "显示前N个进程")
	processCmd.Flags().StringP("sort", "s", "cpu", "排序方式 (cpu, memory, pid, name)")
	processCmd.Flags().StringP("filter", "f", "", "过滤进程名称")
	processCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
}

// ProcessInfo 进程信息结构
type ProcessInfo struct {
	PID     int
	Name    string
	CPU     float64
	Memory  float64
	User    string
	Command string
}

// showProcessHeader 显示进程信息标题
func showProcessHeader() {
	color.Blue.Println("===========================================")
	color.Blue.Println("         🔄 系统进程信息              ")
	color.Blue.Println("===========================================")
	fmt.Println()
}

// showProcessStats 显示进程统计信息
func showProcessStats() {
	color.Blue.Println("📊 进程统计")
	
	if runtime.GOOS == "darwin" {
		// 获取进程总数
		if psOutput := getCommandOutput("ps", "aux"); psOutput != "" {
			lines := strings.Split(psOutput, "\n")
			processCount := len(lines) - 1 // 减去标题行
			if processCount > 0 {
				color.Info.Printf("总进程数: %d\n", processCount)
			}
		}
		
		// 获取运行状态统计
		if psStateOutput := getCommandOutput("ps", "-eo", "stat"); psStateOutput != "" {
			lines := strings.Split(psStateOutput, "\n")
			stateCounts := make(map[string]int)
			
			for i, line := range lines {
				if i == 0 { // 跳过标题行
					continue
				}
				state := strings.TrimSpace(line)
				if state != "" {
					// 取状态的第一个字符
					if len(state) > 0 {
						mainState := string(state[0])
						stateCounts[mainState]++
					}
				}
			}
			
			// 显示状态统计
			stateNames := map[string]string{
				"R": "运行中",
				"S": "睡眠",
				"I": "空闲",
				"T": "停止",
				"Z": "僵尸",
				"U": "不可中断",
			}
			
			for state, count := range stateCounts {
				if name, exists := stateNames[state]; exists {
					color.Info.Printf("%s: %d\n", name, count)
				} else {
					color.Info.Printf("%s: %d\n", state, count)
				}
			}
		}
	}
	
	fmt.Println()
}

// showProcessList 显示进程列表
func showProcessList(top int, sortBy, filter string, verbose bool) {
	color.Blue.Printf("🔍 进程列表 (前 %d 个，按 %s 排序)\n", top, sortBy)
	
	processes := getProcessList()
	if len(processes) == 0 {
		color.Error.Println("❌ 无法获取进程信息")
		return
	}
	
	// 过滤进程
	if filter != "" {
		filteredProcesses := []ProcessInfo{}
		for _, proc := range processes {
			if strings.Contains(strings.ToLower(proc.Name), strings.ToLower(filter)) ||
			   strings.Contains(strings.ToLower(proc.Command), strings.ToLower(filter)) {
				filteredProcesses = append(filteredProcesses, proc)
			}
		}
		processes = filteredProcesses
		color.Info.Printf("过滤结果: %d 个进程\n", len(processes))
	}
	
	// 排序进程
	sortProcesses(processes, sortBy)
	
	// 限制显示数量
	if top > 0 && top < len(processes) {
		processes = processes[:top]
	}
	
	// 显示表头
	fmt.Println()
	if verbose {
		color.Yellow.Printf("%-8s %-20s %-8s %-8s %-10s %s\n", "PID", "进程名", "CPU%", "内存%", "用户", "命令")
		color.Yellow.Println(strings.Repeat("-", 80))
	} else {
		color.Yellow.Printf("%-8s %-25s %-8s %-8s %s\n", "PID", "进程名", "CPU%", "内存%", "用户")
		color.Yellow.Println(strings.Repeat("-", 60))
	}
	
	// 显示进程信息
	for _, proc := range processes {
		// 根据资源使用情况着色
		var cpuColor, memColor func(a ...interface{}) string
		
		if proc.CPU > 50 {
			cpuColor = color.Red.Sprint
		} else if proc.CPU > 20 {
			cpuColor = color.Yellow.Sprint
		} else {
			cpuColor = color.Green.Sprint
		}
		
		if proc.Memory > 10 {
			memColor = color.Red.Sprint
		} else if proc.Memory > 5 {
			memColor = color.Yellow.Sprint
		} else {
			memColor = color.Green.Sprint
		}
		
		if verbose {
			// 截断长命令
			command := proc.Command
			if len(command) > 30 {
				command = command[:27] + "..."
			}
			fmt.Printf("%-8d %-20s %s %-8s %-10s %s\n",
				proc.PID,
				truncateString(proc.Name, 20),
				cpuColor(fmt.Sprintf("%-8.1f", proc.CPU)),
				memColor(fmt.Sprintf("%.1f", proc.Memory)),
				truncateString(proc.User, 10),
				command)
		} else {
			fmt.Printf("%-8d %-25s %s %s %s\n",
				proc.PID,
				truncateString(proc.Name, 25),
				cpuColor(fmt.Sprintf("%-8.1f", proc.CPU)),
				memColor(fmt.Sprintf("%-8.1f", proc.Memory)),
				truncateString(proc.User, 10))
		}
	}
	
	fmt.Println()
}

// getProcessList 获取进程列表
func getProcessList() []ProcessInfo {
	var processes []ProcessInfo
	
	if runtime.GOOS == "darwin" {
		// 使用 ps 命令获取进程信息
		output := getCommandOutput("ps", "aux")
		if output == "" {
			return processes
		}
		
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			if i == 0 { // 跳过标题行
				continue
			}
			
			fields := strings.Fields(line)
			if len(fields) >= 11 {
				pid, _ := strconv.Atoi(fields[1])
				cpu, _ := strconv.ParseFloat(fields[2], 64)
				mem, _ := strconv.ParseFloat(fields[3], 64)
				user := fields[0]
				
				// 进程名通常在第11个字段
				name := fields[10]
				if strings.HasPrefix(name, "[") && strings.HasSuffix(name, "]") {
					// 内核进程
					name = strings.Trim(name, "[]")
				}
				
				// 完整命令
				command := strings.Join(fields[10:], " ")
				
				processes = append(processes, ProcessInfo{
					PID:     pid,
					Name:    name,
					CPU:     cpu,
					Memory:  mem,
					User:    user,
					Command: command,
				})
			}
		}
	}
	
	return processes
}

// sortProcesses 排序进程列表
func sortProcesses(processes []ProcessInfo, sortBy string) {
	switch sortBy {
	case "cpu":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPU > processes[j].CPU
		})
	case "memory", "mem":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].Memory > processes[j].Memory
		})
	case "pid":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].PID < processes[j].PID
		})
	case "name":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].Name < processes[j].Name
		})
	default:
		// 默认按CPU排序
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPU > processes[j].CPU
		})
	}
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}