package xcode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "设置 macOS 代码签名环境",
	Long:  color.Success.Render("\r\n设置和配置 macOS 代码签名环境，包括证书检查、密钥链配置和开发环境验证"),
	Run: func(cmd *cobra.Command, args []string) {
		showCertificates, _ := cmd.Flags().GetBool("show-certificates")
		showKeychain, _ := cmd.Flags().GetBool("show-keychain")
		showProfiles, _ := cmd.Flags().GetBool("show-profiles")
		showAll, _ := cmd.Flags().GetBool("all")
		verbose, _ := cmd.Flags().GetBool("verbose")
		
		// 显示标题
		showSetupHeader()
		
		// 如果指定了 --all，则显示所有信息
		if showAll {
			showCertificates = true
			showKeychain = true
			showProfiles = true
		}
		
		// 如果没有指定任何选项，显示基本信息
		if !showCertificates && !showKeychain && !showProfiles {
			showBasicSetupInfo(verbose)
		} else {
			// 显示指定的信息
			if showCertificates {
				showCertificateInfo(verbose)
			}
			
			if showKeychain {
				showKeychainInfo(verbose)
			}
			
			if showProfiles {
				showProvisioningProfiles(verbose)
			}
		}
		
		// 显示开发路线图
		showDevelopmentRoadmap("setup")
	},
}

func init() {
	setupCmd.Flags().Bool("show-certificates", false, "显示代码签名证书")
	setupCmd.Flags().Bool("show-keychain", false, "显示密钥链信息")
	setupCmd.Flags().Bool("show-profiles", false, "显示配置文件")
	setupCmd.Flags().Bool("all", false, "显示所有信息")
	setupCmd.Flags().BoolP("verbose", "v", false, "详细输出")
}

// showSetupHeader 显示设置标题
func showSetupHeader() {
	color.Blue.Println("===========================================")
	color.Blue.Println("      🔧 macOS 代码签名环境设置         ")
	color.Blue.Println("===========================================")
	fmt.Println()
}

// showBasicSetupInfo 显示基本设置信息
func showBasicSetupInfo(verbose bool) {
	color.Blue.Println("📋 基本环境信息")
	
	// 检查 Xcode
	checkXcodeInstallation(verbose)
	
	// 检查命令行工具
	checkCommandLineTools(verbose)
	
	// 检查代码签名证书
	checkSigningCertificates(verbose)
	
	// 检查密钥链
	checkKeychain(verbose)
	
	fmt.Println()
}

// checkXcodeInstallation 检查 Xcode 安装
func checkXcodeInstallation(verbose bool) {
	color.Info.Println("🔍 检查 Xcode 安装")
	
	// 检查 Xcode 路径
	xcodePath := getCommandOutput("xcode-select", "-p")
	if xcodePath != "" {
		color.Success.Printf("✅ Xcode 路径: %s\n", xcodePath)
		
		// 获取 Xcode 版本
		if xcodeVersion := getCommandOutput("xcodebuild", "-version"); xcodeVersion != "" {
			lines := strings.Split(xcodeVersion, "\n")
			if len(lines) > 0 {
				color.Info.Printf("版本: %s\n", lines[0])
			}
			if len(lines) > 1 {
				color.Info.Printf("构建版本: %s\n", lines[1])
			}
		}
	} else {
		color.Error.Println("❌ 未找到 Xcode 安装")
		color.Info.Println("💡 请从 App Store 安装 Xcode")
	}
	
	fmt.Println()
}

// checkCommandLineTools 检查命令行工具
func checkCommandLineTools(verbose bool) {
	color.Info.Println("🛠️  检查命令行工具")
	
	// 检查是否安装了命令行工具
	cmd := exec.Command("xcode-select", "--install")
	err := cmd.Run()
	if err != nil {
		// 如果返回错误，通常意味着已经安装了
		color.Success.Println("✅ 命令行工具已安装")
	} else {
		color.Yellow.Println("⚠️  命令行工具可能需要安装或更新")
	}
	
	// 检查关键工具
	tools := []string{"codesign", "security", "hdiutil", "plutil", "lipo"}
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err == nil {
			color.Success.Printf("✅ %s 可用\n", tool)
		} else {
			color.Error.Printf("❌ %s 不可用\n", tool)
		}
	}
	
	fmt.Println()
}

