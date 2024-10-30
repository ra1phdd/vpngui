import React, { useState, useEffect } from 'react';
import * as RoutesXrayAPI from "../../../wailsjs/go/xray_api/RoutesXrayAPI.js";
import * as RunXrayAPI from "../../../wailsjs/go/xray_api/RunXrayAPI.js";
import * as Config from "../../../wailsjs/go/config/Config.js";
import * as ConfigRepository from "../../../wailsjs/go/repository/ConfigRepository.js";
import {Countries} from "../../constants/countries.jsx";
import {VpnStatuses} from "../../constants/vpnStatuses.jsx";
import '@styles/pages/home.css'
import {useXray} from "../../contexts/XrayAPI.jsx";

function PageMain() {
    const [isOn, setIsOn] = useState(false);
    const [status, setStatus] = useState(VpnStatuses["off"]);
    const [ip, setIP] = useState("");
    const [countryCode, setCC] = useState("");
    const [isChecked, setIsChecked] = useState(false);
    const { isXrayRunning } = useXray();

    const handleToggle = async () => {
        if (!isOn) {
            setStatus(VpnStatuses["waitOn"]);
            await RunXrayAPI.Run();
            const intervalId = setInterval(async () => {
                const config = await ConfigRepository.GetConfig();

                if (config["ActiveVPN"]) {
                    setStatus(VpnStatuses["on"]);
                    clearInterval(intervalId);
                }
            }, 100);
        } else {
            setStatus(VpnStatuses["waitOff"]);
            await RunXrayAPI.Kill();
            const intervalId = setInterval(async () => {
                const config = await ConfigRepository.GetConfig();
                console.log(config)

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
            const xray = await Config.GetXray();

            if (config["ActiveVPN"] && isXrayRunning) {
                setIsOn(true);
                setStatus(VpnStatuses["on"]);
            }

            if (config["DisableRoutes"]) {
                setIsChecked(true);
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

    const handleCheckboxChange = () => {
        if (isChecked) {
            RoutesXrayAPI.EnableRoutes()
        } else {
            RoutesXrayAPI.DisableRoutes()
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
            <div className="disable-routes">
                <input type="checkbox" checked={isChecked} onChange={handleCheckboxChange}/>
                <span>Отключить маршруты</span>
            </div>
        </>
    );
}

export default PageMain;