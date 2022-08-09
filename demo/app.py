# -*- coding: utf-8 -*-

from flask import Flask

from demo import demo_antigen, demo_ner, demo_qa
from config.settings import BIND_ADDR, BIND_PORT, DEBUG_MODE

app = Flask(__name__)

@app.route('/')
def hello_world():
    return 'Hello World!'

# demo
app.register_blueprint(demo_antigen)
app.register_blueprint(demo_ner)
app.register_blueprint(demo_qa)


if __name__ == '__main__':
    # 外部可见，出错时带调试信息（debug=True）
    # 转生产时，接口需要增减校验机制，避免非授权调用 ！！！！！！
    app.run(host=BIND_ADDR, port=BIND_PORT, debug=DEBUG_MODE)
