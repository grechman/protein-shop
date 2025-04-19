import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, Alert } from 'react-native';
import { createOrder, getProducts } from '../api/api';

const OrderPlacementScreen = ({ route, navigation }) => {
  const { cartItems } = route.params;
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(false);
  const [fetching, setFetching] = useState(true);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const response = await getProducts();
        setProducts(response.data);
      } catch (error) {
        Alert.alert('Ошибка', 'Не удалось загрузить товары');
      } finally {
        setFetching(false);
      }
    };
    fetchProducts();
  }, []);

  const handlePlaceOrder = async () => {
    if (cartItems.length === 0) {
      Alert.alert('Ошибка', 'Корзина пуста');
      return;
    }
    setLoading(true);
    try {
      const response = await createOrder(cartItems);
      Alert.alert('Успех', 'Заказ успешно оформлен!');
      navigation.navigate('Home');
    } catch (error) {
      Alert.alert('Ошибка', error.response?.data?.error || 'Не удалось оформить заказ');
    } finally {
      setLoading(false);
    }
  };

  const total = cartItems.reduce((sum, item) => {
    const product = products.find((p) => p.id === item.product_id);
    return sum + (product?.price || 0) * item.quantity;
  }, 0);

  if (fetching) {
    return <Text style={styles.loading}>Загрузка...</Text>;
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Оформление заказа</Text>
      <FlatList
        data={cartItems}
        keyExtractor={(item, index) => index.toString()}
        renderItem={({ item }) => {
          const product = products.find((p) => p.id === item.product_id);
          return (
            <View style={styles.cartItem}>
              <Text>{product?.name || 'Неизвестно'}</Text>
              <Text>Кол-во: {item.quantity}</Text>
              <Text>${((product?.price || 0) * item.quantity).toFixed(2)}</Text>
            </View>
          );
        }}
      />
      <View style={styles.summary}>
        <Text style={styles.summaryText}>Итого: ${total.toFixed(2)}</Text>
        <Text style={styles.summaryText}>Заработано баллов: {Math.floor(total / 10)}</Text>
      </View>
      <TouchableOpacity style={styles.button} onPress={handlePlaceOrder} disabled={loading}>
        <Text style={styles.buttonText}>{loading ? 'Оформление...' : 'Оформить заказ'}</Text>
      </TouchableOpacity>
      <TouchableOpacity onPress={() => navigation.navigate('Home')}>
        <Text style={styles.link}>Вернуться в магазин</Text>
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, padding: 20 },
  title: { fontSize: 24, fontWeight: 'bold', marginBottom: 20 },
  cartItem: { padding: 10, borderBottomWidth: 1, borderBottomColor: '#ccc' },
  summary: { padding: 10, backgroundColor: '#f0f0f0', borderRadius: 5, marginVertical: 10 },
  summaryText: { fontSize: 16, marginVertical: 5 },
  button: { backgroundColor: '#007AFF', padding: 15, borderRadius: 5, alignItems: 'center', marginVertical: 10 },
  buttonText: { color: '#fff', fontSize: 16 },
  link: { color: '#007AFF', textAlign: 'center' },
  loading: { flex: 1, textAlign: 'center', marginTop: 20, fontSize: 18 },
});

export default OrderPlacementScreen;