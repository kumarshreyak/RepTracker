# SF Pro Fonts Setup

This project uses SF Pro Display and SF Pro Text fonts as specified in the Figma design.

## Font Download Instructions

1. **Download SF Pro fonts from Apple:**
   - Visit: https://developer.apple.com/fonts/
   - Download the SF Pro font family
   - You'll get a package with various font files

2. **Extract and convert fonts:**
   - Extract the downloaded package
   - Look for SF Pro Display and SF Pro Text families
   - Convert TTF/OTF files to WOFF2 format for web optimization

3. **Required font files:**
   Place the following WOFF2 files in this directory:

   **SF Pro Display:**
   - `SFProDisplay-Regular.woff2` (400 weight)
   - `SFProDisplay-Medium.woff2` (500 weight)
   - `SFProDisplay-Semibold.woff2` (600 weight)
   - `SFProDisplay-Bold.woff2` (700 weight)

   **SF Pro Text:**
   - `SFProText-Regular.woff2` (400 weight)
   - `SFProText-Medium.woff2` (500 weight)  
   - `SFProText-Semibold.woff2` (600 weight)
   - `SFProText-Bold.woff2` (700 weight)

## Font Conversion

If you need to convert fonts to WOFF2, you can use online tools or command line:

```bash
# Using woff2 CLI tool
npm install -g woff2
woff2_compress SFProDisplay-Regular.ttf
```

## Fallback Behavior

The fonts are configured with Apple system font fallbacks, so on Apple devices, the system SF Pro fonts will be used automatically even without the font files.

## Usage

The fonts are automatically available through CSS variables:
- `--font-sf-pro-display`
- `--font-sf-pro-text`

And through the Typography component with predefined variants matching the Figma design.

## Typography System

The following typography styles are available:
- `display-large`: SF Pro Display Bold, 52px, line-height 56px
- `display-medium`: SF Pro Display Semibold, 35px
- `heading-xsmall`: SF Pro Text Bold, 15px, line-height 22px  
- `text-large-paragraph`: SF Pro Text Regular, 15px, line-height 22px 