import React from 'react';
import { render, fireEvent } from '@testing-library/react-native';
import LoginScreen from './LoginScreen';
import { AuthContext } from '../utils/AuthContext';

const mockLogin = jest.fn();
const mockNavigation = { navigate: jest.fn() };

test('отображает экран входа и вызывает функцию входа', () => {
  const { getByPlaceholderText, getByText } = render(
    <AuthContext.Provider value={{ handleLogin: mockLogin }}>
      <LoginScreen navigation={mockNavigation} />
    </AuthContext.Provider>
  );

  const emailInput = getByPlaceholderText('user@example.com');
  const passwordInput = getByPlaceholderText('Password');
  const loginButton = getByText('Log In');

  fireEvent.changeText(emailInput, 'test@example.com');
  fireEvent.changeText(passwordInput, 'StrongPass123!');
  fireEvent.press(loginButton);

  expect(mockLogin).toHaveBeenCalledWith('test@example.com', 'StrongPass123!');
});