# RepTracker Standalone Build Guide

This guide explains how to build and deploy the RepTracker mobile app as a standalone application that doesn't require a local development server.

## Overview

By default, React Native/Expo apps run in **development mode**, which requires a local Metro bundler server running on your machine. This guide shows you how to create **production builds** that are completely self-contained and can run independently.

## Prerequisites

- macOS with Xcode installed
- Node.js and npm
- iOS device connected to your Mac
- Apple Developer account (free tier is sufficient for development builds)

## Quick Start

### 1. Build Standalone iOS App

```bash
# Navigate to the frontend directory
cd frontend-native

# Create production build and export as IPA
npm run build:standalone:ios
```

### 2. Install on Your iPhone

**Method A: Using Xcode (Recommended)**
1. Open Xcode
2. Go to **Window > Devices and Simulators**
3. Select your connected iPhone
4. Drag `ios/build/export/RepTracker.ipa` onto your device

**Method B: Using Finder**
1. Connect iPhone to Mac
2. Open Finder and select your iPhone from sidebar
3. Drag the IPA file to your device

## Detailed Build Process

### Available Build Scripts

The following scripts are available in `package.json`:

```json
{
  "scripts": {
    // Development builds (require local server)
    "start": "expo start",
    "ios": "expo run:ios",
    "android": "expo run:android",
    
    // Standalone builds (no server required)
    "ios:release": "expo run:ios --configuration Release --no-bundler",
    "android:release": "expo run:android --variant release --no-bundler",
    
    // Cloud builds (requires paid Apple Developer account)
    "build:ios": "eas build --platform ios",
    "build:android": "eas build --platform android",
    "build:production": "eas build --profile production",
    
    // Bundle generation
    "export": "expo export",
    "bundle:ios": "npx react-native bundle --platform ios --dev false --entry-file index.js --bundle-output ./ios/main.jsbundle"
  }
}
```

### Manual Build Process

If you prefer to build manually or understand the process:

#### Step 1: Create iOS Archive

```bash
# Create the archive (.xcarchive)
xcodebuild -workspace ios/RepTracker.xcworkspace \
           -scheme RepTracker \
           -configuration Release \
           -destination generic/platform=iOS \
           -archivePath ios/build/RepTracker.xcarchive \
           archive
```

#### Step 2: Export as IPA

```bash
# Export the archive as an installable IPA file
xcodebuild -exportArchive \
           -archivePath ios/build/RepTracker.xcarchive \
           -exportPath ios/build/export \
           -exportOptionsPlist ios/RepTracker/ExportOptions.plist
```

The `ExportOptions.plist` file is already configured for development distribution and is stored in the `ios/RepTracker/` directory to prevent it from being cleared by clean commands.

## Build Configurations

### Development Build
- **Purpose**: Testing on physical devices
- **Distribution**: Development team only
- **Signing**: Development certificate
- **File**: `RepTracker.ipa` in `ios/build/export/`

### Production Build (Future)
For App Store distribution, you'll need:
- Paid Apple Developer account ($99/year)
- App Store Connect setup
- Production certificates

## Troubleshooting

### Common Issues

**1. "Waiting on http://localhost:8081"**
- **Problem**: App is still trying to connect to development server
- **Solution**: Use `--no-bundler` flag or manual build process above

**2. Code Signing Errors**
- **Problem**: Invalid certificates or provisioning profiles
- **Solution**: Check Xcode > Preferences > Accounts, ensure your Apple ID is signed in

**3. "No devices found"**
- **Problem**: iPhone not connected or not trusted
- **Solution**: Connect iPhone, trust computer, check `xcrun devicectl list devices`

**4. Build Fails with Pod Errors**
- **Solution**: 
  ```bash
  cd ios
  pod deintegrate
  pod install
  cd ..
  ```

### Verification

To verify your build is truly standalone:

1. **Disconnect from WiFi** on your iPhone
2. **Close any local development servers** on your Mac
3. **Launch the app** - it should work without any network dependency for the core functionality

## File Structure

After building, you'll have:

```
frontend-native/
├── ios/
│   ├── build/
│   │   ├── RepTracker.xcarchive     # Archive file
│   │   └── export/
│   │       ├── RepTracker.ipa       # ← Install this on your device
│   │       ├── DistributionSummary.plist
│   │       └── Packaging.log
│   └── RepTracker/
│       └── ExportOptions.plist          # Export configuration (safe location)
```

## Backend Configuration

Your app is configured to connect to your deployed backend on Google Cloud Services. The backend URL is set via environment variables:

- `EXPO_PUBLIC_API_BASE_URL` - Your GCS backend URL

Make sure your backend is running and accessible before testing the standalone app.

## Android Builds (Future)

For Android standalone builds:

```bash
# Android release build
npm run android:release

# Or manually
cd android
./gradlew assembleRelease
```

The APK will be generated in `android/app/build/outputs/apk/release/`.

## Automation Script

For convenience, you can create a build script:

```bash
#!/bin/bash
# build-standalone.sh

echo "Building RepTracker standalone app..."

# Clean previous builds
rm -rf ios/build/

# Create archive
xcodebuild -workspace ios/RepTracker.xcworkspace \
           -scheme RepTracker \
           -configuration Release \
           -destination generic/platform=iOS \
           -archivePath ios/build/RepTracker.xcarchive \
           archive

# Export IPA
xcodebuild -exportArchive \
           -archivePath ios/build/RepTracker.xcarchive \
           -exportPath ios/build/export \
           -exportOptionsPlist ios/RepTracker/ExportOptions.plist

echo "✅ Build complete! Install ios/build/export/RepTracker.ipa on your device"
```

Make it executable:
```bash
chmod +x build-standalone.sh
./build-standalone.sh
```

## Key Benefits

✅ **No local server dependency** - App runs completely standalone  
✅ **Production performance** - Optimized and minified code  
✅ **Direct backend connection** - Connects to your GCS backend  
✅ **Native performance** - Full React Native performance benefits  
✅ **Offline capable** - Core functionality works without internet  

## Next Steps

1. **Test thoroughly** on your physical device
2. **Set up CI/CD** for automated builds
3. **Prepare for App Store** distribution when ready
4. **Create Android builds** following similar process

---

**Need help?** Check the troubleshooting section or refer to the main project documentation. 