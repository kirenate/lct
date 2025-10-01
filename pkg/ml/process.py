# pip install paddlepaddle
# pip install paddleocr

import time
from dataclasses import dataclass
from pathlib import Path

import cv2
import numpy as np
from loguru import logger
from paddleocr import PaddleOCR
from pydantic import BaseModel


class Prediction(BaseModel): ...


@dataclass
class OCR:
    def __post_init__(self) -> None:
        self._ocr = PaddleOCR(
            text_recognition_model_name="PP-OCRv5_server_rec",  # неизменяемое значение
            text_recognition_model_dir="../.data/inference_2",  # путь до папки с файлами, которые на гугл диске
            use_doc_orientation_classify=False,
            use_doc_unwarping=False,
            use_textline_orientation=False,
        )  # вот этот инстанс один раз создается и всё
        logger.info("ocr.initialized")

    @staticmethod
    def _preprocess_image(raw_image: bytes) -> tuple[cv2.typing.MatLike, float]:
        nparr = np.frombuffer(raw_image, dtype=np.uint8)
        image = cv2.imdecode(nparr, flags=1)

        standard_size = 4000
        gray_img = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        img = cv2.GaussianBlur(gray_img, (5, 5), 0)
        res_img = cv2.adaptiveThreshold(img, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C, cv2.THRESH_BINARY, 25, 10)
        dims_image = cv2.cvtColor(res_img, cv2.COLOR_GRAY2RGB)
        width = dims_image.shape[1]
        height = dims_image.shape[0]
        if height >= width and height > standard_size:
            k_transform = height / standard_size
            dim = (int(width / k_transform), standard_size)
        elif width >= height and width > standard_size:
            k_transform = width / standard_size
            dim = (standard_size, int(height / k_transform))
        else:
            return dims_image, 1
        result_image = cv2.resize(dims_image, dim, interpolation=cv2.INTER_AREA)
        return result_image, k_transform

    def predict(self, raw_image: bytes) -> None:
        t0 = time.monotonic()
        logger.info(f"processing.image len={len(raw_image)/1024/1024:2f}MB")
        img, _ = self._preprocess_image(raw_image)
        logger.info(f"image.reshaped {img.shape=}")

        result = self._ocr.predict(img)  # это уже для каждого файла вызывается
        for res in result:
            res.print()
            res.save_to_img("output")
            res.save_to_json("output")

        result[0]["rec_polys"]  # массив с полигонами
        result[0]["rec_texts"]  # массив строк полученных
        result[0]["rec_scores"]  # массив со значениями уверенности для каждой строки
        logger.info(
            f"image.processed elapsed={time.monotonic() - t0:2f}f "
            f"{result[0]["rec_polys"]}, "
            f"{result[0]["rec_texts"]}, {result[0]["rec_scores"]}"
        )


if __name__ == "__main__":
    with open(Path("test.jpg"), "rb") as file:
        OCR().predict(file.read())
