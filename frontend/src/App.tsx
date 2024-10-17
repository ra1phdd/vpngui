import {useEffect, useState} from 'react';
//import logo from './assets/images/logo-universal.png';
import './App.css';
import { Run, Kill } from "../wailsjs/go/xray_api/XrayAPI";
import { GetJSON } from "../wailsjs/go/config/Config";

function App() {
    const [vpnStatus, setVpnStatus] = useState<string>('VPN выключен');
    const [activeVPN, setActiveVPN] = useState<boolean>(false);

    // Эффект для синхронизации состояния VPN
    useEffect(() => {
        if (activeVPN) {
            setVpnStatus('VPN включен');
        } else {
            setVpnStatus('VPN выключен');
        }
    }, [activeVPN]);

    // Включение VPN
    const enableVPN = async () => {
        console.log("ахуй")
        setVpnStatus('VPN включается...');
        await Run();
        console.log("хуй")

        const checkVPNStatus = setInterval(async () => {
            const config = await GetJSON(); // Предполагаем, что функция возвращает конфиг
            console.log(config)
            if (config["active-vpn"]) {
                console.log("хуйня")
                setActiveVPN(true);
                clearInterval(checkVPNStatus);
            }
        }, 100);
    };

    // Выключение VPN
    const disableVPN = async () => {
        setVpnStatus('VPN выключается...');
        await Kill();

        const checkVPNStatus = setInterval(async () => {
            const config = await GetJSON();
            if (!config["active-vpn"]) {
                setActiveVPN(false);
                clearInterval(checkVPNStatus);
            }
        }, 100);
    };

    // Перезапуск VPN
    const restartVPN = async () => {
        setVpnStatus('VPN перезапускается...');
        await Kill();

        const checkVPNStatus = setInterval(async () => {
            const config = await GetJSON();
            if (!config["active-vpn"]) {
                clearInterval(checkVPNStatus);
                await enableVPN(); // После выключения перезапускаем
            }
        }, 100);
    };

    return (
        <div className="App">
            <h1>{vpnStatus}</h1>
            <div className="button-container">
                <button onClick={enableVPN} disabled={activeVPN}>Включить VPN</button>
                <button onClick={disableVPN} disabled={!activeVPN}>Выключить VPN</button>
                <button onClick={restartVPN}>Перезапустить VPN</button>
            </div>
        </div>
    );
}

export default App
