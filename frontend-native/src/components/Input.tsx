import React, { useState, forwardRef } from 'react';
import { 
  TextInput, 
  View, 
  StyleSheet, 
  TextInputProps, 
  ViewStyle,
  TouchableOpacity,
  Text
} from 'react-native';
import { Typography } from './Typography';
import { SemanticColor, getColor } from './Colors';

export interface InputProps extends Omit<TextInputProps, 'style'> {
  label?: string;
  hint?: string;
  error?: string;
  success?: string;
  variant?: 'default' | 'error' | 'success' | 'disabled';
  disabled?: boolean;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  onRightIconPress?: () => void;
  containerStyle?: ViewStyle;
  inputStyle?: ViewStyle;
  clearable?: boolean;
  onClear?: () => void;
}

export const Input = forwardRef<TextInput, InputProps>(({
  label,
  hint,
  error,
  success,
  variant = 'default',
  disabled = false,
  leftIcon,
  rightIcon,
  onRightIconPress,
  containerStyle,
  inputStyle,
  clearable = false,
  onClear,
  value = '',
  onChangeText,
  ...textInputProps
}, ref) => {
  const [isFocused, setIsFocused] = useState(false);
  
  // Determine the actual variant based on props
  const actualVariant = disabled ? 'disabled' : error ? 'error' : success ? 'success' : variant;
  
  // Determine colors based on state
  const getBorderColor = (): SemanticColor => {
    if (actualVariant === 'disabled') return 'borderOpaque';
    if (actualVariant === 'error') return 'borderNegative';
    if (actualVariant === 'success') return 'borderPositive';
    if (isFocused) return 'borderAccent';
    return 'borderOpaque';
  };
  
  const getBackgroundColor = (): SemanticColor => {
    if (actualVariant === 'disabled') return 'backgroundStateDisabled';
    return 'backgroundPrimary';
  };
  
  const getTextColor = (): SemanticColor => {
    if (actualVariant === 'disabled') return 'contentStateDisabled';
    return 'contentPrimary';
  };
  
  const getHintColor = (): SemanticColor => {
    if (actualVariant === 'error') return 'contentNegative';
    if (actualVariant === 'success') return 'contentPositive';
    return 'contentSecondary';
  };
  
  const shouldShowClearButton = clearable && value.length > 0 && !disabled;
  
  const handleClear = () => {
    if (onClear) {
      onClear();
    } else if (onChangeText) {
      onChangeText('');
    }
  };
  
  return (
    <View style={[styles.container, containerStyle]}>
      {label && (
        <Typography 
          variant="label-small" 
          color={actualVariant === 'disabled' ? 'contentStateDisabled' : 'contentPrimary'}
          style={styles.label}
        >
          {label}
        </Typography>
      )}
      
      <View style={[
        styles.inputContainer,
        {
          backgroundColor: getColor(getBackgroundColor()),
          borderColor: getColor(getBorderColor()),
          borderWidth: isFocused ? 2 : 1,
        },
        actualVariant === 'disabled' && styles.disabledContainer,
      ]}>
        {leftIcon && (
          <View style={styles.leftIconContainer}>
            {leftIcon}
          </View>
        )}
        
        <TextInput
          ref={ref}
          style={[
            styles.input,
            {
              color: getColor(getTextColor()),
              flex: 1,
            },
            leftIcon ? styles.inputWithLeftIcon : null,
            (rightIcon || shouldShowClearButton) ? styles.inputWithRightIcon : null,
            inputStyle,
          ]}
          value={value}
          onChangeText={onChangeText}
          onFocus={(e) => {
            setIsFocused(true);
            textInputProps.onFocus?.(e);
          }}
          onBlur={(e) => {
            setIsFocused(false);
            textInputProps.onBlur?.(e);
          }}
          editable={!disabled}
          placeholderTextColor={getColor(actualVariant === 'disabled' ? 'contentStateDisabled' : 'contentTertiary')}
          {...textInputProps}
        />
        
        {shouldShowClearButton && (
          <TouchableOpacity 
            style={styles.clearButton}
            onPress={handleClear}
            hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
          >
            <View style={styles.clearIcon}>
              <Text style={styles.clearIconText}>×</Text>
            </View>
          </TouchableOpacity>
        )}
        
        {rightIcon && !shouldShowClearButton && (
          <TouchableOpacity 
            style={styles.rightIconContainer}
            onPress={onRightIconPress}
            disabled={!onRightIconPress || disabled}
            hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
          >
            {rightIcon}
          </TouchableOpacity>
        )}
      </View>
      
      {(hint || error || success) && (
        <Typography 
          variant="paragraph-xsmall" 
          color={getHintColor()}
          style={styles.hint}
        >
          {error || success || hint}
        </Typography>
      )}
    </View>
  );
});

Input.displayName = 'Input';

const styles = StyleSheet.create({
  container: {
    marginBottom: 4,
  },
  label: {
    marginBottom: 8,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 4,
    minHeight: 48,
    paddingHorizontal: 16,
  },
  disabledContainer: {
    opacity: 0.6,
  },
  input: {
    fontSize: 14,
    lineHeight: 20,
    paddingVertical: 8,
    flex: 1,
    textAlignVertical: 'center',
    includeFontPadding: false,
  },
  inputWithLeftIcon: {
    marginLeft: 8,
  },
  inputWithRightIcon: {
    marginRight: 8,
  },
  leftIconContainer: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  rightIconContainer: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  clearButton: {
    justifyContent: 'center',
    alignItems: 'center',
    width: 20,
    height: 20,
  },
  clearIcon: {
    width: 16,
    height: 16,
    borderRadius: 8,
    backgroundColor: '#A6A6A6',
    justifyContent: 'center',
    alignItems: 'center',
  },
  clearIconText: {
    color: '#FFFFFF',
    fontSize: 12,
    fontWeight: '600',
    lineHeight: 16,
    textAlign: 'center',
  },
  hint: {
    marginTop: 4,
  },
}); 