
## api 文档

### 1. 全局接口定义

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

| 参数      | 类型    | 说明                                                         | 示例  |
| --------- | ------- | ------------------------------------------------------------ | ----- |
| appId     | string  | 应用渠道编号                                                 |       |
| code      | string  | 接口返回状态代码                                             |       |
| signType  | string  | 签名算法，plain： 不用签名，SM2：使用SM2算法                 | plain |
| encType   | string  | 接口数据加密算法，目前不加密                                 | plain |
| success   | boolean | 成功与否                                                     |       |
| timestamp | int     | unix时间戳                                                   |       |
| data      | json    | 成功时返回结果数据；出错时，data.msg返回错误说明。详见具体接口 |       |

> 成功时：code为0， success为True，data内容见各接口定义；
>
> 出错时：code返回错误代码，具体定义见各接口说明

返回示例

```json
{
    "appId": "19E179E5DC29C05E65B90CDE57A1C7E5", 
    "code": 0, 
    "signType": "plain",
    "encType": "plain",
    "success": true,
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



### 2. 识别文本命令

请求URL

> http://127.0.0.1:5000/talk2ui/text2order

请求方式

> POST

输入参数

| 参数  | 必选 | 类型   | 说明               |
| ----- | ---- | ------ | --------------- |
| text  | 是   | string | 输入的文本        |

请求示例

```json
{
    "text" : "我想查下体检报告"
}
```

返回结果

| 参数       | 必选 | 类型   | 说明                 |
| ---------- | ---- | ------ | -------------------- |
| result     | 是   | string | 识别结果，具体见示例 |

返回示例

```json
{
    "appId": "66A095861BAE55F8735199DBC45D3E8E",
    "code": 0,
    "data": {
        "result": {
            "object": "体检",
            "order": "Check",
            "text": "我想查下体检报告"
        },
        "msg": "success",
        "requestId": "20231127205205120059ba735757e1e3a6b4113e74277e"
    },
    "encType": "plain",
    "signType": "plain",
    "success": true,
    "timestamp": 1701089525
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9101 | 缺少参数                          |



### 3. 识别语音命令

请求URL

> http://127.0.0.1:5000/talk2ui/wav2order

请求方式

> POST

输入参数

| 参数  | 必选 | 类型   | 说明               |
| ----- | ---- | ------ | --------------- |
| wav_data  | 是   | string | 语音数据base64编码串 |

> 语音数据支持 wav 和 mp3 格式

请求示例

```json
{
    "wav_data": "UklGRnBjAQBXQVZFZm10IBAAAAAB ..."
}
```

返回结果

| 参数       | 必选 | 类型   | 说明                 |
| ---------- | ---- | ------ | -------------------- |
| result     | 是   | string | 识别结果，具体见示例 |

返回示例

```json
{
    "appId": "66A095861BAE55F8735199DBC45D3E8E",
    "code": 0,
    "data": {
        "result": {
            "object": "内科",
            "order": "Department",
            "text": "我想预约内科门诊"
        },
        "msg": "success",
        "requestId": "20231127205531dc8aad3bbfb25d3c5a1d32f982f671f9"
    },
    "encType": "plain",
    "signType": "plain",
    "success": true,
    "timestamp": 1701089731
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9101 | 缺少参数                          |
| 9002 | 语音数据太大，base64数据不大于3MB |
| 9901 | base64编码异常                    |