// checkSigningCertificates 检查代码签名证书
func checkSigningCertificates(verbose bool) {
	color.Info.Println("🔐 检查代码签名证书")
	
	// 查找开发证书
	devCerts := getCommandOutput("security", "find-identity", "-v", "-p", "codesigning")
	if devCerts != "" {
		lines := strings.Split(devCerts, "\n")
		devCount := 0
		distCount := 0
		
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.Contains(line, "valid identities found") {
				continue
			}
			
			if strings.Contains(line, "Developer ID Application") {
				distCount++
				if verbose {
					color.Success.Printf("✅ 分发证书: %s\n", extractCertificateName(line))
				}
			} else if strings.Contains(line, "Mac Developer") || strings.Contains(line, "Apple Development") {
				devCount++
				if verbose {
					color.Success.Printf("✅ 开发证书: %s\n", extractCertificateName(line))
				}
			}
		}
		
		color.Info.Printf("开发证书: %d 个\n", devCount)
		color.Info.Printf("分发证书: %d 个\n", distCount)
		
		if devCount == 0 && distCount == 0 {
			color.Yellow.Println("⚠️  未找到有效的代码签名证书")
			color.Info.Println("💡 请在 Xcode 中登录 Apple ID 或导入证书")
		}
	} else {
		color.Error.Println("❌ 无法获取证书信息")
	}
	
	fmt.Println()
}

// checkKeychain 检查密钥链
func checkKeychain(verbose bool) {
	color.Info.Println("🔑 检查密钥链")
	
	// 获取默认密钥链
	defaultKeychain := getCommandOutput("security", "default-keychain")
	if defaultKeychain != "" {
		defaultKeychain = strings.Trim(defaultKeychain, `"`)
		color.Info.Printf("默认密钥链: %s\n", defaultKeychain)
	}
	
	// 列出密钥链搜索列表
	keychainList := getCommandOutput("security", "list-keychains")
	if keychainList != "" && verbose {
		color.Info.Println("密钥链搜索列表:")
		lines := strings.Split(keychainList, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				line = strings.Trim(line, `"`)
				fmt.Printf("  - %s\n", line)
			}
		}
	}
	
	fmt.Println()
}

// showCertificateInfo 显示证书详细信息
func showCertificateInfo(verbose bool) {
	color.Blue.Println("🔐 代码签名证书详情")
	
	// 获取所有代码签名证书
	certOutput := getCommandOutput("security", "find-identity", "-v", "-p", "codesigning")
	if certOutput == "" {
		color.Error.Println("❌ 未找到代码签名证书")
		return
	}
	
	lines := strings.Split(certOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "valid identities found") {
			continue
		}
		
		// 提取证书信息
		if strings.Contains(line, ")") {
			parts := strings.Split(line, ")")
			if len(parts) >= 2 {
				hash := strings.TrimSpace(parts[0])
				hash = strings.TrimPrefix(hash, "1)")
				hash = strings.TrimPrefix(hash, "2)")
				hash = strings.TrimSpace(hash)
				
				certName := strings.TrimSpace(parts[1])
				
				if strings.Contains(certName, "Developer ID Application") {
					color.Success.Printf("📦 分发证书: %s\n", certName)
				} else if strings.Contains(certName, "Mac Developer") || strings.Contains(certName, "Apple Development") {
					color.Info.Printf("🛠️  开发证书: %s\n", certName)
				} else {
					color.Cyan.Printf("🔐 其他证书: %s\n", certName)
				}
				
				if verbose {
					color.Gray.Printf("   SHA-1: %s\n", hash)
				}
			}
		}
	}
	
	fmt.Println()
}

