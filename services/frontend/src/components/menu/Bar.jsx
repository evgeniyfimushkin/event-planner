import Tab from "./Tab.jsx"
import "./Menu.css"

export default function Bar() {
    let tabs = [
        ["Мероприятия", "/"],
        ["Календарь", "/calendar"],
        ["Профиль",],
    ];
    return (
        <div className="bar">
            {tabs.map((e,i) => (
                <Tab key={i} title={e[0]} target={e[1]}/>
            ))}
        </div>
    )
}
