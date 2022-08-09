# Go实现多模型推理和api服务



## 测试



### 编译

```
make
```



### 启动 dispatcher

```
build/multinfer server 0
```



### 启动 http

```
build/multinfer http
```



### 启动 demo

```
cd demo
python3 app.py
```



### 测试脚本

```
python3 test_api.py 127.0.0.1 ner _
```
