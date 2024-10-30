import React, {useState, useEffect, useRef} from 'react';
import '@styles/pages/log.css';
import * as Log from "../../../wailsjs/go/log/Log.js";

function PageLog() {
    const [log, setLog] = useState('');
    const logRef = useRef(null);

    useEffect(() => {
        const intervalId = setInterval(async () => {
            const text = await Log.GetLogs();
            setLog(text);
        }, 1000);

        return () => clearInterval(intervalId);
    }, []);

    useEffect(() => {
        if (logRef.current) {
            logRef.current.scrollTop = logRef.current.scrollHeight;
        }
    }, [log]);

    return (
        <div className="log-program">
            {log ? <p ref={logRef}>{log}</p> : <p style={{textAlign: "center"}}>Логи подгружаются...</p>}
        </div>
    );
}

export default PageLog;