// Color type definitions
export type GrayscaleColor = 'dark' | 'light' | 'white' | 'light-gray-1' | 'light-gray-2' | 'light-gray-3' | 'light-gray-4';

export type ColorFamily = 'yellow' | 'orange' | 'red' | 'pink' | 'purple' | 'blue' | 'cyan' | 'teal' | 'green';

export type AppColor = 
  | GrayscaleColor 
  | ColorFamily
  | `${ColorFamily}-bright`
  | `${ColorFamily}-dark-1`
  | `${ColorFamily}-light-1`
  | `${ColorFamily}-light-2`;

// Color hex values for React Native
export const colors: Record<AppColor, string> = {
  // Grayscale colors
  'dark': '#333333',
  'light': '#757575',
  'white': '#ffffff',
  'light-gray-1': '#fafafa',
  'light-gray-2': '#f2f2f2',
  'light-gray-3': '#e8e8e8',
  'light-gray-4': '#e0e0e0',
  
  // Yellow colors
  'yellow-bright': '#fcb400',
  'yellow': '#e08d00',
  'yellow-dark-1': '#b87503',
  'yellow-light-1': '#ffd66e',
  'yellow-light-2': '#ffeab6',
  
  // Orange colors
  'orange-bright': '#ff6f2c',
  'orange': '#f7653b',
  'orange-dark-1': '#d74d26',
  'orange-light-1': '#ffa981',
  'orange-light-2': '#fee2d5',
  
  // Red colors
  'red-bright': '#f82b60',
  'red': '#ef3061',
  'red-dark-1': '#ba1e45',
  'red-light-1': '#ff9eb7',
  'red-light-2': '#ffdce5',
  
  // Pink colors
  'pink-bright': '#ff08c2',
  'pink': '#e929ba',
  'pink-dark-1': '#b2158b',
  'pink-light-1': '#f99de2',
  'pink-light-2': '#ffdaf6',
  
  // Purple colors
  'purple-bright': '#8b46ff',
  'purple': '#7c39ed',
  'purple-dark-1': '#6b1cb0',
  'purple-light-1': '#cdb0ff',
  'purple-light-2': '#ede3fe',
  
  // Blue colors
  'blue-bright': '#2d7ff9',
  'blue': '#1283da',
  'blue-dark-1': '#2750ae',
  'blue-light-1': '#9cc7ff',
  'blue-light-2': '#cfdfff',
  
  // Cyan colors
  'cyan-bright': '#18bfff',
  'cyan': '#01a9db',
  'cyan-dark-1': '#0b76b7',
  'cyan-light-1': '#77d1f3',
  'cyan-light-2': '#d0f0fd',
  
  // Teal colors
  'teal-bright': '#20d9d2',
  'teal': '#02aaa4',
  'teal-dark-1': '#06a09b',
  'teal-light-1': '#72ddc3',
  'teal-light-2': '#c2f5e9',
  
  // Green colors
  'green-bright': '#20c933',
  'green': '#11af22',
  'green-dark-1': '#338a17',
  'green-light-1': '#93e088',
  'green-light-2': '#d1f7c4',
};

// Color utility function
export const getColor = (color: AppColor): string => {
  return colors[color] || colors.dark;
}; 