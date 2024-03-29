import json

# 生成参数字符串
def gen_param_str(param1):
    param = param1.copy()
    name_list = sorted(param.keys()) # data 按 key 排序, 中文不进行性转义，与go保持一致
    for key in name_list:
        if type(param[key])==type({}):
            param[key] = json.dumps(param[key], sort_keys=True, ensure_ascii=False, separators=(',', ':'))
    return '&'.join(['%s=%s'%(str(i), str(param[i])) for i in name_list if str(param[i])!=''])
