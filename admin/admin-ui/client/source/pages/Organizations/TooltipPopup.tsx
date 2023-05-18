/*
 * Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import React from 'react';
import Tooltip, { TooltipProps, tooltipClasses } from '@mui/material/Tooltip';
import { styled } from '@mui/material/styles';
import Link from '@mui/material/Link';

interface Organization {
    name: string;
    displayName: string;
    organizationClaimValue: string;
    serviceNamespaces: string;
    production: string;
    sandbox: string;
}

interface TooltipPopupProps {
    org: Organization;
}

const CustomWidthTooltip = styled(({ className, ...props }: TooltipProps) => (
    <Tooltip {...props} classes={{ popper: className }} />
))(({ theme }) => ({
    [`& .${tooltipClasses.tooltip}`]: {
        maxWidth: 500,
        backgroundColor: theme.palette.common.white,
        color: 'rgba(0, 0, 0, 0.87)',
        boxShadow: theme.shadows[1],
    },
}));

const TooltipPopup: React.FC<TooltipPopupProps> = ({ org }) => {

    const title = org && (
        <div>
            Name: {org.name}<br />
            Display Name: {org.displayName}<br />
            Organization Claim Value: {org.organizationClaimValue}<br />
            Service Namespaces: {org.serviceNamespaces}<br />
            Production Endpoints: {org.production}<br />
            Sandbox Endpoints: {org.sandbox}<br />
        </div>
    );

    return (
        <div>
            <CustomWidthTooltip
                PopperProps={{
                    disablePortal: true,
                }}
                disableFocusListener
                disableTouchListener
                title={title}
                arrow
                placement="right-start"
            >
                <Link underline="none" style={{ color: 'black', display: 'inline-flex', justifyContent: 'flex-start', marginLeft: '10px' }}>
                    {org.displayName}
                </Link>
            </CustomWidthTooltip>
        </div>
    );
}

export default TooltipPopup;