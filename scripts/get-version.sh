#!/bin/bash

# ====================================
# macOS 应用版本号获取脚本
# ====================================
#
# 这个脚本用于从 Xcode 项目配置文件中获取应用程序的营销版本号（MARKETING_VERSION）。
# 它会自动查找项目中的 .pbxproj 文件，并从中提取版本号信息。
#
# 功能：
# 1. 自动查找项目中的 .pbxproj 文件
# 2. 从配置文件中提取 MARKETING_VERSION
# 3. 输出格式化的版本号（x.y.z）
#
# 使用方法：
# 1. 直接运行（自动查找 .pbxproj）：
#    ./scripts/get-version.sh
#
# 2. 指定 .pbxproj 文件路径：
#    ./scripts/get-version.sh path/to/project.pbxproj
#
# 返回值：
# - 成功：输出版本号（例如：1.0.0）并返回 0
# - 失败：输出错误信息到标准错误并返回非零值
#   * 1: 未找到 .pbxproj 文件
#   * 2: 未找到版本号
#
# 注意事项：
# - 需要在项目根目录或其父目录下运行
# - 会自动过滤掉 Resources 和 temp 目录
# - 如果找到多个 .pbxproj 文件，使用第一个匹配的文件
# ====================================

# 用法: bash scripts/get-version.sh [pbxproj路径]
projectFile=${1:-$(find $(pwd) -maxdepth 2 ! -path "*Resources*" ! -path "*temp*" -type f -name "*.pbxproj" | head -n 1)}
if [ -z "$projectFile" ]; then
  echo "❌ 未找到 .pbxproj 配置文件！" >&2
  exit 1
fi
version=$(grep "MARKETING_VERSION" "$projectFile" | head -n 1 | grep -o '[0-9]\+\.[0-9]\+\.[0-9]\+')
if [ -z "$version" ]; then
  echo "❌ 未找到 MARKETING_VERSION！" >&2
  exit 2
fi
echo "$version"