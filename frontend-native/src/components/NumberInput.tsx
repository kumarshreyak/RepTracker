import React, { useState, useRef, useEffect } from 'react';
import {
  View,
  StyleSheet,
  TextInput,
  Pressable,
  ViewStyle,
  TextStyle,
} from 'react-native';
import { Typography } from './Typography';
import { SemanticColor, getColor } from './Colors';

export interface NumberInputProps {
  /** The current numeric value (will default to 0 if undefined) */
  value: number | undefined;
  /** Callback fired when the value changes */
  onValueChange: (value: number) => void;
  /** Optional unit label displayed below the value (e.g., "kg", "reps", "lbs") */
  unit?: string;
  /** Typography variant for the number display */
  displayVariant?: 'display-large' | 'display-medium' | 'display-small' | 'display-xsmall' | 'heading-xxlarge' | 'heading-xlarge' | 'heading-large' | 'heading-medium' | 'heading-small' | 'heading-xsmall' | 'label-large' | 'label-medium' | 'label-small' | 'label-xsmall' | 'paragraph-large' | 'paragraph-medium' | 'paragraph-small' | 'paragraph-xsmall';
  /** Typography variant for the unit label */
  unitVariant?: 'display-large' | 'display-medium' | 'display-small' | 'display-xsmall' | 'heading-xxlarge' | 'heading-xlarge' | 'heading-large' | 'heading-medium' | 'heading-small' | 'heading-xsmall' | 'label-large' | 'label-medium' | 'label-small' | 'label-xsmall' | 'paragraph-large' | 'paragraph-medium' | 'paragraph-small' | 'paragraph-xsmall';
  /** Semantic color for the number value */
  valueColor?: SemanticColor;
  /** Semantic color for the unit label */
  unitColor?: SemanticColor;
  /** Custom styles for the container */
  containerStyle?: ViewStyle;
  /** Custom styles for the value display/input */
  valueStyle?: TextStyle;
  /** Custom styles for the unit label */
  unitStyle?: TextStyle;
  /** Whether the input is disabled */
  disabled?: boolean;
  /** Minimum allowed value */
  minValue?: number;
  /** Maximum allowed value */
  maxValue?: number;
  /** Step increment for value changes */
  step?: number;
  /** Allow decimal values */
  allowDecimal?: boolean;
  /** Placeholder text when editing */
  placeholder?: string;
  /** Auto-focus when editing starts */
  autoFocus?: boolean;
  /** Select all text when editing starts */
  selectTextOnFocus?: boolean;
}

/**
 * NumberInput Component
 * 
 * A tap-to-edit number input component that displays numbers in a large, readable format
 * and allows inline editing with a single tap. Perfect for workout apps where users need
 * to quickly adjust weights, reps, or other numeric values.
 * 
 * @example
 * ```tsx
 * // Basic usage for workout reps
 * <NumberInput
 *   value={12}
 *   onValueChange={(value) => setReps(value)}
 *   unit="reps"
 *   displayVariant="heading-large"
 *   minValue={1}
 *   maxValue={999}
 * />
 * 
 * // Weight input with decimals
 * <NumberInput
 *   value={135.5}
 *   onValueChange={(value) => setWeight(value)}
 *   unit="lbs"
 *   displayVariant="display-xsmall"
 *   allowDecimal={true}
 *   minValue={0}
 *   maxValue={9999}
 * />
 * 
 * // Time duration in seconds
 * <NumberInput
 *   value={90}
 *   onValueChange={(value) => setRestTime(value)}
 *   unit="sec"
 *   displayVariant="heading-medium"
 *   step={5}
 *   minValue={0}
 *   maxValue={600}
 * />
 * ```
 */
