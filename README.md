# blog
个人博客项目


## 项目按照 convee 大哥的 学习了一边 go 后端编程
- - -
### 附上原文链接 [convee](github.com/convee/goblog)

### 技术栈

* 前端框架：[Bootstrap v3.3.7](http://getbootstrap.com)
* 语言：[go](https://go.dev/)
* 网络库：标准库 net/http
* 配置文件解析库 [Viper](https://github.com/spf13/viper)
* 日志库：[zap](https://github.com/uber-go/zap)
* 搜索引擎：[elasticsearch](https://github.com/olivere/elastic/v7)
* 数据库：[mysql](https://github.com/go-sql-driver/mysql)
* 缓存：[redis](https://github.com/go-redis/redis)
* 文件存储：阿里云 oss、cdn
* markdown 编辑器：[markdown editor](https://github.com/pandao/editor.md)
* pprof 性能调优
* 包管理工具 [Go Modules](https://github.com/golang/go/wiki/Modules)
* 评论插件：[gitalk](https://github.com/gitalk/gitalk)
* 后台登录：cookie
* 使用 make 来管理 Go 工程
* 使用 shell(startup.sh) 脚本来管理进程
* 使用 YAML 文件进行多环境配置
* 优雅退出
* Http 请求 panic 异常捕获
* 错误信息钉钉预警

### 目录结构

```shell
├── Makefile                     # 项目管理文件
├── conf                         # 配置文件统一存放目录
├── internal                     # 业务目录
│   ├── handler                  # http 接口
│   ├── pkg                      # 内部应用程序代码
│   └── routers                  # 业务路由
├── logs                         # 存放日志的目录
├── static                       # 存放静态文件的目录
├── tpl                          # 存放模板的目录
├── main.go                      # 项目入口文件
├── pkg                          # 公共的 package
├── tests                        # 单元测试
└── startup.sh                   # 启动脚本
```

### 功能模块

#### 后台

* 文章管理：文章增删改查
* 页面管理：页面增删改查，可自定义 markdown 页面
* 分类管理：分类增删改查
* 标签管理：标签列表

#### 前台

* 文章列表：倒序展示文章、可置顶
* 内容页面：markdown 内容展示
* 标签页面：按标签文章数量排序
* 关于页面：个人说明
* 阅读清单：个人阅读书籍
* 站内搜索：支持文章标题、描述、内容、分类、标签模糊搜索

## 开发规范

遵循: [Uber Go 语言编码规范](https://github.com/uber-go/guide/blob/master/style.md)

### 常用命令

- make help 查看帮助
- make dep 下载 Go 依赖包
- make build 编译项目
- make tar 打包文件

### 部署流程

* 依赖环境：

  mysql、redis、elasticsearch
  > elasticsearch 可通过配置开启关闭，redis主要考虑到后续加缓存

* 安装部署

```
# 下载安装，可以不用是 GOPATH
git clone git@github.com:Joliya/blog.git

# 进入到下载目录
cd goblog

# 生成环境配置文件
cd conf

# 修改 mysql、redis、elasticsearch 配置

# 导入初始化 sql 结构
mysql -u root -p
> create database blog;
> set names utf8mb4;
> use blog;
> source blog.sql;


# 下载依赖
make dep

# 编译
make build

# 运行
./goblog dev.yml

# 后台运行
nohup ./goblog dev.yml &
```

* docker 启动 es、redis、mysql
```yml
version: '3'
services:
  cerebro:
    image: lmenezes/cerebro:0.8.3
    container_name: cerebro
    ports:
      - "9000:9000"
    command:
      - -Dhosts.0.host=http://elasticsearch:9200
    networks:
      - es7net
  kibana:
    image: docker.elastic.co/kibana/kibana:7.1.0
    container_name: kibana7
    environment:
      - I18N_LOCALE=zh-CN
      - XPACK_GRAPH_ENABLED=true
      - TIMELION_ENABLED=true
      - XPACK_MONITORING_COLLECTION_ENABLED="true"
    ports:
      - "5601:5601"
    networks:
      - es7net
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.1.0
    container_name: es7_01
    environment:
      - cluster.name=geektime
      - node.name=es7_01
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
      - discovery.seed_hosts=es7_01,es7_02
      - cluster.initial_master_nodes=es7_01,es7_02
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - es7data1:/usr/share/elasticsearch/data
    ports:
      - 9200:9200
    networks:
      - es7net
  elasticsearch2:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.1.0
    container_name: es7_02
    environment:
      - cluster.name=geektime
      - node.name=es7_02
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
      - discovery.seed_hosts=es7_01,es7_02
      - cluster.initial_master_nodes=es7_01,es7_02
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - es7data2:/usr/share/elasticsearch/data
    networks:
      - es7net
  mysql:
    container_name: blog-mysql
    restart: always
    platform: linux/x86_64
    image: mysql:5.7
    ports:
      - "3306:3306"
    volumes:
      - ./mysql/conf:/etc/mysql/conf.d
      - ./mysql/logs:/logs
      - ./mysql/data:/var/lib/mysql
    command: [
          'mysqld',
          '--innodb-buffer-pool-size=80M',
          '--character-set-server=utf8mb4',
          '--collation-server=utf8mb4_unicode_ci',
          '--default-time-zone=+8:00',
          '--lower-case-table-names=1'
        ]
    environment:
      MYSQL_ROOT_PASSWORD: 123456
  redis:
    container_name: blog_redis
    image: redis:latest
    volumes:
      - ./redis/redis.conf:/usr/local/etc/redis/redis.conf
      - ./redis/blog_redis/data:/data
    # command: redis-server --requirepass yourpass
    ports:
      - 6379:6379


volumes:
  es7data1:
    driver: local
  es7data2:
    driver: local

networks:
  es7net:
    driver: bridge
```

* supervisord 部署

```
[program:goblog]
directory = /data/modules/blog
command = /data/modules/blog/goblog -c conf/prod.yml
autostart = true
autorestart = true
startsecs = 5
user = root
redirect_stderr = true
stdout_logfile = /data/modules/blog/supervisor.log
```

* 访问首页

http://localhost:9091
