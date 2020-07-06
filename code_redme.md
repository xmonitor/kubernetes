# apiserver 代码分析

## 1 http 请求流程分析

1 

k8s.io/apiserver/pkg/endpoints/filters/requestInfo.go:WithRequestInfo 
comment: 返回一个 http.Handler, 这个 handler 会根据 http.Request req 提取出一个 RequestInfo，然后存入 req.ctx 中【以后想要获取到 reqInfo，可通过request.RequestInfoFrom(req.Context()) 获取到】。

k8s.io/apiserver/pkg/server/filters/priority-and-fairness.go:WithPriorityAndFairness(
	handler http.Handler,
	longRunningRequestCheck apirequest.LongRunningRequestCheck,
	fcIfc utilflowcontrol.Interface,
) 
comment: 
1 从 http.Requst.Context() 中获取 RequestInfo 和 User；
2 检查用户的请求是否是 long-running 类型，是则直接处理，不进行 apf 处理；
3 构造 note 函数：用于构建一个 PriorityAndFairnessClassification；
4 构造一个 execute 函数：获取与 priorityAndFairnessKey 相关的原始 innerCtx，然后执行 handler.ServerHTTP(w, innerReq)；
5 digest := utilflowcontrol.RequestDigest{requestInfo, user}；
6 调用 fcIfc.Handle(ctx, digest, note, execute)

k8s.io/apiserver/pkg/server/config.go:DefaultBuildHandlerChain(apiHandler http.Handler, c *Config) 会调用这两个函数。
