- [sonic ：基于 JIT 技术的开源全场景高性能 JSON 库](#sonic-基于-jit-技术的开源全场景高性能-json-库)
- [简介](#简介)
  - [为什么要自研 JSON 库](#为什么要自研-json-库)
- [开源库 sonic 技术原理](#开源库-sonic-技术原理)
  - [`JIT`](#jit)
- [Lazy-load](#lazy-load)
  - [`sonic-ast`](#sonic-ast)
- [SIMD \& asm2asm](#simd--asm2asm)
- [性能测试](#性能测试)

# sonic ：基于 JIT 技术的开源全场景高性能 JSON 库
项目仓库：https://github.com/bytedance/sonic

版权声明：本文为CSDN博主「 字节跳动技术团队」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。

原文链接：https://blog.csdn.net/ByteDanceTech/article/details/122694591

# 简介
sonic 是字节跳动开源的一款 Golang JSON 库，基于即时编译（Just-In-Time Compilation）与向量化编程（Single Instruction Multiple Data）技术，大幅提升了 Go 程序的 JSON 编解码性能。同时结合 `lazy-load` 设计思想，它也为不同业务场景打造了一套全面高效的 API。

## 为什么要自研 JSON 库
JSON（JavaScript Object Notation） 以其简洁的语法和灵活的自描述能力，被广泛应用于各互联网业务。但是 JSON 由于本质是一种文本协议，且没有类似 Protobuf 的强制模型约束（schema），编解码效率往往十分低下。再加上有些业务开发者对 JSON 库的不恰当选型与使用，最终导致服务性能急剧劣化。


结果显示：目前这些 JSON 库均无法在各场景下都保持最优性能，即使是当前使用最广泛的第三方库 json-iterator，在泛型编解码、大数据量级场景下的性能也满足不了我们的需要。

JSON 库的基准编解码性能固然重要，但是对不同场景的最优匹配更关键 —— 于是我们走上了自研 JSON 库的道路。

# 开源库 sonic 技术原理
由于 JSON 业务场景复杂，指望通过单一算法来优化并不现实。于是在设计 sonic 的过程中，我们借鉴了其他领域/语言的优化思想（不仅限于 JSON），将其融合到各个处理环节中。其中较为核心的技术有三块：JIT、lazy-load 与 SIMD 。

## `JIT`
对于有 schema 的定型编解码场景而言，很多运算其实不需要在“运行时”执行。这里的“运行时”是指程序真正开始解析 JSON 数据的时间段。

举个例子，如果业务模型中确定了某个 JSON key 的值一定是布尔类型，那么我们就可以在序列化阶段直接输出这个对象对应的 JSON 值（‘true’或‘false’），并不需要再检查这个对象的具体类型。

sonic-JIT 的核心思想就是：**将模型解释与数据处理逻辑分离，让前者在“编译期”固定下来**。

这种思想也存在于标准库和某些第三方 JSON 库，如 json-iterator 的函数组装模式：把 Go struct 拆分解释成一个个字段类型的编解码函数，然后组装并缓存为整个对象对应的编解码器（codec），运行时再加载出来处理 JSON。但是这种实现难以避免转化成大量 interface 和 function 调用栈，随着 JSON 数据量级的增长，function-call 开销也成倍放大。只有将模型解释逻辑真正编译出来，实现 stack-less 的执行体，才能最大化 schema 带来的性能收益。

业界实现方式目前主要有两种：`代码生成 code-gen（或模版 template）`和 `即时编译 JIT`。前者的优点是库开发者实现起来相对简单，缺点是增加业务代码的维护成本和局限性，无法做到秒级热更新——这也是代码生成方式的 JSON 库受众并不广泛的原因之一。JIT 则将编译过程移到了程序的加载（或首次解析）阶段，只需要提供 JSON schema 对应的结构体类型信息，就可以一次性编译生成对应的 codec 并高效执行。

sonic-JIT 大致过程如下：

![https://img-blog.csdnimg.cn/img_convert/4e0cb6b7f49b08b9b1999bb4d9ce60de.png](https://img-blog.csdnimg.cn/img_convert/4e0cb6b7f49b08b9b1999bb4d9ce60de.png)

sonic-JIT 体系

- 初次运行时，基于 Go 反射来获取需要编译的 schema 信息； 
- 结合 JSON 编解码算法生成一套自定义的中间代码 OP codes； 
- 将 OP codes 翻译为 Plan9 汇编； 
- 使用第三方库 golang-asm 将 Plan 9 转为机器码； 
- 将生成的二进制码注入到内存 cache 中并封装为 go function； 
- 后续解析，直接根据 type ID （rtype.hash）从 cache 中加载对应的 codec 处理 JSON。

从最终实现的结果来看，`sonic-JIT` 生成的 codec 性能不仅好于 `json-iterator`，甚至超过了代码生成方式的 easyjson（见后文“性能测试”章节）。这一方面跟底层文本处理算子的优化有关（见后文“SIMD & asm2asm”章节），另一方面来自于 sonic-JIT 能控制底层 CPU 指令，在运行时建立了一套独立高效的 ABI（Application Binary Interface）体系：

- 将使用频繁的变量放到固定的寄存器上（如 JSON buffer、结构体指针），尽量避免 memory load & store； 
- 自己维护变量栈（内存池），避免 Go 函数栈扩展； 
- 自动生成跳转表，加速 generic decoding 的分支跳转； 
- 使用寄存器传递参数（当前 Go Assembly 并未支持，见“SIMD & asm2asm”章节）。

# Lazy-load
对于大部分 Go JSON 库，泛型编解码是它们性能表现最差的场景之一，然而由于业务本身需要或业务开发者的选型不当，它往往也是被应用得最频繁的场景。

泛型编解码性能差仅仅是因为没有 schema 吗？其实不然。我们可以对比一下 C++ 的 JSON 库，如 rappidjson、simdjson，它们的解析方式都是泛型的，但性能仍然很好（simdjson 可达 2GB/s 以上）。**标准库泛型解析性能差的根本原因在于它采用了 Go 原生泛型——interface（`map[string]interface{}`）作为 JSON 的编解码对象。**

这其实是一种糟糕的选择：首先是数据反序列化的过程中，map 插入的开销很高；其次在数据序列化过程中，map 遍历也远不如数组高效。

回过头来看，JSON 本身就具有完整的自描述能力，如果我们用一种与 JSON AST 更贴近的数据结构来描述，不但可以让转换过程更加简单，甚至可以实现按需加载（lazy-load）——这便是 sonic-ast 的核心逻辑：**它是一种 JSON 在 Go 中的编解码对象，用 `node {type, length, pointer}` 表示任意一个 JSON 数据节点，并结合树与数组结构描述节点之间的层级关系**。

## `sonic-ast`
sonic-ast 结构示意
![https://img-blog.csdnimg.cn/img_convert/d1eba71c3a30889103fbbfa396e0c052.png](https://img-blog.csdnimg.cn/img_convert/d1eba71c3a30889103fbbfa396e0c052.png)]

sonic-ast 实现了一种有状态、可伸缩的 JSON 解析过程：**当使用者 get 某个 key 时，sonic 采用 skip 计算来轻量化跳过要获取的 key 之前的 json 文本；对于该 key 之后的 JSON 节点，直接不做任何的解析处理；仅使用者真正需要的 key 才完全解析（转为某种 Go 原始类型）**。由于节点转换相比解析 JSON 代价小得多，在并不需要完整数据的业务场景下收益相当可观。

虽然 skip 是一种轻量的文本解析（处理 JSON 控制字符“[”、“{”等），但是使用类似 gjson 这种纯粹的 JSON 查找库时，往往会有相同路径查找导致的重复开销。

针对该问题，sonic 在对于子节点 skip 处理过程增加了一个步骤，将跳过 JSON 的 key、起始位、结束位记录下来，分配一个 Raw-JSON 类型的节点保存下来，这样二次 skip 就可以直接基于节点的 offset 进行。同时 sonic-ast 支持了节点的更新、插入和序列化，甚至支持将任意 Go types 转为节点并保存下来。


换言之，sonic-ast 可以作为一种通用的泛型数据容器替代 Go interface，在协议转换、动态代理等服务场景有巨大潜力。

# SIMD & asm2asm
无论是定型编解码场景还是泛型编解码场景，核心都离不开 JSON 文本的处理与计算。其中一些问题在业界已经有比较成熟高效的解决方案，如浮点数转字符串算法 Ryu，整数转字符串的查表法等，这些都被实现到 sonic 的底层文本算子中。

开发者们会发现这段代码其实是用 C 语言编写的 —— 其实 sonic 中绝大多数文本处理函数都是用 C 实现的：一方面 SIMD 指令集在 C 语言下有较好的封装，实现起来较为容易；另一方面这些 C 代码通过 clang 编译能充分享受其编译优化带来的提升。为此我们开发了一套 x86 汇编转 Plan9 汇编的工具 asm2asm，将 clang 输出的汇编通过 Go Assembly 机制静态嵌入到 sonic 中。同时在 JIT 生成的 codec 中我们利用 asm2asm 工具计算好的 C 函数 PC 值，直接调用 CALL 指令跳转，从而绕过 Go Assembly 不能寄存器传参的限制，压榨最后一丝 CPU 性能。

其它
除了上述提到的技术外，sonic 内部还有很多的细节优化，比如使用 RCU 替换 sync.Map 提升 codec cache 的加载速度，使用内存池减少 encode buffer 的内存分配，等等。这里限于篇幅便不详细展开介绍了，感兴趣的同学可以自行搜索阅读 sonic 源码进行了解。

# 性能测试
![https://img-blog.csdnimg.cn/img_convert/7237eee4173ff5cfa8bb26dd011d5bf3.png](https://img-blog.csdnimg.cn/img_convert/7237eee4173ff5cfa8bb26dd011d5bf3.png)

可以看到 sonic 在几乎所有场景下都处于领先（sonic-ast 由于直接使用了 Go Assembly 导入的 C 函数导致小数据集下有一定性能折损）

平均编码性能较 json-iterator 提升 240% ，平均解码性能较 json-iterator 提升 110% ；

单 key 修改能力较 sjson 提升 75% 。

并且在生产环境中，sonic 中也验证了良好的收益，服务高峰期占用核数减少将近三分之一：

字节某服务在 sonic 上线前后的 CPU 占用（核数）对比

结语
由于底层基于汇编进行开发，sonic 当前仅支持 amd64 架构下的 darwin/linux 平台 ，后续会逐步扩展到其它操作系统及架构。除此之外，我们也考虑将 sonic 在 Go 语言上的成功经验移植到不同语言及序列化协议中。目前 sonic 的 C++ 版本正在开发中，其定位是基于 sonic 核心思想及底层算子实现一套通用的高性能 JSON 编解码接口。

近日，sonic 发布了第一个大版本 v1.0.0，标志着其除了可被企业灵活用于生产环境，也正在积极响应社区需求、拥抱开源生态。我们期待 sonic 未来在使用场景和性能方面可以有更多突破，欢迎开发者们加入进来贡献 PR，一起打造业界最佳的 JSON 库！

相关链接

项目地址：https://github.com/bytedance/sonic

BenchMark：https://github.com/bytedance/sonic/blob/main/bench.sh
