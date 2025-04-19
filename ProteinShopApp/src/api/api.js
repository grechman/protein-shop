import axios from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';

const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use(async (config) => {
  const token = await AsyncStorage.getItem('jwt_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const register = (email, password) =>
  api.post('/auth/register', { email, password });
export const login = (email, password) =>
  api.post('/auth/login', { email, password });
export const getProducts = () => api.get('/products');
export const createOrder = (items) => api.post('/orders', { items });
export const getOrders = () => api.get('/orders');
export const getLoyaltyPoints = () => api.get('/loyalty/points');
export const getProfile = () => api.get('/profile');

export default api;