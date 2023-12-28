package main

import (
	"fmt"

	"github.com/XdpCs/comfyUIclient"
)

func main() {
	var workflow = `{
  "3": {
    "inputs": {
      "seed": 440802502753153,
      "steps": 20,
      "cfg": 2.5,
      "sampler_name": "euler",
      "scheduler": "karras",
      "denoise": 1,
      "model": [
        "14",
        0
      ],
      "positive": [
        "12",
        0
      ],
      "negative": [
        "12",
        1
      ],
      "latent_image": [
        "12",
        2
      ]
    },
    "class_type": "KSampler"
  },
  "8": {
    "inputs": {
      "samples": [
        "3",
        0
      ],
      "vae": [
        "15",
        2
      ]
    },
    "class_type": "VAEDecode"
  },
  "12": {
    "inputs": {
      "width": 512,
      "height": 512,
      "video_frames": 32,
      "motion_bucket_id": 127,
      "fps": 8,
      "augmentation_level": 0,
      "clip_vision": [
        "15",
        1
      ],
      "init_image": [
        "24",
        0
      ],
      "vae": [
        "15",
        2
      ]
    },
    "class_type": "SVD_img2vid_Conditioning"
  },
  "14": {
    "inputs": {
      "min_cfg": 1,
      "model": [
        "15",
        0
      ]
    },
    "class_type": "VideoLinearCFGGuidance"
  },
  "15": {
    "inputs": {
      "ckpt_name": "svd_xt.safetensors"
    },
    "class_type": "ImageOnlyCheckpointLoader"
  },
  "23": {
    "inputs": {
      "frame_rate": 8,
      "loop_count": 0,
      "filename_prefix": "SVD_output",
      "format": "image/gif",
      "pingpong": false,
      "save_image": true,
      "crf": 20,
      "save_metadata": true,
      "videopreview": {
        "hidden": false,
        "paused": false,
        "params": {
          "filename": "SVD_output_00002.gif",
          "subfolder": "",
          "type": "output",
          "format": "image/gif"
        }
      },
      "images": [
        "8",
        0
      ]
    },
    "class_type": "VHS_VideoCombine"
  },
  "24": {
    "inputs": {
      "image": "ComfyUI_temp_pghaf_00089_.png",
      "choose file to upload": "image"
    },
    "class_type": "LoadImage"
  }
}`
	client := comfyUIclient.NewDefaultClient("serverAddress", "port")
	client.ConnectAndListen()
	for !client.IsInitialized() {
	}
	_, err := client.QueuePrompt(workflow)
	if err != nil {
		panic(err)
	}
	for taskStatus := range client.GetTaskStatus() {
		switch taskStatus.Type {
		case comfyUIclient.Status:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataStatus)
			fmt.Printf("Queue remaining: %d", s.Status.ExecInfo.QueueRemaining)
		case comfyUIclient.ExecutionStart:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataExecutionStart)
			fmt.Println(s)
		case comfyUIclient.ExecutionCached:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataExecutionCached)
			fmt.Println(s)
		case comfyUIclient.Executing:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataExecuting)
			fmt.Println(s)
		case comfyUIclient.Progress:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataProgress)
			fmt.Println(s)
		case comfyUIclient.Executed:
			s := taskStatus.Data.(*comfyUIclient.WSMessageDataExecuted)
			fmt.Println(s)
		case comfyUIclient.ExecutionInterrupted:
			s := taskStatus.Data.(*comfyUIclient.WSMessageExecutionInterrupted)
			fmt.Println(s)
		case comfyUIclient.ExecutionError:
			s := taskStatus.Data.(*comfyUIclient.WSMessageExecutionError)
			fmt.Println(s)
		default:
			fmt.Println("unknown message type")
		}
	}
}
