// Package image 统一 OpenAI 兼容生图中继：同步 Helper、异步 task 提交/执行/轮询。
//
// 子模块：
//   - sync.go      同步生图 Helper
//   - execute.go   异步 worker 重放上游
//   - worker.go    队列与 CAS 状态机
//   - fetch.go     GET .../generations|edits/{id} 轮询
//
// 渠道策略与请求补丁见 relay/imagevendor；协议转换见 relay/channel/openai。
package image
