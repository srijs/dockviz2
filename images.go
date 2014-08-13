package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Image struct {
	Id          string
	ParentId    string   `json:",omitempty"`
	RepoTags    []string `json:",omitempty"`
	VirtualSize int64
	Size        int64
	Created     int64
}

type ImagesCommand struct {
	Dot        bool `short:"d" long:"dot" description:"Show image information as Graphviz dot."`
	Tree       bool `short:"t" long:"tree" description:"Show image information as tree."`
	NoTruncate bool `short:"n" long:"no-trunc" description:"Don't truncate the image IDs."`
}

var imagesCommand ImagesCommand

func (x *ImagesCommand) Execute(args []string) error {

	// read in stdin
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading all input", err)
	}

	images, err := parseImagesJSON(stdin)
	if err != nil {
		return err
	}

	if imagesCommand.Dot {
		fmt.Printf(jsonToDot(images))
	} else if imagesCommand.Tree {

		var startImageArg = ""
		if len(args) > 0 {
			startImageArg = args[0]
		}

		fmt.Printf(jsonToTree(images, startImageArg, imagesCommand.NoTruncate))
	} else {
		return fmt.Errorf("Please specify either --dot or --tree")
	}

	return nil
}

func jsonToTree(images *[]Image, startImageArg string, noTrunc bool) string {
	var buffer bytes.Buffer

	var startImage Image

	var roots []Image
	var byParent = make(map[string][]Image)
	for _, image := range *images {
		if image.ParentId == "" {
			roots = append(roots, image)
		} else {
			if children, exists := byParent[image.ParentId]; exists {
				byParent[image.ParentId] = append(children, image)
			} else {
				byParent[image.ParentId] = []Image{image}
			}
		}

		if startImageArg != "" {
			if startImageArg == image.Id || startImageArg == truncate(image.Id) {
				startImage = image
			}

			for _, repotag := range image.RepoTags {
				if repotag == startImageArg {
					startImage = image
				}
			}
		}
	}

	if startImageArg != "" {
		WalkTree(&buffer, noTrunc, []Image{startImage}, byParent, "")
	} else {
		WalkTree(&buffer, noTrunc, roots, byParent, "")
	}

	return buffer.String()
}

func WalkTree(buffer *bytes.Buffer, noTrunc bool, images []Image, byParent map[string][]Image, prefix string) {
	if len(images) > 1 {
		length := len(images)
		for index, image := range images {
			if index+1 == length {
				PrintTreeNode(buffer, noTrunc, image, prefix+"└─")
				if subimages, exists := byParent[image.Id]; exists {
					WalkTree(buffer, noTrunc, subimages, byParent, prefix+"  ")
				}
			} else {
				PrintTreeNode(buffer, noTrunc, image, prefix+"|─")
				if subimages, exists := byParent[image.Id]; exists {
					WalkTree(buffer, noTrunc, subimages, byParent, prefix+"| ")
				}
			}
		}
	} else {
		for _, image := range images {
			PrintTreeNode(buffer, noTrunc, image, prefix+"└─")
			if subimages, exists := byParent[image.Id]; exists {
				WalkTree(buffer, noTrunc, subimages, byParent, prefix+"  ")
			}
		}
	}
}

func PrintTreeNode(buffer *bytes.Buffer, noTrunc bool, image Image, prefix string) {
	var imageID string
	if noTrunc {
		imageID = image.Id
	} else {
		imageID = truncate(image.Id)
	}

	buffer.WriteString(fmt.Sprintf("%s%s Virtual Size: %s", prefix, imageID, humanSize(image.VirtualSize)))
	if image.RepoTags[0] != "<none>:<none>" {
		buffer.WriteString(fmt.Sprintf(" Tags: %s\n", strings.Join(image.RepoTags, ", ")))
	} else {
		buffer.WriteString(fmt.Sprintf("\n"))
	}
}

func humanSize(raw int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB"}

	rawFloat := float64(raw)
	ind := 0

	for {
		if rawFloat < 1000 {
			break
		} else {
			rawFloat = rawFloat / 1000
			ind = ind + 1
		}
	}

	return fmt.Sprintf("%.01f %s", rawFloat, sizes[ind])
}

func truncate(id string) string {
	return id[0:12]
}

func parseImagesJSON(rawJSON []byte) (*[]Image, error) {

	var images []Image
	err := json.Unmarshal(rawJSON, &images)

	if err != nil {
		return nil, fmt.Errorf("Error reading JSON: ", err)
	}

	return &images, nil
}

func jsonToDot(images *[]Image) string {

	var buffer bytes.Buffer
	buffer.WriteString("digraph docker {\n")

	for _, image := range *images {
		if image.ParentId == "" {
			buffer.WriteString(fmt.Sprintf(" base -> \"%s\" [style=invis]\n", truncate(image.Id)))
		} else {
			buffer.WriteString(fmt.Sprintf(" \"%s\" -> \"%s\"\n", truncate(image.ParentId), truncate(image.Id)))
		}
		if image.RepoTags[0] != "<none>:<none>" {
			buffer.WriteString(fmt.Sprintf(" \"%s\" [label=\"%s (+%s) (%s)\\n%s\",shape=box,fillcolor=\"paleturquoise\",style=\"filled,rounded\"];\n", truncate(image.Id), truncate(image.Id), humanSize(image.Size), humanSize(image.VirtualSize), strings.Join(image.RepoTags, "\\n")))
		} else {
			buffer.WriteString(fmt.Sprintf(" \"%s\" [label=\"%s (+%s) (%s)\"];\n", truncate(image.Id), truncate(image.Id), humanSize(image.Size), humanSize(image.VirtualSize)))
    }
	}

	buffer.WriteString(" base [style=invisible]\n}\n")

	return buffer.String()
}

func init() {
	parser.AddCommand("images",
		"Visualize docker images.",
		"",
		&imagesCommand)
}
