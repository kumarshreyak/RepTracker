import React from 'react';
import { TouchableOpacity, View, ViewStyle, TextStyle } from 'react-native';
import { Typography } from './Typography';
import { getColor } from './Colors';

export interface ButtonProps {
  children: React.ReactNode;
  variant?: 'default' | 'primary' | 'danger' | 'secondary' | 'text' | 'success';
  size?: 'small' | 'default' | 'large' | 'xlarge';
  icon?: React.ReactNode;
  iconOnly?: boolean;
  style?: ViewStyle;
  onPress?: () => void;
  disabled?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'default',
  size = 'default',
  icon,
  iconOnly = false,
  style = {},
  disabled = false,
  onPress,
}) => {
  // Size configurations - using proper spacing values
  const sizeConfig = {
    small: {
      paddingVertical: 6,
      paddingHorizontal: 10,
      iconSize: 12,
      gap: 6,
    },
    default: {
      paddingVertical: 8,
      paddingHorizontal: 12,
      iconSize: 14,
      gap: 8,
    },
    large: {
      paddingVertical: 10,
      paddingHorizontal: 16,
      iconSize: 16,
      gap: 10,
    },
    xlarge: {
      paddingVertical: 12,
      paddingHorizontal: 20,
      iconSize: 20,
      gap: 12,
    },
  };

  // Icon only size configurations
  const iconOnlySize = {
    small: { padding: 6 },
    default: { padding: 8 },
    large: { padding: 10 },
    xlarge: { padding: 12 },
  };

  // Variant styles using semantic colors
  const variantStyles = {
    default: {
      backgroundColor: getColor('backgroundPrimary'),
      borderColor: getColor('borderOpaque'),
      textColor: getColor('contentPrimary'),
    },
    primary: {
      backgroundColor: getColor('backgroundAccent'),
      borderColor: getColor('backgroundAccent'),
      textColor: getColor('contentOnColor'),
    },
    danger: {
      backgroundColor: getColor('backgroundNegative'),
      borderColor: getColor('backgroundNegative'),
      textColor: getColor('contentOnColor'),
    },
    secondary: {
      backgroundColor: getColor('backgroundSecondary'),
      borderColor: getColor('borderOpaque'),
      textColor: getColor('contentPrimary'),
    },
    text: {
      backgroundColor: 'transparent',
      borderColor: 'transparent',
      textColor: getColor('contentAccent'),
    },
    success: {
      backgroundColor: getColor('backgroundPositive'),
      borderColor: getColor('backgroundPositive'),
      textColor: getColor('contentOnColor'),
    },
  };

  const currentSize = sizeConfig[size];
  const currentVariant = variantStyles[variant];

  const buttonStyle: ViewStyle = {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 3,
    borderWidth: variant === 'text' ? 0 : 1,
    backgroundColor: currentVariant.backgroundColor,
    borderColor: currentVariant.borderColor,
    opacity: disabled ? 0.5 : 1,
    gap: currentSize.gap,
    ...(iconOnly ? iconOnlySize[size] : {
      paddingVertical: currentSize.paddingVertical,
      paddingHorizontal: currentSize.paddingHorizontal,
    }),
    ...style,
  };

  // Map button sizes to Typography variants
  const textVariant = size === 'small' ? 'label-xsmall' : 
                      size === 'large' ? 'label-medium' :
                      size === 'xlarge' ? 'label-large' : 'label-small';

  // Determine text color based on variant using semantic colors
  const textColor = variant === 'primary' || variant === 'danger' || variant === 'success' 
    ? 'contentOnColor' 
    : variant === 'text' 
    ? 'contentAccent' 
    : 'contentPrimary';

  return (
    <TouchableOpacity 
      style={buttonStyle}
      onPress={onPress}
      disabled={disabled}
      activeOpacity={0.7}
    >
      {icon && (
        <View style={{ width: currentSize.iconSize, height: currentSize.iconSize }}>
          {icon}
        </View>
      )}
      {!iconOnly && (
        <Typography 
          variant={textVariant} 
          color={textColor}
        >
          {children}
        </Typography>
      )}
    </TouchableOpacity>
  );
};

// Edit Icon component for React Native
export const EditIcon: React.FC<{ color?: string; size?: number }> = ({ 
  color = 'currentColor', 
  size = 12 
}) => {
  // For now, we'll use a simple placeholder. In a real app, you'd use react-native-svg
  return (
    <View style={{
      width: size,
      height: size,
      backgroundColor: color === 'currentColor' ? getColor('contentPrimary') : color,
      borderRadius: 2,
    }} />
  );
};

// Rect Primary Button - Based on Uber Design System
export interface RectPrimaryButtonProps {
  children: React.ReactNode;
  onPress?: () => void;
  disabled?: boolean;
  style?: ViewStyle;
  icon?: React.ReactNode;
}

export const RectPrimaryButton: React.FC<RectPrimaryButtonProps> = ({
  children,
  onPress,
  disabled = false,
  style = {},
  icon,
}) => {
  const buttonStyle: ViewStyle = {
    backgroundColor: getColor('primaryA'), // Black background as shown in Uber design
    borderRadius: 0, // Rectangular shape with sharp corners
    paddingVertical: 16,
    paddingHorizontal: 24,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
    opacity: disabled ? 0.5 : 1,
    minHeight: 48, // Consistent minimum height for touch targets
    ...style,
  };

  return (
    <TouchableOpacity 
      style={buttonStyle}
      onPress={onPress}
      disabled={disabled}
      activeOpacity={0.8}
    >
      {icon && (
        <View style={{ width: 16, height: 16 }}>
          {icon}
        </View>
      )}
      <Typography 
        variant="label-medium" 
        color="contentOnColor"
      >
        {children}
      </Typography>
    </TouchableOpacity>
  );
}; 