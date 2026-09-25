from database import Base 
from sqlalchemy import Column, Integer, Float, Date

class ShiftDB(Base):
    __tablename__  = "shifts" 
    shift_id = Column(Integer, primary_key = True)
    work_date = Column(Date)
    hours = Column(Float)
    rate = Column(Float)
