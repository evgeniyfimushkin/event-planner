import Tab from "./Tab.jsx"
import "./Menu.css"

export default function Bar() {
    let tabs = ["Мероприятия", "Календарь", "Профиль"];
    return (
        <div className="bar">
            {tabs.map((e,i) => (
                <Tab key={i} title={e} />
            ))}
        </div>
    )
}
