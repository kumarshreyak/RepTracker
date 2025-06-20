# Uber Base Design System - Color Documentation

## Overview

The GymLog app now uses the complete **Uber Base Design System** semantic color system. This provides a comprehensive, scalable approach to color that ensures consistency and makes it easy to maintain and update the visual design.

## 🎨 Color Architecture

The system is organized into three layers:

### 1. **Core Colors** (Foundation)
These are the foundation colors that drive the entire system:

- `primaryA` - Black (`#000000`)
- `primaryB` - White (`#FFFFFF`) 
- `accent` - Blue 600 (`#266EF1`)
- `negative` - Red 600 (`#DE1135`)
- `warning` - Yellow 300 (`#F6BC2F`)
- `positive` - Green 600 (`#0E8345`)

### 2. **Basic Semantic Colors** (Primary Usage)
Organized by purpose and context:

#### **Background Colors:**
- `backgroundPrimary` - Primary background (white)
- `backgroundSecondary` - Secondary background (gray 50)
- `backgroundTertiary` - Tertiary background (gray 100)
- `backgroundInversePrimary` - Inverse primary (black)
- `backgroundInverseSecondary` - Inverse secondary (gray 800)

#### **Content Colors:**
- `contentPrimary` - Primary text (black)
- `contentSecondary` - Secondary text (gray 800)
- `contentTertiary` - Tertiary text (gray 700)
- `contentInversePrimary` - Inverse primary text (white)
- `contentInverseSecondary` - Inverse secondary text (gray 200)
- `contentInverseTertiary` - Inverse tertiary text (gray 400)

#### **Border Colors:**
- `borderOpaque` - Standard borders (gray 100)
- `borderTransparent` - Transparent borders (8% black)
- `borderSelected` - Selected state borders (black)
- `borderInverseOpaque` - Inverse borders (gray 800)
- `borderInverseTransparent` - Inverse transparent (20% white)
- `borderInverseSelected` - Inverse selected (white)

### 3. **Semantic Extensions** (Specialized Usage)
Extended colors for specific use cases and states:

#### **Background Extensions:**
- State: `backgroundStateDisabled`, `backgroundOverlayArt`, `backgroundOverlayDark`, `backgroundOverlayLight`, `backgroundOverlayElevation`
- Brand: `backgroundAccent`, `backgroundNegative`, `backgroundWarning`, `backgroundPositive`
- Light variants: `backgroundLightAccent`, `backgroundLightNegative`, `backgroundLightWarning`, `backgroundLightPositive`
- Always: `backgroundAlwaysDark`, `backgroundAlwaysLight`

#### **Content Extensions:**
- State: `contentStateDisabled`, `contentOnColor`, `contentOnColorInverse`
- Brand: `contentAccent`, `contentNegative`, `contentWarning`, `contentPositive`

#### **Border Extensions:**
- State: `borderStateDisabled`
- Brand: `borderAccent`, `borderNegative`, `borderWarning`, `borderPositive`
- Light variants: `borderAccentLight`

## 🚀 Usage

### Basic Usage
```tsx
import { getColor } from '@/components';

// Use semantic colors in StyleSheet
const styles = StyleSheet.create({
  container: {
    backgroundColor: getColor('backgroundSecondary'), // Gray 50 background
  },
  card: {
    backgroundColor: getColor('backgroundPrimary'), // White background
    borderColor: getColor('borderOpaque'), // Gray 100 border
  },
  text: {
    color: getColor('contentPrimary'), // Black text
  },
});
```

### Typography Component
```tsx
<Typography variant="heading-large" color="contentPrimary">
  Main Heading
</Typography>
<Typography variant="paragraph-medium" color="contentSecondary">
  Secondary text content
</Typography>
```

