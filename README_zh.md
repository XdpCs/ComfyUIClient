# ComfyUIClient

![GitHub watchers](https://img.shields.io/github/watchers/XdpCs/ComfyUI-Client?style=social)
![GitHub stars](https://img.shields.io/github/stars/XdpCs/ComfyUI-Client?style=social)
![GitHub forks](https://img.shields.io/github/forks/XdpCs/ComfyUI-Client?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/XdpCs/ComfyUI-Client?style=flat-square)
![GitHub repo size](https://img.shields.io/github/repo-size/XdpCs/ComfyUI-Client?style=flat-square)
![GitHub license](https://img.shields.io/github/license/XdpCs/ComfyUI-Client?style=flat-square)

[English](README.md) | [中文](README_zh.md)

## 安装

`go get`

```shell
go get -u github.com/XdpCs/comfyUIclient
```

`go mod`

```shell
require github.com/XdpCs/comfyUIclient
```

## 支持 ComfyUI API

- [x] POST /prompt => func QueuePromptByString, QueuePromptByNodes
- [x] POST /queue => func DeleteAllQueues, DeleteQueueByPromptID
- [x] POST /history => func DeleteAllHistories, DeleteHistoryByPromptID
- [x] POST /interrupt => func InterruptExecution
- [x] POST /upload/image => func UploadImage
- [x] POST /upload/mask => func UploadMask
- [X] GET /embeddings => func GetEmbeddings
- [X] GET /extensions => func GetExtensions
- [X] GET /view => func GetFile
- [X] GET /view_metadata/{folder_name} => func GetViewMetadata
- [X] GET /system_stats => func GetSystemStats
- [X] GET /prompt => func GetQueueRemaining
- [X] GET /history => func GetAllHistories
- [X] GET /history/{prompt_id} => func GetHistoryByPromptID
- [X] GET /queue => func GetQueueInfo
- [X] GET /object_info => func GetObjectInfos
- [X] GET /object_info/{node_class} => func GetObjectInfoByNodeName

## 例子

所有例子都在 `examples` 目录中。

## 许可证

ComfyUI-Client 遵循 [MIT](LICENSE) 许可。请参考 LICENSE 获取更多信息。