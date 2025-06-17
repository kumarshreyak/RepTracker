# React Native Google Sign-in Setup

## Prerequisites

1. **Google Cloud Console Setup**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select an existing one
   - Enable the Google+ API or Google Sign-In API

2. **Create OAuth 2.0 Client IDs**
   
   You'll need to create 3 different client IDs:

   ### Web Client ID
   - Application type: Web application
   - Name: "Web client (auto created by Google Service)"
   - Copy the Client ID (ends with `.googleusercontent.com`)

   ### iOS Client ID  
   - Application type: iOS
   - Name: "iOS client"
   - Bundle ID: Your iOS app bundle ID
   - Copy the Client ID (ends with `.googleusercontent.com`)

   ### Android Client ID
   - Application type: Android
   - Name: "Android client"  
   - Package name: Your Android package name
   - SHA-1 certificate fingerprint: Your debug/release keystore SHA-1
   - Copy the Client ID (ends with `.googleusercontent.com`)

## Configuration

1. **Create Environment File**
   ```bash
   cp .env.example .env
   ```

2. **Update .env File**
   ```
   EXPO_PUBLIC_GOOGLE_WEB_CLIENT_ID=your_web_client_id.googleusercontent.com
   EXPO_PUBLIC_GOOGLE_IOS_CLIENT_ID=your_ios_client_id.googleusercontent.com
   EXPO_PUBLIC_GOOGLE_ANDROID_CLIENT_ID=your_android_client_id.googleusercontent.com
   ```

3. **Update app.json**
   Replace `YOUR_IOS_CLIENT_ID` in the `iosUrlScheme` with your actual iOS client ID:
   ```json
   {
     "expo": {
       "plugins": [
         [
           "@react-native-google-signin/google-signin",
           {
             "iosUrlScheme": "com.googleusercontent.apps.YOUR_ACTUAL_IOS_CLIENT_ID"
           }
         ]
       ]
     }
   }
   ```

## Build and Run

1. **Prebuild the project** (required after configuration changes)
   ```bash
   npx expo prebuild --clean
   ```

2. **Run on device/simulator**
   ```bash
   # For iOS
   npx expo run:ios
   
   # For Android  
   npx expo run:android
   ```

## Troubleshooting

- **iOS**: Make sure the `iosUrlScheme` in app.json matches your iOS client ID
- **Android**: Ensure your SHA-1 fingerprint is correctly added to the Android client ID
- **Both**: Verify that the bundle ID/package name matches what you configured in Google Cloud Console

## Notes

- Google Sign-in requires a physical device or simulator with Google Play Services
- This setup doesn't work in Expo Go - you need a development build
- The sign-in will redirect to home screen after successful authentication 