
## api 文档

### 1. 说明

（1）接口分两大类：人脸识别、定位；人脸特征库的注册和管理。

（2）人脸特征库的人脸数据分组管理，识别时也针对一个组的人脸进行搜索。目的一是将搜索范围控制在比较少的人脸数量，提高人脸搜索性能；二是方便针对不同场景注册不同人脸识别分组，例如针对道闸场景，可以对不同闸机注册不同的人脸分组，识别不同的人群。

（3）人脸注册时，用户user_id和姓名是必填项，user_id要保证唯一。其他为选填项，为了后续能使用对数据进行再利用，建议尽可能填写选填项（性别、年龄）。

（4）最常用的搜索接口是face/search，此接口同时实现了1:1和1:N搜索，当提供人脸图片和用户user_id时就是1:1，当不提供user_id只提供人脸图片时就是1:N。如果同时提供了人脸图片和手机号码后4位，则只有当手机后4位匹配时才返回，相当于双因素认证。

（5）目前传入接口图片大小要求小于2MB，注意这个尺寸不是图片大小而是图片经base64编码后的大小。如果需要，这个最大2MB的限制可以在后台调整。

（6）接口验签接受SHA256和SM2算法。SM2验签算法使用架构部的yhtool-crypto测试通过（如果算法无变化，对yhtool的版本没有具体要求，测试使用yhtool-crypto-1.3.0-RELEASE.jar）。

（7）face/feedback接口用于结果反馈，帮助后台优化识别性能。如果应用场景有机会获知真实的结果，可以用feedback将此次接口调用（用接口返回的request_id标识）是否正确的信息返回给后台。

（8）提高搜索准确度的建议：（A）照片采集尽量正面光照良好，不要有帽子、口罩；（B）对同一人尽可能多注册几张照片；（C）如果确认某人不再需要识别，应从组内删除，控制组内人数，人数越少识别性能会相对更好。

（9）人脸特征数据会在启动时装入内存，新增和删除（/facedb/face/\*）的人脸动态修改内存特征数据，其他对人脸特征的操作需要重启gosearch服务。



### 2. api清单

| url                    | 功能             |
| ---------------------- | ---------------- |
| /face/locate           | 人脸定位         |
| /face/verify           | 人脸对比         |
| /face/search           | 人脸搜索         |
| /face/check            | 是否有人脸       |
| /facedb/face/reg       | 人脸注册         |
| /facedb/face/update    | 人脸更新         |
| /facedb/face/remove    | 人脸删除         |
| /facedb/user/info      | 用户信息查询     |
| /facedb/user/face_list | 获取用户人脸列表 |
| /facedb/user/list      | 获取用户列表     |
| /facedb/user/copy      | 复制用户         |
| /facedb/user/remove    | 删除用户         |
| /facedb/group/new      | 创建用户组       |
| /facedb/group/remove   | 删除用户组       |
| /facedb/group/list     | 获取用户组列表   |
| /facedb/feedback       | 识别结果反馈     |



### 3. 全局接口定义

输入参数

| 参数      | 类型   | 说明                          | 示例        |
| --------- | ------ | ----------------------------- | ----------- |
| appId     | string | 应用渠道编号                  |             |
| version   | string | 版本号                        |             |
| signType  | string | 签名算法，目前使用国密SM2算法 | SM2或SHA256 |
| signData  | string | 签名数据，具体算法见下文      |             |
| encType   | string | 接口数据加密算法，目前不加密  | plain       |
| timestamp | int    | unix时间戳（秒）              |             |
| data      | json   | 接口数据，详见各接口定义      |             |

> 签名/验签算法：
>
> 1. 筛选，获取参数键值对，剔除signData、encData、extra三个参数。data参数按key升序排列进行json序列化。
> 2. 排序，按key升序排序。
> 3. 拼接，按排序好的顺序拼接请求参数
>
> ```key1=value1&key2=value2&...&key=appSecret```，key=appSecret固定拼接在参数串末尾，appSecret需替换成应用渠道所分配的appSecret。
>
> 4. 签名，使用制定的算法进行加签获取二进制字节，使用 16进制进行编码Hex.encode得到签名串，然后base64编码。
> 5. 验签，对收到的参数按1-4步骤签名，比对得到的签名串与提交的签名串是否一致。

