import { useState } from "react";
import axios from "axios";

interface Payment {
    name: string;
    surname: string;
    credit_card_number: string;
}


const Payment = () => {
    const [payment, setPayment] = useState<Payment>({
        name: "",
        surname: "",
        credit_card_number: "",
    });

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        console.log("📡 Sending payment data:", payment);

        try {
            await axios.post("http://localhost:8080/payments", {
                name: payment.name,
                surname: payment.surname,
                credit_card_number: payment.credit_card_number,
            }, {
                headers: { "Content-Type": "application/json" },
            });

            alert("Płatność przetworzona!");
        } catch (error) {
            console.error("Błąd płatności:", error);
        }
    };


    return (
        <div>
            <h2>Formularz płatności</h2>
            <form onSubmit={handleSubmit}>
                <input type="text" placeholder="Imię"
                       onChange={(e) => setPayment({ ...payment, name: e.target.value })} />
                <input type="text" placeholder="Nazwisko"
                       onChange={(e) => setPayment({ ...payment, surname: e.target.value })} />
                <input type="text" placeholder="Numer karty"
                       onChange={(e) => setPayment({ ...payment, credit_card_number: e.target.value })} />

                <button type="submit">Zapłać</button>
            </form>
        </div>
    );
};

export default Payment;
