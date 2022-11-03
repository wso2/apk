// project import
import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faAirbnb, fab } from '@fortawesome/free-brands-svg-icons';
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
                id: 'application-policies',
                title: 'Application Policies',
                type: 'item',
                url: '/application-policies',
                icon: () => <FontAwesomeIcon icon={faBuildingShield} />,
            }, 
            {
                id: 'subscription-policies',
                title: 'Subscription Policies',
                type: 'item',
                url: '/subscription-policies',
                icon: () => <FontAwesomeIcon icon={faRssSquare} />,
            }, 
            {
                id: 'custom-policies',
                title: 'Custom Policies',
                type: 'item',
                url: '/custom-policies',
                icon: () => <FontAwesomeIcon icon={faFilePen} />,
            }, 
            {
                id: 'deny-policies',
                title: 'Deny Policies',
                type: 'item',
                url: '/deny-policies',
                icon: () => <FontAwesomeIcon icon={faBan} />,
            }
        ]
    }, 
    {
        id: 'gateways',
        title: '',
        type: 'group',
        children: [
            {
                id: 'gateways',
                title: 'Gateways',
                type: 'item',
                url: '/gateways',
                icon: () => <FontAwesomeIcon icon={faServer} />,
            }
        ]
    },
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
    {
        id: 'key-managers',
        title: '',
        type: 'group',
        children: [
            {
                id: 'key-managers',
                title: 'Key Managers',
                type: 'item',
                url: '/key-managers',
                icon: () => <FontAwesomeIcon icon={faKey} />,
            }
        ]
    },
    {
        id: 'tasks',
        title: 'Tasks',
        type: 'group',
        children: [
            {
                id: 'user-creation',
                title: 'User Creation',
                type: 'item',
                url: '/user-creation',
                icon: () => <FontAwesomeIcon icon={faUserPlus} />,
            },
            {
                id: 'application-creation',
                title: 'Application Creation',
                type: 'item',
                url: '/application-creation',
                icon: () => <FontAwesomeIcon icon={faSquarePlus} />,
            },
            {
                id: 'application-deletion',
                title: 'Application Deletion',
                type: 'item',
                url: '/application-deletion',
                icon: () => <FontAwesomeIcon icon={faSquareMinus} />,
            },
            {
                id: 'subscription-creation',
                title: 'Subscription Creation',
                type: 'item',
                url: '/subscription-creation',
                icon: () => <FontAwesomeIcon icon={faNewspaper} />,
            },
            {
                id: 'subscription-deletion',
                title: 'Subscription Deletion',
                type: 'item',
                url: '/subscription-deletion',
                icon: () => <FontAwesomeIcon icon={faPersonCircleMinus} />,
            }, 
            {
                id: 'subscription-update',
                title: 'Subscription Update',
                type: 'item',
                url: '/subscription-update',
                icon: () => <FontAwesomeIcon icon={faPersonCircleQuestion} />,
            }, 
            {
                id: 'application-registration',
                title: 'Application Registration',
                type: 'item',
                url: '/application-registration',
                icon: () => <FontAwesomeIcon icon={faHandshake} />,
            }, 
            {
                id: 'api-state-change',
                title: 'API State Change',
                type: 'item',
                url: '/api-state-change',
                icon: () => <FontAwesomeIcon icon={faPenToSquare} />,
            }
        ]
    }, {
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
            {
                id: 'scope-assignments',
                title: 'Scope Assignments',
                type: 'item',
                url: '/scope-assignments',
                icon: () => <FontAwesomeIcon icon={faShapes} />,
            },
            {
                id: 'advanced',
                title: 'Advanced',
                type: 'item',
                url: '/advanced',
                icon: () => <FontAwesomeIcon icon={faGears} />,
            }
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
