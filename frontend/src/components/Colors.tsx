import React from 'react';

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

// Color utility functions
export const getColorClass = (color: AppColor, type: 'text' | 'bg' | 'border' = 'text'): string => {
  const prefix = type === 'text' ? 'text' : type === 'bg' ? 'bg' : 'border';
  
  // Convert color name to CSS variable format
  const cssVariableName = color.replace(/-/g, '-');
  
  return `${prefix}-${cssVariableName}`;
};

// Color swatch component
interface ColorSwatchProps {
  color: AppColor;
  name: string;
  showHex?: boolean;
  size?: 'small' | 'medium' | 'large';
  className?: string;
}

export const ColorSwatch: React.FC<ColorSwatchProps> = ({ 
  color, 
  name, 
  showHex = true,
  size = 'medium',
  className = '' 
}) => {
  const sizeClasses = {
    small: 'w-8 h-8',
    medium: 'w-12 h-12',
    large: 'w-16 h-16'
  };

  // Color hex values for display
  const colorHexMap: Record<AppColor, string> = {
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

  const colorValue = colorHexMap[color] || '#000000';
  
  return (
    <div className={`flex flex-col items-center gap-2 ${className}`}>
      <div className={`${sizeClasses[size]} ${getColorClass(color, 'bg')} rounded-lg border border-gray-200 shadow-sm`} />
      <div className="text-center">
        <div className="text-xs font-medium text-gray-900">{name}</div>
        {showHex && <div className="text-xs text-gray-500 font-mono">{colorValue}</div>}
      </div>
    </div>
  );
};

// Color family demo component
interface ColorFamilyDemoProps {
  family: ColorFamily;
  className?: string;
}

export const ColorFamilyDemo: React.FC<ColorFamilyDemoProps> = ({ family, className = '' }) => {
  const variants = ['bright', 'default', 'dark-1', 'light-1', 'light-2'] as const;
  
  return (
    <div className={`flex flex-col gap-4 ${className}`}>
      <h3 className="text-lg font-semibold capitalize">{family}</h3>
      <div className="flex gap-4">
        {variants.map((variant) => {
          const colorKey = variant === 'default' ? family : `${family}-${variant}`;
          return (
            <ColorSwatch
              key={variant}
              color={colorKey as AppColor}
              name={variant === 'default' ? 'default' : variant}
              size="medium"
            />
          );
        })}
      </div>
    </div>
  );
};

// Complete color demo component
interface ColorDemoProps {
  className?: string;
}

export const ColorDemo: React.FC<ColorDemoProps> = ({ className = '' }) => {
  const colorFamilies: ColorFamily[] = ['yellow', 'orange', 'red', 'pink', 'purple', 'blue', 'cyan', 'teal', 'green'];
  const grayscaleColors: { color: GrayscaleColor; name: string }[] = [
    { color: 'dark', name: 'Dark' },
    { color: 'light', name: 'Light' },
    { color: 'light-gray-4', name: 'Light Gray 4' },
    { color: 'light-gray-3', name: 'Light Gray 3' },
    { color: 'light-gray-2', name: 'Light Gray 2' },
    { color: 'light-gray-1', name: 'Light Gray 1' },
    { color: 'white', name: 'White' },
  ];

  return (
    <div className={`bg-white p-8 ${className}`}>
      <div className="max-w-6xl mx-auto">
        <div className="flex flex-col gap-12">
          {/* Header */}
          <div>
            <h1 className="text-4xl font-bold text-gray-900 mb-4">Color System</h1>
            <p className="text-lg text-gray-600">
              Complete color palette from Figma design with all variants
            </p>
            <div className="bg-gray-200 h-px w-full mt-6" />
          </div>

          {/* Grayscale Colors */}
          <div>
            <h2 className="text-2xl font-semibold text-gray-900 mb-6">Grayscale</h2>
            <div className="flex gap-6 flex-wrap">
              {grayscaleColors.map((item) => (
                <ColorSwatch
                  key={item.color}
                  color={item.color}
                  name={item.name}
                  size="large"
                />
              ))}
            </div>
          </div>

          {/* Color Families */}
          <div>
            <h2 className="text-2xl font-semibold text-gray-900 mb-6">Color Families</h2>
            <div className="flex flex-col gap-8">
              {colorFamilies.map((family) => (
                <ColorFamilyDemo key={family} family={family} />
              ))}
            </div>
          </div>

          {/* Usage Examples */}
          <div>
            <h2 className="text-2xl font-semibold text-gray-900 mb-6">Usage Examples</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {/* Text Colors */}
              <div className="p-6 border border-gray-200 rounded-lg">
                <h3 className="font-semibold mb-4">Text Colors</h3>
                <div className="space-y-2">
                  <p className={getColorClass('dark')}>Dark text</p>
                  <p className={getColorClass('blue')}>Blue text</p>
                  <p className={getColorClass('red')}>Red text</p>
                  <p className={getColorClass('green')}>Green text</p>
                </div>
              </div>

              {/* Background Colors */}
              <div className="p-6 border border-gray-200 rounded-lg">
                <h3 className="font-semibold mb-4">Background Colors</h3>
                <div className="space-y-2">
                  <div className={`p-2 rounded ${getColorClass('blue-light-2', 'bg')}`}>
                    Light blue background
                  </div>
                  <div className={`p-2 rounded ${getColorClass('green-light-2', 'bg')}`}>
                    Light green background
                  </div>
                  <div className={`p-2 rounded ${getColorClass('red-light-2', 'bg')}`}>
                    Light red background
                  </div>
                </div>
              </div>

              {/* Border Colors */}
              <div className="p-6 border border-gray-200 rounded-lg">
                <h3 className="font-semibold mb-4">Border Colors</h3>
                <div className="space-y-2">
                  <div className={`p-2 border-2 rounded ${getColorClass('purple', 'border')}`}>
                    Purple border
                  </div>
                  <div className={`p-2 border-2 rounded ${getColorClass('orange', 'border')}`}>
                    Orange border
                  </div>
                  <div className={`p-2 border-2 rounded ${getColorClass('teal', 'border')}`}>
                    Teal border
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}; 