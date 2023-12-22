
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



### 2. 银行卡号识别

> 检测图片银行卡卡号。
>
> 注意：银行卡图片请按正常文字方向上传，倒置、横置都会影响识别结果。目前不支持竖版银行卡识别。

请求URL

> http://127.0.0.1:5000/ocr/bankcard

请求方式

> POST

输入参数

| 参数  | 必选 | 类型   | 说明               |
| ----- | ---- | ------ | ------------------ |
| image | 是   | string | base64编码图片数据 |

请求示例

```json
{
    "image" : "...."
}
```

返回结果

| 参数        | 必选 | 类型   | 说明             |
| ----------- | ---- | ------ | ---------------- |
| card_number | 是   | string | 检测到的银行卡号 |
| requestId  | 是   | string | 此次请求id       |

返回示例

```json
{
    "data": {
        "card_number": "6217852700001503463",
        "requestId": "4316142b3b7165b11f9ae91f8f4fb5c4",
        "msg": "success"
    }
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |



### 3. 身份证信息识别

> 检测图片身份证信息。
>
> 注意：身份证图片请按正常文字方向上传，倒置、横置都会影响识别结果。

请求URL

> http://127.0.0.1:5000/ocr/idcard

请求方式

> POST

输入参数

| 参数  | 必选 | 类型   | 说明               |
| ----- | ---- | ------ | ------------------ |
| image | 是   | string | base64编码图片数据 |

请求示例

```json
{
    "image" : "...."
}
```

返回结果

| 参数       | 必选 | 类型   | 说明                           |
| ---------- | ---- | ------ | ------------------------------ |
| card_info  | 是   | json   | 身份证信息（未识别项返回空串） |
| requestId | 是   | string | 此次请求id                     |

返回示例

```json
{
    "data": {
        "msg": "success", 
        "card_info": {
            "nation": "汉", /* 民族 */
            "birth": "1990年1月1日", /* 出生年月 */
            "sex": "男", 
            "name": "张三", 
            "idnum": "...",   /* 身份证号码 */
            "addr": "福建省厦门市....",  /* 地址 */
        }, 
        "requestId": "200911171438165173a896444566c5cc557eef817cf0"
    }
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |



### 4. 普通文本识别

> 检测图片中文字信息。
>
> 注意：图片请按正常文字方向上传，倒置、横置都会影响识别结果。

请求URL

> http://127.0.0.1:5000/ocr/text

请求方式

> POST

输入参数

| 参数  | 必选 | 类型   | 说明               |
| ----- | ---- | ------ | ------------------ |
| image | 是   | string | base64编码图片数据 |

请求示例

```json
{
    "image" : "...."
}
```

返回结果

| 参数       | 必选 | 类型   | 说明                           |
| ---------- | ---- | ------ | ------------------------------ |
| text  | 是   | json   | 文字框坐标及文字内容 |
| requestId | 是   | string | 此次请求id                     |

返回示例

```json
{
    "appId": "66A095861BAE55F8735199DBC45D3E8E", 
    "code": 0, 
    "success": true, 
    "signType": "plain", 
    "encType": "plain", 
    "data": {
        "msg": "success", 
        "text": [
            {
                "pos": [221, 13, 1294, 13, 219, 48, 1294, 51], 
                "text": "{1987年股市崩盘可能是现代从众效应引起的第一场危机。金融业普及了动"
            }, 
            {
                "pos": [145, 65, 1294, 65, 147, 109, 1294, 109], 
                "text": "态投资组合保险方法——涉及保护投资者避免投资组合亏损的方法。许多机构提"
            }, 
            {
                "pos": [147, 128, 1294, 130, 145, 167, 1294, 170], 
                "text": "{供这种安全的交易方式:在市场向下时，卖空市场;在市场走高时，做多市场。"
            }, 
            {
                "pos": [145, 188, 998, 188, 147, 232, 998, 232], 
                "text": "当然，上述方法只有市场一小部分人采用时，才行之有效。"
            }, 
        ], 
        "requestId": "210604164252dec951468484ed3048a4b38b5fc37b16"}, 
    "timestamp": 1622796172
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |



### 5. 就诊卡号识别

> 检测图片上号码信息：格式 No. 123456.... 号码目前只限数字形式
>
> 注意：卡片图片请按正常文字方向上传，倒置、横置、倾斜都会影响识别结果。

请求URL

> http://127.0.0.1:5000/ocr/cardnum

请求方式

> POST

输入参数

| 参数  | 必选 | 类型   | 说明               |
| ----- | ---- | ------ | ------------------ |
| image | 是   | string | base64编码图片数据 |

请求示例

```json
{
    "image" : "...."
}
```

返回结果

| 参数       | 必选 | 类型   | 说明                           |
| ---------- | ---- | ------ | ------------------------------ |
| card_info  | 是   | json   | 身份证信息（未识别项返回空串） |
| requestId | 是   | string | 此次请求id                     |

返回示例

```json
{
    "data": {
        "msg": "success", 
        "card_info": {
            "cardnum": "190834120",
        }, 
        "requestId": "200911171438165173a896444566c5cc557eef817cf0"
    }
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |

