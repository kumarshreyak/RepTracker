import { Tabs } from 'expo-router';
import Ionicons from '@expo/vector-icons/Ionicons';
import { getColor } from '../../src/components/Colors';

export default function TabLayout() {
  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: getColor('blue-bright'),
        headerStyle: {
          backgroundColor: getColor('white'),
        },
        headerShadowVisible: true,
        headerTintColor: getColor('dark'),
        tabBarStyle: {
          backgroundColor: getColor('white'),
          borderTopColor: getColor('light-gray-3'),
          borderTopWidth: 1,
        },
        tabBarInactiveTintColor: getColor('light'),
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: 'Home',
          headerShown: false,
          tabBarIcon: ({ color, focused }) => (
            <Ionicons 
              name={focused ? 'home-sharp' : 'home-outline'} 
              color={color} 
              size={24} 
            />
          ),
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: 'Profile',
          headerShown: false,
          tabBarIcon: ({ color, focused }) => (
            <Ionicons 
              name={focused ? 'person' : 'person-outline'} 
              color={color} 
              size={24} 
            />
          ),
        }}
      />
    </Tabs>
  );
} 