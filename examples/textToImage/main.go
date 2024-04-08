package main

import (
	"fmt"
	"log"
	"os"

	"github.com/XdpCs/comfyUIclient"
)

var workflow = `{
  "3": {
    "inputs": {
      "seed": 1114,
      "steps": 20,
      "cfg": 8,
      "sampler_name": "euler",
      "scheduler": "normal",
      "denoise": 1,
      "model": [
        "4",
        0
      ],
      "positive": [
        "6",
        0
      ],
      "negative": [
        "7",
        0
      ],
      "latent_image": [
        "5",
        0
      ]
    },
    "class_type": "KSampler"
  },
  "4": {
    "inputs": {
      "ckpt_name": "CounterfeitV30_v30.safetensors"
    },
    "class_type": "CheckpointLoaderSimple"
  },
  "5": {
    "inputs": {
      "width": 512,
      "height": 512,
      "batch_size": 1
    },
    "class_type": "EmptyLatentImage"
  },
  "6": {
    "inputs": {
      "text": "a beautiful girl",
      "clip": [
        "4",
        1
      ]
    },
    "class_type": "CLIPTextEncode"
  },
  "7": {
    "inputs": {
      "text": "text, watermark",
      "clip": [
        "4",
        1
      ]
    },
    "class_type": "CLIPTextEncode"
  },
  "8": {
    "inputs": {
      "samples": [
        "3",
        0
      ],
      "vae": [
        "4",
        2
      ]
    },
    "class_type": "VAEDecode"
  },
  "9": {
    "inputs": {
      "filename_prefix": "ComfyUI",
      "images": [
        "8",
        0
      ]
    },
    "class_type": "SaveImage"
  }
}`

