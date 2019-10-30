# file-explorer

一个简单的http api文件服务，包含文件的上传、删除、修改名称、移动、下载、查看，文件夹的创建、删除、修改名称、查看

    文件命令都需要附带/user/login返回的token，可附带的位置为query中的token参数或http请求头中的X-REQUEST-TOKEN或Authorization的bearer中

## 命令如下：
* /file
  * /mkdir   
    * 作用：创建文件夹
    * 类型：PUT
    * 参数：
        * name： 文件夹名称(必需)
        * directory_id： 父目录ID(可空，空为根目录)
  * /upload  
    * 作用：上传文件
    * 类型：PUT
    * 参数：
        * file： 上传文件(必需)
        * directory_id： 父目录ID(可空，空为根目录)
  * /:id      
    * 作用：删除文件
    * 类型：DELETE
    * 参数：
        * id[url]： 文件ID(必需)
  * /:id/list         
    * 作用：查看文件列表
    * 类型：GET
    * 参数：
        * id[url]： 文件夹ID(必需)   
  * /:id        
    * 作用：下载文件
    * 类型：GET
    * 参数：
        * id[url]： 文件ID(必需)
  * /:id/info         
    * 作用：获取文件信息
    * 类型：GET
    * 参数：
        * id[url]： 文件ID(必需)    
  * /:id/rename       修改文件名称
    * 作用：修改文件名称
    * 类型：POST
    * 参数：
        * id[url]: 文件ID(必需)
        * new_name： 新文件名称(必需)
  * /:id/move       修改文件位置
    * 作用：修改文件名称
    * 类型：POST
    * 参数：
        * id[url]： 文件ID(必需)
        * directory_id：新目录ID(必需)
* /user
  * /email_code
    * 作用：创建邮箱验证码
    * 类型：POST
    * 参数：
        * email: 电子邮箱(必需)
  * /register
    * 作用：用户注册
    * 类型：POST
    * 参数：
        * username： 用户名(必需)
        * password：密码(必需)
        * email： 电子邮箱(必需)
        * code: 邮箱验证码(必需)
  * /login
    * 作用：用户登录
    * 类型：POST
    * 参数：
        * username： 用户名(必需)
        * password：密码(必需)
  * /password/reset_code
    * 作用：发送修改密码验证码
    * 类型：POST
    * 参数：
        * email: 电子邮箱(必需)
  * /password/reset
    * 作用：修改密码
    * 类型：POST
    * 参数：
        * code: 密码验证码(必需)
        * password: 新密码(必需)
        * email: 电子邮箱(必需)
  * /current
    * 作用：获取当前用户
    * 类型：GET
    * 参数：
        * 无
