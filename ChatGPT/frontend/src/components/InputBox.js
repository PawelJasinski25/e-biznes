import React, { useState } from "react";

const InputBox = ({ sendMessage }) => {
    const [input, setInput] = useState("");

    const handleSend = () => {
        sendMessage(input);
        setInput("");
    };

    return (
        <div className="input-box">
            <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Wpisz wiadomość..."
            />
            <button onClick={handleSend}>🚀 Wyślij</button>
        </div>
    );
};

export default InputBox;
