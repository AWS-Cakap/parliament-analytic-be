import requests
import csv
import os
from datetime import datetime
from transformers import AutoTokenizer, AutoModelForSequenceClassification
import torch


BEARER_TOKEN = 'AAAAAAAAAAAAAAAAAAAAADj90QEAAAAAU4bKrXUSF4jLkUyn20TFlnbpGNM%3D2xv7R8RjHbXXSXTbJ9S11lRte9Ozqir3fldIWjcRBL4ahK4CDQ' 
RAW_CSV = "raw_tweets.csv"
ANALYZED_CSV = "analyzed_tweets.csv"

HEADERS = {
    'Authorization': f'Bearer {BEARER_TOKEN}',
}

tokenizer = AutoTokenizer.from_pretrained("mdhugol/indonesia-bert-sentiment-classification")
model = AutoModelForSequenceClassification.from_pretrained("mdhugol/indonesia-bert-sentiment-classification")

def analyze_sentiment(text):
    inputs = tokenizer(text, return_tensors="pt", truncation=True, padding=True)
    with torch.no_grad():
        outputs = model(**inputs)
        probs = torch.nn.functional.softmax(outputs.logits, dim=-1)
        sentiment_idx = torch.argmax(probs).item()
        labels = ["neg", "neu", "pos"]
        return labels[sentiment_idx]

def save_to_csv(filename, fieldnames, data, mode='a'):
    file_exists = os.path.isfile(filename)
    with open(filename, mode, newline='', encoding='utf-8') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        if not file_exists or mode == 'w':
            writer.writeheader()
        writer.writerow(data)

def crawl_and_analyze(query, max_tweets=15):
    url = 'https://api.twitter.com/2/tweets/search/recent'
    params = {
        'query': f'{query} lang:id -is:retweet',
        'tweet.fields': 'created_at,lang,author_id',
        'expansions': 'author_id',
        'user.fields': 'username',
        'max_results': min(max_tweets, 100)
    }

    response = requests.get(url, headers=HEADERS, params=params)
    data = response.json()

    if 'data' not in data:
        print("Tidak ada tweet ditemukan:", data)
        return

    tweets = data['data']
    users = {user['id']: user['username'] for user in data.get('includes', {}).get('users', [])}

    for tweet in tweets:
        text = tweet['text']
        created_at = tweet['created_at']
        author_id = tweet['author_id']
        username = users.get(author_id, 'unknown')

        save_to_csv(RAW_CSV, ['username', 'text', 'created_at'], {
            'username': username,
            'text': text,
            'created_at': created_at
        })

        sentiment = analyze_sentiment(text)

        analyzed_data = {
            "username": username,
            "text": text,
            "created_at": created_at,
            "sentiment": sentiment
        }
        save_to_csv(ANALYZED_CSV, ['username', 'text', 'created_at', 'sentiment'], analyzed_data)

        print(f"{username} | {sentiment} | {text[:50]}...")


if __name__ == "__main__":
    crawl_and_analyze("Gerindra", max_tweets=15)
