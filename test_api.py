# coding:utf-8

import sys, urllib3, json, base64, time, hashlib
from datetime import datetime
from demo.utils import sm2, helper

urllib3.disable_warnings()


def request(hostname, body, url):
    appid = '66A095861BAE55F8735199DBC45D3E8E'
    unixtime = int(time.time())
    body['timestamp'] = unixtime
    body['appId'] = appid

    param_str = helper.gen_param_str(body)
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

    start_time = datetime.now()
    r = pool.urlopen('POST', url, body=body)
    print('[Time taken: {!s}]'.format(datetime.now() - start_time))

    return r



if __name__ == '__main__':
    if len(sys.argv)<4:
        print("usage: python3 %s <host> <api> <image_path>" % sys.argv[0])
        sys.exit(2)

    hostname = sys.argv[1]
    cate     = sys.argv[2]
    filepath = sys.argv[3]

    host = 'http://%s:5000'%hostname

    body = {
        'version'  : '1',
        #'signType' : 'SHA256', 
        'signType' : 'SM2',
        'encType'  : 'plain',
        'data'     : {},
    }

    if cate=="qa":
        url = host+'/api/bert_qa'
        body['data']['corpus'] = "金字塔（英语：pyramid），在建筑学上是指锥体建筑物，著名的有埃及金字塔，还有玛雅卡斯蒂略金字塔、阿兹特克金字塔（太阳金字塔、月亮金字塔）等。"
        body['data']['question'] = "金字塔是什么？"
    elif cate=="qa2":
        url = host+'/api/albert_qa'
        body['data']['corpus'] = "金字塔（英语：pyramid），在建筑学上是指锥体建筑物，著名的有埃及金字塔，还有玛雅卡斯蒂略金字塔、阿兹特克金字塔（太阳金字塔、月亮金字塔）等。"
        body['data']['question'] = "金字塔是什么？"
    elif cate=="ner":
        url = host+'/ner/ner'
        body['data']['text'] = "1500mg visiable, 2009 年12月底出现黑便,,于当地行胃镜检查并行病理检查示:叒胃体中下部溃疡,叒病理示中分化腺癌,叒无腹胀、泛酸、嗳气、恶心、呕吐、叒无头晕"
        #body['data']['text'] = ",2009年12月底出现黑便,,于当地行胃镜检查并行病理检查示:叒胃体中下部溃疡,叒病理示中分化腺癌,叒无腹胀、泛酸、嗳气、恶心、呕吐、叒无头晕、叒心悸、乏力等症,叒2010年1月13日于我院胃胰科行胃癌根治术,叒2010年1月18日,我院病理:切缘未见癌,叒胃体可见3x2x1cm3溃疡型肿物,叒镜上为中分化腺癌侵及胃壁全层至浆膜层,网膜未见癌,叒肝总动脉旁(0/1)、叒胃大弯(0/1)淋巴结未见癌,叒贲门左(3/3)、叒胃小弯(8/9)、幽门上(2/2)淋巴结可见腺癌转移,,免疫组化:cea(+)、叒p53(+)、叒pr(-)、叒er-b(+)、叒er(+++)、叒共计,ln:叒13/16转移,叒术后于2010年2月-2010年8月行术后化疗6程,叒具体用药为艾素100mg叒静点+叒希罗达1500mg叒bid叒po,2014年6月初出现右侧下上肢活动受限,叒7月份症状逐渐加重,叒7月10日就诊于*****,叒,行mri检查提示:胃癌术后多发脑转移,叒行甘露醇及地塞米松、叒洛赛克治疗后效果不佳。遂于我院就诊,2014-8-5行奥沙利铂150mg叒d1+叒替吉奥叒50mg叒bid叒d1-14化疗一程,2014-08-18开始行三维适形全脑放疗,剂量30gy/10f。2014-09-19始行替吉奥叒50mg叒bid叒d1-14单药化疗一程。本次为行上一程化疗收入我科,叒我科以“胃癌术后脑转移叒rtxnxm1叒iv期”收入,叒入科以来,叒精神饮食尚可,叒无恶心、叒呕吐,二便正常,体重无明显减低。"
    else:
        url = host+'/antigen/check'
        with open(filepath, 'rb') as f:
            img_data = f.read()
        body['data']['image'] = base64.b64encode(img_data).decode('utf-8')

    print("-->", url)

    r = request(hostname, body, url)

    print(r.status)
    if r.status==200:
        print(json.loads(r.data.decode('utf-8')))
    else:
        print(r.data)
