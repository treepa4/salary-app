from fastapi import FastAPI
from pydantic import BaseModel
from datetime import date

class Shift(BaseModel):
    shift_id: int
    work_date: date
    hours: float
    rate: float

app = FastAPI()
shifts = []

@app.post("/shifts")
def add_shift(shift: Shift):
    shift.shift_id = len(shifts) + 1
    shifts.append(shift)
    return shift

@app.get("/shifts")
def get_shifts():
    return shifts 

@app.get("/earnings")
def get_earnings():
    return sum(s.hours*s.rate for s in shifts)

@app.delete("/shifts/{shift_id}")
def del_shift(shift_id: int): 
    for s in shifts: 
        if s.shift_id == shift_id: 
            shifts.remove(s)
            break
    return shifts


@app.get("/")
def read_root():
    return {"message":"My salary app"}