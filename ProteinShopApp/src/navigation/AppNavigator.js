import React, { useContext } from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createBottomTabsNavigator } from '@react-navigation/bottom-tabs';
import { createStackNavigator } from '@react-navigation/stack';
import LoginScreen from '../screens/LoginScreen';
import RegisterScreen from '../screens/RegisterScreen';
import CatalogScreen from '../screens/CatalogScreen';
import PointsScreen from '../screens/PointsScreen';
import ProfileScreen from '../screens/ProfileScreen';
import OrderPlacementScreen from '../screens/OrderPlacementScreen';
import { AuthContext } from '../utils/AuthContext';

const Tab = createBottomTabsNavigator();
const Stack = createStackNavigator();

const MainTabs = () => (
  <Tab.Navigator>
    <Tab.Screen name="Home" component={CatalogScreen} options={{ headerShown: false }} />
    <Tab.Screen name="Points" component={PointsScreen} options={{ headerShown: false }} />
    <Tab.Screen name="Profile" component={ProfileScreen} options={{ headerShown: false }} />
  </Tab.Navigator>
);

const AppNavigator = () => {
  const { isAuthenticated } = useContext(AuthContext);

  if (isAuthenticated === null) {
    return null; // Loading state
  }

  return (
    <NavigationContainer>
      <Stack.Navigator initialRouteName={isAuthenticated ? 'Main' : 'Login'}>
        {!isAuthenticated ? (
          <>
            <Stack.Screen name="Login" component={LoginScreen} options={{ headerShown: false }} />
            <Stack.Screen name="Register" component={RegisterScreen} options={{ headerShown: false }} />
          </>
        ) : (
          <>
            <Stack.Screen name="Main" component={MainTabs} options={{ headerShown: false }} />
            <Stack.Screen name="OrderPlacement" component={OrderPlacementScreen} options={{ title: 'Checkout' }} />
          </>
        )}
      </Stack.Navigator>
    </NavigationContainer>
  );
};

export default AppNavigator;