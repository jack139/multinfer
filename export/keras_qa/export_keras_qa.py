# coding=utf-8

# 导出 locate 模型权重，可用于go tf.LoadSavedModel


import os, shutil
import tensorflow as tf

# ------- from keras_QA 项目
import rc

os.environ["CUDA_VISIBLE_DEVICES"] = '0'

# ------- albert model
rc.model.load_weights("../../../../nlp/keras_QA/outputs/albert_batch64_max512_lr2e-05_F1_82.000/best_model.weights")
rc.model.summary()

save_model_path = "../outputs/keras_qa/saved-model"
if os.path.exists(save_model_path):
    shutil.rmtree(save_model_path) 
os.makedirs(save_model_path)

print('\n'.join([n.name for n in tf.get_default_graph().as_graph_def().node])) # 所有层的名字

# save_model 输出 , for goland 测试
builder = tf.saved_model.builder.SavedModelBuilder(save_model_path)
builder.add_meta_graph_and_variables(rc.K.get_session(), [tf.saved_model.tag_constants.TRAINING], clear_devices=True)
builder.save()  

