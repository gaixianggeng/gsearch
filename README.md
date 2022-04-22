# doraemon

## 项目

### 落盘存储

#### 正排库

#### 倒排库

- 通过 b+数实现

---

## 笔记

### 索引构建方法

#### 静态索引构建

- BSBI
- SPIMI

#### 分布式索引构建

- [MapReduce](https://static.googleusercontent.com/media/research.google.com/zh-CN//archive/bigtable-osdi06.pdf)

#### 动态索引

- 对数合并
- lucene

### 数据库

#### snowflake

https://github.com/bwmarrin/snowflake

#### 正排 b+树

#### 倒排 b+树

#### token 存储 hash 存储

---

- [x] ngaram
- [ ] 测试用例

---

只写入，不修改，每次根据 size 设置的阈值写入数据。
TODO: 段数量过多，需要合并: https://github.com/gaixianggeng/doraemon/blob/fbb74b167eed900ebcb4c11e14b87541dace60e3/internal/index/index.go#L63

写入直接 os.File 写入，记录 offset 不会涉及到 mmap 的页读取，所以直接根据 write 量记录 offset

初始化打开文件时，读取文件，设置 offset 即可

读取通过 mmap 读取
