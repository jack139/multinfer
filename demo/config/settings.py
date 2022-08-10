# -*- coding: utf-8 -*-

############  app server 相关设置

APP_NAME = 'multinfer_demo'

# 参数设置
DEBUG_MODE = False
BIND_ADDR = '0.0.0.0'
BIND_PORT = '8000'

# 图片数据最大尺寸
MAX_IMAGE_SIZE = 1024*1024*3  # 3MB

# port 和 url
DEMO_ANTIGEN = ( 5000, '/antigen/check' )
DEMO_NER_PACK = ( 5000, '/ner/ner' ) # 5001 in ocr-ali
DEMO_KERAS_QA = ( 5000, '/api/albert_qa' ) # 5001 in ocr-ali