### Button Styling
```tsx
// Accent button
<View style={{ backgroundColor: getColor('backgroundAccent') }}>
  <Typography color="contentOnColor">Accent Button</Typography>
</View>

// Negative action
<View style={{ backgroundColor: getColor('backgroundNegative') }}>
  <Typography color="contentOnColor">Delete</Typography>
</View>
```

## 📋 Best Practices

### 1. **Always Use Semantic Colors**
❌ **Don't use primitive colors directly:**
```tsx
backgroundColor: '#FFFFFF' // Wrong
backgroundColor: primitiveColors.white // Wrong
```

✅ **Use semantic colors:**
```tsx
backgroundColor: getColor('backgroundPrimary') // Correct
```

### 2. **Choose the Right Semantic Level**
- Use **Basic Semantic Colors** for 90% of cases
- Use **Core Colors** when you need direct access to foundation colors
- Use **Semantic Extensions** for specific states and specialized use cases

### 3. **Follow Context Patterns**
- **Page backgrounds**: `backgroundSecondary` (gray 50)
- **Card backgrounds**: `backgroundPrimary` (white)
- **Card borders**: `borderOpaque` (gray 100)
- **Primary text**: `contentPrimary` (black)
- **Secondary text**: `contentSecondary` (gray 800)
- **Accent actions**: `backgroundAccent` with `contentOnColor`

### 4. **Maintain Contrast**
- Use `contentPrimary` on light backgrounds
- Use `contentInversePrimary` on dark backgrounds
- Use `contentOnColor` on colored backgrounds

## 🔄 Migration from Legacy System

The legacy color system is temporarily supported for backward compatibility:

```tsx
// Legacy (still works but deprecated)
getColor('primary') // Maps to primitiveColors.blue[500]
getColor('text-primary') // Maps to primitiveColors.gray[900]

// New semantic system (preferred)
getColor('accent') // Core accent color
getColor('contentPrimary') // Primary content color
```

## 🎨 Color Demo Component

To see all available colors in action, you can use the `ColorDemo` component:

```tsx
import { ColorDemo } from '@/components';

// Renders a complete showcase of all semantic colors
<ColorDemo />
```

## 💡 Examples

### Common Layout Pattern
```tsx
const styles = StyleSheet.create({
  // Page container
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'), // Gray 50
  },
  
  // Content card
  card: {
    backgroundColor: getColor('backgroundPrimary'), // White
    borderWidth: 1,
    borderColor: getColor('borderOpaque'), // Gray 100
    borderRadius: 8,
    padding: 16,
  },
  
  // Primary heading
  title: {
    color: getColor('contentPrimary'), // Black
  },
  
  // Secondary text
  subtitle: {
    color: getColor('contentSecondary'), // Gray 800
  },
  
  // Accent button
  button: {
    backgroundColor: getColor('backgroundAccent'), // Blue 600
  },
  
  // Button text
  buttonText: {
    color: getColor('contentOnColor'), // White
  },
});
```

## 🔧 Advanced Usage

### Direct Access to Color Objects
```tsx
import { coreColors, semanticColors, semanticExtensions } from '@/components';

// Access core colors directly
const primaryBlack = coreColors.primaryA;
const accentBlue = coreColors.accent;

// Access semantic color objects
const allBackgroundColors = Object.keys(semanticColors)
  .filter(key => key.startsWith('background'));
```

### Custom Color Utilities
```tsx
// Create custom utilities for specific use cases
const getStateColor = (state: 'normal' | 'error' | 'success') => {
  switch (state) {
    case 'error': return getColor('contentNegative');
    case 'success': return getColor('contentPositive');
    default: return getColor('contentPrimary');
  }
};
```

## 📚 Resources

- **Primitive Colors**: All base colors from the Uber Base design system
- **Color Demo**: Use `<ColorDemo />` component to see all available colors
- **TypeScript Support**: Full type safety with `SemanticColor` type
- **Design System**: Based on official Uber Base design system specifications

---

This semantic color system provides a robust foundation for consistent, maintainable, and scalable color usage throughout the GymLog application. 