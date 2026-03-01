# Uber Base Design System - Typography Documentation

## Overview

The RepTracker app uses the complete **Uber Base Design System** typography system. This provides a comprehensive, scalable approach to typography that ensures consistency, accessibility, and visual hierarchy throughout the application.

## 🎯 Typography Philosophy

The Uber Base typography system is built on these core principles:

### **Hierarchy & Clarity**
- Clear visual hierarchy guides users through content
- Consistent spacing and sizing create rhythm
- Strategic use of weight and size emphasizes important information

### **Accessibility First**
- High contrast ratios for readability
- Appropriate line heights for comfortable reading
- Scalable sizing that works across devices

### **Brand Consistency**
- Uber Move and Uber Move Text fonts maintain brand identity
- Consistent application creates cohesive experience
- Professional, clean aesthetic

## 🎨 Typography Architecture

The system is organized into **four main categories** based on function and hierarchy:

### 1. **Display** (Hero & Large Promotional)
The largest text sizes for impactful, attention-grabbing content:
- `display-large` - 96px/112px line height
- `display-medium` - 52px/64px line height  
- `display-small` - 44px/52px line height
- `display-xsmall` - 36px/44px line height

**Font**: Uber Move Bold (700)  
**Usage**: Hero sections, large promotional content, app branding

### 2. **Headings** (Section & Content Hierarchy)
Bold text for creating clear content structure:
- `heading-xxlarge` - 40px/52px line height
- `heading-xlarge` - 36px/44px line height
- `heading-large` - 32px/40px line height
- `heading-medium` - 28px/36px line height
- `heading-small` - 24px/32px line height
- `heading-xsmall` - 20px/28px line height

**Font**: Uber Move Bold (700)  
**Usage**: Page titles, section headers, card titles, modal headers

### 3. **Labels** (UI Elements & Actions)
Medium weight text for interactive elements and emphasis:
- `label-large` - 18px/24px line height
- `label-medium` - 16px/20px line height
- `label-small` - 14px/16px line height
- `label-xsmall` - 12px/16px line height

**Font**: Uber Move Text Medium (500)  
**Usage**: Button text, form labels, tabs, navigation, badges, chips

### 4. **Paragraphs** (Body & Reading Content)
Regular weight text for comfortable reading:
- `paragraph-large` - 18px/28px line height
- `paragraph-medium` - 16px/24px line height
- `paragraph-small` - 14px/20px line height
- `paragraph-xsmall` - 12px/20px line height

**Font**: Uber Move Text Regular (400)  
**Usage**: Body text, descriptions, captions, helper text

## 🚀 Usage

### Basic Typography Component
```tsx
import { Typography } from '@/components';

<Typography variant="heading-large" color="contentPrimary">
  Workout Overview
</Typography>

<Typography variant="paragraph-medium" color="contentSecondary">
  Track your fitness progress with detailed workout analytics.
</Typography>
```

### Common Patterns
```tsx
// Page Header
<Typography variant="heading-xlarge" color="contentPrimary">
  My Workouts
</Typography>

// Section Header
<Typography variant="heading-small" color="contentPrimary">
  Recent Sessions
</Typography>

// Card Title
<Typography variant="heading-xsmall" color="contentPrimary">
  Upper Body Strength
</Typography>

// Body Text
<Typography variant="paragraph-medium" color="contentSecondary">
  Completed 3 sets of 12 reps with 135 lbs
</Typography>

// Button Text
<Typography variant="label-medium" color="contentOnColor">
  Start Workout
</Typography>

// Form Label
<Typography variant="label-small" color="contentPrimary">
  Exercise Name
</Typography>

// Caption/Meta Text
<Typography variant="paragraph-xsmall" color="contentTertiary">
  Last updated 2 hours ago
</Typography>
```

## 📐 Hierarchy Guidelines

### **When to Use Each Category**

#### **Display Variants**
- ✅ **App splash screens and onboarding**
- ✅ **Empty states with large messaging**
- ✅ **Marketing/promotional sections**
- ❌ Regular page content
- ❌ Within cards or constrained spaces