签名示例：

```json
请求参数：
{
    "appId":"19E179E5DC29C05E65B90CDE57A1C7E5",
    "version": "1",
    "signType": "SM2",
    "signData": "...",
    "encType": "plain",
    "timestamp":1591943910,
    "data": {
    	"user_id":"gt",
    	"face_id":"5ed21b1c262daabe314048f5"
    }
}

密钥：
appSecret="D91CEB11EE62219CD91CEB11EE62219C"
SM2_privateKey="JShsBOJL0RgPAoPttEB1hgtPAvCikOl0V1oTOYL7k5U="

待加签串：
appId=19E179E5DC29C05E65B90CDE57A1C7E5&data={"face_id":"5ed21b1c262daabe314048f5","user_id":"gt"}&encType=plain&signType=SM2&timestamp=1591943910&version=1&key=D91CEB11EE62219CD91CEB11EE62219C

SHA256加签结果：
"2072bd8afb678c03ce9be14202e47b12031aa42a0a8c8593723d7027007ef804"

base64后结果：
"MjA3MmJkOGFmYjY3OGMwM2NlOWJlMTQyMDJlNDdiMTIwMzFhYTQyYTBhOGM4NTkzNzIzZDcwMjcwMDdlZjgwNA=="

SM2加签结果（每次不同）：
"LXgGBQNsXwofSXr+uXYiw0al7MFNNdUl0OyjpxHGKSPjJAr1N5oO6Tq3WL0C8UVX1pmDNH/GZK1Q0h+VvzKiEg=="

```



返回结果

| 参数      | 类型   | 说明                                                         |
| --------- | ------ | ------------------------------------------------------------ |
| code      | string | 接口返回状态代码                                             |
| timestamp | int    | unix时间戳                                                   |
| data      | json   | 成功时返回结果数据。出错时，data.msg返回错误说明。（人脸识别接口有此data.msg字段） |
| msg       | string | 出错时，msg返回错误说明。（特征库接口有此字段。）            |

> 成功时：code为0， msg为"success"，data内容见各接口定义；
>
> 出错时：code返回错误代码，具体定义建各接口说明

返回示例

```json
{
    "code": 0, 
    "timestamp": 1591943910,
    "data": {
       "msg": "success", 
       ...
    }
}
```

全局出错代码

| 编码 | 说明                               |
| ---- | ---------------------------------- |
| 9800 | 无效签名                           |
| 9801 | 签名参数有错误                     |
| 9802 | 调用时间错误，unixtime超出接受范围 |



### 4. 人脸识别

#### （1）人脸定位

> 检测图片中人脸并返回位置

请求URL

> http://127.0.0.1:5001/face/locate

请求方式

> POST

输入参数

| 参数         | 必选 | 类型   | 说明                                              |
| ------------ | ---- | ------ | ------------------------------------------------- |
| image        | 是   | string | base64编码图片数据                                |
| max_face_num | 否   | int    | 最多定位的人脸数量，默认为1，仅检测面积最大的一个 |

请求示例

```json
{
    "image" : "....", 
    "max_face_num" : 5
}
```

返回结果

| 参数                           | 必选 | 类型  | 说明                 |
| ------------------------------ | ---- | ----- | -------------------- |
| face_num                       | 是   | int   | 检测到的图片人脸数量 |
| locations                      | 是   | array | 人脸位置坐标列表     |
| + [ top, right, bottom, left ] | 是   | array | 人脸位置             |

返回示例

