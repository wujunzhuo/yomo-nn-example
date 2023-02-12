# yomo-nn-example

https://github.com/yomorun/yomo#1-install-cli

## 1. zipper

```sh
yomo serve -c workflow.yaml
```

## 2. sink

```sh
cd sink

yomo run --name sink main.go
```

## 3. mobilenet_onnx

```sh
cd mobilenet_onnx

curl -O https://media.githubusercontent.com/media/onnx/models/main/vision/classification/mobilenet/model/mobilenetv2-7.onnx

rustup target add wasm32-wasi

cargo build --release --target wasm32-wasi

yomo run --name mobilenet_onnx target/wasm32-wasi/release/sfn.wasm
```

## 4. source

```sh
cd source

go run main.go
```
