import React, {createContext, useContext, useEffect, useState} from 'react';
import * as ConfigRepository from "../../wailsjs/go/repository/ConfigRepository.js";
import * as RunXrayAPI from "../../wailsjs/go/xray_api/RunXrayAPI.js";

const XrayContext = createContext();

export const XrayProvider = ({ children }) => {
    const [isXrayRunning, setIsXrayRunning] = useState(false);

    useEffect(() => {
        const startXrayIfNeeded = async () => {
            const config = await ConfigRepository.GetConfig();

            if (config["ActiveVPN"] && !isXrayRunning) {
                await RunXrayAPI.Run();
                setIsXrayRunning(true);
            }
        };
        startXrayIfNeeded();
    }, []);

    return (
        <XrayContext.Provider value={{ isXrayRunning, setIsXrayRunning }}>
            {children}
        </XrayContext.Provider>
    );
};

export const useXray = () => useContext(XrayContext);