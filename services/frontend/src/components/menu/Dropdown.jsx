import Tab from "./Tab.jsx"
import "./Menu.css"
import { useLocation } from "react-router-dom";
import { useState } from "react";

export default function Dropdown({items=[]}) {
    const location = useLocation();
    const currentTab = items.find(e=>e[1]===location.pathname)?.[0] || "Меню";
    const [showItems, setShowItems] = useState(false);
    const toggleState = () => {
        setShowItems(!showItems);
    }
    return (
        <div className="dropdown">
            <Tab key={-1} title={currentTab} onClick={toggleState} isOpen={showItems} />
            {showItems && <div className="dropdown-content">
                {items.map((e,i) => (
                    <Tab key={i} title={e[0]} target={e[1]}/>
                ))}
            </div>}
        </div>
    )
}
