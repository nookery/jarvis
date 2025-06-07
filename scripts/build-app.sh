#!/bin/bash

# ====================================
# 通用 iOS/macOS 应用构建脚本
# ====================================
#
# 这个脚本用于构建 iOS/macOS 应用程序，在构建前会显示详细的环境信息，
# 帮助开发者了解当前的构建环境状态，便于调试和问题排查。
#
# 功能：
# 1. 显示系统环境信息（操作系统、架构、主机名等）
# 2. 显示 Xcode 开发环境信息（版本、SDK 路径等）
# 3. 显示 Swift 编译器信息
# 4. 显示 Git 版本控制信息（版本、分支、最新提交等）
# 5. 显示构建环境变量
# 6. 显示构建目标信息（项目、方案、支持的架构等）
# 7. 执行 xcodebuild 构建命令
# 8. 显示构建结果和产物位置
#
# 使用方法：
# 1. 设置必要的环境变量：
#    export SCHEME="YourAppScheme"             # 构建方案名称
#    export BuildPath="/path/to/build"        # 构建输出路径（可选，默认为 ./temp）
#    export ARCH="universal"                  # 目标架构（可选，支持 universal、x86_64、arm64，默认为 universal）
#    export VERBOSE="true"                    # 可选：显示详细构建日志
#
# 2. 在项目根目录运行脚本：
#    ./scripts/build-app.sh
#
# 3. 启用详细日志模式：
#    VERBOSE=true ./scripts/build-app.sh
#
# 注意事项：
# - 需要安装 Xcode 和命令行工具
# - 需要在 Xcode 项目根目录下运行
# - 确保 SCHEME 和 BuildPath 环境变量已正确设置
# - 脚本会执行 clean build，会清除之前的构建缓存
# - 脚本会自动检测项目文件（.xcodeproj 或 .xcworkspace）
#
# 输出：
# - 详细的环境信息报告
# - 构建过程的实时输出
# - 构建结果和产物位置
# - 如果构建失败，脚本会以非零状态码退出
# ====================================

# 设置错误处理
set -e

# 检查必需的环境变量
if [ -z "$SCHEME" ]; then
    printf "\033[31m错误: 未设置 SCHEME 环境变量\033[0m\n"
    
    # 尝试列出项目中可用的 scheme
    printf "\033[33m正在检查项目中可用的 scheme...\033[0m\n"
    
    # 查找项目文件
    PROJECT_FILE=""
    if [ -n "$(find . -maxdepth 1 -name '*.xcworkspace' -type d)" ]; then
        PROJECT_FILE=$(find . -maxdepth 1 -name '*.xcworkspace' -type d | head -n 1)
        PROJECT_TYPE="workspace"
    elif [ -n "$(find . -maxdepth 1 -name '*.xcodeproj' -type d)" ]; then
        PROJECT_FILE=$(find . -maxdepth 1 -name '*.xcodeproj' -type d | head -n 1)
        PROJECT_TYPE="project"
    fi
    
    if [ -n "$PROJECT_FILE" ]; then
        printf "\033[32m在项目 ${PROJECT_FILE} 中找到以下可用的 scheme:\033[0m\n"
        
        if [ "$PROJECT_TYPE" = "workspace" ]; then
            SCHEMES=$(xcodebuild -workspace "$PROJECT_FILE" -list 2>/dev/null | grep -A 100 "Schemes:" | grep -v "Schemes:" | grep -v "^$" | sed 's/^[[:space:]]*//' | head -20)
        else
            SCHEMES=$(xcodebuild -project "$PROJECT_FILE" -list 2>/dev/null | grep -A 100 "Schemes:" | grep -v "Schemes:" | grep -v "^$" | sed 's/^[[:space:]]*//' | head -20)
        fi
        
        if [ -n "$SCHEMES" ]; then
            echo "$SCHEMES" | while read -r scheme; do
                if [ -n "$scheme" ]; then
                    printf "   - %s\n" "$scheme"
                fi
            done
            printf "\n\033[36m请选择一个 scheme 并设置环境变量，例如:\033[0m\n"
            FIRST_SCHEME=$(echo "$SCHEMES" | head -n 1 | sed 's/^[[:space:]]*//')
            if [ -n "$FIRST_SCHEME" ]; then
                printf "export SCHEME=\"%s\"\n" "$FIRST_SCHEME"
            fi
        else
            printf "   \033[31m未找到可用的 scheme\033[0m\n"
            printf "请设置 SCHEME 环境变量，例如: export SCHEME=\"YourAppScheme\"\n"
        fi
    else
        printf "   \033[31m未找到 .xcodeproj 或 .xcworkspace 文件\033[0m\n"
        printf "请设置 SCHEME 环境变量，例如: export SCHEME=\"YourAppScheme\"\n"
    fi
    
    exit 1
fi

# 设置默认构建路径
if [ -z "$BuildPath" ]; then
    BuildPath="./temp"
fi

# 设置默认架构
if [ -z "$ARCH" ]; then
    ARCH="universal"
fi

