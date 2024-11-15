import '../../assets/styles/components/windowControls.css'
import {useState} from "react";
import {WindowHide, WindowMinimise} from "../../../wailsjs/runtime/runtime.js";

function WindowControls() {
    const [isWindows, setIsWindows] = useState(false);

    const closeWindow = () => {
        WindowHide()
    };

    const minimizeWindow = () => {
        WindowMinimise()
    };

    return (
        <div className="controls">
            <div className="header-controls"></div>
            <div className="main-controls">
                <div className="window-controls">
                    <button className="control-button" onClick={minimizeWindow}>–</button>
                    <button className="control-button" onClick={closeWindow}>x</button>
                </div>
            </div>
        </div>
    )
}

export default WindowControls