import React from 'react';
import { View, StyleSheet, ScrollView, SafeAreaView } from 'react-native';
import { Typography } from './Typography';
import { getColor, coreColors, semanticColors, semanticExtensions } from './Colors';

export default function ColorDemo() {
  const renderColorSwatch = (name: string, colorValue: string) => (
    <View key={name} style={styles.swatchContainer}>
      <View style={[styles.swatch, { backgroundColor: colorValue }]} />
      <Typography variant="paragraph-small" color="contentPrimary" style={styles.colorName}>
        {name}
      </Typography>
      <Typography variant="paragraph-xsmall" color="contentSecondary" style={styles.colorValue}>
        {colorValue}
      </Typography>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        <View style={styles.header}>
          <Typography variant="heading-large" color="contentPrimary" style={styles.title}>
            Uber Base Design System
          </Typography>
          <Typography variant="paragraph-medium" color="contentSecondary" style={styles.subtitle}>
            Complete Semantic Color Palette
          </Typography>
        </View>

        {/* Core Colors */}
        <View style={styles.section}>
          <Typography variant="heading-medium" color="contentPrimary" style={styles.sectionTitle}>
            02 Core Colors
          </Typography>
          <Typography variant="paragraph-small" color="contentSecondary" style={styles.sectionDescription}>
            Foundation colors that drive the entire system
          </Typography>
          <View style={styles.swatchGrid}>
            {Object.entries(coreColors).map(([name, value]) => 
              renderColorSwatch(name, value)
            )}
          </View>
        </View>

        {/* Basic Semantic Colors */}
        <View style={styles.section}>
          <Typography variant="heading-medium" color="contentPrimary" style={styles.sectionTitle}>
            03 Basic Semantic Colors
          </Typography>
          <Typography variant="paragraph-small" color="contentSecondary" style={styles.sectionDescription}>
            Primary semantic colors organized by purpose
          </Typography>
          
          {/* Background Colors */}
          <View style={styles.subsection}>
            <Typography variant="label-medium" color="contentPrimary" style={styles.subsectionTitle}>
              Background
            </Typography>
            <View style={styles.swatchGrid}>
              {Object.entries(semanticColors)
                .filter(([name]) => name.startsWith('background'))
                .map(([name, value]) => renderColorSwatch(name, value))}
            </View>
          </View>

          {/* Content Colors */}
          <View style={styles.subsection}>
            <Typography variant="label-medium" color="contentPrimary" style={styles.subsectionTitle}>
              Content
            </Typography>
            <View style={styles.swatchGrid}>
              {Object.entries(semanticColors)
                .filter(([name]) => name.startsWith('content'))
                .map(([name, value]) => renderColorSwatch(name, value))}
            </View>
          </View>

          {/* Border Colors */}
          <View style={styles.subsection}>
            <Typography variant="label-medium" color="contentPrimary" style={styles.subsectionTitle}>
              Border
            </Typography>
            <View style={styles.swatchGrid}>
              {Object.entries(semanticColors)
                .filter(([name]) => name.startsWith('border'))
                .map(([name, value]) => renderColorSwatch(name, value))}
            </View>
          </View>
        </View>

        {/* Semantic Extensions */}
        <View style={styles.section}>
          <Typography variant="heading-medium" color="contentPrimary" style={styles.sectionTitle}>
            03 Semantic Extensions
          </Typography>
          <Typography variant="paragraph-small" color="contentSecondary" style={styles.sectionDescription}>
            Extended semantic colors for specific use cases and states
          </Typography>
          
          {/* Background Extensions */}
          <View style={styles.subsection}>
            <Typography variant="label-medium" color="contentPrimary" style={styles.subsectionTitle}>
              Background Extensions
            </Typography>
            <View style={styles.swatchGrid}>
              {Object.entries(semanticExtensions)
                .filter(([name]) => name.startsWith('background'))
                .map(([name, value]) => renderColorSwatch(name, value))}
            </View>
          </View>

          {/* Content Extensions */}
          <View style={styles.subsection}>
            <Typography variant="label-medium" color="contentPrimary" style={styles.subsectionTitle}>
              Content Extensions
            </Typography>
            <View style={styles.swatchGrid}>
              {Object.entries(semanticExtensions)
                .filter(([name]) => name.startsWith('content'))
                .map(([name, value]) => renderColorSwatch(name, value))}
            </View>
          </View>

          {/* Border Extensions */}
          <View style={styles.subsection}>
            <Typography variant="label-medium" color="contentPrimary" style={styles.subsectionTitle}>
              Border Extensions
            </Typography>
            <View style={styles.swatchGrid}>
              {Object.entries(semanticExtensions)
                .filter(([name]) => name.startsWith('border'))
                .map(([name, value]) => renderColorSwatch(name, value))}
            </View>
          </View>
        </View>

        {/* Usage Examples */}
        <View style={styles.section}>
          <Typography variant="heading-medium" color="contentPrimary" style={styles.sectionTitle}>
            Usage Examples
          </Typography>
          
          <View style={styles.exampleCard}>
            <Typography variant="label-medium" color="contentPrimary" style={styles.exampleTitle}>
              Primary Content Card
            </Typography>
            <Typography variant="paragraph-small" color="contentSecondary">
              Uses backgroundPrimary with borderOpaque
            </Typography>
          </View>
          
          <View style={styles.exampleCardSecondary}>
            <Typography variant="label-medium" color="contentPrimary" style={styles.exampleTitle}>
              Secondary Background
            </Typography>
            <Typography variant="paragraph-small" color="contentSecondary">
              Uses backgroundSecondary for page background
            </Typography>
          </View>
          
          <View style={styles.exampleCardAccent}>
            <Typography variant="label-medium" color="contentOnColor" style={styles.exampleTitle}>
              Accent Action
            </Typography>
            <Typography variant="paragraph-small" color="contentOnColor">
              Uses backgroundAccent with contentOnColor
            </Typography>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
  },
  scrollView: {
    flex: 1,
  },
  header: {
    paddingHorizontal: 16,
    paddingTop: 16,
    paddingBottom: 24,
  },
  title: {
    marginBottom: 8,
  },
  subtitle: {
    marginBottom: 0,
  },
  section: {
    paddingHorizontal: 16,
    marginBottom: 32,
  },
  sectionTitle: {
    marginBottom: 8,
  },
  sectionDescription: {
    marginBottom: 16,
  },
  subsection: {
    marginBottom: 24,
  },
  subsectionTitle: {
    marginBottom: 12,
  },
  swatchGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 12,
  },
  swatchContainer: {
    alignItems: 'center',
    minWidth: 100,
    marginBottom: 8,
  },
  swatch: {
    width: 60,
    height: 60,
    borderRadius: 8,
    marginBottom: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
  },
  colorName: {
    textAlign: 'center',
    fontWeight: '500',
    marginBottom: 2,
  },
  colorValue: {
    textAlign: 'center',
    fontFamily: 'monospace',
  },
  
  // Example cards
  exampleCard: {
    backgroundColor: getColor('backgroundPrimary'),
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    borderRadius: 8,
    padding: 16,
    marginBottom: 12,
  },
  exampleCardSecondary: {
    backgroundColor: getColor('backgroundSecondary'),
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    borderRadius: 8,
    padding: 16,
    marginBottom: 12,
  },
  exampleCardAccent: {
    backgroundColor: getColor('backgroundAccent'),
    borderRadius: 8,
    padding: 16,
    marginBottom: 12,
  },
  exampleTitle: {
    marginBottom: 4,
  },
}); 