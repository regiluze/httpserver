// Copyright 2015 The httpserver Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.
package main

import (
	"flag"
	"fmt"

	"github.com/regiluze/httpserver"
)

const (
	Context = ""
)

func main() {
	port := flag.String("port", "8080", "listen port")
	address := flag.String("address", "0.0.0.0", "server address")
	flag.Parse()
	uploadImageApp := NewImageUploadHandler(Context)
	server := httpserver.NewHttpServer(*address, *port)
	server.SetErrTemplate(ErrorTemplate)
	server.SetNotFoundTemplate(NotFoundTemplate)
	server.Deploy(Context, uploadImageApp)
	error := server.Start()
	if error != nil {
		fmt.Println(error)
	}

}
