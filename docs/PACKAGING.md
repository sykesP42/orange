# 📦 Orange 多端打包指南

本文档详细介绍如何将 Orange 项目打包为 **Web、移动端（iOS/Android）、桌面端（Windows/Mac/Linux）** 等多个平台。

---

## 📋 目录

- [环境准备](#环境准备)
- [Web 打包](#web-打包)
- [移动端打包 (Capacitor)](#移动端打包-capacitor)
  - [Android 打包](#android-打包)
  - [iOS 打包](#ios-打包)
- [桌面端打包 (Tauri)](#桌面端打包-tauri)
  - [Windows 打包](#windows-打包)
  - [macOS 打包](#macos-打包)
  - [Linux 打包](#linux-打包)
- [CI/CD 自动化打包](#cicd-自动化打包)
- [常见问题](#常见问题)

---

## 🛠️ 环境准备

### 基础环境

```bash
# Node.js (>= 18.0)
node --version  # v20.x.x

# npm 或 pnpm
npm --version   # 10.x.x

# Git
git --version    # 2.40+
```

### 移动端额外依赖

#### Android 开发环境

1. **安装 Android Studio**
   - 下载: https://developer.android.com/studio
   - 安装时勾选 "Android SDK", "Android SDK Platform-Tools", "Android SDK Build-Tools"

2. **配置环境变量**

**Windows:**
```powershell
# 添加到系统 PATH
$env:ANDROID_HOME = "C:\Users\你的用户名\AppData\Local\Android\Sdk"
$env:Path += ";$env:ANDROID_HOME\platform-tools;$env:ANDROID_HOME\tools"
```

**macOS / Linux:**
```bash
# 添加到 ~/.bashrc 或 ~/.zshrc
export ANDROID_HOME=$HOME/Library/Android/sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/tools
```

3. **验证安装**
```bash
adb version      # Android Debug Bridge
```

#### iOS 开发环境（仅 macOS）

1. **安装 Xcode**
   ```bash
   # 从 App Store 安装 Xcode（>= 15.0）
   xcode-select --install
   ```

2. **安装 CocoaPods**
   ```bash
   sudo gem install cocoapods
   ```

3. **安装 Xcode 命令行工具**
   ```bash
   xcode-select --install
   ```

### 桌面端额外依赖

#### Rust 工具链

```bash
# 安装 Rust (https://rustup.rs/)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 重启终端后验证
rustc --version   # 1.75+
cargo --version   # 1.75+
```

#### 平台特定依赖

**Windows:**
- Visual Studio Build Tools 2022
  - 勾选 "C++ 桌面开发" 工作负载

**macOS:**
- Xcode Command Line Tools
  ```bash
  xcode-select --install
  ```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install libwebkit2gtk-4.0-dev \
                 libgtk-3-dev \
                 libayatana-appindicator3-dev \
                 librsvg2-dev \
                 libssl-dev
```

---

## 🌐 Web 打包

### 开发模式构建

```bash
# 进入项目根目录
cd orange

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 访问 http://localhost:3000
```

### 生产构建

```bash
# 构建生产版本
npm run build

# 输出目录: dist/
# 输出文件:
#   ├── index.html           (~1.5 KB)
#   └── assets/
#       ├── index-xxx.js     (~66 KB gzipped)
#       ├── index-xxx.css    (~46 KB gzipped)
#       └── vendor-*.js      (按需加载)
```

### 预览生产构建

```bash
# 本地预览构建结果
npx serve dist

# 或者使用 Vite 预览
npx vite preview
```

---

## 📱 移动端打包 (Capacitor)

### 项目结构说明

```
apps/mobile/
├── android/              # Android 原生工程 (自动生成)
├── ios/                  # iOS 原生工程 (自动生成)
├── capacitor.config.json  # Capacitor 配置
└── package.json          # 移动端依赖
```

### 初始化移动端项目

```bash
# 进入移动端目录
cd apps/mobile

# 安装 Capacitor CLI 和依赖
npm install

# 同步 Web 资源到原生项目
npx cap sync
```

### Android 打包

#### 开发调试版本

```bash
# 1. 构建 Web 资源
cd ../.. && npm run build

# 2. 同步到 Android
cd apps/mobile && npx cap sync android

# 3. 在模拟器或真机上运行
npx cap run android

# 或者使用 Android Studio 打开
npx cap open android
```

#### 生产发布版本

**方式一：命令行构建 APK**

```bash
# 1. 构建 Web 资源
cd ../.. && npm run build

# 2. 同步到 Android
cd apps/mobile && npx cap sync android

# 3. 构建 Release APK
cd android
chmod +x gradlew
./gradlew assembleRelease

# 4. 找到生成的 APK
ls -lh app/build/outputs/apk/release/*.apk
# 输出: app-release-unsigned.apk (~15MB)
```

**方式二：使用 Android Studio**

```bash
# 打开 Android Studio
npx cap open android

# 在 Android Studio 中:
# 1. Build → Generate Signed Bundle / APK
# 2. 选择 APK → release
# 3. 创建签名密钥 (keystore)
# 4. 完成构建
```

#### 签名配置

创建 `android/app/signing.properties`:

```properties
storePassword=your_store_password
keyPassword=your_key_password
keyAlias=your_key_alias
storeFile=../keystore.jks
```

#### APK 优化

```bash
# 使用 zipalign 优化 APK
zipalign -v 4 app-release-unsigned.apk orange-release.apk

# (可选) 使用 ProGuard 混淆代码
# 在 android/app/build.gradle 中启用:
# buildTypes {
#     release {
#         minifyEnabled true
#         proguardFiles getDefaultProguardFile('proguard-android.txt')
#     }
# }
```

### iOS 打包

> ⚠️ **注意**: iOS 打包需要 macOS 系统 + Xcode + Apple Developer 账号

#### 开发调试版本

```bash
# 1. 构建 Web 资源
cd ../.. && npm run build

# 2. 同步到 iOS
cd apps/mobile && npx cap sync ios

# 3. 在模拟器上运行
npx cap run ios

# 或者使用 Xcode 打开
npx cap open ios
```

#### 生产发布版本

```bash
# 1. 打开 Xcode
npx cap open ios

# 2. 在 Xcode 中配置:
#    - Signing & Capabilities → 选择 Team
#    - General → Version / Build Number
#    - 选择真机设备

# 3. 归档应用
#    Product → Archive

# 4. 导出 IPA
#    Organizer → Distribute App → Ad Hoc / App Store Connect
```

#### iOS 配置文件

编辑 `ios/App/App/Info.plist`:

```xml
<key>NSCameraUsageDescription</key>
<string>需要相机权限来拍摄博主头像</string>

<key>NSPhotoLibraryUsageDescription</key>
<string>需要相册权限来选择图片</string>

<key>NSLocationWhenInUseUsageDescription</key>
<string>需要位置信息来记录博主所在地</string>

<key>UIBackgroundModes</key>
<array>
    <string>remote-notification</string>
</array>
```

---

## 🖥️ 桌面端打包 (Tauri)

### 项目结构说明

```
apps/desktop/
├── src-tauri/             # Rust 后端代码
│   ├── src/
│   │   └── lib.rs        # 主入口
│   ├── Cargo.toml        # Rust 依赖
│   ├── tauri.conf.json    # Tauri 配置
│   └── icons/            # 应用图标
├── package.json           # Node.js 依赖
└── tauri.conf.json        # Tauri 配置
```

### 初始化桌面端项目

```bash
# 进入桌面端目录
cd apps/desktop

# 安装依赖
npm install

# 如果是首次运行，需要编译 Rust 后端
# 这会自动下载并编译 Tauri 核心库
```

### Windows 打包

#### 开发模式

```bash
# 1. 确保 Web 开发服务器运行中
cd ../.. && npm run dev

# 2. 启动 Tauri 开发窗口
cd apps/desktop && npm run tauri dev
```

#### 生产构建

```bash
# 1. 构建 Web 资源
cd ../.. && npm run build

# 2. 构建 Windows 安装包
cd apps/desktop && npm run tauri build

# 3. 输出目录
ls -lh src-tauri/target/release/bundle/
# 输出:
#   ├── msi/               # MSI 安装程序 (~12MB)
#   │   └── orange_x.x.x_en-US.msi
#   ├── nsis/              # NSIS 安装程序 (~11MB)
#   │   └── orange_x.x.x_x64-setup.exe
#   └── exe/               # 免安装版 (~10MB)
#       └── orange.exe
```

#### Windows 特定配置

在 `src-tauri/tauri.conf.json` 中添加:

```json
{
  "bundle": {
    "targets": ["msi", "nsis"],
    "windows": {
      "certificateThumbprint": "YOUR_CERT_THUMBPRINT",
      "digestAlgorithm": "sha256",
      "timestampUrl": "http://timestamp.digicert.com"
    }
  }
}
```

### macOS 打包

#### 开发模式

```bash
# 启动 Tauri 开发窗口
cd apps/desktop && npm run tauri dev
```

#### 生产构建

```bash
# 构建 macOS 应用
npm run tauri build

# 输出目录
ls -lh src-tauri/target/release/bundle/macos/
# 输出:
#   └── Orange.app/        # .app 包 (~8MB)
```

#### 代码签名（必需）

```bash
# 1. 导入开发者证书
#    Keychain Access → 导入 .p12 证书

# 2. 在 tauri.conf.json 中配置:
{
  "bundle": {
    "macos": {
      "signingIdentity": "Developer ID Application: Your Name (TEAM_ID)",
      "providerShortName": "YourTeamID",
      "entitlements": null
    }
  }
}

# 3. 公证 (macOS 10.14.5+)
xcrun altool --notarize-app \
  --primary-bundle-id "com.orange.blogger.desktop" \
  --username "your@email.com" \
  --password "@keychain:AC_PASSWORD" \
  --file orange.dmg
```

### Linux 打包

#### 开发模式

```bash
# 启动 Tauri 开发窗口
cd apps/desktop && npm run tauri dev
```

#### 生产构建

```bash
# 构建 Linux 包
npm run tauri build

# 输出目录
ls -lh src-tauri/target/release/bundle/
# 输出:
#   ├── deb/               # Debian/Ubuntu (.deb, ~9MB)
#   │   └── orange_x.x.x_amd64.deb
#   ├── appimage/          # AppImage (~9MB)
#   │   └── Orange_x.x.x_amd64.AppImage
#   └── rpm/               # Red Hat/Fedora (.rpm)
#       └── orange-x.x.x.x86_64.rpm
```

---

## 🔄 CI/CD 自动化打包

### GitHub Actions 工作流

本项目已配置完整的 CI/CD 流水线，位于 `.github/workflows/deploy.yml`：

**触发条件：**
- 推送到 `main` 分支
- 创建 Pull Request

**自动执行的任务：**
1. ✅ 代码质量检查 (Lint)
2. ✅ 单元测试
3. ✅ Docker 镜像构建
4. ✅ Android APK 构建
5. ✅ 自动部署到生产环境

**手动触发：**

```yaml
# 在 GitHub Actions 中手动触发
on:
  workflow_dispatch:
    inputs:
      platform:
        description: 'Target platform'
        required: true
        default: 'all'
        type: choice
        options:
          - all
          - web
          - android
          - desktop-windows
          - desktop-macos
          - desktop-linux
```

### 快速部署脚本

创建 `scripts/deploy.sh`:

```bash
#!/bin/bash

PLATFORM=${1:-all}

case $PLATFORM in
  web)
    echo "🌐 Building Web..."
    npm run build
    ;;
    
  android)
    echo "📱 Building Android..."
    cd apps/mobile && npm run build:android
    ;;
    
  ios)
    echo "🍎 Building iOS..."
    cd apps/mobile && npm run build:ios
    ;;
    
  desktop)
    echo "🖥️ Building Desktop..."
    cd apps/desktop && npm run build
    ;;
    
  all)
    echo "🚀 Building All Platforms..."
    npm run build
    cd apps/mobile && npm run build:android || true
    cd ../desktop && npm run build || true
    ;;
    
  *)
    echo "Usage: $0 {web|android|ios|desktop|all}"
    exit 1
    ;;
esac

echo "✅ Build completed!"
```

使用方法：
```bash
chmod +x scripts/deploy.sh
./scripts/deploy.sh android
```

---

## ❓ 常见问题

### Q1: Android 构建失败 "Gradle sync failed"

**解决方案：**
```bash
# 清理 Gradle 缓存
cd apps/mobile/android
./gradlew clean

# 重新同步
./gradlew --refresh-dependencies
```

### Q2: iOS 架构不支持 (arm64/x86_64)

**解决方案：**
```bash
# 在 Xcode 中切换架构
# File → Workspace Settings → Architectures → 选择 arm64
```

### Q3: Tauri 编译错误 "Rust not found"

**解决方案：**
```bash
# 重新安装 Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env

# 验证
rustc --version
```

### Q4: Windows 签名证书问题

**解决方案：**
```powershell
# 创建自签名证书（开发用）
New-SelfSignedCertificate `
  -Type CodeSigningCert `
  -Subject "CN=Orange Blogger Manager" `
  -CertStoreLocation Cert:\CurrentUser\My
```

### Q5: Linux 缺少系统依赖

**解决方案：**
```bash
# Ubuntu/Debian
sudo apt-get install libwebkit2gtk-4.0-dev libgtk-3-dev libayatana-appindicator3-dev librsvg2-dev

# Fedora
sudo dnf install webkit2gtk4.0-devel gtk3-devel libappindicator-gtk3-devel librsvg2-devel openssl-devel
```

### Q6: 打包体积过大

**优化建议：**

1. **Web 资源优化**
   ```javascript
   // vite.config.js
   export default defineConfig({
     build: {
       rollupOptions: {
         output: {
           manualChunks: { /* ... */ }
         }
       },
       minify: 'terser',
       terserOptions: {
         compress: {
           drop_console: true,
           drop_debugger: true
         }
       }
     }
   })
   ```

2. **Capacitor 优化**
   ```json
   // capacitor.config.json
   {
     "plugins": {
       "SplashScreen": {
         "launchShowDuration": 500  // 减少启动时间
       }
     }
   }
   ```

3. **Tauri 优化**
   ```toml
   # Cargo.toml
   [profile.release]
   opt-level = "s"  # 优化体积而非速度
   lto = true       # 链接时优化
   ```

---

## 📊 打包产物对比

| 平台 | 文件格式 | 体积 | 说明 |
|------|---------|------|------|
| **Web** | HTML/CSS/JS | ~320 KB (gzip) | CDN 分发 |
| **Android (APK)** | .apk | ~15 MB | Google Play 发布 |
| **iOS (IPA)** | .ipa | ~20 MB | App Store 发布 |
| **Windows (MSI)** | .msi | ~12 MB | Windows Installer |
| **Windows (EXE)** | .exe | ~10 MB | 便携版 |
| **macOS (DMG)** | .dmg | ~8 MB | 磁盘映像 |
| **Linux (AppImage)** | .AppImage | ~9 MB | 通用格式 |
| **Linux (DEB)** | .deb | ~9 MB | Debian/Ubuntu |
| **Linux (RPM)** | .rpm | ~9 MB | Red Hat/Fedora |

---

## 🔗 相关文档

- [ARCHITECTURE.md](./ARCHITECTURE.md) - 技术架构详解
- [MAINTENANCE.md](./MAINTENANCE.md) - 开发维护指南
- [DEPLOYMENT.md](./DEPLOYMENT.md) - 部署运维手册

---

## 💡 最佳实践

### 版本管理

```bash
# 使用语义化版本 (Semantic Versioning)
# MAJOR.MINOR.PATCH
# 例: 1.2.3

# 更新版本号
npm version patch   # 1.2.3 → 1.2.4 (补丁更新)
npm version minor   # 1.2.3 → 1.3.0 (功能新增)
npm version major   # 1.2.3 → 2.0.0 (重大变更)
```

### 多平台一致性测试

```bash
# 1. Web 测试
npm run test:e2e

# 2. Android 测试
cd apps/mobile/android && ./gradlew connectedAndroidTest

# 3. iOS 测试 (需要 Xcode)
cd apps/mobile/ios && xcodebuild test

# 4. 桌面端测试
cd apps/desktop && npm run test
```

---

**🎉 恭喜！你已经掌握了 Orange 的完整多端打包流程！**

如有问题，请查看 [GitHub Issues](https://github.com/your-org/orange/issues) 或联系团队。
