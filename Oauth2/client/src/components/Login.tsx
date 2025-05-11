import { type ChangeEvent, type FormEvent, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import API from '../api';

interface LoginData {
    email: string;
    password: string;
}

const Login = () => {
    const [loginData, setLoginData] = useState<LoginData>({ email: '', password: '' });
    const navigate = useNavigate();

    const handleLogin = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        try {
            const response = await API.post('/login', loginData);
            console.log(`Zalogowano: ${response.data.user.email}`);
            navigate('/home');
        } catch (error) {
            alert('Błąd logowania. Wprowadź poprawny email lub hasło');
            console.error('Login error:', error);
        }
    };

    const handleGoogleLogin = () => {
        window.location.href = "http://localhost:8080/auth/google/login";
    };

    return (
        <div style={{ maxWidth: '400px', margin: 'auto', paddingTop: '50px' }}>
            <form onSubmit={handleLogin}>
                <h2>Logowanie</h2>
                <input
                    type="email"
                    placeholder="Email"
                    value={loginData.email}
                    onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        setLoginData({ ...loginData, email: e.target.value })
                    }
                    required
                />
                <input
                    type="password"
                    placeholder="Hasło"
                    value={loginData.password}
                    onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        setLoginData({ ...loginData, password: e.target.value })
                    }
                    required
                />
                <button type="submit">Zaloguj się</button>
            </form>
            <p>
                <button onClick={handleGoogleLogin} style={{ marginLeft: '10px' }}>
                    Zaloguj przez Google
                </button>
            </p>
            <p>
                Nie masz konta? <Link to="/register">Zarejestruj się</Link>
            </p>
        </div>
    );
};

export default Login;