var extraDataString string = `{
                    "workflow": {
                        "last_node_id": 9,
                        "last_link_id": 9,
                        "nodes": [
                            {
                                "id": 7,
                                "type": "CLIPTextEncode",
                                "pos": [
                                    413,
                                    389
                                ],
                                "size": {
                                    "0": 425.27801513671875,
                                    "1": 180.6060791015625
                                },
                                "flags": {},
                                "order": 3,
                                "mode": 0,
                                "inputs": [
                                    {
                                        "name": "clip",
                                        "type": "CLIP",
                                        "link": 5
                                    }
                                ],
                                "outputs": [
                                    {
                                        "name": "CONDITIONING",
                                        "type": "CONDITIONING",
                                        "links": [
                                            6
                                        ],
                                        "slot_index": 0
                                    }
                                ],
                                "properties": {
                                    "Node name for S&R": "CLIPTextEncode"
                                },
                                "widgets_values": [
                                    "text, watermark"
                                ]
                            },
                            {
                                "id": 6,
                                "type": "CLIPTextEncode",
                                "pos": [
                                    415,
                                    186
                                ],
                                "size": {
                                    "0": 422.84503173828125,
                                    "1": 164.31304931640625
                                },
                                "flags": {},
                                "order": 2,
                                "mode": 0,
                                "inputs": [
                                    {
                                        "name": "clip",
                                        "type": "CLIP",
                                        "link": 3
                                    }
                                ],
                                "outputs": [
                                    {
                                        "name": "CONDITIONING",
                                        "type": "CONDITIONING",
                                        "links": [
                                            4
                                        ],
                                        "slot_index": 0
                                    }
                                ],
                                "properties": {
                                    "Node name for S&R": "CLIPTextEncode"
                                },
                                "widgets_values": [
                                    "beautiful scenery nature glass bottle landscape, , purple galaxy bottle,"
                                ]
                            },
                            {
                                "id": 5,
                                "type": "EmptyLatentImage",
                                "pos": [
                                    473,
                                    609
                                ],
                                "size": {
                                    "0": 315,
                                    "1": 106
                                },
                                "flags": {},
                                "order": 0,
                                "mode": 0,
                                "outputs": [
                                    {
                                        "name": "LATENT",
                                        "type": "LATENT",
                                        "links": [
                                            2
                                        ],
                                        "slot_index": 0
                                    }
                                ],
                                "properties": {
                                    "Node name for S&R": "EmptyLatentImage"
                                },
                                "widgets_values": [
                                    512,
                                    512,
                                    1
                                ]
                            },
                            {
                                "id": 3,
                                "type": "KSampler",
                                "pos": [
                                    863,
                                    186
                                ],
                                "size": {
                                    "0": 315,
                                    "1": 262
                                },
                                "flags": {},
                                "order": 4,
                                "mode": 0,
                                "inputs": [
                                    {
                                        "name": "model",
                                        "type": "MODEL",
                                        "link": 1
                                    },
                                    {
                                        "name": "positive",
                                        "type": "CONDITIONING",
                                        "link": 4
                                    },
                                    {
                                        "name": "negative",
                                        "type": "CONDITIONING",
                                        "link": 6
                                    },
                                    {
                                        "name": "latent_image",
                                        "type": "LATENT",
                                        "link": 2
                                    }
                                ],
                                "outputs": [
                                    {
                                        "name": "LATENT",
                                        "type": "LATENT",
                                        "links": [
                                            7
                                        ],
                                        "slot_index": 0
                                    }
                                ],
                                "properties": {
                                    "Node name for S&R": "KSampler"
                                },
                                "widgets_values": [
                                    156680208700286,
                                    "randomize",
                                    20,
                                    8,
                                    "euler",
                                    "normal",
                                    1
                                ]
                            },
                            {
                                "id": 8,
                                "type": "VAEDecode",
                                "pos": [
                                    1209,
                                    188
                                ],
                                "size": {
                                    "0": 210,
                                    "1": 46
                                },
                                "flags": {},
                                "order": 5,
                                "mode": 0,
                                "inputs": [
                                    {
                                        "name": "samples",
                                        "type": "LATENT",
                                        "link": 7
                                    },
                                    {
                                        "name": "vae",
                                        "type": "VAE",
                                        "link": 8
                                    }
                                ],
                                "outputs": [
                                    {
                                        "name": "IMAGE",
                                        "type": "IMAGE",
                                        "links": [
                                            9
                                        ],
                                        "slot_index": 0
                                    }
                                ],
                                "properties": {
                                    "Node name for S&R": "VAEDecode"
                                }
                            },
                            {
                                "id": 9,
                                "type": "SaveImage",
                                "pos": [
                                    1451,
                                    189
                                ],
                                "size": {
                                    "0": 210,
                                    "1": 58
                                },
                                "flags": {},
                                "order": 6,
                                "mode": 0,
                                "inputs": [
                                    {
                                        "name": "images",
                                        "type": "IMAGE",
                                        "link": 9
                                    }
                                ],
                                "properties": {},
                                "widgets_values": [
                                    "ComfyUI"
                                ]
                            },
                            {
                                "id": 4,
                                "type": "CheckpointLoaderSimple",
                                "pos": [
                                    26,
                                    474
                                ],
                                "size": {
                                    "0": 315,
                                    "1": 98
                                },
                                "flags": {},
                                "order": 1,
                                "mode": 0,
                                "outputs": [
                                    {
                                        "name": "MODEL",
                                        "type": "MODEL",
                                        "links": [
                                            1
                                        ],
                                        "slot_index": 0
                                    },
                                    {
                                        "name": "CLIP",
                                        "type": "CLIP",
                                        "links": [
                                            3,
                                            5
                                        ],
                                        "slot_index": 1
                                    },
                                    {
                                        "name": "VAE",
                                        "type": "VAE",
                                        "links": [
                                            8
                                        ],
                                        "slot_index": 2
                                    }
                                ],
                                "properties": {
                                    "Node name for S&R": "CheckpointLoaderSimple"
                                },
                                "widgets_values": [
                                    "CounterfeitV30_v30.safetensors"
                                ]
                            }
                        ],
                        "links": [
                            [
                                1,
                                4,
                                0,
                                3,
                                0,
                                "MODEL"
                            ],
                            [
                                2,
                                5,
                                0,
                                3,
                                3,
                                "LATENT"
                            ],
                            [
                                3,
                                4,
                                1,
                                6,
                                0,
                                "CLIP"
                            ],
                            [
                                4,
                                6,
                                0,
                                3,
                                1,
                                "CONDITIONING"
                            ],
                            [
                                5,
                                4,
                                1,
                                7,
                                0,
                                "CLIP"
                            ],
                            [
                                6,
                                7,
                                0,
                                3,
                                2,
                                "CONDITIONING"
                            ],
                            [
                                7,
                                3,
                                0,
                                8,
                                0,
                                "LATENT"
                            ],
                            [
                                8,
                                4,
                                2,
                                8,
                                1,
                                "VAE"
                            ],
                            [
                                9,
                                8,
                                0,
                                9,
                                0,
                                "IMAGE"
                            ]
                        ],
                        "groups": [],
                        "config": {},
                        "extra": {
							"xdp": "xdp loves fyy."
						},
                        "version": 0.4
                    }
                }`

