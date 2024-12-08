import React, {useEffect, useState} from 'react';
import '@styles/pages/routes.css';
import ModeSwitch from "../../components/specific/ModeListSwitch.jsx";
import {RoutesSection} from "../../components/specific/RoutesComponents.jsx";
import {
    DisableBlackList, EnableBlackList,
    GetDomain, GetPort, GetIP,
} from "../../../bindings/vpngui/internal/app/xray-core/routesxrayapi.js";
import {RoutesXrayCore} from "../../../bindings/vpngui/internal/app/xray-core/index.js";

function PageRoutes() {
    const [isChecked, setIsChecked] = useState(false);
    const [data, setData] = useState({
        blacklistDomains: '',
        blacklistIPs: '',
        blacklistPorts: '',
        whitelistDomains: '',
        whitelistIPs: '',
        whitelistPorts: '',
    });
    const [inputs, setInputs] = useState({
        blacklistDomains: '',
        blacklistIPs: '',
        blacklistPorts: '',
        whitelistDomains: '',
        whitelistIPs: '',
        whitelistPorts: '',
    });

    const fetchList = async (listType) => {
        const [domains, ips, ports] = await Promise.all([
            GetDomain(listType),
            GetIP(listType),
            GetPort(listType),
        ]);

        setData((prev) => ({
            ...prev,
            [`${listType}Domain`]: domains,
            [`${listType}IP`]: ips,
            [`${listType}Port`]: ports,
        }));
    };

    useEffect(() => {
        void fetchList('blacklist');
        void fetchList('whitelist');
    }, []);

    const handleCheckboxChange = async () => {
        if (isChecked) {
            await EnableBlackList();
        } else {
            await DisableBlackList();
        }
        setIsChecked(!isChecked);
    };

    const handleInputChange = (type, value) => {
        setInputs((prev) => ({ ...prev, [type]: value }));
    };

    const handleUpdateList = async (action, listType, itemType) => {
        const inputKey = `${listType}${itemType}`;
        await RoutesXrayCore[`${action}${itemType}`](listType, inputs[inputKey]);
        setInputs((prev) => ({ ...prev, [inputKey]: '' }));
        await fetchList(listType);
    };

    return (
        <>
            <ModeSwitch isChecked={isChecked} handleCheckboxChange={handleCheckboxChange} />
            <div className="routes">
                {['Domain', 'IP', 'Port'].map((item) => (
                    <div key={item.toLowerCase()} className={`routes-${item.toLowerCase()}`}>
                        <RoutesSection
                            isChecked={isChecked}
                            mode="blacklist"
                            placeholder={`Здесь будут находиться ваши ${item}s...`}
                            value={data[`blacklist${item}`]}
                            inputValue={inputs[`blacklist${item}`]}
                            handleInputChange={(e) => handleInputChange(`blacklist${item}`, e.target.value)}
                            handleAdd={() => handleUpdateList('Add', 'blacklist', item)}
                            handleDelete={() => handleUpdateList('Del', 'blacklist', item)}
                        />
                        <RoutesSection
                            isChecked={isChecked}
                            mode="whitelist"
                            placeholder={`Здесь будут находиться ваши ${item}...`}
                            value={data[`whitelist${item}`]}
                            inputValue={inputs[`whitelist${item}`]}
                            handleInputChange={(e) => handleInputChange(`whitelist${item}`, e.target.value)}
                            handleAdd={() => handleUpdateList('Add', 'whitelist', item)}
                            handleDelete={() => handleUpdateList('Del', 'whitelist', item)}
                        />
                    </div>
                ))}
            </div>
        </>
    );
}

export default PageRoutes
