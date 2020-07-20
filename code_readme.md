# apiserver 代码分析

## 1 http 请求流程分析

以 apf 的处理流程为例，把 http 的处理流程串起来。

apf 相关参数 "--feature-gates=APIPriorityAndFairness=true --runtime-config=flowcontrol.apiserver.k8s.io/v1alpha1=true"。

1 构建 filter 链

staging/src/k8s.io/apiserver/pkg/authentication/request/headerrequest/requestheader.go requestHeaderAuthRequestHandler:AuthenticateRequest

验证用户，并获取 user 信息。

k8s.io/apiserver/pkg/endpoints/filters/requestInfo.go:WithRequestInfo 
comment: 返回一个 http.Handler, 这个 handler 会根据 http.Request req 提取出一个 RequestInfo，然后存入 req.ctx 中【以后想要获取到 reqInfo，可通过request.RequestInfoFrom(req.Context()) 获取到】。

k8s.io/apiserver/pkg/server/filters/priority-and-fairness.go:WithPriorityAndFairness(
	handler http.Handler,
	longRunningRequestCheck apirequest.LongRunningRequestCheck,
	fcIfc utilflowcontrol.Interface,
) 
comment: 

* 1 从 http.Requst.Context() 中获取 RequestInfo 和 User；
    User 心思是调用 AuthenticateRequest 后获取的。
* 2 检查用户的请求是否是 long-running 类型，是则直接处理，不进行 apf 处理；
* 3 构造 note 函数：用于构建一个 PriorityAndFairnessClassification；
* 4 构造一个 execute 函数：获取与 priorityAndFairnessKey 相关的原始 innerCtx，然后执行 handler.ServerHTTP(w, innerReq)；
* 5 digest := utilflowcontrol.RequestDigest{requestInfo, user}；
* 6 调用 fcIfc.Handle(ctx, digest, note, execute)

k8s.io/apiserver/pkg/server/config.go:DefaultBuildHandlerChain(apiHandler http.Handler, c *Config) 会调用这两个函数。
DefaultBuildHandlerChain 中，构建了 filter 链，先注册的后调用。

## authn & authz & admission

* authn: 用户合法性 [如用户名/密码] 认证
* authz: 验证用户对其 http 请求中的 资源是否有访问权限
* admission: 通过如 webhook 方式让第三方实现的验证用户是否对某资源有访问权限

## apiserver 的代码目录 

staging/src/k8s.io/apiserver/pkg/apis 目录下各个子 目录 如下：

* admission      第三方认证
* apis           对外提供的接口
* audit          用户动作审计
* authentication 用户认证，认证时获取到的 User\Group\Org 等信息存入 RequestInfo
* authorization  用户授权，主要是 RBAC
* endpoints      用户访问 filter 链路等
* features       as 接口列表
* registry       etcd 的使用接口
* server         apiserver 自身的逻辑
* storage        etcd 数据库访问
* util           flowcontrol 等扩展功能
