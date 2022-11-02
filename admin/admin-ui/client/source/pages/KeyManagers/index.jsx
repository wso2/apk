import React from 'react';
import { Route, Routes } from 'react-router-dom';
import ResourceNotFound from 'AppComponents/Base/Errors/ResourceNotFound';
import ListKeyManagers from './ListKeyManagers';
import AddEditKeyManager from './AddEditKeyManager';

/**
 * Render a list
 * @returns {JSX} Header AppBar components.
 */
function KeyManagers() {
    return (
        <Routes>
            <Route exact path='/settings/key-managers' component={ListKeyManagers} />
            <Route exact path='/settings/key-managers/create' component={AddEditKeyManager} />
            <Route exact path='/settings/key-managers/:id' component={AddEditKeyManager} />
            <Route component={ResourceNotFound} />
        </Routes>
    );
}

export default KeyManagers;
