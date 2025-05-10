import { type ChangeEvent, type FormEvent, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import API from '../api';

interface RegisterData {
    email: string;
    name: string;
    surname: string;
    password: string;
}

const Register = () => {
    const [registerData, setRegisterData] = useState<RegisterData>({
        email: '',
        name: '',
        surname: '',
        password: '',
    });
    const navigate = useNavigate();

    const handleRegister = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        try {
            await API.post('/register', registerData);
            alert('Rejestracja zakończona powodzeniem!');
            navigate('/login');
        } catch (error) {
            alert('Błąd rejestracji. Wprowadź poprawne dane i spróbuj ponownie.');
            console.error('Registration error:', error);
        }
    };

    return (
        <div style={{ maxWidth: '400px', margin: 'auto', paddingTop: '50px' }}>
            <form onSubmit={handleRegister}>
                <h2>Rejestracja</h2>
                <input
                    type="email"
                    placeholder="Email"
                    value={registerData.email}
                    onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        setRegisterData({ ...registerData, email: e.target.value })
                    }
                    required
                />
                <input
                    type="text"
                    placeholder="Imię"
                    value={registerData.name}
                    onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        setRegisterData({ ...registerData, name: e.target.value })
                    }
                    required
                />
                <input
                    type="text"
                    placeholder="Nazwisko"
                    value={registerData.surname}
                    onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        setRegisterData({ ...registerData, surname: e.target.value })
                    }
                    required
                />
                <input
                    type="password"
                    placeholder="Hasło"
                    value={registerData.password}
                    onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        setRegisterData({ ...registerData, password: e.target.value })
                    }
                    required
                />
                <button type="submit">Zarejestruj się</button>
            </form>
            <p>
                Masz już konto? <Link to="/login">Zaloguj się</Link>
            </p>
        </div>
    );
};

export default Register;
