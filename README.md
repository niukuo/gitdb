# gitdb
impl backend in golang for git2go
一开始是觉得git2go里面的后端只支持c实现，想要基于go实现mysql后端。于是实现了这些代码。里面各种go-c-go-c调用，基本上把cgo中能遇到的问题都遇到了。
后来考虑使用go-git实现可能简单一些，这些代码就先留在这里了。
