import '../../assets/styles/components/windowControls.css'
import {useEffect, useState} from "react";
import * as App from "../../../wailsjs/go/app/App.js";
import {WindowHide, WindowMinimise} from "../../../wailsjs/runtime/runtime.js";
import {toast} from "react-toastify";

function WindowControls() {
    const [isWindows, setIsWindows] = useState(false);

    useEffect(async () => {

        try {
            const result = await App.IsGOOSWindows();
            if (result !== null) {
                throw new Error(result);
            }
            setIsWindows(result);
        } catch (error) {
            toast.error(`Ошибка: ${error}`);
        }
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