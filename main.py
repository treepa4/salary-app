from database import engine, get_db
from models import Base, ShiftDB
from fastapi import FastAPI, Depends
from pydantic import BaseModel
from datetime import date
from sqlalchemy.orm import Session

class Shift(BaseModel):
    shift_id: int
    work_date: date
    hours: float
    rate: float

app = FastAPI()
Base.metadata.create_all(bind=engine)

@app.post("/shifts")
def add_shift(shift: Shift, db: Session = Depends(get_db)):
    shiftDb = ShiftDB(rate = shift.rate, work_date = shift.work_date, hours = shift.hours)
    db.add(shiftDb)
    db.commit()
    db.refresh(shiftDb)
    return shiftDb

@app.get("/shifts")
def get_shifts(db: Session = Depends(get_db)):
    return db.query(ShiftDB).all()

@app.get("/earnings")
def get_earnings(db: Session = Depends(get_db)):
    shifts = db.query(ShiftDB).all()
    return sum(s.hours*s.rate for s in shifts)

@app.delete("/shifts/{shift_id}")
def del_shift(shift_id: int, db: Session = Depends(get_db)): 
    shifts = db.query(ShiftDB).filter(ShiftDB.shift_id == shift_id).first()
    if shifts is None: 
        return {"error":"id not in database"}
    db.delete(shifts)
    db.commit()
    return {"message":"Succesfull deleted!"}

    

@app.get("/")
def read_root():
    return {"message":"My salary app"}