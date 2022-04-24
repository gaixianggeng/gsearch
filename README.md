# doraemon

## 项目

### 落盘存储

> 使用b+数结构的bolt开源库存储

#### 正排库

* 正排文件，通过bolt进行kv存储 docid维度

#### 倒排库

* term文件 bolt存储，key为term，value为term对应的doccount、倒排索引文件的offset和size

* 倒排索引文件，os.File写入，mmap读取，存储docid、positions等倒排信息，存入的位置信息存入term文件

### 文件类型

* x.term存储term文件
* x.forword存储正排文件
* x.inverted存储倒排索引文件
* segments.gen 存储segment元数据信息，包括上述文件属性

---

## 笔记

### 索引构建方法

#### 静态索引构建

* BSBI
* SPIMI

#### 分布式索引构建

* [MapReduce](https://static.googleusercontent.com/media/research.google.com/zh-CN//archive/bigtable-osdi06.pdf)

#### 动态索引

* 对数合并
* lucene

### 数据库

---

### 备忘
* [snowflake](https://github.com/bwmarrin/snowflake)





* 只写入，不修改，每次根据 size 设置的阈值写入数据。

* TODO: 段数量过多，需要合并: <https://github.com/gaixianggeng/doraemon/blob/fbb74b167eed900ebcb4c11e14b87541dace60e3/internal/index/index.go#L63>

* 写入直接 os.File 写入，记录 offset 不会涉及到 mmap 的页读取，所以直接根据 write 量记录 offset, 初始化打开文件时，读取文件，设置 offset 即可

* 读取通过 mmap 读取
