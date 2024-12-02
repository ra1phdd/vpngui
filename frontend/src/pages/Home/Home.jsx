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
import ToggleSwitch from "../../components/specific/ToggleSwitch.jsx";
import TrafficMonitor from "../../components/specific/TrafficMonitor.jsx";
import RouteCheckbox from "../../components/specific/DisableRoutesCheckbox.jsx";
import {formatBytes} from "../../utils/formatBytes.js";

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
            await toggleOn();
        } else {
            await toggleOff();
        }
        setIsOn((prev) => !prev);
    };

    const toggleOn = async () => {
        setStatus(VpnStatuses["waitOn"]);
        try {
            const result = await RunXrayAPI.Run();
            if (result !== null) throw new Error(result);

            await waitForConfig(true);
        } catch (error) {
            toast.error(`Ошибка: ${error}`);
        }
    };

    const toggleOff = async () => {
        setStatus(VpnStatuses["waitOff"]);
        try {
            const result = await RunXrayAPI.Kill();
            if (result !== null) throw new Error(result);

            await waitForConfig(false);
        } catch (error) {
            toast.error(`Ошибка: ${error}`);
        }
    };

    const waitForConfig = async (shouldBeActive) => {
        const intervalId = setInterval(async () => {
            const config = await ConfigRepository.GetConfig();
            if (config["ActiveVPN"] === shouldBeActive) {
                setStatus(shouldBeActive ? VpnStatuses["on"] : VpnStatuses["off"]);
                clearInterval(intervalId);
            }
        }, 100);
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

    const handleCheckboxChange = async () => {
        try {
            const result = isChecked
                ? await RoutesXrayAPI.EnableRoutes()
                : await RoutesXrayAPI.DisableRoutes();
            if (result !== null) throw new Error(result);
        } catch (error) {
            toast.error(`Ошибка: ${error}`);
        }
        setIsChecked((prev) => !prev);
    };

    return (
        <>
            <ToggleSwitch
                isOn={isOn}
                status={status}
                onToggle={handleToggle}
                country={Countries[countryCode]}
                ip={ip}
            />
            <TrafficMonitor
                uplink={isTrafficProxyUplink}
                downlink={isTrafficProxyDownlink}
            />
            <RouteCheckbox isChecked={isChecked} onChange={handleCheckboxChange} />
        </>
    );
}

export default PageMain;