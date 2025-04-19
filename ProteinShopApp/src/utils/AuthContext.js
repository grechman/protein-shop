import React, { createContext, useState, useEffect } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { login, register } from '../api/api';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(null);
  const [user, setUser] = useState(null);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const token = await AsyncStorage.getItem('jwt_token');
        setIsAuthenticated(!!token);
      } catch (error) {
        console.error('Error checking auth:', error);
        setIsAuthenticated(false);
      }
    };
    checkAuth();
  }, []);

  const handleLogin = async (email, password) => {
    try {
      const response = await login(email, password);
      const { token } = response.data;
      await AsyncStorage.setItem('jwt_token', token);
      setIsAuthenticated(true);
      return { success: true };
    } catch (error) {
      console.error('Login error:', error);
      return { success: false, error: error.response?.data?.error || 'Login failed' };
    }
  };

  const handleRegister = async (email, password) => {
    try {
      const response = await register(email, password);
      const { user_id } = response.data;
      const loginResponse = await login(email, password);
      const { token } = loginResponse.data;
      await AsyncStorage.setItem('jwt_token', token);
      setIsAuthenticated(true);
      setUser({ id: user_id, email });
      return { success: true };
    } catch (error) {
      console.error('Register error:', error);
      return { success: false, error: error.response?.data?.error || 'Registration failed' };
    }
  };

  const logout = async () => {
    await AsyncStorage.removeItem('jwt_token');
    setIsAuthenticated(false);
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, handleLogin, handleRegister, logout }}>
      {children}
    </AuthContext.Provider>
  );
};