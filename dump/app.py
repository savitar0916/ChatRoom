from flask import Flask, jsonify, render_template
import requests

app = Flask(__name__)

# 将这里的URL替换为ngrok为Golang服务器提供的URL
GOLANG_SERVER_URL = 'http://localhost:8000/'

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/get_chat_log', methods=['GET'])
def get_chat_log():
    response = requests.get(f'{GOLANG_SERVER_URL}/messages')
    if response.status_code == 200:
        return jsonify(response.json())
    else:
        return jsonify({'error': 'Failed to fetch chat log'}), response.status_code

if __name__ == '__main__':
    app.run(port=5000)
