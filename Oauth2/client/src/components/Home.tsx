import { useNavigate } from 'react-router-dom';

const Home = () => {
    const navigate = useNavigate();

    const handleLogout = () => {
        navigate('/login');
    };

    return (
        <div style={{ textAlign: 'center', paddingTop: '50px' }}>
            <p>Zalogowano pomyślnie. Witamy w aplikacji!</p>
            <button onClick={handleLogout}>Wyloguj się</button>
        </div>
    );
};

export default Home;
