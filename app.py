from flask import Flask, request

app = Flask(__name__)


@app.route('/', methods=['GET'])
def index():
    print(request.args)
    print(request.json)
    return 'hello 1'

@app.route('/public', methods=['GET'])
def public():
    return 'public'

@app.route('/你好世界/', methods=['GET'])
def hello_world_ch():
    return '你好世界'

@app.route('/admin', methods=['GET'])
def admin():
    return 'admin api'


if __name__ == '__main__':
    app.run(port=5000)
