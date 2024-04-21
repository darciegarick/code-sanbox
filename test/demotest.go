package main

import (
	"fmt"
	"jassue-gin/sanbox"
)

func main() {
	code := `public class Main {
	public void main(String[] args) {
		int a =  Integer.parseInt(args[0]);
		int b =  Integer.parseInt(args[1]);
		System.out.println(a + b);
	}
}
	`

	file1, path, err := sanbox.SaveCodeToFile(code)
	if err != nil {
		return
	}
	fmt.Println("!!!!!!!!!!!!!!!!file1", file1)
	fmt.Println("!!!!!!!!!!!!!!!!path", path)
}
