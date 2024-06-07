import requests # type: ignore

def dump_chat():
    url = 'http://localhost:8000/dump_data'  # 改用新的路由来获取数据
    response = requests.get(url)

    if response.status_code == 200:
        with open('chat_dump.json', 'wb') as file:
            file.write(response.content)
        print("Chat log dumped to chat_dump.json")
    else:
        print("Failed to dump chat log")

if __name__ == '__main__':
    dump_chat()
