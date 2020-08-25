# oa
## 验证
> **/se-token** 
>
> **用于获取记录状态的值**
>
> 响应头的authentication的值为sessionid，用于记录用户的会话
>
> 返回json，其中的token用于登陆校验，一次性数据



> **/se-login**
>
> **用于登陆检测**
>
> 用**post**方法发送请求，请求主体需要username，password，token等数据
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值
>
> 返回sessionid和status，其中status用于检验登陆是否成功 



>**/se-course**
>
>**用于获取课程**
>
> 用**post**方法发送请求，请求主体需要type,limit等数据
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
  "infolist" : {
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
}
```



> **/se-blog**
>
> **用于搜索栏查找，获取已发布的博客**
>
> 用**post**方法发送请求，请求主体需要search,limit等数据
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
    "bloglist": {
        "关键字" : [
        	0: {
            	"id":,
            	"keyword":,
            	"title":,
            	"content":,
            	"author":,
            	"record":,
            	"publicstatus":,
            	"publictime":,
        	},
        	...
    	], 
        ...
    }   
}
```



> **/se-upload**
>
> **用于上传图片**
>
> 用**post**方法发送请求，请求主体为单张图片文件，允许的类型有.png , .jpg , .jpeg
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值

```shell
{
	"url" : ,
	"status" : ,
}
```

> 需要注意的是：返回的url格式汝下
>
> /static/0e3c4adb-b74d-496b-9f21-0557c1e89777.png