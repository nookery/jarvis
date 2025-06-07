#!/bin/bash

# ====================================
# macOS 应用代码签名脚本 (重构版)
# ====================================
#
# 这个脚本用于对 macOS 应用程序进行代码签名，包括 Sparkle 框架的各个组件。
# 脚本采用结构化的四阶段执行流程，便于维护和调试。
#
# 执行流程：
# 1. 环境检查和信息输出
# 2. 检测可用资源（应用方案、代码签名、应用路径）并存储
# 3. 参数验证，不满足要求则输出建议
# 4. 满足要求则执行正常签名
#
# 使用方法：
# 1. 设置必要的环境变量：
#    export SCHEME="YourAppScheme"             # 应用方案名称
#    export SIGNING_IDENTITY="Developer ID"   # 代码签名身份
#    export BuildPath="/path/to/build"        # 构建输出路径（可选，默认为 ./temp）
#    export VERBOSE="true"                    # 可选：显示详细签名日志
#
# 2. 在项目根目录运行脚本：
#    ./scripts/codesign-app-refactored.sh
#
# 注意事项：
# - 需要有效的 Apple 开发者证书
# - 需要在 macOS 系统上运行
# - 应用程序必须已经构建完成
# ====================================

# 设置错误处理
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 全局变量存储检测结果
AVAILABLE_SCHEMES=""
AVAILABLE_IDENTITIES=""
AVAILABLE_APP_PATHS=()
PROJECT_FILE=""
PROJECT_TYPE=""

# ====================================
# 工具函数
# ====================================

# 打印成功信息
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# 打印错误信息
print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 打印警告信息
print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# 打印信息
print_info() {
    echo -e "${BLUE}ℹ️  $1: $2${NC}"
}

# 打印分隔线
print_separator() {
    echo -e "${BLUE}===========================================${NC}"
}

# 打印标题
print_title() {
    echo -e "${PURPLE}$1${NC}"
}

# 执行命令并显示结果
execute_command() {
    local cmd="$1"
    local description="$2"
    
    echo -e "${PURPLE}🔧 $description${NC}"
    
    if [ "$VERBOSE" = "true" ]; then
        echo -e "${BLUE}执行命令: $cmd${NC}"
    fi
    
    if eval "$cmd"; then
        print_success "$description 完成"
    else
        print_error "$description 失败"
        exit 1
    fi
    echo
}

# ====================================
# 阶段 1: 环境检查和信息输出
# ====================================

check_environment() {
    print_separator
    print_title "         🔐 代码签名脚本启动 (重构版)         "
    print_separator
    echo
    
    # 设置默认值
    BuildPath=${BuildPath:-"./temp"}
    VERBOSE=${VERBOSE:-"false"}
    
    # 显示基本环境信息
    printf "${GREEN}📋 环境信息:${NC}\n"
    echo "   操作系统: $(uname -s)"
    echo "   架构: $(uname -m)"
    echo "   主机名: $(hostname)"
    echo "   工作目录: $(pwd)"
    echo "   当前时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "   详细日志: ${VERBOSE}"
    echo
    
    # 显示用户提供的参数
    printf "${GREEN}📋 用户参数:${NC}\n"
    echo "   应用方案: ${SCHEME:-'未设置'}"
    echo "   构建路径: ${BuildPath}"
    echo "   签名身份: ${SIGNING_IDENTITY:-'未设置'}"
    echo
}

# ====================================
# 阶段 2: 检测可用资源并存储
# ====================================

detect_available_resources() {
    printf "${GREEN}🔍 检测可用资源...${NC}\n"
    echo
    
    # 检测项目文件
    detect_project_file
    
    # 检测可用的 Schemes
    detect_available_schemes
    
    # 检测可用的代码签名证书
    detect_available_identities
    
    # 检测可能的应用程序路径
    detect_available_app_paths
    
    # 显示检测结果摘要
    show_detection_summary
}

