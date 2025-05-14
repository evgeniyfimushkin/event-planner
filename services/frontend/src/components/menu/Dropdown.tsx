import "./Menu.css"

import React from "react";
import Tab from "./Tab.js"
import { useLocation } from "react-router-dom";
import { useState } from "react";
import { Relation } from "../../assets/Relations.js";

export default function Dropdown({items=[]}: {
    items: Array<Relation>
}) {
    const location = useLocation();
    const currentTab: string = items.find(e=>e.to===location.pathname)?.[0] || "Меню";
    const [showItems, setShowItems] = useState<boolean>(false);
    const toggleState = () => {
        setShowItems(!showItems);
    }
    
    return (
        <div className="dropdown">
            <Tab key={-1} title={currentTab} onClick={toggleState} isOpen={showItems} />
            {showItems && <div className="dropdown-content">
                {items.map((e,i) => (
                    <Tab key={i} title={e.from} target={e.to}/>
                ))}
            </div>}
        </div>
    )
}
