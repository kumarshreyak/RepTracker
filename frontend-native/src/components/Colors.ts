// Uber Base Design System - Primitive Colors
// Based on the "Primitive Scale - NEW" from Figma

// Gray Scale (Black to White)
export const primitiveColors = {
  // Gray scale
  gray: {
    50: '#F3F3F3',
    100: '#E8E8E8', 
    200: '#DDDDDD',
    300: '#BBBBBB',
    400: '#A6A6A6',
    500: '#868686',
    600: '#727272',
    700: '#5E5E5E',
    800: '#4B4B4B',
    900: '#282828',
  },
  
  // Red scale
  red: {
    50: '#FFF0EE',
    100: '#FFE1DE',
    200: '#FFD2CD',
    300: '#FFB2AB',
    400: '#FC7F79',
    500: '#F83446',
    600: '#DE1135',
    700: '#BB032A',
    800: '#950F22',
    900: '#520810',
  },
  
  // Orange scale
  orange: {
    50: '#FFF0E9',
    100: '#FEE2D4',
    200: '#FFD3BC',
    300: '#FFB48C',
    400: '#FC823A',
    500: '#E65300',
    600: '#C54600',
    700: '#A33B04',
    800: '#823006',
    900: '#461A00',
  },
  
  // Amber scale
  amber: {
    50: '#FFF1E1',
    100: '#FFE4B7',
    200: '#FFD5A1',
    300: '#FFB749',
    400: '#DF9500',
    500: '#C46E00',
    600: '#A95F03',
    700: '#904A07',
    800: '#763A00',
    900: '#3E2000',
  },
  
  // Yellow scale
  yellow: {
    50: '#FDF2DC',
    100: '#FBE5B6',
    200: '#FFD688',
    300: '#F6BC2F',
    400: '#D79900',
    500: '#B97502',
    600: '#9F6402',
    700: '#845201',
    800: '#6B4100',
    900: '#392300',
  },
  
  // Lime scale
  lime: {
    50: '#EEF6E3',
    100: '#DEEEC6',
    200: '#CAE6A0',
    300: '#A6D467',
    400: '#77B71C',
    500: '#5B9500',
    600: '#4F7F06',
    700: '#3F6900',
    800: '#365310',
    900: '#1B2D00',
  },
  
  // Green scale
  green: {
    50: '#EAF6ED',
    100: '#D3EFDA',
    200: '#B1EAC2',
    300: '#7FD99A',
    400: '#06C167',
    500: '#009A51',
    600: '#0E8345',
    700: '#166C3B',
    800: '#0D572D',
    900: '#002F14',
  },
  
  // Teal scale
  teal: {
    50: '#E2F8FB',
    100: '#CDEEF3',
    200: '#B0E7EF',
    300: '#77D5E3',
    400: '#01B8CA',
    500: '#0095A4',
    600: '#007F8C',
    700: '#016974',
    800: '#1A535A',
    900: '#002D33',
  },
  
  // Blue scale
  blue: {
    50: '#EFF4FE',
    100: '#DEE9FE',
    200: '#CDDEFF',
    300: '#A9C9FF',
    400: '#6DAAFB',
    500: '#068BEE',
    600: '#266EF1',
    700: '#175BCC',
    800: '#1948A3',
    900: '#002661',
  },
  
  // Purple scale
  purple: {
    50: '#F9F1FF',
    100: '#F2E3FF',
    200: '#EBD5FF',
    300: '#DDB9FF',
    400: '#C490F9',
    500: '#A964F7',
    600: '#944DE7',
    700: '#7C3EC3',
    800: '#633495',
    900: '#3A1659',
  },
  
  // Magenta scale
  magenta: {
    50: '#FEEFF9',
    100: '#FEDFF3',
    200: '#FFCEF2',
    300: '#FFACE5',
    400: '#F877D2',
    500: '#E142BC',
    600: '#CA26A5',
    700: '#A91A90',
    800: '#891869',
    900: '#50003F',
  },
  
  // Special colors
  black: '#000000',
  white: '#FFFFFF',
} as const;

// === UBER BASE SEMANTIC COLOR SYSTEM ===

// 02 Core Colors
export const coreColors = {
  primaryA: primitiveColors.black,
  primaryB: primitiveColors.white,
  accent: primitiveColors.blue[600],
  negative: primitiveColors.red[600],
  warning: primitiveColors.yellow[300],
  positive: primitiveColors.green[600],
} as const;