# 检测项目文件
detect_project_file() {
    if ls *.xcworkspace 1> /dev/null 2>&1; then
        PROJECT_FILE=$(ls -d *.xcworkspace | head -1)
        PROJECT_TYPE="-workspace"
        print_info "项目文件" "$PROJECT_FILE (workspace)"
    elif ls *.xcodeproj 1> /dev/null 2>&1; then
        PROJECT_FILE=$(ls -d *.xcodeproj | head -1)
        PROJECT_TYPE="-project"
        print_info "项目文件" "$PROJECT_FILE (project)"
    else
        print_error "未找到 Xcode 项目文件"
        exit 1
    fi
}

# 检测可用的 Schemes
detect_available_schemes() {
    printf "${CYAN}📋 检测可用 Scheme...${NC}\n"
    
    AVAILABLE_SCHEMES=$(xcodebuild $PROJECT_TYPE "$PROJECT_FILE" -list 2>/dev/null | sed -n '/Schemes:/,/^$/p' | grep -v 'Schemes:' | grep -v '^$' | sed 's/^[[:space:]]*//' | sort -u)
    
    if [ -n "$AVAILABLE_SCHEMES" ]; then
        echo "$AVAILABLE_SCHEMES" | while read -r scheme; do
            [ -n "$scheme" ] && echo "  - $scheme"
        done
    else
        print_warning "未检测到可用的 Scheme"
    fi
    echo
}

# 检测可用的代码签名证书
detect_available_identities() {
    printf "${CYAN}📋 检测可用代码签名证书...${NC}\n"
    
    AVAILABLE_IDENTITIES=$(security find-identity -v -p codesigning | grep -E "(Developer ID Application|Apple Development|iPhone Developer|Mac Developer)" | head -5)
    
    if [ -n "$AVAILABLE_IDENTITIES" ]; then
        echo "$AVAILABLE_IDENTITIES" | while IFS= read -r line; do
            # 提取证书名称
            CERT_NAME=$(echo "$line" | sed 's/.*"\(.*\)"/\1/')
            # 根据证书类型添加说明
            if [[ "$CERT_NAME" == *"Developer ID Application"* ]]; then
                echo "  - $CERT_NAME [分发证书 - 可公开分发]"
            elif [[ "$CERT_NAME" == *"Apple Development"* ]]; then
                echo "  - $CERT_NAME [开发证书 - 仅限开发测试]"
            elif [[ "$CERT_NAME" == *"Mac Developer"* ]]; then
                echo "  - $CERT_NAME [开发证书 - 仅限开发测试]"
            elif [[ "$CERT_NAME" == *"iPhone Developer"* ]]; then
                echo "  - $CERT_NAME [开发证书 - 仅限开发测试]"
            else
                echo "  - $CERT_NAME"
            fi
        done
    else
        print_warning "未检测到可用的代码签名证书"
    fi
    echo
}

