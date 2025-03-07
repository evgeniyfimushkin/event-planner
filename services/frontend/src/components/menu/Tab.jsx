import "./Menu.css"

export default function Tab({title, target}) {
    return (
        <div className="tab">
            <a href={target}>{title}</a>
        </div>
    )
}
