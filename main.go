package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/foobaz/go-zopfli/zopfli"
	"io/ioutil"
	"os"
)

func main() {
	umlText, err := input()
	if err != nil {
		fmt.Println("error[input]", err)
		os.Exit(1)
	}
	bs, err := deflate(umlText)
	if err != nil {
		fmt.Println("error[compress]", err)
	}
	path := base64Encoding(bs)
	fmt.Printf("http://www.plantuml.com/plantuml/png/%s", path)
}

func input() (string, error) {
	if len(os.Args) == 1 {
		return fromStdIn(), nil
	} else if len(os.Args) == 2 {
		return fromFile(os.Args[1])
	} else {
		return "", fmt.Errorf("invalid params [%v]", os.Args[1:])
	}
}

func fromStdIn() string {
	buffer := bytes.Buffer{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func fromFile(file string) (string, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func deflate(umlText string) ([]byte, error) {
	var buffer bytes.Buffer
	options := zopfli.DefaultOptions()
	err := zopfli.ZlibCompress(&options, []byte(umlText), &buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

const mapper = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

func base64Encoding(input []byte) string {
	var buffer bytes.Buffer
	inputLength := len(input)
	for i := 0; i < 3 - inputLength % 3; i++ {
		input = append(input, byte(0))
	}

	for i := 0; i < inputLength; i += 3 {
		b1, b2, b3, b4 := input[i], input[i+1], input[i+2], byte(0)

		b4 = b3 & 0x3f
		b3 = ((b2 & 0xf) << 2) | (b3 >> 6)
		b2 = ((b1 & 0x3) << 4) | (b2 >> 4)
		b1 = b1 >> 2

		for _,b := range []byte{b1,b2,b3,b4} {
			buffer.WriteByte(byte(mapper[b]))
		}
	}
	return string(buffer.Bytes())
}

