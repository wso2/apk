/*
 * Copyright (c) 2020, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
import { FormattedMessage, useIntl } from 'react-intl';
import { Link as RouterLink } from 'react-router-dom';
import { Card } from '@mui/material';
import Box from '@mui/material/Box';
import CardContent from '@mui/material/CardContent';
import Divider from '@mui/material/Divider';
import Link from '@mui/material/Link';
import { makeStyles } from 'tss-react/mui';
import Typography from '@mui/material/Typography';
import PolicyIcon from '@ant-design/icons/GroupOutlined';
import AssignmentTurnedInIcon from '@ant-design/icons/GroupOutlined';
import AssignmentIcon from '@ant-design/icons/GroupOutlined';
import PhonelinkSetupIcon from '@ant-design/icons/GroupOutlined';
import { useAppContext } from 'context/AppContext';

const useStyles = makeStyles()(() => {
    return {
        root: {
            minWidth: 275,
            minHeight: 270,
            textAlign: 'center',
        },
        title: {
            fontSize: 20,
            fontWeight: 'fontWeightBold',
        },
        pos: {
            marginBottom: 12,
        },
    }
});

/**
 * Render progress inside a container centering in the container.
 * @returns {JSX} Loading animation.
 */
export default function RateLimitingCard() {
    const {classes} = useStyles();
    const { isSuperTenant } = useAppContext();
    const intl = useIntl();
    const selectedRateLimitingPolicies = [
        {
            name: intl.formatMessage({
                id: 'Dashboard.rateLimiting.card.advancedPolicies.name',
                defaultMessage: 'Advanced Policies',
            }),
            description: intl.formatMessage({
                id: 'Dashboard.rateLimiting.card.advancedPolicies.description',
                defaultMessage: 'Control access per API or API resource using advanced rules',
            }),
            icon: <PolicyIcon color='inherit' fontSize='small' />,
            path: '/throttling/advanced',
        },
        {
            name: intl.formatMessage({
                id: 'Dashboard.rateLimiting.card.applicationPolicies.name',
                defaultMessage: 'Application Policies',
            }),
            description: intl.formatMessage({
                id: 'Dashboard.rateLimiting.card.applicationPolicies.description',
                defaultMessage: 'Applicable per access token generated for an application',
            }),
            icon: <PhonelinkSetupIcon color='inherit' fontSize='small' />,
            path: '/throttling/application',
        },
        {
            name: intl.formatMessage({
                id: 'Dashboard.rateLimiting.card.subscriptionPolicies.name',
                defaultMessage: 'Subscription Policies',
            }),
            description: intl.formatMessage({
                id: 'Dashboard.rateLimiting.card.subscriptionPolicies.description',
                defaultMessage: 'Control access per Subscription',
            }),
            icon: <AssignmentTurnedInIcon color='inherit' fontSize='small' />,
            path: '/throttling/subscription',
        },
        {
            name: intl.formatMessage({
                id: 'Dashboard.rateLimiting.card.customPolicies.name',
                defaultMessage: 'Custom Policies',
            }),
            description: intl.formatMessage({
                id: 'Dashboard.rateLimiting.card.customPolicies.description',
                defaultMessage: 'Allows system administrators to define dynamic '
                    + 'rules for specific use cases, which are applied globally across all tenants.',
            }),
            icon: <AssignmentIcon color='inherit' fontSize='small' />,
            path: '/throttling/custom',
            id: 'Custom Policies',
        },
    ];
    let selectedPolicies = selectedRateLimitingPolicies;
    if (!isSuperTenant) {
        selectedPolicies = selectedRateLimitingPolicies.filter((item) => item.id !== 'Custom Policies');
    }

    return (
        <Card className={classes.root} style={{ textAlign: 'left' }}>
            <CardContent>
                <Typography className={classes.title} gutterBottom>
                    <FormattedMessage
                        id='Dashboard.rateLimiting.card.title'
                        defaultMessage='Rate Limiting'
                    />
                </Typography>

                <Divider light />
                <Box mt={1} mb={-2}>
                    {selectedPolicies.map((policy) => {
                        return (
                            <Box display='flex'>
                                <Box mx={1} mt={0.5}>
                                    {policy.icon}
                                </Box>
                                <Box flexGrow={1}>
                                    <Link component={RouterLink} to={policy.path} color='inherit'>
                                        <Typography
                                            variant='body1'
                                            style={{ fontWeight: 'bold' }}
                                        >
                                            {policy.name}
                                        </Typography>
                                    </Link>
                                    <Typography variant='body2' gutterBottom>
                                        {policy.description}
                                    </Typography>
                                </Box>
                            </Box>
                        );
                    })}
                </Box>
            </CardContent>
        </Card>
    );
}
