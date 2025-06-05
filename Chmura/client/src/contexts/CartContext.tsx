import { createContext, useContext, useState, useEffect } from "react";
import axios from "axios";

interface CartItem {
    product_id: number;
    quantity: number;
}

interface CartContextType {
    cart: CartItem[];
    addToCart: (productId: number) => void;
    fetchCart: () => void;
}

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

const CartContext = createContext<CartContextType | undefined>(undefined);

export const CartProvider = ({ children }: { children: React.ReactNode }) => {
    const [cart, setCart] = useState<CartItem[]>([]);

    const fetchCart = async () => {
        try {
            const response = await axios.get<{ items: CartItem[] }>(`${API_URL}/cart`);
            setCart(response.data.items);
        } catch (error) {
            console.error("Błąd pobierania koszyka:", error);
        }
    };

    useEffect(() => {
        fetchCart();
    }, []);

    const addToCart = async (productId: number) => {
        try {
            await axios.post(`${API_URL}/cart`, {
                product_id: productId,
                quantity: 1,
            });
            fetchCart();
        } catch (error) {
            console.error("Error adding to cart:", error);
        }
    };

    return (
        <CartContext.Provider value={{ cart, addToCart, fetchCart }}>
            {children}
        </CartContext.Provider>
    );
};

export const useCart = () => {
    const context = useContext(CartContext);
    if (!context) {
        throw new Error("useCart must be used within a CartProvider");
    }
    return context;
};
