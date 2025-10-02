from typing import Any

import uvicorn
from fastapi import FastAPI, UploadFile
from process import OCR


def create_app() -> FastAPI:
    app = FastAPI()
    ocr = OCR()

    @app.post("/process")
    async def process(image: UploadFile) -> tuple[Any, Any, Any]:
        return ocr.predict(await image.read())

    return app


def main() -> None:
    app = create_app()
    uvicorn.run(app, host="0.0.0.0", port=8080, log_level="debug")


if __name__ == "__main__":
    main()