#### **Heading Variants**
- ✅ **Page titles** (`heading-xlarge` or `heading-large`)
- ✅ **Section headers** (`heading-medium` or `heading-small`)
- ✅ **Card/modal titles** (`heading-small` or `heading-xsmall`)
- ✅ **List item titles** (`heading-xsmall`)
- ❌ Body content
- ❌ Button text

#### **Label Variants**
- ✅ **Button and action text** (`label-medium` or `label-small`)
- ✅ **Form field labels** (`label-small`)
- ✅ **Navigation items** (`label-medium`)
- ✅ **Tabs and segmented controls** (`label-medium`)
- ✅ **Badges and status indicators** (`label-xsmall`)
- ❌ Long-form reading content
- ❌ Multi-line descriptions

#### **Paragraph Variants**
- ✅ **Body text and descriptions** (`paragraph-medium`)
- ✅ **List item descriptions** (`paragraph-small`)
- ✅ **Helper text and instructions** (`paragraph-small`)
- ✅ **Captions and metadata** (`paragraph-xsmall`)
- ❌ Headings or titles
- ❌ Interactive element text

## 🎯 Context-Specific Usage

### **Workout Screens**
```tsx
// Workout title
<Typography variant="heading-large" color="contentPrimary">
  Push Day Workout
</Typography>

// Exercise name
<Typography variant="heading-xsmall" color="contentPrimary">
  Bench Press
</Typography>

// Set information
<Typography variant="paragraph-medium" color="contentSecondary">
  Set 1: 12 reps × 185 lbs
</Typography>

// Timer display
<Typography variant="label-large" color="contentPrimary">
  02:30
</Typography>

// Action button
<Typography variant="label-medium" color="contentOnColor">
  Complete Set
</Typography>
```

### **List Items & Cards**
```tsx
// Card header
<Typography variant="heading-small" color="contentPrimary">
  Recent Workouts
</Typography>

// List item title
<Typography variant="heading-xsmall" color="contentPrimary">
  Upper Body Strength
</Typography>

// List item subtitle
<Typography variant="paragraph-small" color="contentSecondary">
  45 minutes • 8 exercises
</Typography>

// List item meta
<Typography variant="paragraph-xsmall" color="contentTertiary">
  Completed today
</Typography>
```

### **Forms & Input**
```tsx
// Form section title
<Typography variant="heading-small" color="contentPrimary">
  Exercise Details
</Typography>

// Field label
<Typography variant="label-small" color="contentPrimary">
  Weight (lbs)
</Typography>

// Helper text
<Typography variant="paragraph-xsmall" color="contentSecondary">
  Enter the weight used for this set
</Typography>

// Error message
<Typography variant="paragraph-xsmall" color="contentNegative">
  Please enter a valid weight
</Typography>
```

## 📱 Responsive Considerations

### **Mobile-First Approach**
- Prefer smaller variants on mobile (`heading-small` vs `heading-large`)
- Use `paragraph-medium` as the standard body text size
- Ensure touch targets with `label-medium` have at least 44px height

### **Scaling Recommendations**
```tsx
// Mobile preference
<Typography variant="heading-small" color="contentPrimary">
  Section Title
</Typography>

// Tablet/larger screens could use
<Typography variant="heading-medium" color="contentPrimary">
  Section Title
</Typography>
```

## 🎨 Color Pairing Guidelines

### **Primary Content Hierarchy**
```tsx
// Most important text
<Typography variant="heading-large" color="contentPrimary">
  Primary Title
</Typography>

// Secondary information
<Typography variant="paragraph-medium" color="contentSecondary">
  Supporting description
</Typography>

// Tertiary/meta information
<Typography variant="paragraph-small" color="contentTertiary">
  Additional context
</Typography>
```

### **Colored Backgrounds**
```tsx
// On accent/colored backgrounds
<Typography variant="label-medium" color="contentOnColor">
  Button Text
</Typography>

// On inverse/dark backgrounds
<Typography variant="heading-medium" color="contentInversePrimary">
  Dark Theme Title
</Typography>
```

