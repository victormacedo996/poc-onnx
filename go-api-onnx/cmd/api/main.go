package main

import (
	"github.com/victormacedo996/poc-onnx/internal/config"
	"github.com/victormacedo996/poc-onnx/internal/onnx"
	"github.com/victormacedo996/poc-onnx/internal/webserver"
	"github.com/victormacedo996/poc-onnx/internal/webserver/rest/controllers"
	"github.com/yalue/onnxruntime_go"
)

// https://github.com/microsoft/onnxruntime/releases/download/v1.20.1/onnxruntime-linux-x64-1.20.1.tgz
func main() {
	c := config.NewConfig()
	onnx.InitializeEnvironment(c.SHARED_LIB_PATH)
	defer onnxruntime_go.DestroyEnvironment()

	sess, err := onnx.CreateSession(c.MODEL_PATH_ONNX)
	if err != nil {
		panic(err)
	}
	defer sess.Destroy()

	router := webserver.NewRouter(*c)
	health_controller := controllers.NewHealthController()
	model_controller := controllers.NewModelController(sess)
	router.Start(*health_controller, *model_controller)
}
