import React, { useState, useEffect } from 'react';
import * as RoutesXrayAPI from "../../../wailsjs/go/xray_api/RoutesXrayAPI.js";
import * as RunXrayAPI from "../../../wailsjs/go/xray_api/RunXrayAPI.js";
import * as Config from "../../../wailsjs/go/config/Config.js";
import * as ConfigRepository from "../../../wailsjs/go/repository/ConfigRepository.js";
import * as Traffic from "../../../wailsjs/go/stats/Traffic.js";
import {Countries} from "../../constants/countries.jsx";
import {VpnStatuses} from "../../constants/vpnStatuses.jsx";
import '@styles/pages/home.css'
import {useXray} from "../../contexts/XrayAPI.jsx";
import { toast } from 'react-toastify';

function PageMain() {
    const [isOn, setIsOn] = useState(false);
    const [status, setStatus] = useState(VpnStatuses["off"]);
    const [ip, setIP] = useState("");
    const [countryCode, setCC] = useState("");
    const [isChecked, setIsChecked] = useState(false);
    const [isTrafficProxyUplink, setIsTrafficProxyUplink] = useState("0.00 Kb/s");
    const [isTrafficProxyDownlink, setIsTrafficProxyDownlink] = useState("0.00 Kb/s");
    const { isXrayRunning } = useXray();

    const handleToggle = async () => {
        if (!isOn) {
            setStatus(VpnStatuses["waitOn"]);
            try {
                const result = await RunXrayAPI.Run();
                if (result !== null) {
                    throw new Error(result);
                }
            } catch (error) {
                toast.error(`Ошибка: ${error}`);
            }

            const intervalId = setInterval(async () => {
                const config = await ConfigRepository.GetConfig();

                if (config["ActiveVPN"]) {
                    setStatus(VpnStatuses["on"]);
                    clearInterval(intervalId);
                }
            }, 100);
        } else {
            setStatus(VpnStatuses["waitOff"]);
            try {
                const result = await RunXrayAPI.Kill();
                if (result !== null) {
                    throw new Error(result);
                }
            } catch (error) {
                toast.error(`Ошибка: ${error}`);
            }

            const intervalId = setInterval(async () => {
                const config = await ConfigRepository.GetConfig();

                if (!config["ActiveVPN"]) {
                    setStatus(VpnStatuses["off"]);
                    clearInterval(intervalId);
                }
            }, 100)
        }

        setIsOn(prev => !prev);
    };

    useEffect(() => {
        const checkVPNStatus = async () => {
            const config = await ConfigRepository.GetConfig();
            const xray = await Config.Get();

            if (config["ActiveVPN"] && isXrayRunning) {
                setIsOn(true);
                setStatus(VpnStatuses["on"]);
            }

            if (config["DisableRoutes"]) {
                setIsChecked(false);
            }

            for (const outbound of xray["outbounds"]) {
                if (outbound["tag"] === "proxy") {
                    setIP(outbound["settings"]["vnext"][0]["address"]);
                    setCC(outbound["settings"]["vnext"][0]["country_code"]);
                }
            }
        };

        checkVPNStatus();
    }, []);

    useEffect(() => {
        const fetchTrafficData = async () => {
            const proxyUplink = await Traffic.GetTraffic("proxy", "uplink");
            const proxyDownlink = await Traffic.GetTraffic("proxy", "downlink");

            setIsTrafficProxyUplink(formatBytes(proxyUplink));
            setIsTrafficProxyDownlink(formatBytes(proxyDownlink));
        };

        const intervalId = setInterval(fetchTrafficData, 2000);
        return () => clearInterval(intervalId);
    }, []);

    const formatBytes = (bytes) => {
        if (bytes < 0) return `0.00 Kb/s`;
        const kb = bytes / 1024;
        if (kb < 1024) return `${kb.toFixed(2)} Kb/s`;
        const mb = kb / 1024;
        if (mb < 1024) return `${mb.toFixed(2)} Mb/s`;
        const gb = mb / 1024;
        return `${gb.toFixed(2)} Gb/s`;
    };

    const handleCheckboxChange = async () => {
        if (isChecked) {
            try {
                const result = await RoutesXrayAPI.EnableRoutes();
                if (result !== null) {
                    throw new Error(result);
                }
            } catch (error) {
                toast.error(`Ошибка: ${error}`);
            }
        } else {
            try {
                const result = await RoutesXrayAPI.DisableRoutes();
                if (result !== null) {
                    throw new Error(result);
                }
            } catch (error) {
                toast.error(`Ошибка: ${error}`);
            }
        }

        setIsChecked(!isChecked);
    };

    return (
        <>
            <div className="main-controls">
                <p className={`toggle-status ${isOn ? 'on' : 'off'}`}>{status}</p>
                <div className={`toggle-switch ${isOn ? 'on' : 'off'}`} onClick={handleToggle}>
                    <div className="toggle-knob">
                        <span className="icon">⏻</span>
                    </div>
                </div>
                <p className="toggle-country">{isOn && countryCode !== "" ? <>{Countries[countryCode]}</> : <>&nbsp;</>}</p>
                <p className="toggle-ip">&nbsp;{isOn && ip !== "" ? <>IP: {ip}</> : <>&nbsp;</>}</p>
            </div>
            <div className="traffic">
                <div className="traffic-proxy">
                    <div className="traffic-proxy_download">
                        <p className="traffic-title">Загрузка</p>
                        <p className="traffic-speed">▲ {isTrafficProxyUplink}</p>
                    </div>
                    <div className="traffic-proxy_upload">
                        <p className="traffic-title">Отдача</p>
                        <p className="traffic-speed">▼ {isTrafficProxyDownlink}</p>
                    </div>
                </div>
            </div>
            <div className="disable-routes">
                <label className="checkbox-routes">
                    <input type="checkbox" checked={isChecked} onChange={handleCheckboxChange}/>
                    {/*<span className="checkbox"></span>*/ }
                    Отключить маршруты
                </label>
            </div>
        </>
    );
}

export default PageMain;