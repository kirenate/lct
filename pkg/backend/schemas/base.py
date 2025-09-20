from typing import Any

import humps
from fastapi.encoders import jsonable_encoder
from pydantic import BaseModel, ConfigDict


# * Pure pydantic model without any alias generator
class PureBaseModel(BaseModel):
    def jsonable_encoder(self, **kwargs: Any) -> Any:  # noqa: ANN401
        return jsonable_encoder(self, **kwargs)

    model_config = ConfigDict(populate_by_name=True, from_attributes=True)


# * Camel alias generator model
class CamelizedBaseModel(PureBaseModel):
    model_config = ConfigDict(populate_by_name=True, from_attributes=True, alias_generator=humps.camelize)
