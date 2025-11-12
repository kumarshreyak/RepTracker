# Add project specific ProGuard rules here.
# By default, the flags in this file are appended to flags specified
# in /usr/local/Cellar/android-sdk/24.3.3/tools/proguard/proguard-android.txt
# You can edit the include path and order by changing the proguardFiles
# directive in build.gradle.
#
# For more details, see
#   http://developer.android.com/guide/developing/tools/proguard.html

# react-native-reanimated
-keep class com.swmansion.reanimated.** { *; }
-keep class com.facebook.react.turbomodule.** { *; }

# Add any project specific keep options here:
# Google Play Services
-keep class com.google.android.gms.** { *; }
-dontwarn com.google.android.gms.**

# Google Sign-In
-keep class com.google.android.gms.auth.** { *; }
-keep class com.google.android.gms.common.** { *; }
-keep class com.google.android.gms.internal.auth.** { *; }

# Google API Client Library
-keep class com.google.api.client.** { *; }
-dontwarn com.google.api.client.**

# If using specific Google APIs (e.g., Drive, Calendar), add rules for those as well.
# Example for Google Drive API:
#-keep class com.google.api.services.drive.** { *; }
#-dontwarn com.google.api.services.drive.**

# Keep all custom classes that extend from Google Play Services classes
-keep class * extends com.google.android.gms.common.api.GoogleApiClient$Builder { *; }
-keep class * extends com.google.android.gms.common.api.GoogleApiClient$ConnectionCallbacks { *; }
-keep class * extends com.google.android.gms.common.api.GoogleApiClient$OnConnectionFailedListener { *; }

# General rules for annotations and interfaces used by Google Play Services
-keepattributes Signature
-keepattributes SourceFile,LineNumberTable
-keep public interface com.google.android.gms.common.api.Api$ApiOptions { *; }