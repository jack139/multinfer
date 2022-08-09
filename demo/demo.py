# coding:utf-8

import os
from flask import Flask, Blueprint, render_template, request
import urllib3, json, base64, time, hashlib
from datetime import datetime
from utils import helper, sm2
from config.settings import MAX_IMAGE_SIZE


ALLOWED_EXTENSIONS = set(['png', 'jpg', 'jpeg'])
def allowed_file(filename):
    return '.' in filename and \
           filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS


# 接口演示 antigen

demo_antigen = Blueprint('demo_antigen', __name__)


@demo_antigen.route("/demo/antigen", methods=["GET"])
def demo_get():
    return render_template('demo_antigen.html')

@demo_antigen.route("/demo/antigen", methods=["POST"])
def demo_post():
    file = request.files['file']
    if file and allowed_file(file.filename):
        #file.save(os.path.join(os.getcwd(), file.filename))
        image_data = file.stream.read()
        if len(image_data)>MAX_IMAGE_SIZE:
            return "image file size exceeds 3MB"
        body_data = { 'image' : base64.b64encode(image_data).decode('utf-8') }
        api_url, params, status, rdata, timespan = call_api("antigen", body_data)
        return render_template('result.html', 
            result=rdata, status=status, 
            timespan=timespan, params=params, api_url=api_url)
    else:
        return "not allowed image"


# 接口演示 NER

demo_ner = Blueprint('demo_ner', __name__)

@demo_ner.route("/demo/ner", methods=["GET"])
def demo_get():
    return render_template('demo_ner.html')

@demo_ner.route("/demo/ner", methods=["POST"])
def demo_post():
    text = request.form['text']
    body_data = { 'text' : text }
    api_url, params, status, rdata, timespan = call_api("ner", body_data)
    return render_template('result.html', 
        result=rdata, status=status, 
        timespan=timespan, params=params, api_url=api_url)


# 接口演示 QA

demo_qa = Blueprint('demo_qa', __name__)

@demo_qa.route("/demo/qa", methods=["GET"])
def demo_get():
    return render_template('demo_qa.html')

@demo_qa.route("/demo/qa", methods=["POST"])
def demo_post():
    corpus = request.form['corpus']
    question = request.form['question']
    body_data = { 'corpus' : corpus, 'question' : question }
    api_url, params, status, rdata, timespan = call_api("keras_qa", body_data)
    return render_template('result.html', 
        result=rdata, status=status, 
        timespan=timespan, params=params, api_url=api_url)


# 调用接口
def call_api(cate, body_data):
    hostname = '127.0.0.1'

    body = {
        'version'  : '1',
        'signType' : 'SHA256', 
        #'signType' : 'SM2',
        'encType'  : 'plain',
        'data'     : body_data
    }

    appid = '66A095861BAE55F8735199DBC45D3E8E'
    unixtime = int(time.time())
    body['timestamp'] = unixtime
    body['appId'] = appid

    param_str = helper.gen_param_str(body)
    sign_str = '%s&key=%s' % (param_str, '43E554621FF7BF4756F8C1ADF17F209C')

    if body['signType'] == 'SHA256':
        signature_str =  base64.b64encode(hashlib.sha256(sign_str.encode('utf-8')).hexdigest().encode('utf-8')).decode('utf-8')
    else: # SM2
        signature_str = sm2.SM2withSM3_sign_base64(sign_str)

    #print(sign_str)

    body['signData'] = signature_str

    body_str = json.dumps(body)
    #print(body)

    pool = urllib3.PoolManager(num_pools=2, timeout=180, retries=False)

    host = 'http://%s:5000'%hostname

    if cate=='antigen':
        url = host+'/antigen/check'
    elif cate=='ner':
        url = host+'/ner/ner'
    elif cate=='keras_qa':
        url = host+'/api/albert_qa'
    else:
        url = host+'/'

    start_time = datetime.now()
    r = pool.urlopen('POST', url, body=body_str)
    #print('[Time taken: {!s}]'.format(datetime.now() - start_time))

    print(r.status)
    if r.status==200:
        rdata = json.dumps(json.loads(r.data.decode('utf-8')), ensure_ascii=False, indent=4)
    else:
        rdata = r.data

    # 截短一下 image 字段显示内容
    if cate=='antigen':
        body['data']['image'] = body['data']['image'][:20]+' ... ' + body['data']['image'][-20:]

    body2 = json.dumps(body, ensure_ascii=False, indent=4)
    return url, body2, r.status, rdata, \
        '{!s}s'.format(datetime.now() - start_time)
