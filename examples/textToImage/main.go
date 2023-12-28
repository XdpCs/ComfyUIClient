package main

import (
	"fmt"
	"log"
	"os"

	"github.com/XdpCs/comfyUIclient"
)

func main() {
	var workflow = `{
    "3": {
        "class_type": "KSampler",
        "inputs": {
            "cfg": 8,
            "denoise": 1,
            "latent_image": [
                "5",
                0
            ],
            "model": [
                "4",
                0
            ],
            "negative": [
                "7",
                0
            ],
            "positive": [
                "6",
                0
            ],
            "sampler_name": "euler",
            "scheduler": "normal",
            "seed": 8566257,
            "steps": 20
        }
    },
    "4": {
        "class_type": "CheckpointLoaderSimple",
        "inputs": {
            "ckpt_name": "TV002.ckpt"
        }
    },
    "5": {
        "class_type": "EmptyLatentImage",
        "inputs": {
            "batch_size": 1,
            "height": 512,
            "width": 512
        }
    },
    "6": {
        "class_type": "CLIPTextEncode",
        "inputs": {
            "clip": [
                "4",
                1
            ],
            "text": "masterpiece best quality girl"
        }
    },
    "7": {
        "class_type": "CLIPTextEncode",
        "inputs": {
            "clip": [
                "4",
                1
            ],
            "text": "bad hands"
        }
    },
    "8": {
        "class_type": "VAEDecode",
        "inputs": {
            "samples": [
                "3",
                0
            ],
            "vae": [
                "4",
                2
            ]
        }
    },
    "9": {
        "class_type": "SaveImage",
        "inputs": {
            "filename_prefix": "ComfyUI",
            "images": [
                "8",
                0
            ]
        }
    }
}`
	client := comfyUIclient.NewDefaultClient("hz-t2.matpool.com", "29900")
	client.ConnectAndListen()
	for !client.IsInitialized() {
	}
	_, err := client.QueuePrompt(workflow)
	if err != nil {
		panic(err)
	}
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
		case comfyUIclient.ExecutionInterrupted:
			s := taskStatus.Data.(*comfyUIclient.WSMessageExecutionInterrupted)
			fmt.Printf("Type: %v, Data:%+v\n", comfyUIclient.ExecutionInterrupted, s)
		case comfyUIclient.ExecutionError:
			s := taskStatus.Data.(*comfyUIclient.WSMessageExecutionError)
			fmt.Printf("Type: %v, Data:%+v\n", comfyUIclient.ExecutionError, s)
		default:
			fmt.Println("unknown message type")
		}
	}
}
