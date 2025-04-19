import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TouchableOpacity, TextInput, StyleSheet, Alert } from 'react-native';
import { getProducts } from '../api/api';

const CatalogScreen = ({ navigation }) => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [cart, setCart] = useState({});

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const response = await getProducts();
        setProducts(response.data);
      } catch (error) {
        Alert.alert('Error', 'Failed to load products');
      } finally {
        setLoading(false);
      }
    };
    fetchProducts();
  }, []);

  const addToCart = (productId) => {
    setCart((prev) => ({
      ...prev,
      [productId]: (prev[productId] || 0) + 1,
    }));
  };

  const goToCheckout = () => {
    const items = Object.keys(cart).map((productId) => ({
      product_id: productId,
      quantity: cart[productId],
    }));
    navigation.navigate('OrderPlacement', { cartItems: items });
  };

  if (loading) {
    return <Text style={styles.loading}>Loading...</Text>;
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Shop</Text>
        <TouchableOpacity onPress={goToCheckout}>
          <Text style={styles.cart}>Cart ({Object.values(cart).reduce((a, b) => a + b, 0)})</Text>
        </TouchableOpacity>
      </View>
      <TextInput style={styles.search} placeholder="Search protein bars, supplements..." />
      <FlatList
        data={products}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <View style={styles.productCard}>
            <Text style={styles.productName}>{item.name}</Text>
            <Text style={styles.productDescription}>{item.description}</Text>
            <Text style={styles.productPrice}>${item.price.toFixed(2)}</Text>
            <TouchableOpacity style={styles.addButton} onPress={() => addToCart(item.id)}>
              <Text style={styles.addButtonText}>Add to Cart</Text>
            </TouchableOpacity>
          </View>
        )}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, padding: 10 },
  header: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  title: { fontSize: 24, fontWeight: 'bold' },
  cart: { fontSize: 16, color: '#007AFF' },
  search: { borderWidth: 1, borderColor: '#ccc', padding: 10, marginVertical: 10, borderRadius: 5 },
  productCard: { padding: 10, borderBottomWidth: 1, borderBottomColor: '#ccc' },
  productName: { fontSize: 18, fontWeight: 'bold' },
  productDescription: { fontSize: 14, color: '#666' },
  productPrice: { fontSize: 16, marginVertical: 5 },
  addButton: { backgroundColor: '#007AFF', padding: 10, borderRadius: 5, alignItems: 'center' },
  addButtonText: { color: '#fff', fontSize: 14 },
  loading: { flex: 1, textAlign: 'center', marginTop: 20, fontSize: 18 },
});

export default CatalogScreen;