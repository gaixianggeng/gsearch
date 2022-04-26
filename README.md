# doraemon

## 项目

### 落盘存储

> 使用b+数结构的[bolt](https://github.com/boltdb/bolt)开源库存储term和正排数据
> 本来想自己手撸一个bptree，但是看了bolt的源码，我确定写不出来人家这种水平的代码，遂放弃，决定先聚焦核心功能。

#### 正排库

* 正排文件，通过bolt进行kv存储 docid主键

#### 倒排库

* term文件 bolt存储，key为term，value为term对应的doccount、倒排索引文件的offset和size

* 倒排索引文件，os.File写入，mmap读取，存储docid、positions等倒排信息，存入的位置信息存入term文件

#### engine对象

> engine是recall召回和index索引的控制模块

* 通过engine mode区分是查询还是写入，主要需要标识出要处理的segment，recall使用cur_seg_id，index使用next_seg_id

* 召回和索引是不同的engine

* 召回上层有多个engine

* recall 召回时会在engine上层进行多segment合并

### 文件类型

* x.term存储term文件
* x.forward存储正排文件
* x.inverted存储倒排索引文件
* segments.gen 存储segment元数据信息，包括上述文件属性

---

## NOTE

### 功能list

* 索引

  * [x] 创建term、正排、倒排

  * [x] 分segment写入  

  * [ ] segment merge

  * [ ] 删除

* 分词

  * [x] ngaram

  * [ ] ik(准备接入开源)

* 召回

  * [x] 短语召回

  * [x] 100%match召回
  
  * [x] 单segment召回

  * [ ] 多segment召回结果合并

* 相关性

  * [ ] tfidf

  * [ ] bm25

  * [ ] 加入词向量(看看实现难度吧...)

* 效果

  * [ ] and

  * [ ] or

  * [ ] 排序

  * [ ] 分页

### 记一下

#### 索引构建方法

静态索引构建

* BSBI
* SPIMI

分布式索引构建

* [MapReduce](https://static.googleusercontent.com/media/research.google.com/zh-CN//archive/bigtable-osdi06.pdf)

动态索引

* 对数合并
* lucene

笔记

* [snowflake](https://github.com/bwmarrin/snowflake)

* 只写入，不修改，每次根据 size 设置的阈值写入数据。

* 写入直接 os.File 写入，记录 offset 不会涉及到 mmap 的页读取，所以直接根据 write 量记录 offset, 初始化打开文件时，读取文件，设置 offset 即可

* 读取通过 mmap 读取
