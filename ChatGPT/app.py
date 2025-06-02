from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from gpt4all import GPT4All

app = FastAPI()

MODEL_PATH = "mistral-7b-instruct-v0.1.Q4_0.gguf"
gpt4all_model = GPT4All(MODEL_PATH)


class ChatRequest(BaseModel):
    user_message: str

@app.post("/chat")
async def chat(chat_req: ChatRequest):
    user_message = chat_req.user_message.strip()

    try:
        response_text = gpt4all_model.generate(user_message, max_tokens=100, temp=0.7)

        return {"response": response_text.strip()}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("app:app", host="0.0.0.0", port=8000, reload=True)
