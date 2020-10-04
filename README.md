# oa
## 树图

|--- **/se-token** //一次性校验码

|--- **/se-login** //登陆

|

|--- **/se-course** //查询课程

|--- **/se-blog** //搜索博客

|--- **/se-uploadimage** //上传图片

|--- **/se-uploadfile** //上传文件

|--- **/se-blog-exec** //博客操作

|--- **/se-article** //博客单篇查询

|--- **/se-download** //文件下载

## 使用

### **/se-token** 

> **用于获取记录状态的值**
>
> 响应头的authentication的值为sessionid，用于记录用户的会话
>
> 返回json，其中的token用于登陆校验，一次性数据



### **/se-login**

> **用于登陆检测**
>
> 用**post**方法发送请求，请求主体需要username，password，token等数据
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值
>
> 返回sessionid和status，其中status用于检验登陆是否成功 



### **/se-course**

>**用于获取课程**
>
>用**post**方法发送请求，请求主体需要type,limit等数据
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
              "sessionid": "XXXXXXXXXXXXX",
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



### **/se-blog**

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
    keyword:"",
    "bloglist": [
        	"number":{
            	"id":"",
            	"keyword":"",
            	"title":"",
            	"content":"",//这个字段为空
        		"summary":"",
            	"author":"",
            	"record":"",
            	"publicstatus":"",
            	"publictime":"",
        		"isshow":""
        	}
        	...
		]
       ...
}
```



### **/se-uploadimage**

> **用于上传图片**
>
> 用**post**方法发送请求，请求主体为单张图片文件，允许的类型有.png , .jpg , .jpeg
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值

```json
{
	"url":"",
	"status":"OK/not login sessionId"
}
```

> 需要注意的是：返回的url格式汝下
>
> /static/0e3c4adb-b74d-496b-9f21-0557c1e89777.png



### /se-uploadfile

> **用于上传文件压缩包**
>
> > **压缩文件**
> >
> > 用**post**方法发送请求，允许的类型有**.zip** , **.7z** , **.rar** ,  **.gz** , **.tar** 。
> >
> > 字段要求：**name = “file“**
>
> > **封面**
> >
> > 用**post**方法发送请求，允许的类型有**.png , .jpg , .jpeg** 。
> >
> > 字段要求：**name = “cover“**
>
> >**简介**
> >
> >字段要求：**name = ”brief“**
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值

```json
{
	"sessionid":"",
	"status":"OK/not login sessionId"
}
```





### **/se-blog-exec**

> **用于新增博客或者修改博客**
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值
>
> 用**post**方法发送请求，请求主体为json数据，要求格式如下

```json
{
    "id":"",//文章id，如果是新建则不需要
    "keyword":"",//关键词，各个关键词之间用 ; 分开
    "title":"",
    "content":"",
    "status":"",//决定文章是处于发布状态或者草稿状态，发布状态为 release
    "command":"add/update" //控制操作
}
```

> 返回的数据如下

```json
{
    "status":"not login sessionID/OK",
}
```



### **/se-article**

> **用于查找单篇文章的所有数据**
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值
>
> 用**post**方法发送请求，请求主体为表单数据 **id**
>
> 返回单篇文章的数据，如下

```json
{		
				"id":"",
            	"keyword":"",
            	"title":"",
            	"content":"",//这个字段为空
        		"summary":"",
            	"author":"",
            	"record":"",
            	"publicstatus":"",
            	"publictime":"",
        		"isshow":""
}
```

### /se-logout

> **用于登出账号**
>
> 请求头里需要authentication字段，作为代替cookie，将从 **/se-token**获取的sessionid作为值
>
> 操作失败时返回**500状态码**

### /se-download

> 用于下载文件
>
> 字段要求：**在链接后面带上名为url的参数**，如下
>
> ```
> /se-download?url=b8068f27-ef84-4de9-be9c-ce07d11ef8ba.7z
> ```
>
> 