# 根据架构设置构建目标和路径
case "$ARCH" in
    "x86_64")
        DESTINATION="platform=macOS,arch=x86_64"
        BuildPath="${BuildPath}/x86_64"
        ;;
    "arm64")
        DESTINATION="platform=macOS,arch=arm64"
        BuildPath="${BuildPath}/arm64"
        ;;
    "universal")
        DESTINATION="platform=macOS"
        ARCHS="x86_64 arm64"
        BuildPath="${BuildPath}/universal"
        ;;
    *)
        printf "${RED}错误: 不支持的架构 '$ARCH'。支持的架构: universal, x86_64, arm64${NC}\n"
        exit 1
        ;;
esac

# 输出颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # 无颜色

# 显示开发路线图
show_development_roadmap() {
    local current_step="$1"
    
    echo
    printf "${PURPLE}===========================================${NC}\n"
    printf "${PURPLE}         🗺️  开发分发路线图                ${NC}\n"
    printf "${PURPLE}===========================================${NC}\n"
    echo
    
    # 定义路线图步骤
    local steps=(
        "build:🔨 构建应用:编译源代码，生成可执行文件"
        "codesign:🔐 代码签名:为应用添加数字签名，确保安全性"
        "package:📦 打包分发:创建 DMG 安装包"
        "notarize:✅ 公证验证:Apple 官方验证（可选）"
        "distribute:🚀 发布分发:上传到分发平台或直接分发"
    )
    
    printf "${CYAN}📍 当前位置: "
    case "$current_step" in
        "build") printf "${GREEN}构建应用${NC}\n" ;;
        "codesign") printf "${GREEN}代码签名${NC}\n" ;;
        "package") printf "${GREEN}打包分发${NC}\n" ;;
        "notarize") printf "${GREEN}公证验证${NC}\n" ;;
        "distribute") printf "${GREEN}发布分发${NC}\n" ;;
        *) printf "${YELLOW}未知步骤${NC}\n" ;;
    esac
    echo
    
    # 显示路线图
    for step in "${steps[@]}"; do
        local step_id=$(echo "$step" | cut -d':' -f1)
        local step_icon=$(echo "$step" | cut -d':' -f2)
        local step_desc=$(echo "$step" | cut -d':' -f3)
        
        if [ "$step_id" = "$current_step" ]; then
            printf "${GREEN}▶ %s %s${NC}\n" "$step_icon" "$step_desc"
        else
            printf "  %s %s\n" "$step_icon" "$step_desc"
        fi
    done
    
    echo
    printf "${YELLOW}💡 下一步建议:${NC}\n"
    case "$current_step" in
        "build")
            printf "   运行代码签名: ${CYAN}./scripts/codesign-app.sh${NC}\n"
            ;;
        "codesign")
            printf "   创建安装包: ${CYAN}./scripts/create-dmg.sh${NC}\n"
            ;;
        "package")
            printf "   进行公证验证或直接分发应用\n"
            ;;
        "notarize")
            printf "   发布到分发平台或提供下载链接\n"
            ;;
        "distribute")
            printf "   🎉 开发分发流程已完成！\n"
            ;;
    esac
    
    echo
    printf "${PURPLE}===========================================${NC}\n"
}

printf "${BLUE}===========================================${NC}\n"
printf "${BLUE}         应用构建环境信息                ${NC}\n"
printf "${BLUE}===========================================${NC}\n"
echo

# 系统信息
printf "${GREEN}📱 系统信息:${NC}\n"
echo "   操作系统: $(uname -s) $(uname -r)"
echo "   系统架构: $(uname -m)"
echo "   主机名称: $(hostname)"
echo

# Xcode 信息
printf "${GREEN}🔨 Xcode 开发环境:${NC}\n"
if command -v xcodebuild &> /dev/null; then
    echo "   Xcode 版本: $(xcodebuild -version | head -n 1)"
    echo "   构建版本: $(xcodebuild -version | tail -n 1)"
    echo "   SDK 路径: $(xcrun --show-sdk-path)"
    echo "   开发者目录: $(xcode-select -p)"
else
    printf "   ${RED}❌ 未找到 Xcode${NC}\n"
    exit 1
fi
echo

# Swift 信息
printf "${GREEN}🚀 Swift 编译器:${NC}\n"
if command -v swift &> /dev/null; then
    SWIFT_VERSION=$(swift --version 2>/dev/null | grep -o 'Swift version [0-9]\+\.[0-9]\+\.[0-9]\+' | cut -d' ' -f3)
    echo "   Swift 版本: ${SWIFT_VERSION}"
else
    printf "   ${RED}❌ 未找到 Swift${NC}\n"
fi
echo

# Git 信息
printf "${GREEN}📝 Git 版本控制:${NC}\n"
if command -v git &> /dev/null; then
    echo "   Git 版本: $(git --version)"
    if git rev-parse --git-dir > /dev/null 2>&1; then
        echo "   当前分支: $(git branch --show-current)"
        echo "   最新提交: $(git log -1 --pretty=format:'%h - %s (%an, %ar)')"
    fi
else
    printf "   ${RED}❌ 未找到 Git${NC}\n"
fi
echo