```json
{
    "appId": "66A095861BAE55F8735199DBC45D3E8E", 
    "code": 0, 
    "data": {
        "face_num": 2, 
        "locations": [
            [1145, 364, 1335, 607], 
            [764, 391, 947, 641]
        ], 
        "msg": "success", 
        "requestId": "2022090115f1a6e30df9ad3cbf51d5b6e2de3ada"
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1662013631
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |



#### （2）人脸对比

> 比对两张照片中人脸的相似度（1:1），返回相似度分值

请求URL

> http://127.0.0.1:5001/face/verify

请求方式

> POST

输入参数

| 参数   | 必选 | 类型   | 说明               |
| ------ | ---- | ------ | ------------------ |
| image1 | 是   | string | base64编码图片数据 |
| image2 | 是   | string | base64编码图片数据 |

请求示例

```json
{
    "image1": "....", 
    "image2": "....", 
}
```

返回结果

| 参数     | 必选 | 类型    | 说明                           |
| -------- | ---- | ------- | ------------------------------ |
| is_match | 是   | boolean | 是否同一人，TRUE 或 FALSE      |
| score    | 是   | float   | 相似度得分（值越小相似度越高） |

返回示例

```json
{
    "appId": "66A095861BAE55F8735199DBC45D3E8E", 
    "code": 0, 
    "data": {
        "is_match": true, 
        "msg": "success", 
        "requestId": "2022090100b89d6fac623c729f4a8d0911377c35", 
        "score": 0.06643416019927911
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1662014458
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |





#### （3）人脸搜索

> 1:N识别，在指定人脸用户分组中，找到最相似的人脸；当指定user_id时，进行1:1验证

请求URL

> http://127.0.0.1:5001/face/search

请求方式

> POST

输入参数

| 参数        | 必选 | 类型   | 说明                                                      |
| ----------- | ---- | ------ | --------------------------------------------------------- |
| image       | 是   | string | base64编码图片数据                                        |
| group_id    | 否   | string | 在指定分组内搜索，默认为"DEFAULT"分组                     |
| user_id     | 否   | string | 如果提供，则与指定user_id的用户进行比对，相当于1:1验证    |
| mobile_tail | 否   | string | 手机后4位，如果提供，相当于双因素验证（人脸+手机号后4位） |

请求示例

```json
{
    "image": "....", 
    "group_id": "test", 
}
```

返回结果 

(1) 1:N 时 识别时，入参只提供图片，未提供user_id和mobile_tail

| 参数          | 必选 | 类型   | 说明                                  |
| ------------- | ---- | ------ | ------------------------------------- |
| request_id    | 是   | string | 本次查询请求id（用于结果反馈）        |
| user_list     | 是   | string | 匹配到的用户列表                      |
| + user_id     | 是   | string | 用户id                                |
| + name        | 是   | string | 用户姓名，如果注册时有提供            |
| + mobile_tail | 是   | string | 手机号后4位，如果注册了手机号         |
| + score       | 是   | float  | 相似度得分，越小越相似                 |
| + location    | 是   | array  | 人脸的位置 (top, right, bottom, left) |
| fake          | 是   | bool   | bool值为true是假人脸，false是真人脸   |
| fake_score | 是 | float | 为假人脸的概率 |

返回示例

```json
{
    "appId": "66A095861BAE55F8735199DBC45D3E8E", 
    "code": 0, 
    "data": {
        "msg": "success", 
        "requestId": "20220901587d7ae1492e16f25dfb6aee91738d5f", 
        "user_list": [
            {
                "fake": false, 
                "fake_score": 0.09242594,
                "location": [371, 85, 609, 422], 
                "mobile_tail": "4665", 
                "name": "obama", 
                "score": 0.22684336, 
                "user_id": "obama"
            }
        ]
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1662014938
}
```


(2) 双因素验证时，入参还提供了mobile_tail

| 参数          | 必选 | 类型   | 说明                                  |
| ------------- | ---- | ------ | ------------------------------------- |
| request_id    | 是   | string | 本次查询请求id（用于结果反馈）        |
| user_list     | 是   | string | 匹配到的用户列表                      |
| + user_id     | 是   | string | 用户id                                |
| + name        | 是   | string | 用户姓名，如果注册时有提供            |
| + score       | 是   | float  | 相似度得分，越小越相似                 |
| + location    | 是   | array  | 人脸的位置 (top, right, bottom, left) |

返回示例

```json
{
    "appId": "66A095861BAE55F8735199DBC45D3E8E", 
    "code": 0, 
    "data": {
        "msg": "success", 
        "requestId": "202209011f6fc08d0d419489b47b798b316abc4a", 
        "user_list": [
            {
                "location": [724, 611, 1767, 2081], 
                "name": "obama", 
                "score": 0.282884, 
                "user_id": "obama"
            }
        ]
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1662015225
}
```


(3) 1:1 验证时，入参还提供了user_id

| 参数       | 必选 | 类型    | 说明                                  |
| ---------- | ---- | ------- | ------------------------------------- |
| request_id | 是   | string  | 本次查询请求id（用于结果反馈）        |
| is_match   | 是   | boolean | 是否同一人，TRUE 或 FALSE             |
| location   | 是   | array   | 人脸的位置 (top, right, bottom, left) |
| score      | 是   | float   | 相似度得分                            |

返回示例

```json
{
    "appId": "66A095861BAE55F8735199DBC45D3E8E", 
    "code": 0, 
    "data": {
        "is_match": true, 
        "location": [371, 85, 609, 422], 
        "msg": "success", 
        "requestId": "202209014e0f11b9bd15fe1f18e1eb70514c06b5", 
        "score": 0.22684327
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1662015355
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9003 | 参数格式错误                      |
| 9901 | base64编码异常                    |



#### （4）是否有人脸

> 检测图片中是否有人脸，并进行反欺骗检测

请求URL

> http://127.0.0.1:5001/face/check

请求方式

> POST

输入参数

| 参数  | 必选 | 类型   | 说明               |
| ----- | ---- | ------ | ------------------ |
| image | 是   | string | base64编码图片数据 |

请求示例

```json
{
    "image" : "....", 
}
```

返回结果

| 参数     | 必选 | 类型 | 说明                          |
| -------- | ---- | ---- | ----------------------------- |
| has_face | 是   | bool | 检测到人脸返回true，否则false |
| fake      | 否  | bool | bool值为true是假人脸，false是真人脸 |
| fake_score | 否 | float | 为假人脸的概率 |


返回示例

```json
{
    "appId": "66A095861BAE55F8735199DBC45D3E8E", 
    "code": 0, 
    "data": {
        "fake": false, 
        "fake_score": 0.09242594,
        "has_face": true, 
        "msg": "success", 
        "requestId": "202209019ccc9d370815f986612a3e81281bcc3f"
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1662015758
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |



### 5. 特征库管理

#### （1）特征库结构

```
|- 特征库
   |- 分组一（group_id）
      |- 用户01（user_id）
         |- 人脸（face_id）
      |- 用户02（user_d）
         |- 人脸（face_id）
         |- 人脸（face_id）
         ....
       ....
   |- 分组二（group_id）
   |- 分组三（group_id）
   ....
```



#### （2）人脸注册

> 向特征库中添加人脸，当user_id在库中已经存在时，新注册的图片会追加到该user_id下

请求URL

> http://127.0.0.1:5001/facedb/face/reg

请求方式

> POST

输入参数

| 参数     | 必选 | 类型   | 说明                                                         |
| -------- | ---- | ------ | ------------------------------------------------------------ |
| image    | 是   | string | base64编码图片数据                                           |
| group_id | 否   | string | 用户分组id，默认为"DEFAULT"分组                              |
| user_id  | 是   | string | 用户id（由数字、字母、下划线组成），必须唯一，建议使用身份证号码 |
| name     | 是   | string | 用户姓名                                                     |
| mobile   | 否   | string | 手机号码                                                     |
| gender   | 否   | int    | 性别（ 1 男， 2 女）                                         |
| age      | 否   | int    | 年龄                                                         |

请求示例

```json
{
    "image": "...",
    "group_id": "group1",
    "user_id": "gt2",
    "name": "gt",
    "gender": 1,
    "age": 18
}
```

返回结果

| 参数    | 必选   | 类型             | 说明 |
| ------- | ------ | ---------------- | ---- |
| face_id | string | 人脸特征唯一标识 |      |

返回示例

```json
{
    "code": 0, 
    "msg": "success",
    "data": {
        "face_id": "5ed7725e0d72875f136cdbbe"
    }, 
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |
| 9003 | 未定位到人脸                      |
| 9004 | 用户组不存在                      |
| 9005 | user_id已存在                     |



#### （3）人脸更新

> 更新特征库中指定用户下的人脸信息，使用新图替换库中该user_id下所有图片，若user_id不存在则报错

请求URL

> http://127.0.0.1:5001/facedb/face/update

请求方式

> POST

输入参数

| 参数     | 必选 | 类型   | 说明                                 |
| -------- | ---- | ------ | ------------------------------------ |
| image    | 否   | string | base64编码图片数据，不提供则不修改   |
| group_id | 否   | string | 用户分组id，默认为"DEFAULT"分组      |
| user_id  | 是   | string | 用户id                               |
| mobile   | 否   | string | 手机号码，不提供则不修改             |
| name     | 否   | string | 用户姓名，不提供则不修改             |
| gender   | 否   | int    | 性别（ 1 男， 2 女），不提供则不修改 |
| age      | 否   | int    | 年龄，不提供则不修改                 |

请求示例

```json
{
    "image": "...",
    "group_id": "group1",
    "user_id": "gt2",
}
```

返回结果

| 参数    | 必选   | 类型             | 说明                                |
| ------- | ------ | ---------------- | ----------------------------------- |
| face_id | string | 人脸特征唯一标识 | 只更新用户信息，不更新图片，则返回0 |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "face_id": "5ed77468b643e4aa5b27cf49"
    }
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |
| 9003 | 未定位到人脸                      |
| 9004 | user_id不存在                     |



#### （4）人脸删除

> 删除指定用户的某张人脸特征数据

请求URL

> http://127.0.0.1:5001/facedb/face/remove

请求方式

> POST

输入参数

| 参数     | 必选 | 类型   | 说明                            |
| -------- | ---- | ------ | ------------------------------- |
| group_id | 否   | string | 用户分组id，默认为"DEFAULT"分组 |
| user_id  | 是   | string | 用户id                          |
| face_id  | 是   | string | 人脸特征标识                    |

请求示例

```shell
curl -X POST --data "{"group_id":"group1","user_id":"gt","face_id":"5ed21b1c262daabe314048f5"}" http://127.0.0.1:5001/facedb/face/remove
```

返回结果

| 参数       | 必选 | 类型   | 说明                          |
| ---------- | ---- | ------ | ----------------------------- |
| type | 是   | string | 成功返回SUCCESS |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "type": "SUCCESS"
    }
}
```

出错代码

| 编码 | 说明                    |
| ---- | ----------------------- |
| 9001 | 缺少参数                |
| 9002 | user_id不存在           |
| 9003 | 用户没有face_id人脸数据 |



#### （5）用户信息查询

> 查询特征库中某个用户的详细信息

请求URL

> http://127.0.0.1:5001/facedb/user/info

请求方式

> POST

输入参数

| 参数     | 必选 | 类型   | 说明                            |
| -------- | ---- | ------ | ------------------------------- |
| group_id | 否   | string | 用户分组id，默认为"DEFAULT"分组 |
| user_id  | 是   | string | 用户id                          |

请求示例

```shell
curl -X POST --data "{"group_id":"test","user_id":"obama"}" http://127.0.0.1:5001/facedb/user/info
```

返回结果

| 参数      | 必选 | 类型   | 说明                        |
| --------- | ---- | ------ | --------------------------- |
| group_id  | 是   | string | 分组id                      |
| user_id   | 是   | string | 用户id                      |
| name      | 是   | string | 用户姓名                    |
| mobile    | 否   | string | 手机号码，如果注册时有提供  |
| memo      | 否   | int    | 性别，返回0表示无此数据     |
| age       | 否   | int    | 年龄，返回0表示无此数据     |
| image_num | 是   | int    | 此user_id下已注册的照片数量 |
| ctime     | 是   | string | 注册时间                    |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "image_num": 1, 
        "group_id": "test", 
        "user_id": "obama", 
        "name": "obama", 
        "gender": 0,
        "age": 18,
        "ctime": "2020-05-30 16:36:40", 
        "mobile": ""
    }
}
```

出错代码

| 编码 | 说明          |
| ---- | ------------- |
| 9001 | 缺少参数      |
| 9002 | user_id不存在 |



#### （6）获取用户人脸列表

> 获取某个用户的全部人脸列表

请求URL

> http://127.0.0.1:5001/facedb/user/face_list

请求方式

> POST

输入参数

| 参数     | 必选 | 类型   | 说明                            |
| -------- | ---- | ------ | ------------------------------- |
| group_id | 否   | string | 用户分组id，默认为"DEFAULT"分组 |
| user_id  | 是   | string | 用户id                          |

请求示例

```shell
curl -X POST --data "{"group_id":"test","user_id":"gt"}" http://127.0.0.1:5001/facedb/user/face_list
```

返回结果

| 参数      | 必选 | 类型  | 说明                  |
| --------- | ---- | ----- | --------------------- |
| face_list | 是   | array | 人脸特征face_id的列表 |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "face_list": [
            "5ed21b1c262daabe314048f5", 
            "5ed21b1d262daabe314048f6"
        ]
    }
}
```

出错代码

| 编码 | 说明          |
| ---- | ------------- |
| 9001 | 缺少参数      |
| 9002 | user_id不存在 |



#### （7）获取用户列表

> 查询指定用户组中的用户列表

请求URL

> http://127.0.0.1:5001/facedb/user/list

请求方式

> POST

输入参数

| 参数     | 必选 | 类型   | 说明                            |
| -------- | ---- | ------ | ------------------------------- |
| group_id | 否   | string | 用户分组id，默认为"DEFAULT"分组 |
| start    | 否   | int    | 起始位置，默认为0               |
| length   | 否   | int    | 返回数量，默认100，最大1000     |

请求示例

```shell
curl -X POST --data "{"group_id":"test"}" http://127.0.0.1:5001/facedb/user/list
```

返回结果

| 参数      | 必选 | 类型  | 说明            |
| --------- | ---- | ----- | --------------- |
| user_list | 是   | array | 用户user_id列表 |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "user_list": [
            "biden", 
            "obama", 
            "alex_lacamoire", 
            "gt", 
            "zhiqiang", 
            "obama2", 
            "kit_harington", 
            "obama1", 
            "rose_leslie"
        ]
    }
}
```

出错代码

| 编码 | 说明 |
| ---- | ---- |
|      |      |



#### （8）复制用户

> 将指定用户复制到另外的人脸组

请求URL

> http://127.0.0.1:5001/facedb/user/copy

请求方式

> POST

输入参数

| 参数         | 必选 | 类型   | 说明                 |
| ------------ | ---- | ------ | -------------------- |
| user_id      | 是   | string | 用户id               |
| src_group_id | 是   | string | 从指定组里复制信息   |
| dst_group_id | 是   | string | 需要添加用户的分组id |

请求示例

```shell
curl -X POST --data "{"user_id":"gt","src_group_id":"test","dst_group_id":"group1"}" http://127.0.0.1:5001/facedb/user/copy
```

返回结果

| 参数 | 必选 | 类型   | 说明            |
| ---- | ---- | ------ | --------------- |
| type | 是   | string | 成功返回SUCCESS |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "type": "SUCCESS"
    }
}
```

