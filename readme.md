# SHUCourseProxy

访问[上海大学选课网站](http://xk.autoisp.shu.edu.cn)及与其使用同一套Auth系统的网站时使用的的代理服务。

## 使用方法

假设服务部署在： http://localhost:8086

则可以向 http://localhost:8086/login 发送下面的JSON进行模拟登录：

```json
{
	"from_url": "http://xk.autoisp.shu.edu.cn:8080/",
	"username": "你的学生证",
	"password": "你的密码"
}
```

后台会模拟登录，然后储存这个学生证号在这个网站上的Cookie。

会返回一个JWT，可以使用这个JWT（通过设置Authorization头）来通过这个Proxy服务来请求登录后可以拿到的页面。

如，向 http://localhost:8086/get 发送Post请求：

header:
```
Authorization="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdHVkZW50SWQiOiIxNzEyMDIzOCJ9.ri5Q1zgGPD1ohGF72Tx4Kdf3iLndtNuPVWcOI9X2YJ4"
```

body:
```json
{
	"url": "http://xk.autoisp.shu.edu.cn:8080"
}
```
即会返回以该JWT头所代表的用户身份去Get http://xk.autoisp.shu.edu.cn:8080 得到的结果。