// showKeychainInfo 显示密钥链详细信息
func showKeychainInfo(verbose bool) {
	color.Blue.Println("🔑 密钥链详细信息")
	
	// 显示默认密钥链
	defaultKeychain := getCommandOutput("security", "default-keychain")
	if defaultKeychain != "" {
		defaultKeychain = strings.Trim(defaultKeychain, `"`)
		color.Info.Printf("默认密钥链: %s\n", defaultKeychain)
		
		// 检查密钥链状态
		if _, err := os.Stat(defaultKeychain); err == nil {
			color.Success.Println("✅ 密钥链文件存在")
		} else {
			color.Error.Println("❌ 密钥链文件不存在")
		}
	}
	
	// 显示密钥链搜索列表
	color.Info.Println("\n密钥链搜索列表:")
	keychainList := getCommandOutput("security", "list-keychains")
	if keychainList != "" {
		lines := strings.Split(keychainList, "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				line = strings.Trim(line, `"`)
				if i == 0 {
					color.Success.Printf("  1. %s (默认)\n", line)
				} else {
					color.Info.Printf("  %d. %s\n", i+1, line)
				}
				
				// 检查密钥链状态
				if verbose {
					if _, err := os.Stat(line); err == nil {
						color.Gray.Println("     ✅ 可访问")
					} else {
						color.Gray.Println("     ❌ 不可访问")
					}
				}
			}
		}
	} else {
		color.Error.Println("❌ 无法获取密钥链列表")
	}
	
	fmt.Println()
}

// showProvisioningProfiles 显示配置文件
func showProvisioningProfiles(verbose bool) {
	color.Blue.Println("📄 配置文件信息")
	
	// 配置文件路径
	profileDir := os.ExpandEnv("$HOME/Library/MobileDevice/Provisioning Profiles")
	
	if _, err := os.Stat(profileDir); os.IsNotExist(err) {
		color.Yellow.Println("⚠️  配置文件目录不存在")
		color.Info.Printf("路径: %s\n", profileDir)
		return
	}
	
	// 列出配置文件
	files, err := os.ReadDir(profileDir)
	if err != nil {
		color.Error.Printf("❌ 无法读取配置文件目录: %v\n", err)
		return
	}
	
	profileCount := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mobileprovision") {
			profileCount++
			if verbose {
				profilePath := filepath.Join(profileDir, file.Name())
				showProvisioningProfileDetails(profilePath)
			}
		}
	}
	
	color.Info.Printf("配置文件总数: %d\n", profileCount)
	
	if profileCount == 0 {
		color.Yellow.Println("⚠️  未找到配置文件")
		color.Info.Println("💡 请在 Xcode 中下载配置文件")
	}
	
	fmt.Println()
}

// showProvisioningProfileDetails 显示配置文件详情
func showProvisioningProfileDetails(profilePath string) {
	// 使用 security 命令解析配置文件
	profileInfo := getCommandOutput("security", "cms", "-D", "-i", profilePath)
	if profileInfo == "" {
		return
	}
	
	// 提取基本信息
	lines := strings.Split(profileInfo, "\n")
	profileName := ""
	teamName := ""
	appID := ""
	
	for _, line := range lines {
		if strings.Contains(line, "<key>Name</key>") {
			// 下一行包含名称
			continue
		}
		if strings.Contains(line, "<string>") && profileName == "" {
			profileName = extractStringValue(line)
		}
		if strings.Contains(line, "<key>TeamName</key>") {
			// 下一行包含团队名称
			continue
		}
		if strings.Contains(line, "<key>application-identifier</key>") {
			// 下一行包含应用 ID
			continue
		}
	}
	
	filename := filepath.Base(profilePath)
	color.Cyan.Printf("📄 %s\n", filename)
	if profileName != "" {
		color.Info.Printf("   名称: %s\n", profileName)
	}
	if teamName != "" {
		color.Info.Printf("   团队: %s\n", teamName)
	}
	if appID != "" {
		color.Info.Printf("   应用ID: %s\n", appID)
	}
}

// extractCertificateName 提取证书名称
func extractCertificateName(line string) string {
	if idx := strings.Index(line, ")"); idx != -1 && idx+1 < len(line) {
		return strings.TrimSpace(line[idx+1:])
	}
	return line
}

// extractStringValue 提取字符串值
func extractStringValue(line string) string {
	start := strings.Index(line, "<string>")
	end := strings.Index(line, "</string>")
	if start != -1 && end != -1 && start+8 < end {
		return line[start+8 : end]
	}
	return ""
}