出错代码

| 编码 | 说明                      |
| ---- | ------------------------- |
| 9001 | 缺少参数                  |
| 9002 | user_id不存在             |
| 9003 | user_id在目的用户组已存在 |



#### （9）删除用户

> 删除指定用户

请求URL

> http://127.0.0.1:5001/facedb/user/remove

请求方式

> POST

输入参数

| 参数     | 必选 | 类型   | 说明                            |
| -------- | ---- | ------ | ------------------------------- |
| group_id | 否   | string | 用户分组id，默认为"DEFAULT"分组 |
| user_id  | 是   | string | 用户id                          |

请求示例

```shell
curl -X POST --data "{"user_id":"gt","group_id":"group1"}" http://127.0.0.1:5001/facedb/user/remove
```

返回结果

| 参数 | 必选 | 类型   | 说明            |
| ---- | ---- | ------ | --------------- |
| type | 是   | string | 成功返回SUCCESS |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "type": "SUCCESS"
    }
}
```

出错代码

| 编码 | 说明          |
| ---- | ------------- |
| 9001 | 缺少参数      |
| 9002 | user_id不存在 |



#### （10）创建用户组

> 创建一个新的用户组，如果用户组已存在 则返回错误

请求URL

> http://127.0.0.1:5001/facedb/group/new

请求方式

> POST

输入参数

| 参数     | 必选 | 类型   | 说明                                               |
| -------- | ---- | ------ | -------------------------------------------------- |
| group_id | 是   | string | 用户组id，标识一组用户（由数字、字母、下划线组成） |

> ```DEFAULT``` 和 ```__RECYCLE__``` 为默认组，不可再创建

请求示例

```shell
curl -X POST --data "{"group_id":"group1"}" http://127.0.0.1:5001/facedb/group/new
```

返回结果

| 参数 | 必选 | 类型   | 说明            |
| ---- | ---- | ------ | --------------- |
| type | 是   | string | 成功返回SUCCESS |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "type": "SUCCESS"
    }
}
```

