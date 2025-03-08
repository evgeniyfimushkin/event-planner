import "./Menu.css";

export default function Tab({title, target, onClick, isOpen}) {
    return (
        <div className={"tab " + (isOpen && "open" || "")} onClick={onClick || (()=>{})}>
            {target
            && <a href={target}>{title}</a>
            || <span>{title}</span>}
        </div>
    )
}
