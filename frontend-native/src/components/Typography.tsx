import React from 'react';
import { Text, TextStyle } from 'react-native';
import { SemanticColor, getColor } from './Colors';

// New Uber Base Design System Typography Variants
export type TypographyVariant = 
  // Paragraph variants (Uber Move Text 400 - Regular)
  | 'paragraph-large'
  | 'paragraph-medium'
  | 'paragraph-small'
  | 'paragraph-xsmall'
  // Label variants (Uber Move Text 500 - Medium)
  | 'label-large'
  | 'label-medium'
  | 'label-small'
  | 'label-xsmall'
  // Heading variants (Uber Move 700 - Bold)
  | 'heading-xxlarge'
  | 'heading-xlarge'
  | 'heading-large'
  | 'heading-medium'
  | 'heading-small'
  | 'heading-xsmall'
  // Display variants (Uber Move 700 - Bold)
  | 'display-large'
  | 'display-medium'
  | 'display-small'
  | 'display-xsmall';

interface TypographyProps {
  children: React.ReactNode;
  variant?: TypographyVariant;
  color?: SemanticColor;
  style?: TextStyle;
}

export const Typography: React.FC<TypographyProps> = ({ 
  children, 
  variant = 'paragraph-medium', 
  color = 'contentPrimary',
  style = {}
}) => {
  const variantStyles: Record<TypographyVariant, TextStyle> = {
    // Paragraph variants - Uber Move Text Regular (400)
    'paragraph-large': {
      fontFamily: 'Uber Move Text',
      fontSize: 18,
      lineHeight: 28,
      fontWeight: '400',
    },
    'paragraph-medium': {
      fontFamily: 'Uber Move Text',
      fontSize: 16,
      lineHeight: 24,
      fontWeight: '400',
    },
    'paragraph-small': {
      fontFamily: 'Uber Move Text',
      fontSize: 14,
      lineHeight: 20,
      fontWeight: '400',
    },
    'paragraph-xsmall': {
      fontFamily: 'Uber Move Text',
      fontSize: 12,
      lineHeight: 20,
      fontWeight: '400',
    },
    
    // Label variants - Uber Move Text Medium (500)
    'label-large': {
      fontFamily: 'Uber Move Text',
      fontSize: 18,
      lineHeight: 24,
      fontWeight: '500',
    },
    'label-medium': {
      fontFamily: 'Uber Move Text',
      fontSize: 16,
      lineHeight: 20,
      fontWeight: '500',
    },
    'label-small': {
      fontFamily: 'Uber Move Text',
      fontSize: 14,
      lineHeight: 16,
      fontWeight: '500',
    },
    'label-xsmall': {
      fontFamily: 'Uber Move Text',
      fontSize: 12,
      lineHeight: 16,
      fontWeight: '500',
    },
    
    // Heading variants - Uber Move Bold (700)
    'heading-xxlarge': {
      fontFamily: 'Uber Move',
      fontSize: 40,
      lineHeight: 52,
      fontWeight: '700',
    },
    'heading-xlarge': {
      fontFamily: 'Uber Move',
      fontSize: 36,
      lineHeight: 44,
      fontWeight: '700',
    },
    'heading-large': {
      fontFamily: 'Uber Move',
      fontSize: 32,
      lineHeight: 40,
      fontWeight: '700',
    },
    'heading-medium': {
      fontFamily: 'Uber Move',
      fontSize: 28,
      lineHeight: 36,
      fontWeight: '700',
    },
    'heading-small': {
      fontFamily: 'Uber Move',
      fontSize: 24,
      lineHeight: 32,
      fontWeight: '700',
    },
    'heading-xsmall': {
      fontFamily: 'Uber Move',
      fontSize: 20,
      lineHeight: 28,
      fontWeight: '700',
    },
    
    // Display variants - Uber Move Bold (700)
    'display-large': {
      fontFamily: 'Uber Move',
      fontSize: 96,
      lineHeight: 112,
      fontWeight: '700',
    },
    'display-medium': {
      fontFamily: 'Uber Move',
      fontSize: 52,
      lineHeight: 64,
      fontWeight: '700',
    },
    'display-small': {
      fontFamily: 'Uber Move',
      fontSize: 44,
      lineHeight: 52,
      fontWeight: '700',
    },
    'display-xsmall': {
      fontFamily: 'Uber Move',
      fontSize: 36,
      lineHeight: 44,
      fontWeight: '700',
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