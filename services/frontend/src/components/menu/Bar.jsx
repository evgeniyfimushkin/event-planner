import Tab from "./Tab.jsx"
import "./Menu.css"

export default function Bar({items=[]}) {
    return (
        <div className="bar">
            {items.map((e,i) => (
                <Tab key={i} title={e[0]} target={e[1]}/>
            ))}
        </div>
    )
}
