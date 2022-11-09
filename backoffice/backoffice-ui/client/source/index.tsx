
import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
// scroll bar
import 'simplebar/src/simplebar.css';
// project import
import App from './App';

// ==============================|| MAIN - REACT DOM RENDER  ||============================== //

const container = document.getElementById('root');
const root = createRoot(container!);
root.render(
    <BrowserRouter>
        <App />
    </BrowserRouter>
);