export const NumberInput: React.FC<NumberInputProps> = ({
  value,
  onValueChange,
  unit,
  displayVariant = 'display-xsmall',
  unitVariant = 'label-small',
  valueColor = 'contentPrimary',
  unitColor = 'contentTertiary',
  containerStyle,
  valueStyle,
  unitStyle,
  disabled = false,
  minValue,
  maxValue,
  step = 1,
  allowDecimal = false,
  placeholder = '0',
  autoFocus = true,
  selectTextOnFocus = true,
}) => {
  // Ensure value is never undefined
  const safeValue = value ?? 0;
  
  const [isEditing, setIsEditing] = useState(false);
  const [inputValue, setInputValue] = useState(safeValue.toString());
  const inputRef = useRef<TextInput>(null);

  // Update input value when prop value changes
  useEffect(() => {
    if (!isEditing) {
      setInputValue(safeValue.toString());
    }
  }, [safeValue, isEditing]);

  const handlePress = () => {
    if (disabled) return;
    setIsEditing(true);
    setInputValue(safeValue.toString());
  };

  const handleSubmit = () => {
    setIsEditing(false);
    
    let newValue: number;
    
    if (allowDecimal) {
      newValue = parseFloat(inputValue) || 0;
    } else {
      newValue = parseInt(inputValue) || 0;
    }

    // Apply constraints
    if (minValue !== undefined && newValue < minValue) {
      newValue = minValue;
    }
    if (maxValue !== undefined && newValue > maxValue) {
      newValue = maxValue;
    }

    // Apply step rounding if needed
    if (step !== 1 && !allowDecimal) {
      newValue = Math.round(newValue / step) * step;
    }

    onValueChange(newValue);
  };

  const handleTextChange = (text: string) => {
    // Filter input based on allowDecimal
    if (allowDecimal) {
      // Allow numbers and single decimal point
      const filtered = text.replace(/[^0-9.]/g, '');
      const parts = filtered.split('.');
      if (parts.length > 2) {
        // Only allow one decimal point
        setInputValue(parts[0] + '.' + parts.slice(1).join(''));
      } else {
        setInputValue(filtered);
      }
    } else {
      // Only allow integers
      const filtered = text.replace(/[^0-9]/g, '');
      setInputValue(filtered);
    }
  };

  const handleBlur = () => {
    handleSubmit();
  };

  // Auto-focus when editing starts
  useEffect(() => {
    if (isEditing && autoFocus) {
      setTimeout(() => {
        inputRef.current?.focus();
        if (selectTextOnFocus) {
          inputRef.current?.setSelection(0, inputValue.length);
        }
      }, 100);
    }
  }, [isEditing, autoFocus, selectTextOnFocus, inputValue.length]);

  const getInputStyle = (): TextStyle => {
    // Base styles that match the display variant
    const baseStyle: TextStyle = {
      color: getColor(valueColor),
      textAlign: 'center',
      borderBottomWidth: 2,
      borderBottomColor: getColor('accent'),
      paddingVertical: 4,
      minWidth: 80,
      backgroundColor: 'transparent',
    };

    // Map display variants to font sizes (approximating Typography component styles)
    switch (displayVariant) {
      case 'display-large':
        return { ...baseStyle, fontSize: 96, fontWeight: '700' };
      case 'display-medium':
        return { ...baseStyle, fontSize: 52, fontWeight: '700' };
      case 'display-small':
        return { ...baseStyle, fontSize: 44, fontWeight: '700' };
      case 'display-xsmall':
        return { ...baseStyle, fontSize: 36, fontWeight: '700' };
      case 'heading-xxlarge':
        return { ...baseStyle, fontSize: 40, fontWeight: '700' };
      case 'heading-xlarge':
        return { ...baseStyle, fontSize: 36, fontWeight: '700' };
      case 'heading-large':
        return { ...baseStyle, fontSize: 32, fontWeight: '700' };
      case 'heading-medium':
        return { ...baseStyle, fontSize: 28, fontWeight: '700' };
      case 'heading-small':
        return { ...baseStyle, fontSize: 24, fontWeight: '700' };
      case 'heading-xsmall':
        return { ...baseStyle, fontSize: 20, fontWeight: '700' };
      case 'label-large':
        return { ...baseStyle, fontSize: 18, fontWeight: '500' };
      case 'label-medium':
        return { ...baseStyle, fontSize: 16, fontWeight: '500' };
      case 'label-small':
        return { ...baseStyle, fontSize: 14, fontWeight: '500' };
      case 'label-xsmall':
        return { ...baseStyle, fontSize: 12, fontWeight: '500' };
      case 'paragraph-large':
        return { ...baseStyle, fontSize: 18, fontWeight: '400' };
      case 'paragraph-medium':
        return { ...baseStyle, fontSize: 16, fontWeight: '400' };
      case 'paragraph-small':
        return { ...baseStyle, fontSize: 14, fontWeight: '400' };
      case 'paragraph-xsmall':
        return { ...baseStyle, fontSize: 12, fontWeight: '400' };
      default:
        return { ...baseStyle, fontSize: 36, fontWeight: '700' };
    }
  };

  return (
    <View style={[styles.container, containerStyle]}>
      <Pressable 
        style={styles.valueContainer}
        onPress={handlePress}
        disabled={disabled}
      >
        {isEditing ? (
          <TextInput
            ref={inputRef}
            style={[getInputStyle(), valueStyle]}
            value={inputValue}
            onChangeText={handleTextChange}
            onBlur={handleBlur}
            onSubmitEditing={handleSubmit}
            keyboardType={allowDecimal ? 'decimal-pad' : 'numeric'}
            returnKeyType="done"
            placeholder={placeholder}
            placeholderTextColor={getColor('contentTertiary')}
            selectTextOnFocus={selectTextOnFocus}
          />
        ) : (
          <Typography 
            variant={displayVariant} 
            color={disabled ? 'contentStateDisabled' : valueColor}
            style={valueStyle}
          >
            {safeValue.toString()}
          </Typography>
        )}
      </Pressable>
      
      {unit && (
        <Typography 
          variant={unitVariant} 
          color={disabled ? 'contentStateDisabled' : unitColor}
          style={{ ...styles.unit, ...unitStyle }}
        >
          {unit}
        </Typography>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
  },
  valueContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 44, // Ensure good touch target
    paddingHorizontal: 8,
  },
  unit: {
    marginTop: 4,
  },
}); 