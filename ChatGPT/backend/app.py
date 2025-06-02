from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from gpt4all import GPT4All
from starlette.middleware.cors import CORSMiddleware
import random

app = FastAPI()

MODEL_PATH = "mistral-7b-instruct-v0.1.Q4_0.gguf"
gpt4all_model = GPT4All(MODEL_PATH)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class ChatRequest(BaseModel):
    user_message: str

# List of predefined conversation openers and closers
OPENERS = [
    "Hello. How can I assist you today?",
    "Welcome to our store. Are you looking for a specific item?",
    "Hi. I have great fashion recommendations for you.",
    "How can I help? Do you have a style in mind?",
    "Good to see you. Are you looking for something casual or more formal?"
]

CLOSERS = [
    "Thank you for chatting. Feel free to return anytime.",
    "I'm glad I could assist. If you want to see more outfits, check our collection.",
    "I hope you found what you were looking for. Have a great day.",
    "That was a great conversation. Visit our store to see our latest arrivals.",
    "Thanks for chatting. A well-chosen outfit makes all the difference."
]

GREETINGS = ["hello", "hi", "hey", "good morning", "good afternoon", "good evening"]

@app.get("/welcome")
async def welcome():
    """ Returns a random greeting when the chat starts. """
    return {"response": random.choice(OPENERS)}

@app.post("/chat")
async def chat(chat_req: ChatRequest):
    user_message = chat_req.user_message.strip().lower()

    try:
        if user_message in GREETINGS:
            response_text = random.choice(OPENERS)
        else:
            response_text = gpt4all_model.generate(user_message, max_tokens=100, temp=0.7)

        if user_message in ["bye", "goodbye", "thanks"]:
            response_text += "\n\n" + random.choice(CLOSERS)

        return {"response": response_text.strip()}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("app:app", host="0.0.0.0", port=8000, reload=True)