# 环境变量
printf "${GREEN}🌍 构建环境变量:${NC}\n"
echo "   构建方案: ${SCHEME}"
echo "   构建路径: ${BuildPath}"
echo "   目标架构: ${ARCH}"
echo "   构建目标: ${DESTINATION}"
if [ -n "$ARCHS" ]; then
    echo "   支持架构: ${ARCHS}"
fi
echo "   构建配置: Release"
echo "   详细日志: ${VERBOSE:-'false'}"
echo "   工作目录: $(pwd)"
echo

# 构建目标信息
printf "${GREEN}🎯 构建目标信息:${NC}\n"

# 自动检测项目文件
PROJECT_FILE=""
if [ -n "$(find . -maxdepth 1 -name '*.xcworkspace' -type d)" ]; then
    PROJECT_FILE=$(find . -maxdepth 1 -name '*.xcworkspace' -type d | head -n 1)
    PROJECT_TYPE="workspace"
    echo "   项目文件: ${PROJECT_FILE}"
    echo "   项目类型: Xcode Workspace"
elif [ -n "$(find . -maxdepth 1 -name '*.xcodeproj' -type d)" ]; then
    PROJECT_FILE=$(find . -maxdepth 1 -name '*.xcodeproj' -type d | head -n 1)
    PROJECT_TYPE="project"
    echo "   项目文件: ${PROJECT_FILE}"
    echo "   项目类型: Xcode Project"
else
    printf "   ${RED}❌ 未找到 .xcodeproj 或 .xcworkspace 文件${NC}\n"
    exit 1
fi

echo "   构建方案: ${SCHEME}"

# 显示支持的架构
if [ "${PROJECT_TYPE}" = "workspace" ]; then
    PROJECT_ARCHS=$(xcodebuild -workspace "${PROJECT_FILE}" -scheme "${SCHEME}" -showBuildSettings -configuration Release 2>/dev/null | grep 'ARCHS =' | head -n 1 | cut -d'=' -f2 | xargs || echo '无法确定')
else
    PROJECT_ARCHS=$(xcodebuild -project "${PROJECT_FILE}" -scheme "${SCHEME}" -showBuildSettings -configuration Release 2>/dev/null | grep 'ARCHS =' | head -n 1 | cut -d'=' -f2 | xargs || echo '无法确定')
fi
echo "   项目支持架构: ${PROJECT_ARCHS}"
if [ -n "$ARCHS" ]; then
    echo "   构建目标架构: ${ARCHS}"
else
    echo "   构建目标架构: ${ARCH}"
fi
echo

printf "${BLUE}===========================================${NC}\n"
printf "${YELLOW}🚀 开始构建过程...${NC}\n"
printf "${BLUE}===========================================${NC}\n"
echo

# 开始构建
printf "${GREEN}正在构建应用(VERBOSE=${VERBOSE:-false})...${NC}\n"

# 构建命令
if [ "${VERBOSE}" = "true" ]; then
    QUIET_FLAG=""
else
    QUIET_FLAG="-quiet"
fi

# 构建通用的 xcodebuild 参数
BASE_ARGS="-scheme \"${SCHEME}\" -configuration Release -derivedDataPath \"${BuildPath}\""
if [ "${PROJECT_TYPE}" = "workspace" ]; then
    BASE_ARGS="-workspace \"${PROJECT_FILE}\" ${BASE_ARGS}"
else
    BASE_ARGS="-project \"${PROJECT_FILE}\" ${BASE_ARGS}"
fi

# 添加目标和架构参数
BUILD_ARGS="${BASE_ARGS} -destination \"${DESTINATION}\""
if [ -n "$ARCHS" ]; then
    BUILD_ARGS="${BUILD_ARGS} ARCHS=\"${ARCHS}\" ONLY_ACTIVE_ARCH=NO"
fi

# 添加静默参数
if [ "${VERBOSE}" != "true" ]; then
    BUILD_ARGS="${BUILD_ARGS} -quiet"
fi

# 显示完整的执行命令（包含架构参数）
echo "xcodebuild ${BUILD_ARGS}"
echo

# 执行构建命令
printf "${YELLOW}正在清理之前的构建...${NC}\n"
eval "xcodebuild ${BUILD_ARGS} clean"

if [ "$ARCH" = "universal" ] && [ -n "$ARCHS" ]; then
    printf "${YELLOW}开始构建应用 (通用二进制: ${ARCHS})...${NC}\n"
elif [ "$ARCH" != "universal" ]; then
    printf "${YELLOW}开始构建应用 (架构: ${ARCH})...${NC}\n"
else
    printf "${YELLOW}开始构建应用 (架构: ${ARCH})...${NC}\n"
fi
eval "xcodebuild ${BUILD_ARGS} build"

# 检查构建结果
if [ $? -eq 0 ]; then
    echo
    printf "${GREEN}✅ 构建成功完成！${NC}\n"
    printf "${GREEN}📦 构建产物位置: ${BuildPath}/Build/Products/Release/${NC}\n"
else
    echo
    printf "${RED}❌ 构建失败！${NC}\n"
    exit 1
fi

# 显示开发路线图
show_development_roadmap "build"