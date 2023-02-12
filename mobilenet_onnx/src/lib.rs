use std::{fs::OpenOptions, sync::Mutex};

use anyhow::Result;
use image::{
    imageops::{resize, FilterType},
    load_from_memory,
};
use serde::Serialize;
use tract_onnx::{
    onnx,
    prelude::{
        tract_ndarray, tvec, Framework, Graph, InferenceModelExt, RunnableModel, Tensor, TypedFact,
        TypedOp,
    },
};

const LOG_PREFIX: &str = "\u{001B}[34m[mobilenet_onnx]\u{001B}[0m ";

static mut GLOBAL_MODEL: Mutex<
    Option<RunnableModel<TypedFact, Box<dyn TypedOp>, Graph<TypedFact, Box<dyn TypedOp>>>>,
> = Mutex::new(None);

#[yomo::init]
fn init() -> Result<Vec<u32>> {
    println!("{}successufully loading model from file", LOG_PREFIX);
    // https://github.com/onnx/models/blob/main/vision/classification/mobilenet/model/mobilenetv2-7.onnx
    let mut file = OpenOptions::new().read(true).open("./mobilenetv2-7.onnx")?;

    let model = onnx()
        .model_for_read(&mut file)?
        .into_optimized()?
        .into_runnable()?;
    unsafe {
        GLOBAL_MODEL = Some(model).into();
    }
    println!("{}successufully loaded model", LOG_PREFIX);

    Ok(vec![0x33])
}

#[derive(Serialize)]
struct ImageResult {
    score: f32,
    class: i32,
}

#[yomo::handler]
fn handler(input: &[u8]) -> Result<(u32, Vec<u8>)> {
    println!("{}processing image", LOG_PREFIX);
    let image = load_from_memory(input)?.to_rgb8();
    let resized = resize(&image, 224, 224, FilterType::Triangle);
    let image: Tensor = tract_ndarray::Array4::from_shape_fn((1, 3, 224, 224), |(_, c, y, x)| {
        let mean = [0.485, 0.456, 0.406][c];
        let std = [0.229, 0.224, 0.225][c];
        (resized[(x as _, y as _)][c] as f32 / 255.0 - mean) / std
    })
    .into();

    println!("{}running model inference", LOG_PREFIX);
    let model = unsafe { GLOBAL_MODEL.lock().unwrap() };
    let pred = model.as_ref().unwrap().run(tvec!(image.into()))?;
    let best = pred[0]
        .to_array_view::<f32>()
        .unwrap()
        .iter()
        .cloned()
        .zip(2..)
        .max_by(|a, b| a.0.partial_cmp(&b.0).unwrap())
        .unwrap();

    let result = ImageResult {
        score: best.0,
        class: best.1,
    };
    let output = serde_json::to_vec(&result)?;
    println!("{}finished", LOG_PREFIX);

    Ok((0x34, output))
}
