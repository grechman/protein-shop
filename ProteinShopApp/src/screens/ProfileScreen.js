import React, { useState, useEffect, useContext } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Alert } from 'react-native';
import { getProfile } from '../api/api';
import { AuthContext } from '../utils/AuthContext';

const ProfileScreen = () => {
  const { logout } = useContext(AuthContext);
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await getProfile();
        setProfile(response.data);
      } catch (error) {
        Alert.alert('Error', 'Failed to load profile');
      } finally {
        setLoading(false);
      }
    };
    fetchProfile();
  }, []);

  if (loading) {
    return <Text style={styles.loading}>Loading...</Text>;
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Profile</Text>
      {profile && (
        <View style={styles.infoSection}>
          <Text>Email: {profile.email}</Text>
          <Text>Member Since: {new Date(profile.created_at).toLocaleDateString()}</Text>
        </View>
      )}
      <TouchableOpacity style={styles.button} onPress={logout}>
        <Text style={styles.buttonText}>Log Out</Text>
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, padding: 20 },
  title: { fontSize: 24, fontWeight: 'bold', marginBottom: 20 },
  infoSection: { padding: 10, backgroundColor: '#f0f0f0', borderRadius: 5, marginBottom: 20 },
  button: { backgroundColor: '#FF3B30', padding: 15, borderRadius: 5, alignItems: 'center' },
  buttonText: { color: '#fff', fontSize: 16 },
  loading: { flex: 1, textAlign: 'center', marginTop: 20, fontSize: 18 },
});

export default ProfileScreen;