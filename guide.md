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
    5. index 计数

> segment 层
> 
1. Token2PostingsLists
    1. 初始化 token 维度的倒排索引
    2. CreateNewInvertedIndex
    3. CreateNewPostingsList
2. MergeInvertedIndex
3. Flush