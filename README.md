# ComfyUI-Client

![GitHub watchers](https://img.shields.io/github/watchers/XdpCs/ComfyUI-Client?style=social)
![GitHub stars](https://img.shields.io/github/stars/XdpCs/ComfyUI-Client?style=social)
![GitHub forks](https://img.shields.io/github/forks/XdpCs/ComfyUI-Client?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/XdpCs/ComfyUI-Client?style=flat-square)
![GitHub repo size](https://img.shields.io/github/repo-size/XdpCs/ComfyUI-Client?style=flat-square)
![GitHub license](https://img.shields.io/github/license/XdpCs/ComfyUI-Client?style=flat-square)

## Install

`go get`

```shell
go get -u github.com/XdpCs/comfyUIclient
```

`go mod`

```shell
require github.com/XdpCs/comfyUIclient
```

## Support the ComfyUI API

- [x] POST /prompt => func QueuePromptByString, QueuePromptByNodes
- [x] POST /queue => func DeleteAllQueues, DeleteQueueByPromptID
- [x] POST /history => func DeleteAllHistories, DeleteHistoryByPromptID
- [x] POST /interrupt => func InterruptExecution
- [ ] POST /upload/image
- [ ] POST /upload/mask
- [X] GET /embeddings => func GetEmbeddings
- [X] GET /extensions => func GetExtensions
- [X] GET /view => func GetImage
- [X] GET /view_metadata/{folder_name} => func GetViewMetadata
- [X] GET /system_stats => func GetSystemStats
- [X] GET /prompt => func GetQueueRemaining
- [X] GET /history => func GetAllHistories
- [X] GET /prompt/{prompt_id} => func GetHistoryByPromptID
- [ ] GET /queue => func GetQueueInfo
- [ ] GET /object_info => func GetObjectInfo
- [ ] GET /object_info/{node_class} => func GetObjectInfoByNodeClass

## Examples

All examples are in the `examples` directory.

## License

ComfyUI-Client is under the [MIT](LICENSE). Please refer to LICENSE for more information.
