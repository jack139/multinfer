#! -*- coding: utf-8 -*-
# 用GlobalPointer做中文命名实体识别

# 模型来自 medicla_ner 项目

import os

from bert4keras.layers import EfficientGlobalPointer as GlobalPointer
from bert4keras.models import build_transformer_model
from keras.models import Model

# 标签： 要与训练时顺序一致
categories = ['检验和检查', '治疗和手术', '疾病和诊断', '症状和体征', '药物', '解剖部位']

config_path = '../../../../nlp/nlp_model/chinese_bert_L-12_H-768_A-12/bert_config.json'

model = build_transformer_model(config_path)
output = GlobalPointer(len(categories), 64)(model.output)

model = Model(model.input, output)
#model.summary()
