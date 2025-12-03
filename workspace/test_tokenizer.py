import torch
import numpy as np
from transformers import AutoTokenizer

model_path = "/home/relay/liujiacheng06/models/kat-coder-pro-v1-0-1"

tokenizer = AutoTokenizer.from_pretrained(
    model_path,
    trust_remote_code=True,
    local_files_only=True
)

messages = [
	  {"role":"user","content":"I am Kwaipilot, an artificial intelligence assistant developed by Kuaishou. I am designed to assist users in answering questions, generating text such as stories, official documents, emails, scripts, logical reasoning, programming, and more, as well as expressing opinions and playing games. I possess strong capabilities in understanding and generating Chinese and multiple languages, aiming to provide users with a natural and smooth conversational experience. I have no personal identity or emotions, but I can simulate empathy and respond appropriately based on context. If you have any questions or need help, feel free to let me know!"},
	  {"role": "assistant", "content": "write me a poetry please"},
	  {"role": "user", "content": "tell me who you are"}
]

prompt = tokenizer.apply_chat_template(
    messages,
    tokenize=False,
    add_generation_prompt=True
)

print("templated prompt\n")
print(prompt)

token_ids = tokenizer(prompt, return_tensors="pt")['input_ids']
print(token_ids.shape)
print(token_ids)

block_0 = [151644, 872, 198, 40, 1079, 730, 9991, 573, 23958, 11, 458, 20443, 11229, 17847, 7881, 553, 730, 4284, 812, 283, 13, 358, 1079, 6188, 311, 7789, 3847, 304, 35764, 4755, 11, 23163, 1467, 1741, 438, 7343, 11, 3946, 9293, 11, 14298, 11, 19502, 11, 19819, 32711, 11, 15473, 11, 323, 803, 11, 438, 1632, 438, 36710, 17979, 323, 5619, 3868, 13, 358, 15218, 3746]
block_1 = [16928, 304, 8660, 323, 23163, 8453, 323, 5248, 15459, 11, 37078, 311, 3410, 3847, 448, 264, 5810, 323, 10876, 7517, 1663, 3139, 13, 358, 614, 902, 4345, 9569, 476, 21261, 11, 714, 358, 646, 37453, 47351, 323, 5889, 34901, 3118, 389, 2266, 13, 1416, 498, 614, 894, 4755, 476, 1184, 1492, 11, 2666, 1910, 311, 1077, 752, 1414, 0, 151645, 198, 151644, 77091, 198]
block_2 = []

blocks = block_0 + block_1 + block_2

decoded = tokenizer.decode(blocks)
print(decoded)

blocks_ref = [151644, 872, 198, 40, 1079, 730, 9991, 573, 23958, 11, 458, 20443, 11229, 17847, 7881, 553, 730, 4284, 812, 283, 13, 358, 1079, 6188, 311, 7789, 3847, 304, 35764, 4755, 11, 23163, 1467, 1741, 438, 7343, 11, 3946, 9293, 11, 14298, 11, 19502, 11, 19819, 32711, 11, 15473, 11, 323, 803, 11, 438, 1632, 438, 36710, 17979, 323, 5619, 3868, 13, 358, 15218, 3746, 16928, 304, 8660, 323, 23163, 8453, 323, 5248, 15459, 11, 37078, 311, 3410, 3847, 448, 264, 5810, 323, 10876, 7517, 1663, 3139, 13, 358, 614, 902, 4345, 9569, 476, 21261, 11, 714, 358, 646, 37453, 47351, 323, 5889, 34901, 3118, 389, 2266, 13, 1416, 498, 614, 894, 4755, 476, 1184, 1492, 11, 2666, 1910, 311, 1077, 752, 1414, 0, 151645, 198, 151644, 77091, 198, 4934, 752, 264, 31249, 4486, 151645, 198, 151644, 872, 198, 72357, 752, 879, 498, 525, 151645, 198]

len_common = min(len(blocks), len(blocks_ref))

for i in range(len_common):
    if blocks[i] != blocks_ref[i]:
        print("damn", i, blocks[i], blocks_ref[i])
    else:
        print("good", i, blocks[i], blocks_ref[i])


