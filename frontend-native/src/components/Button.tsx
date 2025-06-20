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
  // Size configurations
  const sizeConfig = {
    small: {
      paddingVertical: 6,
      paddingHorizontal: 10,
      fontSize: 11,
      iconSize: 12,
      gap: 6,
    },
    default: {
      paddingVertical: 8,
      paddingHorizontal: 12,
      fontSize: 13,
      iconSize: 14,
      gap: 8,
    },
    large: {
      paddingVertical: 10,
      paddingHorizontal: 16,
      fontSize: 15,
      iconSize: 16,
      gap: 10,
    },
    xlarge: {
      paddingVertical: 12,
      paddingHorizontal: 20,
      fontSize: 17,
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

  // Variant styles using new semantic colors
  const variantStyles = {
    default: {
      backgroundColor: getColor('background'),
      borderColor: getColor('border-medium'),
      textColor: getColor('text-primary'),
    },
    primary: {
      backgroundColor: getColor('primary'),
      borderColor: getColor('primary'),
      textColor: getColor('text-inverse'),
    },
    danger: {
      backgroundColor: getColor('danger'),
      borderColor: getColor('danger'),
      textColor: getColor('text-inverse'),
    },
    secondary: {
      backgroundColor: getColor('surface'),
      borderColor: getColor('border-medium'),
      textColor: getColor('text-primary'),
    },
    text: {
      backgroundColor: 'transparent',
      borderColor: 'transparent',
      textColor: getColor('primary'),
    },
    success: {
      backgroundColor: getColor('success'),
      borderColor: getColor('success'),
      textColor: getColor('text-inverse'),
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

  const textVariant = size === 'small' ? 'label-xsmall' : 
                      size === 'large' ? 'label-large' :
                      size === 'xlarge' ? 'label-large' : 'label-medium';

  // Determine text color based on variant
  const textColor = variant === 'primary' || variant === 'danger' || variant === 'success' 
    ? 'text-inverse' 
    : variant === 'text' 
    ? 'primary' 
    : 'text-primary';

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
          style={{ fontWeight: '500' }}
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
      backgroundColor: color === 'currentColor' ? getColor('text-primary') : color,
      borderRadius: 2,
    }} />
  );
}; 