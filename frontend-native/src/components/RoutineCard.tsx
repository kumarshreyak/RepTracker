import React from 'react';
import {
  View,
  StyleSheet,
  Pressable,
  Alert,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { Typography } from './Typography';
import { getColor } from './Colors';

interface Exercise {
  id: string;
  name: string;
  muscleGroup: string;
  equipmentType: string;
}

interface Routine {
  id: string;
  userId: string;
  name: string;
  description: string;
  exercises: Exercise[];
  createdAt: string;
  updatedAt: string;
}

interface RoutineCardProps {
  routine: Routine;
  onStart: (routineId: string) => void;
  onEdit: (routineId: string) => void;
  onDelete: (routineId: string) => void;
}

export const RoutineCard: React.FC<RoutineCardProps> = ({
  routine,
  onStart,
  onEdit,
  onDelete,
}) => {
  const handleDelete = () => {
    Alert.alert(
      'Delete Routine',
      `Are you sure you want to delete "${routine.name}"? This action cannot be undone.`,
      [
        {
          text: 'Cancel',
          style: 'cancel',
        },
        {
          text: 'Delete',
          style: 'destructive',
          onPress: () => onDelete(routine.id),
        },
      ]
    );
  };

  const handleStart = () => {
    onStart(routine.id);
  };

  const handleEdit = () => {
    onEdit(routine.id);
  };

  const totalExercises = routine.exercises?.length || 0;

  return (
    <View style={styles.routineCard}>
      {/* Header with title and exercise count */}
      <View style={styles.header}>
        <View style={styles.titleSection}>
          <Typography variant="label-medium" color="contentPrimary">
            {routine.name}
          </Typography>
          <Typography variant="label-small" color="contentSecondary">
            {totalExercises} exercise{totalExercises !== 1 ? 's' : ''}
          </Typography>
        </View>

        {/* Action buttons */}
        <View style={styles.actionButtons}>
          <Pressable 
            style={styles.iconButton}
            onPress={handleEdit}
          >
            <MaterialIcons 
              name="edit" 
              size={20} 
              color={getColor('contentSecondary')} 
            />
          </Pressable>
          
          <Pressable 
            style={styles.iconButton}
            onPress={handleDelete}
          >
            <MaterialIcons 
              name="delete" 
              size={20} 
              color={getColor('contentNegative')} 
            />
          </Pressable>

          <Pressable 
            style={styles.playButton}
            onPress={handleStart}
          >
            <MaterialIcons 
              name="play-arrow" 
              size={16} 
              color={getColor('contentOnColor')} 
            />
          </Pressable>
        </View>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  routineCard: {
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    padding: 16,
    marginBottom: 12,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  titleSection: {
    flex: 1,
    marginRight: 12,
    gap: 4,
  },
  actionButtons: {
    flexDirection: 'row',
    gap: 8,
  },
  iconButton: {
    padding: 8,
    minWidth: 32,
    alignItems: 'center',
  },
  playButton: {
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: getColor('backgroundAccent'),
    justifyContent: 'center',
    alignItems: 'center',
  },
});