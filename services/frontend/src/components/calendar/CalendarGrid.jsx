import Minicard from "../cards/Minicard";
import "./Calendar.css";

export default function CalendarGrid({events=[]}) {

    const getCurrentDate = () => {
        return new Date();
    }

    const isEventNow = (event, day, hour) => {
        const eventDate = new Date(event.start_time);
        const date = getCurrentDate();
        const start = new Date(date.getFullYear(), date.getMonth(), day, hour);
        const end = new Date(date.getFullYear(), date.getMonth(), day, hour+1);
        const isNow = eventDate > start && eventDate < end;
        return isNow
    }
    const isToday = (day) => {
        const date = getCurrentDate();
        return date.getDate() === day;
    }
    const isNow = (hour) => {
        const date = getCurrentDate();
        return date.getHours() === hour;
    }

    const currentDate = getCurrentDate();
    const firstDay = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1);
    const lastDay = new Date(currentDate.getFullYear(), currentDate.getMonth()+1, 0);
    const firstHour = 0;
    const lastHour = 23;
    const rangeDate = [];
    for (let day = firstDay.getDate(); day <= lastDay.getDate(); day++) {
        rangeDate.push(day);
    }
    const rangeHour = [];
    for (let hour = firstHour; hour <= lastHour; hour++) {
        rangeHour.push(hour);
    }

    const grid = <div className="calendar">
        <table className="calendar">
            <caption>{getCurrentDate().toLocaleString("ru", {month: "long", year: "numeric"})}</caption>
            <thead>
                <tr>
                    <td>\</td>
                    {rangeDate.map((day,di)=><td key={di}>{day}</td>)}
                </tr>
            </thead>
            <tbody>
                {rangeHour.map((hour,hi) => <tr key={hi}>
                    <td>{hour}:00</td>
                    {rangeDate.map((day,di) => <td key={di} className={(isToday(day) && "today" || "") + " " + (isNow(hour) && "now" || "")}>
                        {events.filter(e=>isEventNow(e,day,hour)).map(e=><Minicard event={e} subscribed={true} />)}
                    </td>)}
                </tr>)}
            </tbody>
        </table>
    </div>

    return (
        <>
            {grid}
        </>
    );
}
