import React, {useEffect, useState} from 'react';
import '@styles/pages/routes.css';
import * as RoutesXrayAPI from "../../../wailsjs/go/xray_api/RoutesXrayAPI.js";

function PageRoutes() {
    const [isChecked, setIsChecked] = useState(false);
    const [blacklistDomains, setBlacklistDomains] = useState('');
    const [blacklistIPs, setBlacklistIPs] = useState('');
    const [blacklistPorts, setBlacklistPorts] = useState('');
    const [whilelistDomains, setWhitelistDomains] = useState('');
    const [whilelistIPs, setWhitelistIPs] = useState('');
    const [whilelistPorts, setWhitelistPorts] = useState('');

    const [inputBlacklistDomains, setInputBlacklistDomains] = useState('');
    const [inputBlacklistIPs, setInputBlacklistIPs] = useState('');
    const [inputBlacklistPorts, setInputBlacklistPorts] = useState('');
    const [inputWhitelistDomains, setInputWhitelistDomains] = useState('');
    const [inputWhitelistIPs, setInputWhitelistIPs] = useState('');
    const [inputWhitelistPorts, setInputWhitelistPorts] = useState('');

    const getBlackList = async () => {
        const domains = await RoutesXrayAPI.GetDomain("blacklist");
        const ips = await RoutesXrayAPI.GetIP("blacklist");
        const ports = await RoutesXrayAPI.GetPort("blacklist");

        setBlacklistDomains(domains);
        setBlacklistIPs(ips);
        setBlacklistPorts(ports);
    };

    const getWhiteList = async () => {
        const domains = await RoutesXrayAPI.GetDomain("whitelist");
        const ips = await RoutesXrayAPI.GetIP("whitelist");
        const ports = await RoutesXrayAPI.GetPort("whitelist");

        setWhitelistDomains(domains);
        setWhitelistIPs(ips);
        setWhitelistPorts(ports);
    };

    useEffect(() => {
        getBlackList();
        getWhiteList();
    }, []);

    const handleCheckboxChange = () => {
        if (isChecked) {
            RoutesXrayAPI.EnableBlackList();
        } else {
            RoutesXrayAPI.DisableBlackList();
        }

        setIsChecked(!isChecked);
    };

    const handleInputBlacklistDomains = (event) => {
        setInputBlacklistDomains(event.target.value);
    };
    const handleInputBlacklistIPs = (event) => {
        setInputBlacklistIPs(event.target.value);
    };
    const handleInputBlacklistPorts = (event) => {
        setInputBlacklistPorts(event.target.value);
    };
    const handleInputWhitelistDomains = (event) => {
        setInputWhitelistDomains(event.target.value);
    };
    const handleInputWhitelistIPs = (event) => {
        setInputWhitelistIPs(event.target.value);
    };
    const handleInputWhitelistPorts = (event) => {
        setInputWhitelistPorts(event.target.value);
    };

    const handleAddDomainInBlacklist = async () => {
        await RoutesXrayAPI.AddDomain("blacklist", inputBlacklistDomains);
        setInputBlacklistDomains("");

        getBlackList();
    };
    const handleAddIPInBlacklist = async () => {
        await RoutesXrayAPI.AddIP("blacklist", inputBlacklistIPs);
        setInputBlacklistIPs("");

        getBlackList();
    };
    const handleAddPortInBlacklist = async () => {
        await RoutesXrayAPI.AddPort("blacklist", inputBlacklistPorts);
        setInputBlacklistPorts("");

        getBlackList();
    };
    const handleAddDomainInWhitelist = async () => {
        await RoutesXrayAPI.AddDomain("whitelist", inputWhitelistDomains);
        setInputWhitelistDomains("");

        getWhiteList();
    };
    const handleAddIPInWhitelist = async () => {
        await RoutesXrayAPI.AddIP("whitelist", inputWhitelistIPs);
        setInputWhitelistIPs("");

        getWhiteList();
    };
    const handleAddPortInWhitelist = async () => {
        await RoutesXrayAPI.AddPort("whitelist", inputWhitelistPorts);
        setInputWhitelistPorts("");

        getWhiteList();
    };

    const handleDelDomainInBlacklist = async () => {
        await RoutesXrayAPI.DelDomain("blacklist", inputBlacklistDomains);
        setInputBlacklistDomains("");

        getBlackList();
    };
    const handleDelIPInBlacklist = async () => {
        await RoutesXrayAPI.DelIP("blacklist", inputBlacklistIPs);
        setInputBlacklistIPs("");

        getBlackList();
    };
    const handleDelPortInBlacklist = async () => {
        await RoutesXrayAPI.DelPort("blacklist", inputBlacklistPorts);
        setInputBlacklistPorts("");

        getBlackList();
    };
    const handleDelDomainInWhitelist = async () => {
        await RoutesXrayAPI.DelDomain("whitelist", inputWhitelistDomains);
        setInputWhitelistDomains("");

        getWhiteList();
    };
    const handleDelIPInWhitelist = async () => {
        await RoutesXrayAPI.DelIP("whitelist", inputWhitelistIPs);
        setInputWhitelistIPs("");

        getWhiteList();
    };
    const handleDelPortInWhitelist = async () => {
        await RoutesXrayAPI.DelPort("whitelist", inputWhitelistPorts);
        setInputWhitelistPorts("");

        getWhiteList();
    };

    return (
        <>
            <div className="mode-list">
                <span className={`black-list ${isChecked ? 'inactive' : 'active'}`}>черный список</span>
                <input
                    className="react-switch-checkbox"
                    id={`react-switch-new`}
                    type="checkbox"
                    checked={isChecked}
                    onChange={handleCheckboxChange}
                />
                <label
                    className="react-switch-label"
                    htmlFor={`react-switch-new`}
                >
                    <span className={`react-switch-button`}/>
                </label>
                <span className={`white-list ${isChecked ? 'active' : 'inactive'}`}>белый список</span>
            </div>
            <div className="routes">
                <div className="routes-domain">
                    <section className={`routes-blacklist ${isChecked ? 'inactive' : 'active'}`}>
                        <textarea placeholder="Здесь будут находиться ваши домены..." value={blacklistDomains} disabled={true}></textarea>
                        <div className="routes-action">
                            <input placeholder="Добавить домен..." disabled={isChecked} value={inputBlacklistDomains} onChange={handleInputBlacklistDomains}/>
                            <a className={`routes-add ${isChecked ? 'inactive' : 'active'}`} onClick={handleAddDomainInBlacklist}>+</a>
                            <a className={`routes-del ${isChecked ? 'inactive' : 'active'}`} onClick={handleDelDomainInBlacklist}>-</a>
                        </div>
                    </section>
                    <section className={`routes-whitelist ${isChecked ? 'active' : 'inactive'}`}>
                    <textarea placeholder="Здесь будут находиться ваши домены..." value={whilelistDomains} disabled={true}></textarea>
                        <div className="routes-action">
                            <input placeholder="Добавить домен..." disabled={!isChecked} value={inputWhitelistDomains} onChange={handleInputWhitelistDomains}/>
                            <a className={`routes-add ${isChecked ? 'active' : 'inactive'}`} onClick={handleAddDomainInWhitelist}>+</a>
                            <a className={`routes-del ${isChecked ? 'active' : 'inactive'}`} onClick={handleDelDomainInWhitelist}>-</a>
                        </div>
                    </section>
                </div>
                <div className="routes-ip">
                    <section className={`routes-blacklist ${isChecked ? 'inactive' : 'active'}`}>
                        <textarea placeholder="Здесь будут находиться ваши IP-адреса..." value={blacklistIPs} disabled={true}></textarea>
                        <div className="routes-action">
                            <input placeholder="Добавить IP-адрес..." disabled={isChecked} value={inputBlacklistIPs} onChange={handleInputBlacklistIPs}/>
                            <a className={`routes-add ${isChecked ? 'inactive' : 'active'}`} onClick={handleAddIPInBlacklist}>+</a>
                            <a className={`routes-del ${isChecked ? 'inactive' : 'active'}`} onClick={handleDelIPInBlacklist}>-</a>
                        </div>
                    </section>
                    <section className={`routes-whitelist ${isChecked ? 'active' : 'inactive'}`}>
                        <textarea placeholder="Здесь будут находиться ваши IP-адреса..." value={whilelistIPs} disabled={true}></textarea>
                        <div className="routes-action">
                            <input placeholder="Добавить IP-адрес..." disabled={!isChecked} value={inputWhitelistIPs} onChange={handleInputWhitelistIPs}/>
                            <a className={`routes-add ${isChecked ? 'active' : 'inactive'}`} onClick={handleAddIPInWhitelist}>+</a>
                            <a className={`routes-del ${isChecked ? 'active' : 'inactive'}`} onClick={handleDelIPInWhitelist}>-</a>
                        </div>
                    </section>
                </div>
                <div className="routes-port">
                    <section className={`routes-blacklist ${isChecked ? 'inactive' : 'active'}`}>
                        <textarea className="routes-textarea" placeholder="Здесь будут находиться ваши порты..." value={blacklistPorts} disabled={true}></textarea>
                        <div className="routes-action">
                            <input placeholder="Добавить порт..." disabled={isChecked} value={inputBlacklistPorts} onChange={handleInputBlacklistPorts}/>
                            <a className={`routes-add ${isChecked ? 'inactive' : 'active'}`} onClick={handleAddPortInBlacklist}>+</a>
                            <a className={`routes-del ${isChecked ? 'inactive' : 'active'}`} onClick={handleDelPortInBlacklist}>-</a>
                        </div>
                    </section>
                    <section className={`routes-whitelist ${isChecked ? 'active' : 'inactive'}`}>
                        <textarea className="routes-textarea" placeholder="Здесь будут находиться ваши порты..." value={whilelistPorts} disabled={true}></textarea>
                        <div className="routes-action">
                            <input placeholder="Добавить порт..." disabled={!isChecked} value={inputWhitelistPorts} onChange={handleInputWhitelistPorts}/>
                            <a className={`routes-add ${isChecked ? 'active' : 'inactive'}`} onClick={handleAddPortInWhitelist}>+</a>
                            <a className={`routes-del ${isChecked ? 'active' : 'inactive'}`} onClick={handleDelPortInWhitelist}>-</a>
                        </div>
                    </section>
                </div>
            </div>
        </>
    );
}

export default PageRoutes
