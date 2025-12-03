import requests

url = f"http://127.0.0.1:30000/flush_cache"
response = requests.post(url)
print(response)
