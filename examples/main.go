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
      "seed": 1118,
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

func main() {
	client := comfyUIclient.NewDefaultClient("serverAddress", "port")
	client.ConnectAndListen()
	for !client.IsInitialized() {
	}

	// if you use the same seed, you will get the same result, so comfyUI will not give you result.
	go func() {
		_, err := client.QueuePromptByString(workflow)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	go func() {
		_, err := client.QueuePromptByNodes(getNodes())
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
					imageData, err := client.GetImage(image)
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
			"seed":         1114,
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
