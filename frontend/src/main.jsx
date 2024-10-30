import React from 'react'
import {createRoot} from 'react-dom/client'
import App from './App.jsx'
import { XrayProvider } from './contexts/XrayAPI.jsx';

const container = document.getElementById('root')

const root = createRoot(container)

root.render(
    <XrayProvider>
        <App />
    </XrayProvider>
)