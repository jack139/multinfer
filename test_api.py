# coding:utf-8

import sys, urllib3, json, base64, time, hashlib
from datetime import datetime
from utils import sm2

urllib3.disable_warnings()

# 生成参数字符串
def gen_param_str(param1):
    param = param1.copy()
    name_list = sorted(param.keys())
    if 'data' in name_list: # data 按 key 排序, 中文不进行性转义，与go保持一致
        param['data'] = json.dumps(param['data'], sort_keys=True, ensure_ascii=False, separators=(',', ':'))
    return '&'.join(['%s=%s'%(str(i), str(param[i])) for i in name_list if str(param[i])!=''])


def request(hostname, body, cate):
    appid = '66A095861BAE55F8735199DBC45D3E8E'
    unixtime = int(time.time())
    body['timestamp'] = unixtime
    body['appId'] = appid

    param_str = gen_param_str(body)
    sign_str = '%s&key=%s' % (param_str, '43E554621FF7BF4756F8C1ADF17F209C')

    #print(sign_str)

    if body['signType'] == 'SHA256':
        signature_str =  base64.b64encode(hashlib.sha256(sign_str.encode('utf-8')).hexdigest().encode('utf-8')).decode('utf-8')
    else: # SM2
        signature_str = sm2.SM2withSM3_sign_base64(sign_str)

    body['signData'] = signature_str

    body = json.dumps(body)
    #print(body)

    pool = urllib3.PoolManager(num_pools=2, timeout=180, retries=False)

    host = 'http://%s:5000'%hostname
    

    if cate=="qa":
        url = host+'/api/bert_qa'
    elif cate=="qa2":
        url = host+'/api/albert_qa'
    else:
        url = host+'/antigen/check'

    print("-->", url)

    start_time = datetime.now()
    r = pool.urlopen('POST', url, body=body)
    print('[Time taken: {!s}]'.format(datetime.now() - start_time))

    return r



if __name__ == '__main__':
    if len(sys.argv)<4:
        print("usage: python3 %s <host> <api> <image_path>" % sys.argv[0])
        sys.exit(2)

    hostname = sys.argv[1]
    filepath = sys.argv[3]

    #with open(filepath, 'rb') as f:
    #    img_data = f.read()

    body = {
        'version'  : '1',
        #'signType' : 'SHA256', 
        'signType' : 'SM2',
        'encType'  : 'plain',
        'data'     : {
            #'image'    : base64.b64encode(img_data).decode('utf-8'),
            'corpus'   : "金字塔（英语：pyramid），在建筑学上是指锥体建筑物，著名的有埃及金字塔，还有玛雅卡斯蒂略金字塔、阿兹特克金字塔（太阳金字塔、月亮金字塔）等。",
            'question' : "金字塔是什么？",
            'text'     : "测试测试",
        }
    }

    r = request(hostname, body, sys.argv[2])

    print(r.status)
    if r.status==200:
        print(json.loads(r.data.decode('utf-8')))
    else:
        print(r.data)
