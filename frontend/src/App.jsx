import React, {lazy, Suspense, useEffect, useState} from 'react';
import Header from "./components/common/Header.jsx";
import {Routes, Route, BrowserRouter} from "react-router-dom";
import WindowControls from "./components/common/windowControls.jsx";
import '@styles/main.css';
import {ToastContainer} from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';

const PageHome = lazy(() => import("./pages/Home/Home.jsx"));
const PageRoutes = lazy(() => import("./pages/Routes/Routes.jsx"));
const PageLog = lazy(() => import("./pages/Log/Log.jsx"));
const PageAccounts = lazy(() => import("./pages/Accounts/Accounts.jsx"));

function App() {
    useEffect(() => {
        const handleSelectStart = (event) => {
            event.preventDefault();
        };

        document.addEventListener('selectstart', handleSelectStart);

        return () => {
            document.removeEventListener('selectstart', handleSelectStart);
        };
    }, []);

    return (
        <BrowserRouter>
            <WindowControls />
            <div className="container">
                <Header />
                <main>
                    <Suspense fallback={<div style={{textAlign: "center"}}>Загрузка...</div>}>
                        <Routes>
                            <Route path="/" element={<PageHome />} />
                            <Route path="/routes" element={<PageRoutes />} />
                            <Route path="/log" element={<PageLog />} />
                            <Route path="/accounts" element={<PageAccounts />} />
                        </Routes>
                    </Suspense>
                </main>
            </div>
            <ToastContainer />
        </BrowserRouter>
    );
}

export default App
