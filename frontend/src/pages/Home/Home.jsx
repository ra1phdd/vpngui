import React, { useState, useEffect } from 'react';
import {EnableRoutes, DisableRoutes} from "../../../bindings/vpngui/internal/app/xray-core/routesxrayapi.js";
import {Run, Kill} from "../../../bindings/vpngui/internal/app/xray-core/runxraycore.js";
import {Get} from "../../../bindings/vpngui/internal/app/config/config.js";
import {GetConfig} from "../../../bindings/vpngui/internal/app/repository/configrepository.js";
import {GetTraffic} from "../../../bindings/vpngui/internal/app/stats/traffic.js";
import {Countries} from "../../constants/countries.jsx";
import {VpnStatuses} from "../../constants/vpnStatuses.jsx";
import '@styles/pages/home.css'
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

    const handleToggle = async () => {
        if (!isOn) {
            setStatus(VpnStatuses["waitOn"]);
            await toggleOn();
            setStatus(VpnStatuses["on"]);
        } else {
            setStatus(VpnStatuses["waitOff"]);
            await toggleOff();
            setStatus(VpnStatuses["off"]);
        }
        setIsOn((prev) => !prev);
    };

    const toggleOn = async () => {
        try {
            const result = await Run();
            if (result !== null) throw new Error(result);
        } catch (error) {
            if (error.error !== undefined) {
                toast.error(`Ошибка: ${error.error}`);
            }
        }
    };

    const toggleOff = async () => {
        try {
            const result = await Kill(true);
            if (result !== null) throw new Error(result);
        } catch (error) {
            if (error.error !== undefined) {
                toast.error(`Ошибка: ${error.error}`);
            }
        }
    };

    useEffect(() => {
        const checkVPNStatus = async () => {
            const config = await GetConfig();
            const xray = await Get();

            if (config["ActiveVPN"]) {
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
            const proxyUplink = await GetTraffic("proxy", "uplink");
            const proxyDownlink = await GetTraffic("proxy", "downlink");

            setIsTrafficProxyUplink(formatBytes(proxyUplink));
            setIsTrafficProxyDownlink(formatBytes(proxyDownlink));
        };

        const intervalId = setInterval(fetchTrafficData, 2000);
        return () => clearInterval(intervalId);
    }, []);

    const handleCheckboxChange = async () => {
        try {
            const result = isChecked
                ? await EnableRoutes()
                : await DisableRoutes();
            if (result !== null) throw new Error(result);
        } catch (error) {
            if (error.error !== undefined) {
                toast.error(`Ошибка: ${error.error}`);
            }
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