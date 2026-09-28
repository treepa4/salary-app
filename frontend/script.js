let calendar;
document.getElementById('shift-form').addEventListener('submit',function(event) {
    event.preventDefault();
    // console.log('ok');
    const workDate = document.getElementById('work-date').value;
    const hours = document.getElementById('hours').value;
    const rate = document.getElementById('rate').value;

    fetch('http://localhost:8000/shifts', {
        method: 'POST',
        headers: {
            'Content-Type' : 'application/json'
                },
        body: JSON.stringify({
            work_date: workDate,
            hours: parseFloat(hours),
            rate: parseFloat(rate)
        })
    })
    .then(response => response.json())
    .then(data => {
        console.log("added: ", data)
        calendar.refetchEvents()
        loadShifts()  
        loadEarnings()  
    })
    .catch(err => console.error("error: ", err));
    });

function loadShifts() {
    fetch('http://localhost:8000/shifts')
    .then(response => response.json())
    .then(data => {
        const list = document.getElementById('shifts-list');
        list.innerHTML = '';
        data.forEach(shift => {
            let li = document.createElement('li');
            li.className = 'flex justify-between items-center border border-line rounded-lg px-4 py-2 font-mono text-sm';
            li.textContent = `#${shift.shift_id} · ${shift.work_date} · ${shift.hours} ч × ${shift.rate} Рупий/ч`;                        
            let button = document.createElement('button');
            button.className = 'text-red-500 hover:text-red-700 text-xs ml-2';
            button.textContent = 'del';
            button.addEventListener('click', function() {
                fetch(
                    `http://localhost:8000/shifts/${shift.shift_id}`, {method: 'DELETE'}
                    )
                .then(() => {
                    loadShifts()
                    calendar.refetchEvents()
                });
            });
        li.appendChild(button);
        list.appendChild(li);
        });
    })}
loadShifts()
function loadEarnings() {
    let earn = document.getElementById('earnings');
    fetch('http://localhost:8000/earnings')
    .then(data => data.json())
    .then(response => {earn.textContent = `earnings: ${response}`;});
}
loadEarnings()

document.addEventListener('DOMContentLoaded', function() {
    const calendarEl = document.getElementById('calendar');
    calendar = new FullCalendar.Calendar(calendarEl, {
        initialView: 'dayGridMonth',
        locale: 'ru',
        events: function(fetchInfo, successCallback, failureCallback) {
            fetch('http://localhost:8000/shifts')
            .then(response => response.json())
            .then(data => {
                const events = data.map(shift => ({
                    title: `${shift.hours}ч.`,
                    start: shift.work_date
                }));
                successCallback(events);
            })
            .catch(failureCallback);
        }

    });
    calendar.render();
});

