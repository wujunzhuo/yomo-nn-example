# yomo-nn-example

## 1. Install YoMo Cli

https://yomo.run/docs#install-cli

## 2. run YoMo Zipper

```sh
yomo serve -c config.yaml
```

## 2. mobilenet_onnx: the serverless function to make AI inference

```sh
# download the ONNX model
curl -O https://media.githubusercontent.com/media/onnx/models/main/vision/classification/mobilenet/model/mobilenetv2-7.onnx

# compile the serverless function to WebAssembly file
rustup target add wasm32-wasi
cargo build --release --target wasm32-wasi --manifest-path=mobilenet_onnx/Cargo.toml

# run the serverless function
yomo run mobilenet_onnx/target/wasm32-wasi/release/sfn.wasm
```

## 3. send image and receive the infernce result

```sh
go build -o cli ./cmd

./cli sample.png
```
