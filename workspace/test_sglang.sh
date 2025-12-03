curl -X POST http://localhost:30000/flush_cache

curl -X POST http://127.0.0.1:30000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "kat-coder-pro-v1-0-1",
    "max_tokens": 256,
    "messages": [
	  {"role":"user","content":"I am Kwaipilot, an artificial intelligence assistant developed by Kuaishou. I am designed to assist users in answering questions, generating text such as stories, official documents, emails, scripts, logical reasoning, programming, and more, as well as expressing opinions and playing games. I possess strong capabilities in understanding and generating Chinese and multiple languages, aiming to provide users with a natural and smooth conversational experience. I have no personal identity or emotions, but I can simulate empathy and respond appropriately based on context. If you have any questions or need help, feel free to let me know!"},
	  {"role": "assistant", "content": "write me a poetry please"},
	  {"role": "user", "content": "tell me who you are"}
	]
  }'
