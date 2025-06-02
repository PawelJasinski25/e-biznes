import React, { useState } from "react";
import axios from "axios";
import Message from "./components/Message";
import InputBox from "./components/InputBox";
import "./index.css";

const App = () => {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);

  const sendMessage = async (input) => {
    if (!input.trim()) return;
    const userMessage = { role: "user", content: input };

    setMessages([...messages, userMessage]);
    setLoading(true);

    try {
      const response = await axios.post("http://localhost:8000/chat", {
        user_message: input,
      });

      const botMessage = { role: "assistant", content: response.data.response };
      setMessages([...messages, userMessage, botMessage]);
    } catch (error) {
      const botMessage = { role: "assistant", content: "❌ Błąd serwera!" };
      setMessages([...messages, userMessage, botMessage]);
    }
    setLoading(false);
  };

  return (
      <div className="chat-container">
        <h1>💬 Chatbot Sklepu</h1>
        <div className="chat-box">
          {messages.map((msg, index) => (
              <Message key={index} role={msg.role} content={msg.content} />
          ))}
          {loading && <p className="loading">⏳ Generowanie odpowiedzi...</p>}
        </div>
        <InputBox sendMessage={sendMessage} />
      </div>
  );
};

export default App;
