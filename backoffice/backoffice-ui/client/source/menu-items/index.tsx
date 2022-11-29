// project import
import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faShapes } from '@fortawesome/free-solid-svg-icons';

// ==============================|| MENU ITEMS ||============================== //

const menuItems = {
    items: [{
        id: 'group-manage-apis',
        title: 'MANAGE APIS',
        type: 'group',
        children: [
            {
                id: 'manage-apis',
                title: 'Manage APIs',
                type: 'item',
                url: '/manage-apis',
                icon: () => <FontAwesomeIcon icon={faShapes} />,
            }
        ]
    }] 
};

export default menuItems;