出错代码

| 编码 | 说明           |
| ---- | -------------- |
| 9001 | 缺少参数       |
| 9002 | group_id已存在 |



#### （12）删除用户组

> 删除指定用户组，如果组不存在 则返回错误。**组内用户将一同删除，需谨慎操作！**

请求URL

> http://127.0.0.1:5001/facedb/group/remove

请求方式

> POST

输入参数

| 参数     | 必选 | 类型   | 说明       |
| -------- | ---- | ------ | ---------- |
| group_id | 是   | string | 用户分组id |

请求示例

```shell
curl -X POST --data "{"group_id":"group1"}" http://127.0.0.1:5001/facedb/group/remove
```

返回结果

| 参数 | 必选 | 类型   | 说明            |
| ---- | ---- | ------ | --------------- |
| type | 是   | string | 成功返回SUCCESS |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "type": "SUCCESS"
    }
}
```

出错代码

| 编码 | 说明           |
| ---- | -------------- |
| 9001 | 缺少参数       |
| 9002 | group_id不存在 |



#### （13）获取用户组列表

> 查询特征库中用户组的列表

请求URL

> http://127.0.0.1:5001/facedb/group/list

请求方式

> POST

输入参数

| 参数   | 必选 | 类型 | 说明                        |
| ------ | ---- | ---- | --------------------------- |
| start  | 否   | int  | 起始位置，默认为0           |
| length | 否   | int  | 返回数量，默认100，最大1000 |

请求示例

```shell
curl -X POST --data "{}" http://127.0.0.1:5001/facedb/group/list
```

返回结果

| 参数       | 必选 | 类型  | 说明               |
| ---------- | ---- | ----- | ------------------ |
| group_list | 是   | array | 用户组group_id列表 |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
        "user_list": [
            "test", 
            "train3", 
            "train2", 
            "test2", 
            "test3"
        ]
    }
}
```

出错代码

| 编码 | 说明 |
| ---- | ---- |
|      |      |



#### （14）识别结果反馈

> 反馈最近搜索识别的结果是否正确，用于帮助系统提高识别准确率

请求URL

> http://127.0.0.1:5001/facedb/feedback

请求方式

> POST

输入参数

| 参数       | 必选 | 类型   | 说明                       |
| ---------- | ---- | ------ | -------------------------- |
| request_id | 是   | string | search接口返回的request_id |
| is_correct | 是   | int    | 1 结果正确； 0 结果不正确  |

请求示例

```json
{
    "request_id": "4316142b3b7165b11f9ae91f8f4fb5c4", 
    "is_correct": 1, 
}
```

返回示例

```json
{ 
    "code": 0, 
    "msg": "success"
}
```

出错代码

| 编码 | 说明             |
| ---- | ---------------- |
| 9001 | 缺少参数         |
| 9002 | request_id不存在 |

