from ast import With
from shared.ulid import ulid_as_uuid
from sqlalchemy import Column, MetaData
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.ext.declarative import declarative_base

metadata_obj = MetaData(schema="public")

Base = declarative_base()
Base.metadata = metadata_obj
