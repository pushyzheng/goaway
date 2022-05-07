from flask import Flask, request

app = Flask(__name__)


@app.route('/', methods=['GET'])
def index():
    print(request.args)
    print(request.json)
    return 'hello 1'

@app.route('/admin', methods=['GET'])
def index():
    return 'admin api'


if __name__ == '__main__':
    app.run(port=5000)
