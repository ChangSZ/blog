# blog

 [点击访问博客网站](http://water-melon.top)



## 主要功能

1. 文章发布和修改
2. 文章回收站和撤回
3. 文章分类
4. 文章标签
5. 文章支持markdown
6. 网站静态文件可自由配置`本地`或 `CDN`
7. 可上传图片至服务器,同时支持上传至 `七牛`
8. 自由添加友链和管理友链顺序
9.  采用`gitalk`功能作为评论系统,界面优美且方便其他用户留言和通知
10. 定时备份数据和静态资源并发送至指定邮箱
11. 日志支持`trace.id`追踪
12. 网站信息自由设置

</br>

## 技术栈

主要代码是 `golang`+`vue`+`HTML`+`CSS`+`MySQL`
>   - [博客管理后台](https://github.com/ChangSZ/blog-admin)是基于`vue`的`iview`UI组件开发的, 
>   - 前台是基于`HTML+CSS`展示[静态页面](http://water-melon.top)
>   - 缓存用的`redis`
>   - 数据库用的是 `MySQL`
>   - 配置文件用的 `yaml`

</br>

## 运行方法
### 前置配置
- 需在mysql中创建`blog`数据库，并将`/common/sql.sql`导入
- 需将`env.example.yaml`配置文件拷贝一份，dev及本地测试环境命名为`env.dev.yaml`，prod环境命名为`env.prod.yaml`，并将其中的db、redis配置完全

### 启动方式1
```golang
   go run main.go
   // 然后访问http://127.0.0.1:8081
```

### 启动方式2：Docker启动
```bash
   docker build -f .\Dockerfile -t blog:v0.0.2 .

   # 可导出镜像
   docker save -o blog.tar blog:v0.0.2

   # 通过ftp上传至服务器，然后执行导入
   docker load -i blog.tar

   # 运行
   docker run -idt --name blog --network host blog:v0.0.2

   # 然后访问http://服务器IP:8081
```
</br>

 ## nginx部署
 参考nginx.conf配置及说明
