# oa
## 验证
> **/se-token** 
>
>响应头的authentication的值为sessionid，用于记录用户的会话
>
>返回json，其中的token用于登陆校验，一次性数据

> **/se-login**
>
> 用**post**方法发送请求，请求主体需要username，password，token等数据
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值
>
> 返回sessionid和status，其中status用于检验登陆是否成功 

> **/se-course**
> 用**post**方法发送请求，请求主体需要type数据
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值
> 
> 返回json，格式如下
```json
{
  "userinfo": {
        "sesionid": "XXXXXXXXXXXXX",
        "status": "success/fail"
  },
  "list": [
      {
        "id": "",
        "type" : "",
        "url" : "",
        "title" : "",
        "content" : ""
      }
      ...
  ]
}
```
>
>
>**/se-course-limit**
>
> 用**post**方法发送请求，请求主体需要type,page等数据
>
> **需要注意的是页面的下标是从0开始的**
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值
>
> 返回json，格式如下
```json
{
  "queryinfo": {
        "length": "x", //总长度
        "PageNum" : "x", //当前页面
        "error": "" //错误信息
  },
  "list": [
      {
        "id": "",
        "type" : "",
        "url" : "",
        "title" : "",
        "content" : ""
      }
      ...
  ]
}
```
