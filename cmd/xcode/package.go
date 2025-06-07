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

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "创建 DMG 安装包",
	Long:  color.Success.Render("\r\n为 macOS 应用程序创建 DMG 安装包，支持自动检测应用程序路径和自定义输出名称"),
	Run: func(cmd *cobra.Command, args []string) {
		scheme, _ := cmd.Flags().GetString("scheme")
		buildPath, _ := cmd.Flags().GetString("build-path")
		outputDir, _ := cmd.Flags().GetString("output")
		dmgName, _ := cmd.Flags().GetString("name")
		includeArch, _ := cmd.Flags().GetBool("include-arch")
		verbose, _ := cmd.Flags().GetBool("verbose")
		useCreateDmg, _ := cmd.Flags().GetBool("use-create-dmg")
		
		// 显示配置信息
		showPackageConfig(scheme, buildPath, outputDir, dmgName, includeArch, verbose)
		
		// 自动检测 SCHEME
		if scheme == "" {
			scheme = detectScheme()
			if scheme == "" {
				color.Error.Println("❌ 错误: 未设置 SCHEME 且无法自动检测")
				showAvailableSchemes()
				os.Exit(1)
			}
		}
		
		// 设置默认值
		if buildPath == "" {
			buildPath = "./temp/Build/Products/Release"
		}
		if outputDir == "" {
			outputDir = "./temp"
		}
		
		// 检查依赖
		err := checkDependencies(useCreateDmg)
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(1)
		}
		
		// 检查应用程序
		appPath, err := checkApplication(buildPath, scheme)
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(1)
		}
		
		// 创建 DMG 文件
		dmgFiles, err := createDMGFile(appPath, outputDir, dmgName, scheme, includeArch, useCreateDmg, verbose)
		if err != nil {
			color.Error.Printf("❌ 创建 DMG 失败: %s\n", err.Error())
			os.Exit(1)
		}
		
		// 显示结果
		showResults(dmgFiles)
		
		// 显示开发路线图
		showDevelopmentRoadmap("package")
	},
}

func init() {
	packageCmd.Flags().StringP("scheme", "s", "", "应用程序方案名称")
	packageCmd.Flags().StringP("build-path", "b", "./temp/Build/Products/Release", "构建产物路径")
	packageCmd.Flags().StringP("output", "o", "./temp", "DMG 输出目录")
	packageCmd.Flags().StringP("name", "n", "", "DMG 文件名称")
	packageCmd.Flags().Bool("include-arch", true, "是否在文件名中包含架构信息")
	packageCmd.Flags().BoolP("verbose", "v", false, "详细日志输出")
	packageCmd.Flags().Bool("use-create-dmg", false, "使用 create-dmg 工具（需要 npm 安装）")
}

// showPackageConfig 显示配置信息
func showPackageConfig(scheme, buildPath, outputDir, dmgName string, includeArch, verbose bool) {
	color.Blue.Println("===========================================")
	color.Blue.Println("         🚀 DMG 创建脚本                ")
	color.Blue.Println("===========================================")
	fmt.Println()
	
	color.Blue.Println("⚙️  配置信息")
	color.Info.Printf("应用方案: %s\n", scheme)
	color.Info.Printf("构建路径: %s\n", buildPath)
	color.Info.Printf("输出目录: %s\n", outputDir)
	if dmgName != "" {
		color.Info.Printf("DMG 名称: %s\n", dmgName)
	} else {
		color.Info.Println("DMG 名称: 自动生成")
	}
	color.Info.Printf("包含架构: %t\n", includeArch)
	color.Info.Printf("详细日志: %t\n", verbose)
	fmt.Println()
}

// checkDependencies 检查依赖
func checkDependencies(useCreateDmg bool) error {
	color.Blue.Println("🔍 检查依赖工具")
	
	// 检查 hdiutil（macOS 原生工具）
	if _, err := exec.LookPath("hdiutil"); err != nil {
		return fmt.Errorf("未找到 hdiutil 工具")
	}
	color.Success.Println("✅ hdiutil 可用")
	
	// 如果使用 create-dmg，检查是否安装
	if useCreateDmg {
		if _, err := exec.LookPath("create-dmg"); err != nil {
			color.Yellow.Println("⚠️  create-dmg 未安装，将使用 hdiutil")
			color.Info.Println("💡 安装 create-dmg: npm install -g create-dmg")
		} else {
			color.Success.Println("✅ create-dmg 可用")
		}
	}
	
	fmt.Println()
	return nil
}

// checkApplication 检查应用程序
func checkApplication(buildPath, scheme string) (string, error) {
	color.Blue.Println("📱 检查应用程序")
	
	appPath := filepath.Join(buildPath, scheme+".app")
	
	// 检查应用是否存在
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		// 搜索可能的应用程序位置
		foundApps := searchForApps(scheme)
		
		if len(foundApps) == 0 {
			return "", fmt.Errorf("未找到应用程序: %s", appPath)
		}
		
		color.Info.Printf("📍 发现 %d 个可能的应用程序:\n", len(foundApps))
		for i, app := range foundApps {
			appSize := "未知"
			if sizeOutput := getCommandOutput("du", "-sh", app); sizeOutput != "" {
				parts := strings.Fields(sizeOutput)
				if len(parts) > 0 {
					appSize = parts[0]
				}
			}
			fmt.Printf("   %d. %s (%s)\n", i+1, app, appSize)
		}
		fmt.Println()
		color.Info.Println("💡 建议: 请设置 BuildPath 环境变量指向正确的构建目录，例如:")
		fmt.Println()
		for _, app := range foundApps {
			buildDir := filepath.Dir(app)
			fmt.Printf(" go run main.go xcode package --build-path '%s'\n", buildDir)
		}
		fmt.Println()
		
		return "", fmt.Errorf("请先运行构建脚本: go run main.go xcode build")
	}
	
	// 显示应用信息
	showAppInfo(appPath, scheme)
	
	// 检测架构
	detectArchitecture(appPath)
	
	return appPath, nil
}

// searchForApps 搜索应用程序
func searchForApps(scheme string) []string {
	possiblePaths := []string{
		fmt.Sprintf("./temp/Build/Products/Release/%s.app", scheme),
		fmt.Sprintf("./temp/Build/Products/Debug/%s.app", scheme),
		fmt.Sprintf("./Build/Products/Release/%s.app", scheme),
		fmt.Sprintf("./Build/Products/Debug/%s.app", scheme),
		fmt.Sprintf("./build/Release/%s.app", scheme),
		fmt.Sprintf("./build/Debug/%s.app", scheme),
		fmt.Sprintf("./DerivedData/Build/Products/Release/%s.app", scheme),
		fmt.Sprintf("./DerivedData/Build/Products/Debug/%s.app", scheme),
	}
	
	foundApps := []string{}
	
	// 检查预定义路径
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			foundApps = append(foundApps, path)
		}
	}
	
	// 使用 find 命令搜索更多可能的位置
	cmd := exec.Command("find", ".", "-name", scheme+".app", "-type", "d", "-not", "-path", "*/.*")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				// 避免重复添加
				alreadyFound := false
				for _, existing := range foundApps {
					if existing == line {
						alreadyFound = true
						break
					}
				}
				if !alreadyFound && len(foundApps) < 20 {
					foundApps = append(foundApps, line)
				}
			}
		}
	}
	
	return foundApps
}

// detectArchitecture 检测架构
func detectArchitecture(appPath string) {
	executablePath := filepath.Join(appPath, "Contents/MacOS")
	
	// 查找可执行文件
	files, err := os.ReadDir(executablePath)
	if err != nil {
		return
	}
	
	for _, file := range files {
		if !file.IsDir() {
			execFile := filepath.Join(executablePath, file.Name())
			if archOutput := getCommandOutput("lipo", "-archs", execFile); archOutput != "" {
				color.Info.Printf("应用架构: %s\n", archOutput)
				break
			}
		}
	}
}

// createDMGFile 创建 DMG 文件
func createDMGFile(appPath, outputDir, dmgName, scheme string, includeArch, useCreateDmg, verbose bool) ([]string, error) {
	color.Blue.Println("📦 创建 DMG 安装包")
	
	// 设置输出目录
	if outputDir != "." {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("无法创建输出目录: %v", err)
		}
		
		// 切换到输出目录
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(outputDir)
		appPath = "../" + appPath
	}
	
	// 选择创建方法
	if useCreateDmg {
		if _, err := exec.LookPath("create-dmg"); err == nil {
			color.Info.Println("创建方法: create-dmg (npm)")
			return createDMGWithCreateDmg(appPath, dmgName, scheme, includeArch, verbose)
		}
	}
	
	color.Info.Println("创建方法: hdiutil (原生)")
	return createDMGWithHdiutil(appPath, dmgName, scheme, includeArch, verbose)
}

// createDMGWithHdiutil 使用 hdiutil 创建 DMG
func createDMGWithHdiutil(appPath, dmgName, scheme string, includeArch, verbose bool) ([]string, error) {
	finalDMG := generateDMGFilename(dmgName, scheme, includeArch, appPath)
	tempDMG := strings.Replace(finalDMG, ".dmg", "-temp.dmg", 1)
	
	// 创建临时 DMG
	args := []string{"create", "-srcfolder", appPath, "-format", "UDRW", "-volname", scheme, tempDMG}
	err := executeCommand("hdiutil", args, "创建临时 DMG", verbose)
	if err != nil {
		return nil, err
	}
	
	// 挂载 DMG
	mountOutput, err := exec.Command("hdiutil", "attach", tempDMG, "-readwrite", "-noverify", "-noautoopen").Output()
	if err != nil {
		return nil, fmt.Errorf("挂载 DMG 失败: %v", err)
	}
	
	// 解析挂载点
	mountPoint := ""
	lines := strings.Split(string(mountOutput), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/Volumes/") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "/Volumes/") {
					mountPoint = part
					break
				}
			}
			break
		}
	}
	
	if mountPoint == "" {
		return nil, fmt.Errorf("无法找到挂载点")
	}
	
	// 创建应用程序快捷方式
	err = executeCommand("ln", []string{"-s", "/Applications", filepath.Join(mountPoint, "Applications")}, "创建 Applications 快捷方式", verbose)
	if err != nil {
		// 卸载 DMG
		exec.Command("hdiutil", "detach", mountPoint).Run()
		return nil, err
	}
	
	// 卸载 DMG
	err = executeCommand("hdiutil", []string{"detach", mountPoint}, "卸载 DMG", verbose)
	if err != nil {
		return nil, err
	}
	
	// 压缩为最终文件名
	err = executeCommand("hdiutil", []string{"convert", tempDMG, "-format", "UDZO", "-imagekey", "zlib-level=9", "-o", finalDMG}, "压缩 DMG", verbose)
	if err != nil {
		return nil, err
	}
	
	// 删除临时文件
	os.Remove(tempDMG)
	
	return []string{finalDMG}, nil
}

