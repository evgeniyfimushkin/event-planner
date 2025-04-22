import "./Menu.css"

import React from "react"
import Bar from "./Bar"
import Dropdown from "./Dropdown"
import { menuItems } from "../../assets/Relations"

export default function Menu() {
    return (
        <div className="menu">
            <Bar items={menuItems} />
            <Dropdown items={menuItems} />
        </div>
    )
}
