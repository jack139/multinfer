# coding=utf-8

# 导出 locate 模型权重，可用于go tf.LoadSavedModel


import os, shutil
from keras import backend as K
import tensorflow as tf

# ------- from medical_ner 项目
from ner_model import model

os.environ["CUDA_VISIBLE_DEVICES"] = '0'

# ------- albert model
model.load_weights("../../../../nlp/medical_ner/pack_best_f1_0.82966.weights")
model.summary()

save_model_path = "../outputs/ner_pack/saved-model"
if os.path.exists(save_model_path):
    shutil.rmtree(save_model_path) 
os.makedirs(save_model_path)

print('\n'.join([n.name for n in tf.get_default_graph().as_graph_def().node])) # 所有层的名字

# save_model 输出 , for goland 测试
builder = tf.saved_model.builder.SavedModelBuilder(save_model_path)
builder.add_meta_graph_and_variables(K.get_session(), [tf.saved_model.tag_constants.TRAINING], clear_devices=True)
builder.save()  