func main() {
	endPoint := comfyUIclient.NewEndPoint("https", "serverAddress", "port")
	client := comfyUIclient.NewDefaultClient(endPoint)
	client.ConnectAndListen()
	for !client.IsInitialized() {
	}

	// if you use the same seed, you will get the same result, so comfyUI will not give you result.
	go func() {
		_, err := client.QueuePromptByString(workflow, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	go func() {
		_, err := client.QueuePromptByNodes(getNodes(), extraDataString)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	count := 0
	for taskStatus := range client.GetTaskStatus() {
		switch taskStatus.Type {
		case comfyUIclient.ExecutionStart:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataExecutionStart)
			fmt.Printf("Type: %v, Data:%+v\n", comfyUIclient.ExecutionStart, s)
		case comfyUIclient.ExecutionCached:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataExecutionCached)
			fmt.Printf("Type: %v, Data:%+v\n", comfyUIclient.ExecutionCached, s)
		case comfyUIclient.Executing:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataExecuting)
			fmt.Printf("Type: %v, Data:%+v\n", comfyUIclient.Executing, s)
		case comfyUIclient.Progress:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataProgress)
			fmt.Printf("Type: %v, Data:%+v\n", comfyUIclient.Progress, s)
		case comfyUIclient.Executed:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataExecuted)
			fmt.Printf("Type: %v, Data:%+v\n", comfyUIclient.Executed, s)
			for _, images := range s.Output {
				for _, image := range images {
					imageData, err := client.GetFile(image)
					if err != nil {
						panic(err)
					}
					f, err := os.Create(image.Filename)
					if err != nil {
						log.Println("Failed to write image:", err)
						os.Exit(1)
					}
					f.Write(*imageData)
					f.Close()
				}
			}
			count++
			IsEndQueuePrompt(count, 2)
		case comfyUIclient.ExecutionInterrupted:
			s := taskStatus.Data.(*comfyUIclient.WSMessageExecutionInterrupted)
			fmt.Printf("Type: %v, Data:%+v\n", comfyUIclient.ExecutionInterrupted, s)
			count++
			IsEndQueuePrompt(count, 2)
		case comfyUIclient.ExecutionError:
			s := taskStatus.Data.(*comfyUIclient.WSMessageExecutionError)
			fmt.Printf("Type: %v, Data:%+v\n", comfyUIclient.ExecutionError, s)
			count++
			IsEndQueuePrompt(count, 2)
		default:
			fmt.Println("unknown message type")
		}
	}
}

func IsEndQueuePrompt(count int, num int) {
	if count >= num {
		fmt.Println("end queue prompt")
		os.Exit(0)
	}
}

func getNodes() map[string]comfyUIclient.PromptNode {
	nodes := map[string]comfyUIclient.PromptNode{}
	nodes["3"] = comfyUIclient.PromptNode{
		Inputs: map[string]interface{}{
			"seed":         1118,
			"steps":        20,
			"cfg":          8,
			"sampler_name": "euler",
			"scheduler":    "normal",
			"denoise":      1,
			"model":        []interface{}{"4", 0},
			"positive":     []interface{}{"6", 0},
			"negative":     []interface{}{"7", 0},
			"latent_image": []interface{}{"5", 0},
		},
		ClassType: "KSampler",
	}
	nodes["4"] = comfyUIclient.PromptNode{
		Inputs: map[string]interface{}{
			"ckpt_name": "CounterfeitV30_v30.safetensors",
		},
		ClassType: "CheckpointLoaderSimple",
	}

	nodes["5"] = comfyUIclient.PromptNode{
		Inputs: map[string]interface{}{
			"width":      512,
			"height":     512,
			"batch_size": 1,
		},
		ClassType: "EmptyLatentImage",
	}

	nodes["6"] = comfyUIclient.PromptNode{
		Inputs: map[string]interface{}{
			"text": "a beautiful girl",
			"clip": []interface{}{"4", 1},
		},
		ClassType: "CLIPTextEncode",
	}

	nodes["7"] = comfyUIclient.PromptNode{
		Inputs: map[string]interface{}{
			"text": "text, watermark",
			"clip": []interface{}{"4", 1},
		},
		ClassType: "CLIPTextEncode",
	}

	nodes["8"] = comfyUIclient.PromptNode{
		Inputs: map[string]interface{}{
			"samples": []interface{}{"3", 0},
			"vae":     []interface{}{"4", 2},
		},
		ClassType: "VAEDecode",
	}

	nodes["9"] = comfyUIclient.PromptNode{
		Inputs: map[string]interface{}{
			"filename_prefix": "ComfyUI",
			"images":          []interface{}{"8", 0},
		},
		ClassType: "SaveImage",
	}
	return nodes
}