# 检测可能的应用程序路径
detect_available_app_paths() {
    printf "${CYAN}📋 检测可能的应用程序路径...${NC}\n"
    
    # 如果用户提供了 SCHEME，则搜索该 SCHEME 的应用程序
    local search_scheme="${SCHEME:-*}"
    
    # 搜索可能的路径
    local possible_paths=(
        "./temp/Build/Products/Debug/${search_scheme}.app"
        "./temp/Build/Products/Release/${search_scheme}.app"
        "./temp/arm64/Build/Products/Release/${search_scheme}.app"
        "./temp/arm64/Build/Products/Debug/${search_scheme}.app"
        "./temp/x86_64/Build/Products/Release/${search_scheme}.app"
        "./temp/x86_64/Build/Products/Debug/${search_scheme}.app"
        "./temp/universal/Build/Products/Release/${search_scheme}.app"
        "./temp/universal/Build/Products/Debug/${search_scheme}.app"
        "./Build/Products/Release/${search_scheme}.app"
        "./Build/Products/Debug/${search_scheme}.app"
        "./build/Release/${search_scheme}.app"
        "./build/Debug/${search_scheme}.app"
        "./DerivedData/Build/Products/Release/${search_scheme}.app"
        "./DerivedData/Build/Products/Debug/${search_scheme}.app"
    )
    
    AVAILABLE_APP_PATHS=()
    
    # 检查预定义路径
    for path in "${possible_paths[@]}"; do
        # 使用通配符展开路径
        for expanded_path in $path; do
            if [ -d "$expanded_path" ]; then
                # 检查应用程序大小，过滤掉大小为0的应用程序
                local app_size_bytes=$(du -s "$expanded_path" 2>/dev/null | cut -f1 || echo "0")
                if [ "$app_size_bytes" -gt 0 ]; then
                    AVAILABLE_APP_PATHS+=("$expanded_path")
                fi
            fi
        done
    done
    
    # 使用 find 命令搜索更多可能的位置
    while IFS= read -r -d '' app_path; do
        local app_size_bytes=$(du -s "$app_path" 2>/dev/null | cut -f1 || echo "0")
        if [ "$app_size_bytes" -gt 0 ]; then
            # 检查是否已经在列表中
            local already_exists=false
            for existing_path in "${AVAILABLE_APP_PATHS[@]}"; do
                if [ "$existing_path" = "$app_path" ]; then
                    already_exists=true
                    break
                fi
            done
            if [ "$already_exists" = false ]; then
                AVAILABLE_APP_PATHS+=("$app_path")
            fi
        fi
    done < <(find . -name "${search_scheme}.app" -type d -print0 2>/dev/null | head -20)
    
    # 显示找到的应用程序
    if [ ${#AVAILABLE_APP_PATHS[@]} -gt 0 ]; then
        for i in "${!AVAILABLE_APP_PATHS[@]}"; do
            local app_path="${AVAILABLE_APP_PATHS[$i]}"
            local app_size=$(du -sh "$app_path" 2>/dev/null | cut -f1 || echo "未知")
            echo "  $((i+1)). $app_path ($app_size)"
        done
    else
        print_warning "未找到任何应用程序"
    fi
    echo
}

# 显示检测结果摘要
show_detection_summary() {
    printf "${GREEN}📊 资源检测摘要:${NC}\n"
    
    local scheme_count=$(echo "$AVAILABLE_SCHEMES" | grep -c . || echo "0")
    local identity_count=$(echo "$AVAILABLE_IDENTITIES" | grep -c . || echo "0")
    local app_count=${#AVAILABLE_APP_PATHS[@]}
    
    echo "   可用 Scheme: $scheme_count 个"
    echo "   可用签名证书: $identity_count 个"
    echo "   可用应用程序: $app_count 个"
    echo
}

# ====================================
# 阶段 3: 参数验证和建议输出
# ====================================

validate_parameters_and_suggest() {
    printf "${GREEN}🔍 验证用户参数...${NC}\n"
    echo
    
    local missing_vars=""
    local invalid_vars=""
    
    # 检查 SCHEME
    if [ -z "$SCHEME" ]; then
        missing_vars="${missing_vars}SCHEME "
    elif [ -n "$AVAILABLE_SCHEMES" ] && ! echo "$AVAILABLE_SCHEMES" | grep -q "^$SCHEME$"; then
        invalid_vars="${invalid_vars}SCHEME(不在可用列表中) "
    fi
    
    # 检查 SIGNING_IDENTITY
    if [ -z "$SIGNING_IDENTITY" ]; then
        missing_vars="${missing_vars}SIGNING_IDENTITY "
    elif [ -n "$AVAILABLE_IDENTITIES" ] && ! echo "$AVAILABLE_IDENTITIES" | grep -q "$SIGNING_IDENTITY"; then
        invalid_vars="${invalid_vars}SIGNING_IDENTITY(不在可用列表中) "
    fi
    
    # 如果有缺失或无效的参数，提供建议
    if [ -n "$missing_vars" ] || [ -n "$invalid_vars" ]; then
        if [ -n "$missing_vars" ]; then
            print_error "以下环境变量未设置: $missing_vars"
        fi
        if [ -n "$invalid_vars" ]; then
            print_error "以下环境变量值无效: $invalid_vars"
        fi
        echo
        
        generate_suggestions
        exit 1
    else
        print_success "参数验证通过"
        echo
    fi
}

# 生成建议命令
generate_suggestions() {
    printf "${GREEN}💡 建议使用以下命令进行代码签名:${NC}\n"
    echo
    
    # 生成所有可能的组合建议
    if [ -n "$AVAILABLE_SCHEMES" ] && [ -n "$AVAILABLE_IDENTITIES" ] && [ ${#AVAILABLE_APP_PATHS[@]} -gt 0 ]; then
        # 将schemes转换为数组
        local schemes_array=()
        while IFS= read -r scheme; do
            [ -n "$scheme" ] && schemes_array+=("$scheme")
        done <<< "$AVAILABLE_SCHEMES"
        
        # 将identities转换为数组
        local identities_array=()
        while IFS= read -r line; do
            local cert_name=$(echo "$line" | sed 's/.*"\(.*\)"/\1/')
            [ -n "$cert_name" ] && identities_array+=("$cert_name")
        done <<< "$AVAILABLE_IDENTITIES"
        
        local command_count=0
        
        # 生成所有组合
        for scheme in "${schemes_array[@]}"; do
            for identity in "${identities_array[@]}"; do
                for app_path in "${AVAILABLE_APP_PATHS[@]}"; do
                    local build_path=$(dirname "$app_path")
                    # 转换为绝对路径
                    local abs_build_path=$(cd "$build_path" 2>/dev/null && pwd || echo "$build_path")
                    local abs_script_path=$(cd "$(dirname "$0")" && pwd)/$(basename "$0")
                    echo " SCHEME='$scheme' SIGNING_IDENTITY='$identity' BuildPath='$abs_build_path' '$abs_script_path'"
                    echo
                    command_count=$((command_count + 1))
                done
            done
        done
        
        printf "${GREEN}📊 总共生成了 ${command_count} 个命令建议 (${#schemes_array[@]} 个 Scheme × ${#identities_array[@]} 个签名证书 × ${#AVAILABLE_APP_PATHS[@]} 个应用程序位置)${NC}\n"
    else
        # 简化建议
        if [ -n "$AVAILABLE_SCHEMES" ]; then
            echo "$AVAILABLE_SCHEMES" | while read -r scheme; do
                if [ -n "$scheme" ]; then
                    echo " SCHEME='$scheme' SIGNING_IDENTITY='YOUR_SIGNING_IDENTITY' ./scripts/codesign-app-refactored.sh"
                    echo
                fi
            done
            echo "注意: 请将 YOUR_SIGNING_IDENTITY 替换为您的实际代码签名身份"
        fi
    fi
    
    echo
    printf "${GREEN}📋 证书类型说明:${NC}\n"
    echo "   🟢 Developer ID Application: 用于 Mac App Store 外分发，可被所有用户安装"
    echo "   🟡 Apple Development: 用于开发测试，仅限开发团队内部使用"
    echo "   🔴 Mac App Store: 用于 App Store 上架（需单独申请）"
    echo
}

# ====================================
# 阶段 4: 执行代码签名
# ====================================

perform_code_signing() {
    printf "${GREEN}🔐 开始执行代码签名...${NC}\n"
    echo
    
    # 确定应用程序路径
    local app_path=""
    
    # 检查 BuildPath 是否已经包含 Build/Products 路径
    if [[ "$BuildPath" == *"/Build/Products/"* ]]; then
        # 如果已经包含，直接使用
        app_path="$BuildPath/$SCHEME.app"
    else
        # 如果不包含，添加标准路径
        app_path="$BuildPath/Build/Products/Release/$SCHEME.app"
    fi
    
    # 验证应用程序是否存在
    if [ ! -d "$app_path" ]; then
        print_error "应用程序不存在: $app_path"
        
        # 尝试从检测到的路径中找到匹配的应用程序
        local found_match=false
        for available_path in "${AVAILABLE_APP_PATHS[@]}"; do
            if [[ "$available_path" == *"/$SCHEME.app" ]]; then
                app_path="$available_path"
                found_match=true
                print_info "使用检测到的应用程序" "$app_path"
                break
            fi
        done
        
        if [ "$found_match" = false ]; then
            print_error "无法找到匹配的应用程序"
            exit 1
        fi
    fi
    
    # 显示签名信息
    printf "${GREEN}📋 签名信息:${NC}\n"
    echo "   应用程序: $app_path"
    echo "   签名身份: $SIGNING_IDENTITY"
    echo "   详细日志: $VERBOSE"
    echo
    
    # 执行 Sparkle 框架签名
    sign_sparkle_framework "$app_path"
    
    # 执行主应用程序签名
    sign_main_application "$app_path"
    
    # 验证签名
    verify_code_signature "$app_path"
    
    # 显示完成信息
    show_completion_info
}

# 签名 Sparkle 框架
sign_sparkle_framework() {
    local app_path="$1"
    local sparkle_framework="$app_path/Contents/Frameworks/Sparkle.framework"
    
    if [ -d "$sparkle_framework" ]; then
        printf "${PURPLE}🔧 签名 Sparkle 框架组件...${NC}\n"
        
        # Sparkle 框架内的组件
        local sparkle_components=(
            "$sparkle_framework/Versions/B/Resources/Autoupdate.app/Contents/MacOS/Autoupdate"
            "$sparkle_framework/Versions/B/Resources/Autoupdate.app"
            "$sparkle_framework/Versions/B/Resources/Updater.app/Contents/MacOS/Updater"
            "$sparkle_framework/Versions/B/Resources/Updater.app"
            "$sparkle_framework/Versions/B/Sparkle"
            "$sparkle_framework"
        )
        
        for component in "${sparkle_components[@]}"; do
            if [ -e "$component" ]; then
                local component_name=$(basename "$component")
                execute_command "codesign --force --options runtime --sign '$SIGNING_IDENTITY' '$component'" "签名 $component_name"
            fi
        done
    else
        print_info "Sparkle 框架" "未找到，跳过"
    fi
}

# 签名主应用程序
sign_main_application() {
    local app_path="$1"
    
    printf "${PURPLE}🔧 签名主应用程序...${NC}\n"
    execute_command "codesign --force --options runtime --sign '$SIGNING_IDENTITY' '$app_path'" "签名主应用程序"
}

# 验证代码签名
verify_code_signature() {
    local app_path="$1"
    
    printf "${PURPLE}🔍 验证代码签名...${NC}\n"
    
    if codesign --verify --deep --strict "$app_path" 2>/dev/null; then
        print_success "代码签名验证通过"
    else
        print_error "代码签名验证失败"
        exit 1
    fi
    
    # 显示签名信息
    if [ "$VERBOSE" = "true" ]; then
        echo -e "${BLUE}签名详细信息:${NC}"
        codesign -dv "$app_path" 2>&1 | sed 's/^/   /'
    fi
    echo
}

# 显示完成信息
show_completion_info() {
    print_separator
    print_title "         🗺️  开发分发路线图                "
    print_separator
    print_separator
    echo
    
    printf "${GREEN}📍 当前位置: 代码签名${NC}\n"
    echo
    echo "  🔨 构建应用 编译源代码，生成可执行文件"
    echo "▶ 🔐 代码签名 为应用添加数字签名，确保安全性"
    echo "  📦 打包分发 创建 DMG 安装包"
    echo "  ✅ 公证验证 Apple 官方验证（可选）"
    echo "  🚀 发布分发 上传到分发平台或直接分发"
    echo
    echo
    printf "${GREEN}💡 下一步建议:${NC}\n"
    echo "   创建安装包: ./scripts/create-dmg.sh"
    echo
    echo
    print_separator
}

# ====================================
# 主执行流程
# ====================================

main() {
    # 阶段 1: 环境检查和信息输出
    check_environment
    
    # 阶段 2: 检测可用资源并存储
    detect_available_resources
    
    # 阶段 3: 参数验证和建议输出
    validate_parameters_and_suggest
    
    # 阶段 4: 执行代码签名
    perform_code_signing
}

# 执行主函数
main "$@"