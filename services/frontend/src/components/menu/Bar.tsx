import "./Menu.css"

import React from "react"
import Tab from "./Tab.jsx"
import { Relation } from "../../assets/Relations.js"

export default function Bar({items=[]}: {
    items: Array<Relation>
}) {
    return (
        <div className="bar">
            {items.map((e,i) => (
                <Tab key={i} title={e.from} target={e.to}/>
            ))}
        </div>
    )
}
