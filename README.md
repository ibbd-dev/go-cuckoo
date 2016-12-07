# go-cuckoo

Golang Cuckoo Filter

## Install 

```sh
go get -u github.com/ibbd-dev/go-cuckoo
```

## Example


## 测试数据

对于100w的key，分别使用以下两种hash算法：

- `fnv-1`: 每个hash table的buckets的数量为2^22
- `fnv-1a`: 每个hash table的buckets的数量为2^22 （空间效率比`fnv-1`稍低）
- `md5`: 每个hash table的buckets的数量为2^20

ps: md5的散列性能比fnv好很多
