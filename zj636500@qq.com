[33mcommit ed033f92802cae57700fc8d33f4ade769af83307[m[33m ([m[1;36mHEAD -> [m[1;32mmaster[m[33m)[m
Author: ZJ-p007 <zj636500@qq.com>
Date:   Thu Mar 4 10:46:09 2021 +0800

    优化代码将寻找nonce和hash计算封装在一起

[33mcommit 1bdee64d95a1cac01535bf8746654f76cf422de8[m
Author: ZJ-p007 <zj636500@qq.com>
Date:   Thu Mar 4 10:32:43 2021 +0800

    解决包循环引用问题，定义区块数据接口，等并进行数据接口；解决pow bug

[33mcommit 281ef5d6946b1de2f97fb6e9fba56a4ccb4c0133[m
Author: ZJ-p007 <zj636500@qq.com>
Date:   Thu Mar 4 09:25:13 2021 +0800

    引入共识机制模块，采用pow实现共识算法

[33mcommit c8960756f64134d240daa64f818163642995f1d5[m[33m ([m[1;31mXianfengChain04/master[m[33m)[m
Author: ZJ-p007 <zj636500@qq.com>
Date:   Wed Mar 3 19:27:37 2021 +0800

    生成区块，创世区块、新区块的功能封装，区块哈希值计算

[33mcommit fe63d2761cd54a3f7714b01a7a27dbb223114dc1[m
Author: ZJ-p007 <zj636500@qq.com>
Date:   Wed Mar 3 19:23:25 2021 +0800

    生成一个区块，创世区块、新区快生成功能封装，区块哈希值计算

[33mcommit 2a94a6a504c6ac502e90aec89a171acd29463373[m
Author: ZJ-p007 <zj636500@qq.com>
Date:   Wed Mar 3 15:09:41 2021 +0800

    封装区块哈希值计算，已测试

[33mcommit c07bf3c6d6d7abb04c1a4ca59476d9b405fef2fa[m
Author: ZJ-p007 <zj636500@qq.com>
Date:   Wed Mar 3 14:40:20 2021 +0800

    首次提交，完成创世区块结构体定义，完成创世区块的封装，已测试
