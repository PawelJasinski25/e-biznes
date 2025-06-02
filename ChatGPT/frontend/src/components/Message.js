import React from "react";

const Message = ({ role, content }) => {
    return (
        <div className={`message-container ${role === "user" ? "left" : "right"}`}>
            <div className={`message ${role === "user" ? "user-message" : "assistant-message"}`}>
                {content}
            </div>
        </div>
    );
};

export default Message;
