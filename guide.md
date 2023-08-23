# 项目解释
## 相关名词
```
forward 正排
segment 段
InvertedIndexValue 倒排索引
PostingsList 倒排列表
TermValue 存储的doc_count、offset、size
```

## 项目结构
### 写入流程

1. 解析读取 mete.Profile
    1. 存储的是索引的元数据，包括 segment 的信息 

> index 层
> 
1. 初始化 engine
2. 读取源文件，解析源文件，获取 id、标题和正文
3. 单独启动一个协程进行 merge

> engine 层
> 
1. 添加到正排数据
2. Text2PostingsLists 文本和id 转为倒排索引
    1. 分词获取token list
    2. 将每个 token转为倒排列表
    3. 合并token 相同，doc 不同的数据
    4. 落盘 flush 操作
        1. 更新meta 的信息
    5. index 计数

> segment 层
> 
1. Token2PostingsLists
    1. 初始化 token 维度的倒排索引
    2. CreateNewInvertedIndex 创建新的倒排索引
    3. CreateNewPostingsList 创建新的倒排列表，赋值给b 中新建的索引
2. MergeInvertedIndex 合并两个 map[string]*InvertedIndexValue类型的倒排索引
3. 是否达到阈值，达到的话 Flush
4. storagePostings 存储InvertedIndexValue
    1. 先编码 存储的编码结构 ？？ 待确定
    2. 再写入InvertedDB.StoragePostings ？？inverted_db 层
5. indexToCount 对索引的 doc 计数

> merge 流程
> 
- index 索引流程启动的时候会启动 channel 来接收带 merge 的segment
- 每次 flush 的时候会计算是否需要 merge
- 当前merge 的算法实现还没弄