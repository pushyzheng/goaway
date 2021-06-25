from flask import Flask, request

app = Flask(__name__)


@app.route('/', methods=['POST'])
def index():
    import time
    time.sleep(2)

    print(request.args)
    print(request.json)

    return 'hello 1'


if __name__ == '__main__':
    app.run(port=5000)
