package main

import (
	"context"
	"fmt"
	"github.com/codeuniversity/smag-mvp/faces/proto"
	"github.com/codeuniversity/smag-mvp/imgproxy"
	"google.golang.org/grpc"
)

func main() {
	con, err := grpc.Dial("localhost:6666", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := proto.NewFaceRecognizerClient(con)
	response, err := c.RecognizeFaces(context.Background(), &proto.RecognizeRequest{
		Url: "https://cdn.madaracosmetics.com/media/catalog/category/FACE_OK_3.jpg",
	})
	if err != nil {
		panic(err)
	}
	p, err := imgproxy.New("localhost:8080", "5800c215e5cd5110365c390e83752526fa40758efa4dcc406e3a4fdd6e22877c", "520f986b998545b4785e0defbc4f3c1203f22de2374a3d53cb7a7fe9fea309c5")
	if err != nil {
		panic(err)
	}
	faces := response.Faces
	for _, face := range faces {
		fmt.Println(face)
		x := int(face.X)
		y := int(face.Y)
		width := int(face.Width)
		height := int(face.Height)
		url := p.GetCropURL(x, y, width, height, "https://cdn.madaracosmetics.com/media/catalog/category/FACE_OK_3.jpg")
		fmt.Println(url)
	}
}
