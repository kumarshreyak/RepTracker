import { Tabs } from 'expo-router';
import Ionicons from '@expo/vector-icons/Ionicons';
import { getColor } from '../../src/components/Colors';

export default function TabLayout() {
  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: getColor('accent'),
        headerStyle: {
          backgroundColor: getColor('backgroundPrimary'),
        },
        headerShadowVisible: true,
        headerTintColor: getColor('contentPrimary'),
        tabBarStyle: {
          backgroundColor: getColor('backgroundPrimary'),
          borderTopColor: getColor('borderOpaque'),
          borderTopWidth: 1,
        },
        tabBarInactiveTintColor: getColor('contentSecondary'),
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