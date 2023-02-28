// project import
import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faGaugeSimple, faShieldAlt, 
    faBuildingShield, 
    faRssSquare, 
    faFilePen, 
    faServer,
    faBan, 
    faObjectGroup,
    faKey,
    faUserPlus,
    faSquarePlus,
    faNewspaper,
    faSquareMinus,
    faPersonCircleMinus,
    faPersonCircleQuestion,
    faFilm,
    faHandshake,
    faPenToSquare,
    faScrewdriverWrench,
    faShapes,
    faGears} from '@fortawesome/free-solid-svg-icons';

// ==============================|| MENU ITEMS ||============================== //

const menuItems = {
    items: [{
        id: 'group-dashboard',
        title: '',
        type: 'group',
        children: [
            {
                id: 'dashboard',
                title: 'Dashboard',
                type: 'item',
                url: '/dashboard',
                icon: () => <FontAwesomeIcon icon={faGaugeSimple} />,
            }
        ]
    }, {
        id: 'rate-limiting-policies',
        title: 'Rate Limiting Policies',
        type: 'group',
        children: [
            {
                id: 'advanced-policies',
                title: 'Advanced Policies',
                type: 'item',
                url: '/advanced-policies',
                icon: () => <FontAwesomeIcon icon={faShieldAlt} />,
            },
            {
                id: 'application-rate-plans',
                title: 'Application Rate Plans',
                type: 'item',
                url: '/application-rate-plans',
                icon: () => <FontAwesomeIcon icon={faBuildingShield} />,
            }, 
            {
                id: 'business-plans',
                title: 'Business Plans',
                type: 'item',
                url: '/business-plans',
                icon: () => <FontAwesomeIcon icon={faRssSquare} />,
            }, 
            // {
            //     id: 'custom-policies',
            //     title: 'Custom Policies',
            //     type: 'item',
            //     url: '/custom-policies',
            //     icon: () => <FontAwesomeIcon icon={faFilePen} />,
            // }, 
            {
                id: 'deny-policies',
                title: 'Deny Policies',
                type: 'item',
                url: '/deny-policies',
                icon: () => <FontAwesomeIcon icon={faBan} />,
            }
        ]
    }, 
    // {
    //     id: 'gateways',
    //     title: '',
    //     type: 'group',
    //     children: [
    //         {
    //             id: 'gateways',
    //             title: 'Gateways',
    //             type: 'item',
    //             url: '/gateways',
    //             icon: () => <FontAwesomeIcon icon={faServer} />,
    //         }
    //     ]
    // },
    {
        id: 'api-categories',
        title: '',
        type: 'group',
        children: [
            {
                id: 'api-categories',
                title: 'API Categories',
                type: 'item',
                url: '/api-categories',
                icon: () => <FontAwesomeIcon icon={faObjectGroup} />,
            }
        ]
    },
    // {
    //     id: 'key-managers',
    //     title: '',
    //     type: 'group',
    //     children: [
    //         {
    //             id: 'key-managers',
    //             title: 'Key Managers',
    //             type: 'item',
    //             url: '/key-managers',
    //             icon: () => <FontAwesomeIcon icon={faKey} />,
    //         }
    //     ]
    // },
    {
        id: 'settings',
        title: 'Settings',
        type: 'group',
        children: [
            {
                id: 'applications',
                title: 'Applications',
                type: 'item',
                url: '/applications',
                icon: () => <FontAwesomeIcon icon={faScrewdriverWrench} />,
            },
            // {
            //     id: 'scope-assignments',
            //     title: 'Documentation',
            //     type: 'item',
            //     url: '/scope-assignments',
            //     icon: StarOutlined,
            //     external: true,
            //     target: true
            // }
        ]
    }]
};

export default menuItems;
