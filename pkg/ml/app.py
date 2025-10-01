import uvicorn
from fastapi import FastAPI


def create_app() -> FastAPI:
    app = FastAPI()

    @app.post("/process")
    def process() -> None:
        return None

    return app


def main() -> None:
    app = create_app()
    uvicorn.run(app, host="0.0.0.0", port=8080, log_level="debug")


if __name__ == "__main__":
    main()
