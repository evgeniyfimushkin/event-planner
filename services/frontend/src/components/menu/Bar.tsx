import React from "react"

import Tab from "./Tab.jsx"
import "./Menu.css"

export default function Bar({items=[]}: {
    items: Array<[string, string?]>
}) {
    return (
        <div className="bar">
            {items.map((e,i) => (
                <Tab key={i} title={e[0]} target={e[1]}/>
            ))}
        </div>
    )
}
