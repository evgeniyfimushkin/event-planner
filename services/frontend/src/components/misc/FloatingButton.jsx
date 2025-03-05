import "./FloatingButton.css"
import { useNavigate } from "react-router-dom";

export default function FloatingButton({text, target}) {
    const navigate = useNavigate();
    const redirect = () => {
        navigate(target);
    }
    return (
        <input className="floatingButton" type="button" value={text} onClick={redirect} />
    );
}
