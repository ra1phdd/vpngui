import '../../assets/styles/components/windowControls.css'
import {useEffect, useState} from "react";
import * as App from "../../../wailsjs/go/app/App.js";
import {WindowHide, WindowMinimise} from "../../../wailsjs/runtime/runtime.js";

function WindowControls() {
    const [isWindows, setIsWindows] = useState(false);

    useEffect(async () => {
        const result = await App.IsGOOSWindows();
        setIsWindows(result);
    }, []);

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
                {isWindows &&
                    <div className="window-controls">
                        <button className="control-button" onClick={minimizeWindow}>–</button>
                        <button className="control-button" onClick={closeWindow}>x</button>
                    </div>
                }
            </div>
        </div>
    )
}

export default WindowControls