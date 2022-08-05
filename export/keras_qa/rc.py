#! -*- coding: utf-8 -*-

import os

# AMP要使用 tf.keras 
os.environ["TF_KERAS"] = "1"


import json, shutil
import numpy as np
from bert4keras.backend import keras, K, tf
from bert4keras.models import build_transformer_model
from bert4keras.optimizers import Adam
from keras.layers import Layer, Dense, Permute
from keras.models import Model

# 基本信息
maxlen = 512
model_type = 'albert'


'''
bert4keras 支持的 BERT model_type
    'bert': BERT,
    'albert': ALBERT,
    'albert_unshared': ALBERT_Unshared,
    'roberta': BERT,
'''
if model_type=='bert':
    # bert配置
    config_path = '../../../../nlp/nlp_model/chinese_bert_L-12_H-768_A-12/bert_config.json'
    checkpoint_path = '../../../../nlp/nlp_model/chinese_bert_L-12_H-768_A-12/bert_model.ckpt'
    dict_path = '../../../../nlp/nlp_model/chinese_bert_L-12_H-768_A-12/vocab.txt'
elif model_type=='albert':
    # albert配置
    config_path = '../../../../nlp/nlp_model/albert_zh_base/albert_config.json'
    checkpoint_path = '../../../../nlp/nlp_model/albert_zh_base/model.ckpt-best'
    dict_path = '../../../../nlp/nlp_model/albert_zh_base/vocab_chinese.txt'
else:
    print('unknow model type.')
    sys.exit(1)


# 建立模型，载入权重
class MaskedSoftmax(Layer):
    """在序列长度那一维进行softmax，并mask掉padding部分
    """
    def compute_mask(self, inputs, mask=None):
        return None

    def call(self, inputs, mask=None):
        if mask is not None:
            mask = K.cast(mask, K.floatx())
            mask = K.expand_dims(mask, 2)
            inputs = inputs - (1.0 - mask) * 1e12
        return K.softmax(inputs, 1)


model = build_transformer_model(
    config_path=config_path,
    checkpoint_path=checkpoint_path,
    model=model_type
)


output = Dense(2)(model.output)
output = MaskedSoftmax()(output)
output = Permute((2, 1))(output)

model = Model(model.input, output)
#model.summary()


'''
    model.load_weights('outputs/albert_batch64_max512_lr2e-05_F1_82.000/best_model.weights')
    corpus = "深度学习（英语：deep learning）是机器学习的分支，是一种以人工神经网络为架构，对资料进行表征学习\
的算法。深度学习是机器学习中一种基于对数据进行表征学习的算法。观测值（例如一幅图像）可以使用多种方式来表示，如\
每个像素强度值的向量，或者更抽象地表示成一系列边、特定形状的区域等。而使用某些特定的表示方法更容易从实例中学习\
任务（例如，人脸识别或面部表情识别）。深度学习的好处是用非监督式或半监督式的特征学习和分层特征提取高效算法\
来替代手工获取特征。"
    ans = extract_answer("什么是深度学习？", corpus)
    print(ans)
'''