import { useState, useEffect } from "react";
import axios from "axios";

interface CartItem {
    product_id: number;
    quantity: number;
}

interface Product {
    id: number;
    name: string;
    price: number;
}

interface Cart {
    items: CartItem[];
}

const Cart = () => {
    const [cart, setCart] = useState<Cart | null>(null);
    const [products, setProducts] = useState<Product[]>([]);

    useEffect(() => {
        axios.get<Cart>("http://localhost:8080/cart")
            .then(response => setCart(response.data))
            .catch(error => console.error("Error fetching cart:", error));

        axios.get<Product[]>("http://localhost:8080/products")
            .then(response => setProducts(response.data))
            .catch(error => console.error("Error fetching products:", error));
    }, []);

    const calculateTotalPrice = () => {
        return cart?.items.reduce((total, item) => {
            const product = products.find(p => p.id === item.product_id);
            return total + (product?.price || 0) * item.quantity;
        }, 0) || 0;
    };

    return (
        <div>
            <h2>Twój koszyk</h2>
            {cart?.items?.length ? (
                <>
                    <ul>
                        {cart.items.map((item) => {
                            const product = products.find(p => p.id === item.product_id);
                            const totalPrice = (product?.price || 0) * item.quantity;

                            return (
                                <li key={item.product_id}>
                                    {product?.name} - Ilość: {item.quantity} - {totalPrice.toFixed(2)} zł
                                </li>
                            );
                        })}
                    </ul>
                    <h3>Całkowita wartość koszyka: {calculateTotalPrice().toFixed(2)} zł</h3>
                </>
            ) : (
                <p>Koszyk jest pusty! Dodaj produkty do koszyka.</p>
            )}
        </div>
    );


};

export default Cart;
