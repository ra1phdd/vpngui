import React, { useState, useEffect } from 'react';
import './App.css';
import * as xray_api from "../wailsjs/go/xray_api/XrayAPI.js";
import * as Config from "../wailsjs/go/config/Config.js";

function App() {
    const [isVPNActive, setVPNActive] = useState(false);
    const [status, setStatus] = useState("VPN выключен");

    useEffect(() => {
        // Инициализация состояния VPN из конфигурации
        const checkVPNStatus = async () => {
            const config = await Config.GetJSON();
            setVPNActive(config["active-vpn"]);
        };

        checkVPNStatus();
    }, []);

    useEffect(() => {
        if (isVPNActive) {
            setStatus("VPN включен");
        } else {
            setStatus("VPN выключен");
        }
    }, [isVPNActive]);

    const toggleVPN = async () => {
        if (isVPNActive) {
            setStatus("VPN выключается...");
            await xray_api.Kill();
            const intervalId = setInterval(async () => {
                const config = await Config.GetJSON();

                if (!config["active-vpn"]) {
                    setVPNActive(false);
                    clearInterval(intervalId);
                }
            }, 100);
        } else {
            setStatus("VPN включается...");
            await xray_api.Run();
            const intervalId = setInterval(async () => {
                const config = await Config.GetJSON();

                if (config["active-vpn"]) {
                    setVPNActive(true);
                    clearInterval(intervalId);
                }
            }, 100);
        }
    };

    const restartVPN = async () => {
        setStatus("VPN перезапускается...");
        await xray_api.Kill();
        while (true) {
            const config = await Config.GetJSON();
            if (!config["active-vpn"]) {
                break;
            }
            await new Promise((resolve) => setTimeout(resolve, 100));
        }

        await xray_api.Run();
        while (true) {
            const config = await Config.GetJSON();
            if (config["active-vpn"]) {
                setStatus("VPN включен");
                break;
            }
            await new Promise((resolve) => setTimeout(resolve, 100));
        }
    };

    return (
        <div style={{textAlign: 'center'}}>
            <h1>{status}</h1>
            <button onClick={toggleVPN}>
                {isVPNActive ? "Подключено" : "Отключено"}
            </button>
            <button onClick={restartVPN}>
                Перезапустить VPN
            </button>
        </div>
    );
}

export default App
