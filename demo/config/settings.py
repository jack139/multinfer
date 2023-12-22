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
DEMO_NER_PACK = ( 5000, '/ner/ner' )
DEMO_KERAS_QA = ( 5000, '/api/albert_qa' )
DEMO_TEXT2ORDER = ( 5000, '/talk2ui/text2order' )
DEMO_WAV2ORDER = ( 5000, '/talk2ui/wav2order' )
DEMO_WAV2TEXT = ( 5000, '/talk2ui/wav2text' )
DEMO_OCR_IDCARD = ( 5000, '/ocr2/id_card' )
DEMO_OCR_BANKCARD = ( 5000, '/ocr2/bank_card' )
DEMO_OCR_CARDNO = ( 5000, '/ocr2/card_no' )
DEMO_OCR_TEXT = ( 5000, '/ocr2/ocr_text' )