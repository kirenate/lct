from dataclasses import dataclass
from fastapi import FastAPI, File, UploadFile
import services.batch_processing_service
from loguru import logger



@dataclass
class Presentation:
    service: services.batch_processing_service.Service

    async def upload_batch(self, file: UploadFile) -> Exception | None:
        try:
            await self.service.save_batch(bytes(await file.read()), str(file.filename))
            await file.seek(0)
        except Exception as e:
            logger.error(e)
            return e
        return

    async def delete_batch(self, batch_name :str) -> Exception | None:
        try:
            await self.service.delete_batch(batch_name)
        except Exception as e:
            logger.error(e)
            return e
        return