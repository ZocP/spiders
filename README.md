# 羊驼太极查询
## 控制台简单指令：
1. `spider update` 更新爬虫
2. `spider reload` 从文件中重新读取qa文件
3. `spider inspect_all_titles` 获取在内存中的全部文章标题
4. `server stop` os.EXIT(1) //人懒

## 文件目录
所有: `./files/`
配置文件: `./files/cfg.json`
生成的QA文件: `./files/spider/`


## API
GET 
(url)/v1/spider/find?key=<要搜索的QA>