// createDMGWithCreateDmg 使用 create-dmg 创建 DMG
func createDMGWithCreateDmg(appPath, dmgName, scheme string, includeArch, verbose bool) ([]string, error) {
	finalDMG := generateDMGFilename(dmgName, scheme, includeArch, appPath)
	
	// 替换空格为连字符
	finalDMG = strings.ReplaceAll(finalDMG, " ", "-")
	
	// 使用 --overwrite 参数创建 DMG，避免 "Target already exists" 错误
	err := executeCommand("create-dmg", []string{"--overwrite", appPath}, "生成 DMG 文件", verbose)
	if err != nil {
		return nil, err
	}
	
	// 查找生成的 DMG 文件并重命名
	files, err := filepath.Glob("*.dmg")
	if err != nil {
		return nil, fmt.Errorf("查找 DMG 文件失败: %v", err)
	}
	
	for _, file := range files {
		if file != finalDMG {
			err = os.Rename(file, finalDMG)
			if err != nil {
				return nil, fmt.Errorf("重命名 DMG 文件失败: %v", err)
			}
			break
		}
	}
	
	return []string{finalDMG}, nil
}

// generateDMGFilename 生成 DMG 文件名
func generateDMGFilename(dmgName, scheme string, includeArch bool, appPath string) string {
	if dmgName != "" {
		if !strings.HasSuffix(dmgName, ".dmg") {
			dmgName += ".dmg"
		}
		return dmgName
	}
	
	// 获取版本信息
	version := ""
	infoPath := filepath.Join(appPath, "Contents/Info.plist")
	if versionOutput := getCommandOutput("plutil", "-p", infoPath); versionOutput != "" {
		lines := strings.Split(versionOutput, "\n")
		for _, line := range lines {
			if strings.Contains(line, "CFBundleShortVersionString") {
				parts := strings.Split(line, `"`)
				if len(parts) >= 4 {
					version = parts[3]
					break
				}
			}
		}
	}
	
	// 获取架构信息
	arch := ""
	if includeArch {
		executablePath := filepath.Join(appPath, "Contents/MacOS")
		files, err := os.ReadDir(executablePath)
		if err == nil {
			for _, file := range files {
				if !file.IsDir() {
					execFile := filepath.Join(executablePath, file.Name())
					if archOutput := getCommandOutput("lipo", "-archs", execFile); archOutput != "" {
						if strings.Contains(archOutput, "x86_64") && strings.Contains(archOutput, "arm64") {
							arch = "universal"
						} else {
							arch = strings.TrimSpace(archOutput)
						}
						break
					}
				}
			}
		}
	}
	
	// 构建文件名
	filename := scheme
	if version != "" {
		filename += "-" + version
	}
	if arch != "" {
		filename += "-" + arch
	}
	filename += ".dmg"
	
	return filename
}

// executeCommand 执行命令
func executeCommand(command string, args []string, description string, verbose bool) error {
	if verbose {
		color.Blue.Printf("🔧 %s\n", description)
		color.Cyan.Printf("命令: %s %s\n", command, strings.Join(args, " "))
	}
	
	cmd := exec.Command(command, args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s失败: %v", description, err)
	}
	
	if verbose {
		color.Success.Printf("✅ %s 完成\n", description)
	}
	
	return nil
}

// showResults 显示结果
func showResults(dmgFiles []string) {
	color.Blue.Println("📋 DMG 创建结果")
	
	for _, dmgFile := range dmgFiles {
		if _, err := os.Stat(dmgFile); err == nil {
			fileSize := "未知"
			if sizeOutput := getCommandOutput("ls", "-lh", dmgFile); sizeOutput != "" {
				parts := strings.Fields(sizeOutput)
				if len(parts) >= 5 {
					fileSize = parts[4]
				}
			}
			color.Info.Printf("%s: %s\n", dmgFile, fileSize)
		}
	}
	
	fmt.Println()
	color.Success.Println("✅ DMG 安装包创建完成！")
}