// 03 Basic Semantic Colors
export const semanticColors = {
  // Background
  backgroundPrimary: coreColors.primaryB, // white
  backgroundSecondary: primitiveColors.gray[50],
  backgroundTertiary: primitiveColors.gray[100],
  backgroundInversePrimary: coreColors.primaryA, // black
  backgroundInverseSecondary: primitiveColors.gray[800],
  
  // Content
  contentPrimary: primitiveColors.black,
  contentSecondary: primitiveColors.gray[800],
  contentTertiary: primitiveColors.gray[700],
  contentInversePrimary: primitiveColors.white,
  contentInverseSecondary: primitiveColors.gray[200],
  contentInverseTertiary: primitiveColors.gray[400],
  
  // Border
  borderOpaque: primitiveColors.gray[100],
  borderTransparent: primitiveColors.black + '14', // 8% black (14 in hex = ~8% of 255)
  borderSelected: coreColors.primaryA, // black
  borderInverseOpaque: primitiveColors.gray[800],
  borderInverseTransparent: primitiveColors.white + '33', // 20% white (33 in hex = ~20% of 255)
  borderInverseSelected: coreColors.primaryB, // white
} as const;

// 03 Semantic Extensions
export const semanticExtensions = {
  // Background Extensions
  backgroundStateDisabled: primitiveColors.gray[50],
  backgroundOverlayArt: primitiveColors.black + '00', // 0% black
  backgroundOverlayDark: primitiveColors.black + '80', // 50% black
  backgroundOverlayLight: primitiveColors.black + '14', // 8% black
  backgroundOverlayElevation: primitiveColors.black + '00', // 0% black
  backgroundAccent: primitiveColors.blue[600],
  backgroundNegative: primitiveColors.red[600],
  backgroundWarning: primitiveColors.yellow[300],
  backgroundPositive: primitiveColors.green[600],
  backgroundLightAccent: primitiveColors.blue[50],
  backgroundLightNegative: primitiveColors.red[50],
  backgroundLightWarning: primitiveColors.yellow[50],
  backgroundLightPositive: primitiveColors.green[50],
  backgroundAlwaysDark: primitiveColors.black,
  backgroundAlwaysLight: primitiveColors.white,
  
  // Content Extensions
  contentStateDisabled: primitiveColors.gray[400],
  contentOnColor: primitiveColors.white,
  contentOnColorInverse: primitiveColors.black,
  contentAccent: primitiveColors.blue[600],
  contentNegative: primitiveColors.red[600],
  contentWarning: primitiveColors.yellow[600],
  contentPositive: primitiveColors.green[600],
  
  // Border Extensions
  borderStateDisabled: primitiveColors.gray[50],
  borderAccent: primitiveColors.blue[600],
  borderNegative: primitiveColors.red[600],
  borderWarning: primitiveColors.yellow[600],
  borderPositive: primitiveColors.green[600],
  borderAccentLight: primitiveColors.blue[200],
} as const;

// Combined semantic color type
export type SemanticColor = 
  // Core colors
  | keyof typeof coreColors
  // Basic semantic colors
  | keyof typeof semanticColors
  // Semantic extensions
  | keyof typeof semanticExtensions
  // Legacy compatibility (temporary)
  | 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'info'
  | 'surface' | 'background' | 'text-primary' | 'text-secondary' | 'text-tertiary' | 'text-inverse'
  | 'border-light' | 'border-medium' | 'border-strong';

// Legacy compatibility mappings (temporary - to be removed)
const legacyCompatibility = {
  primary: primitiveColors.blue[500],
  secondary: primitiveColors.gray[500],
  success: primitiveColors.green[500],
  warning: primitiveColors.amber[500],
  danger: primitiveColors.red[500],
  info: primitiveColors.blue[400],
  surface: primitiveColors.white,
  background: primitiveColors.gray[50],
  'text-primary': primitiveColors.gray[900],
  'text-secondary': primitiveColors.gray[700],
  'text-tertiary': primitiveColors.gray[500],
  'text-inverse': primitiveColors.white,
  'border-light': primitiveColors.gray[200],
  'border-medium': primitiveColors.gray[300],
  'border-strong': primitiveColors.gray[400],
} as const;

// Main color lookup function
export const getColor = (color: SemanticColor): string => {
  // Check core colors
  if (color in coreColors) {
    return coreColors[color as keyof typeof coreColors];
  }
  
  // Check basic semantic colors
  if (color in semanticColors) {
    return semanticColors[color as keyof typeof semanticColors];
  }
  
  // Check semantic extensions
  if (color in semanticExtensions) {
    return semanticExtensions[color as keyof typeof semanticExtensions];
  }
  
  // Legacy compatibility (temporary)
  if (color in legacyCompatibility) {
    return legacyCompatibility[color as keyof typeof legacyCompatibility];
  }
  
  // Default fallback
  return primitiveColors.gray[900];
};

// Export primitives for direct access when needed
export { primitiveColors as colors };

// Color utility functions
export const getPrimitiveColor = (family: keyof typeof primitiveColors, shade: keyof typeof primitiveColors.gray): string => {
  if (family === 'black' || family === 'white') {
    return primitiveColors[family];
  }
  return (primitiveColors[family] as any)?.[shade] || primitiveColors.gray[900];
};

// Export types for TypeScript
export type PrimitiveColorFamily = keyof typeof primitiveColors;
export type ColorShade = keyof typeof primitiveColors.gray; 