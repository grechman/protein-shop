import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, StyleSheet, Alert } from 'react-native';
import { getOrders, getLoyaltyPoints } from '../api/api';

const PointsScreen = () => {
  const [points, setPoints] = useState({ balance: 0, history: [] });
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [pointsResponse, ordersResponse] = await Promise.all([
          getLoyaltyPoints(),
          getOrders(),
        ]);
        setPoints(pointsResponse.data);
        setOrders(ordersResponse.data);
      } catch (error) {
        Alert.alert('Error', 'Failed to load data');
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  if (loading) {
    return <Text style={styles.loading}>Loading...</Text>;
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>My Rewards</Text>
      <View style={styles.pointsSection}>
        <Text style={styles.pointsText}>Balance: {points.balance} points</Text>
      </View>
      <Text style={styles.subtitle}>Order History</Text>
      <FlatList
        data={orders}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <View style={styles.orderCard}>
            <Text>Order #{item.id.slice(0, 8)}</Text>
            <Text>Date: {new Date(item.created_at).toLocaleDateString()}</Text>
            <Text>Total: ${item.total.toFixed(2)}</Text>
            <Text>Status: {item.status}</Text>
          </View>
        )}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, padding: 10 },
  title: { fontSize: 24, fontWeight: 'bold', marginBottom: 10 },
  pointsSection: { padding: 10, backgroundColor: '#f0f0f0', borderRadius: 5 },
  pointsText: { fontSize: 18, fontWeight: 'bold' },
  subtitle: { fontSize: 20, fontWeight: 'bold', marginVertical: 10 },
  orderCard: { padding: 10, borderBottomWidth: 1, borderBottomColor: '#ccc' },
  loading: { flex: 1, textAlign: 'center', marginTop: 20, fontSize: 18 },
});

export default PointsScreen;