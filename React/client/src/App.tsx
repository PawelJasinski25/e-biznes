import { Routes, Route, Link } from "react-router-dom";
import Products from "./components/Product";
import Payment from "./components/Payment";

const App = () => {
    return (
        <div>
            <nav>
                <Link to="/">Home</Link> |
                <Link to="/products">Produkty</Link> |
                <Link to="/payment">Płatności</Link>
            </nav>

            <Routes>
                <Route path="/" element={<h1>Witaj w sklepie</h1>} />
                <Route path="/products" element={<Products />} />
                <Route path="/payment" element={<Payment />} />
            </Routes>
        </div>
    );
};

export default App;
