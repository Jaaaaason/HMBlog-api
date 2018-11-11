# HMBlog
一个博客 api 系统，遵循 RESTful API 规范

### 已完成的 api
Method |            URL Path         | Description
------ | --------------------------- | ------------------------
POST   | /admin/login                | 登录后台，获取 JWT Token
GET    | /categories                 | 以访客身份获取所有分类
GET    | /categories/:id             | 以访客身份获取某个分类
GET    | /admin/categories           | 以后台用户身份获取所有分类
GET    | /admin/categories/:id       | 以后台用户身份获取某个分类
POST   | /admin/categories/:id       | 以后台用户身份创建一个新的分类
PUT    | /admin/categories/:id       | 以后台用户身份修改某个分类
PATCH  | /admin/categories/:id       | 以后台用户身份修改某个分类
DELETE | /admin/categories/:id       | 以后台用户身份删除某个分类
GET    | /posts                      | 以访客身份获取所有博文
GET    | /posts/:id                  | 以访客身份获取某个博文
GET    | /categories/:id/posts       | 以访客身份获取某个分类下所有博文
GET    | /admin/posts                | 以后台用户身份获取所有博文
GET    | /admin/posts/:id            | 以后台用户身份获取某个博文
GET    | /admin/categories/:id/posts | 以后台用户身份获取某个分类下的所有博文
POST   | /admin/posts/:id            | 以后台用户身份创建一个新的博文
PUT    | /admin/posts/:id            | 以后台用户身份修改某个博文
PATCH  | /admin/posts/:id            | 以后台用户身份修改某个博文
DELETE | /admin/posts/:id            | 以后台用户身份删除某个博文

详细的 api 文档请移步 [HMBlog Api Doc](https://app.swaggerhub.com/apis-docs/Jaaaaason/hmblog/1.0.0)
