import React from 'react';
import { Text, TextStyle } from 'react-native';
import { AppColor, getColor } from './Colors';

interface TypographyProps {
  children: React.ReactNode;
  variant?: 
    | 'text-small' 
    | 'text-default' 
    | 'text-large' 
    | 'text-xlarge'
    | 'text-small-paragraph'
    | 'text-default-paragraph'
    | 'text-large-paragraph'
    | 'text-xlarge-paragraph'
    | 'heading-xsmall'
    | 'heading-small'
    | 'heading-default'
    | 'heading-large'
    | 'heading-xlarge'
    | 'heading-xxlarge'
    | 'heading-xsmall-caps'
    | 'heading-small-caps'
    | 'heading-default-caps';
  color?: AppColor;
  style?: TextStyle;
}

export const Typography: React.FC<TypographyProps> = ({ 
  children, 
  variant = 'text-default', 
  color = 'dark',
  style = {}
}) => {
  const variantStyles: Record<string, TextStyle> = {
    // Text variants
    'text-small': {
      fontSize: 11,
      lineHeight: 16,
      fontWeight: '400',
    },
    'text-default': {
      fontSize: 13,
      lineHeight: 18,
      fontWeight: '400',
    },
    'text-large': {
      fontSize: 15,
      lineHeight: 20,
      fontWeight: '400',
    },
    'text-xlarge': {
      fontSize: 17,
      lineHeight: 24,
      fontWeight: '400',
    },
    
    // Paragraph variants
    'text-small-paragraph': {
      fontSize: 11,
      lineHeight: 16,
      fontWeight: '400',
    },
    'text-default-paragraph': {
      fontSize: 13,
      lineHeight: 18,
      fontWeight: '400',
    },
    'text-large-paragraph': {
      fontSize: 15,
      lineHeight: 20,
      fontWeight: '400',
    },
    'text-xlarge-paragraph': {
      fontSize: 17,
      lineHeight: 24,
      fontWeight: '400',
    },
    
    // Heading variants
    'heading-xsmall': {
      fontSize: 11,
      lineHeight: 16,
      fontWeight: '600',
    },
    'heading-small': {
      fontSize: 13,
      lineHeight: 18,
      fontWeight: '600',
    },
    'heading-default': {
      fontSize: 15,
      lineHeight: 20,
      fontWeight: '600',
    },
    'heading-large': {
      fontSize: 17,
      lineHeight: 24,
      fontWeight: '600',
    },
    'heading-xlarge': {
      fontSize: 21,
      lineHeight: 28,
      fontWeight: '700',
    },
    'heading-xxlarge': {
      fontSize: 28,
      lineHeight: 36,
      fontWeight: '700',
    },
    
    // All caps heading variants
    'heading-xsmall-caps': {
      fontSize: 11,
      lineHeight: 16,
      fontWeight: '600',
      textTransform: 'uppercase',
      letterSpacing: 0.5,
    },
    'heading-small-caps': {
      fontSize: 13,
      lineHeight: 18,
      fontWeight: '600',
      textTransform: 'uppercase',
      letterSpacing: 0.5,
    },
    'heading-default-caps': {
      fontSize: 15,
      lineHeight: 20,
      fontWeight: '600',
      textTransform: 'uppercase',
      letterSpacing: 0.5,
    },
  };

  const combinedStyle: TextStyle = {
    ...variantStyles[variant],
    color: getColor(color),
    ...style,
  };

  return (
    <Text style={combinedStyle}>
      {children}
    </Text>
  );
}; 