## ⚠️ Common Mistakes to Avoid

### **❌ Don't Mix Categories Incorrectly**
```tsx
// Wrong: Using paragraph for button text
<Typography variant="paragraph-medium" color="contentOnColor">
  Button Text
</Typography>

// Correct: Use label for UI elements
<Typography variant="label-medium" color="contentOnColor">
  Button Text
</Typography>
```

### **❌ Don't Skip Hierarchy Levels**
```tsx
// Wrong: Jumping from large to xsmall
<Typography variant="heading-large">Title</Typography>
<Typography variant="heading-xsmall">Subtitle</Typography>

// Better: Gradual hierarchy
<Typography variant="heading-large">Title</Typography>
<Typography variant="heading-medium">Subtitle</Typography>
```

### **❌ Don't Use Wrong Colors**
```tsx
// Wrong: Using legacy color names
<Typography variant="heading-large" color="text-primary">
  Title
</Typography>

// Correct: Use semantic colors
<Typography variant="heading-large" color="contentPrimary">
  Title
</Typography>
```

## 📏 Spacing & Layout

### **Recommended Margins**
```tsx
const styles = StyleSheet.create({
  // Between different hierarchy levels
  titleSpacing: {
    marginBottom: 8, // heading to paragraph
  },
  
  // Between same-level elements
  paragraphSpacing: {
    marginBottom: 16, // paragraph to paragraph
  },
  
  // Section spacing
  sectionSpacing: {
    marginBottom: 24, // between major sections
  },
});
```

### **Line Height Considerations**
- All line heights are carefully calculated for optimal readability
- Don't override line heights unless absolutely necessary
- The built-in line heights provide proper vertical rhythm

## 🔧 Advanced Usage

### **Custom Styling**
```tsx
// Override specific styles while maintaining typography
<Typography 
  variant="heading-medium" 
  color="contentPrimary"
  style={{ textAlign: 'center', marginBottom: 20 }}
>
  Centered Title
</Typography>

// Multiple style overrides
<Typography 
  variant="paragraph-medium" 
  color="contentSecondary"
  style={{ 
    fontStyle: 'italic', 
    textAlign: 'center',
    maxWidth: 280 
  }}
>
  Italic centered text with max width
</Typography>
```

### **Dynamic Typography**
```tsx
const getTitleVariant = (screenSize: 'mobile' | 'tablet'): TypographyVariant => {
  return screenSize === 'mobile' ? 'heading-medium' : 'heading-large';
};

<Typography variant={getTitleVariant(screenSize)} color="contentPrimary">
  Responsive Title
</Typography>
```

## 📚 Quick Reference

### **Hierarchy Scale (Largest to Smallest)**
1. **Display**: `display-large` → `display-medium` → `display-small` → `display-xsmall`
2. **Headings**: `heading-xxlarge` → `heading-xlarge` → `heading-large` → `heading-medium` → `heading-small` → `heading-xsmall`
3. **Labels**: `label-large` → `label-medium` → `label-small` → `label-xsmall`
4. **Paragraphs**: `paragraph-large` → `paragraph-medium` → `paragraph-small` → `paragraph-xsmall`

### **Most Common Combinations**
- **Page Title**: `heading-xlarge` + `contentPrimary`
- **Section Header**: `heading-small` + `contentPrimary`
- **Body Text**: `paragraph-medium` + `contentSecondary`
- **Button Text**: `label-medium` + `contentOnColor`
- **Form Label**: `label-small` + `contentPrimary`
- **Caption**: `paragraph-xsmall` + `contentTertiary`

### **Font Families Summary**
- **Uber Move** (Bold 700): Display and Heading variants
- **Uber Move Text** (Medium 500): Label variants  
- **Uber Move Text** (Regular 400): Paragraph variants

---

This typography system provides a robust foundation for consistent, readable, and hierarchical text throughout the RepTracker application. Always use the Typography component and semantic color system for the best results. 