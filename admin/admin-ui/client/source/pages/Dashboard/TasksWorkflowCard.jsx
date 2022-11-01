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

import React, { useEffect, useState } from 'react';
import { FormattedMessage, useIntl } from 'react-intl';
import { Link as RouterLink } from 'react-router-dom';
import { Card } from '@mui/material';
import Avatar from '@mui/material/Avatar';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import CardContent from '@mui/material/CardContent';
import Divider from '@mui/material/Divider';
import Link from '@mui/material/Link';
import { makeStyles } from 'tss-react/mui';
import Typography from '@mui/material/Typography';
import DeviceHubIcon from '@ant-design/icons/GroupOutlined';
import DnsRoundedIcon from '@ant-design/icons/GroupOutlined';
import PeopleIcon from '@ant-design/icons/GroupOutlined';
import PermMediaOutlinedIcon from '@ant-design/icons/GroupOutlined';
import PublicIcon from '@ant-design/icons/GroupOutlined';
import SettingsEthernetIcon from '@ant-design/icons/GroupOutlined';
import Alert from 'components/Alert';
import moment from 'moment';
import useAxiosPromise from 'components/hooks/useAxiosPromise';

const useStyles = makeStyles()((theme) => {
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
        avatar: {
            width: theme.spacing(4),
            height: theme.spacing(4),
        },
        approveButton: {
            textDecoration: 'none',
            backgroundColor: theme.palette.success.light,
            margin: theme.spacing(0.5),
        },
        rejectButton: {
            textDecoration: 'none',
            backgroundColor: theme.palette.error.light,
            margin: theme.spacing(0.5),
        },
    }
});

/**
 * Render progress inside a container centering in the container.
 * @returns {JSX} Loading animation.
 */
export default function TasksWorkflowCard() {
    const {classes} = useStyles();
    const intl = useIntl();
    const [allTasksSet, setAllTasksSet] = useState({});

    /**
    * Calculate total task count
    * @returns {int} total task count
    */
    function getAllTaskCount() {
        let counter = 0;
        for (const task in allTasksSet) {
            if (allTasksSet[task]) {
                counter += allTasksSet[task].length;
            }
        }
        return counter;
    }

    // Fetch all workflow tasks
    const fetchAllWorkFlows = () => {
        const promiseUserSign = useAxiosPromise({url: 'workflows?workflowType=AM_USER_SIGNUP'});
        const promiseStateChange = useAxiosPromise({url: 'workflows?workflowType=AM_API_STATE'});
        const promiseApiProductStateChange = useAxiosPromise({url: 'workflows?workflowType=AM_API_PRODUCT_STATE'});
        const promiseAppCreation = useAxiosPromise({url: 'workflows?workflowType=AM_APPLICATION_CREATION'});
        const promiseAppDeletion = useAxiosPromise({url: 'workflows?workflowType=AM_APPLICATION_DELETION'});
        const promiseSubCreation = useAxiosPromise({url: 'workflows?workflowType=AM_SUBSCRIPTION_CREATION'});
        const promiseSubDeletion = useAxiosPromise({url: 'workflows?workflowType=AM_SUBSCRIPTION_DELETION'});
        const promiseSubUpdate = useAxiosPromise({url: 'workflows?workflowType=AM_SUBSCRIPTION_UPDATE'});
        const promiseRegProd = useAxiosPromise({url: 'workflows?workflowType=AM_APPLICATION_REGISTRATION_PRODUCTION'});
        const promiseRegSb = useAxiosPromise({url: 'workflows?workflowType=AM_APPLICATION_REGISTRATION_SANDBOX'});
        Promise.all([promiseUserSign, promiseStateChange, promiseAppCreation, promiseAppDeletion, promiseSubCreation,
            promiseSubDeletion, promiseSubUpdate, promiseRegProd, promiseRegSb, promiseApiProductStateChange])
            .then(([resultUserSign, resultStateChange, resultAppCreation, resultAppDeletion, resultSubCreation,
                resultSubDeletion, resultSubUpdate, resultRegProd, resultRegSb, resultApiProductStateChange]) => {
                const userCreation = resultUserSign.body.list;
                const stateChange = resultStateChange.body.list;
                const productStateChange = resultApiProductStateChange.body.list;
                const applicationCreation = resultAppCreation.body.list;
                const applicationDeletion = resultAppDeletion.body.list;
                const subscriptionCreation = resultSubCreation.body.list;
                const subscriptionDeletion = resultSubDeletion.body.list;
                const subscriptionUpdate = resultSubUpdate.body.list;
                const registration = resultRegProd.body.list.concat(resultRegSb.body.list);
                setAllTasksSet({
                    userCreation,
                    stateChange,
                    applicationCreation,
                    applicationDeletion,
                    subscriptionCreation,
                    subscriptionDeletion,
                    subscriptionUpdate,
                    registration,
                    productStateChange,
                });
            });
    };

    useEffect(() => {
        fetchAllWorkFlows();
    }, []);

    // Component to be displayed when there's no task available
    // Note: When workflow is not enabled, this will be displayed
    const noTasksCard = (
        <Card className={classes.root}>
            <CardContent>
                <Box mt={2}>
                    <DeviceHubIcon color='secondary' style={{ fontSize: 60 }} />
                </Box>

                <Typography className={classes.title} gutterBottom>
                    <FormattedMessage
                        id='Dashboard.tasksWorkflow.noTasks.card.title'
                        defaultMessage='All the pending tasks completed'
                    />
                </Typography>

                <Typography variant='body2' component='p'>
                    <FormattedMessage
                        id='Dashboard.tasksWorkflow.noTasks.card.description'
                        defaultMessage='Manage workflow tasks, increase productivity and enhance
                        competitiveness by enabling developers to easily deploy
                        business processes and models.'
                    />
                </Typography>
            </CardContent>
        </Card>
    );

    // Compact task card component's individual category component
    const getCompactTaskComponent = (IconComponent, path, name, numberOfTasks) => {
        return (
            <Box alignItems='center' display='flex' width='50%' my='1%'>
                <Box mx={1}>
                    <Avatar className={classes.avatar}>
                        <IconComponent fontSize='inherit' />
                    </Avatar>
                </Box>
                <Box flexGrow={1}>
                    <Link component={RouterLink} to={path} color='inherit'>
                        <Typography>
                            {name}
                        </Typography>
                    </Link>
                    <Typography variant='body2' gutterBottom>
                        {numberOfTasks + ' '}
                        {numberOfTasks === 1
                            ? (
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.compactTasks.card.numberOfPendingTasks.postFix.singular'
                                    defaultMessage=' Pending task'
                                />
                            ) : (
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.compactTasks.card.numberOfPendingTasks.postFix.plural'
                                    defaultMessage=' Pending tasks'
                                />
                            )}
                    </Typography>
                </Box>
            </Box>
        );
    };

    // Component to be displayed when there are more than 4 tasks available
    // Renders the total task count, each task category remaining task count and links
    const compactTasksCard = () => {
        const compactTaskComponentDetails = [
            {
                icon: PeopleIcon,
                path: '/tasks/user-creation',
                name: intl.formatMessage({
                    id: 'Dashboard.tasksWorkflow.compactTasks.userCreation.name',
                    defaultMessage: 'User Creation',
                }),
                count: allTasksSet.userCreation.length,
            },
            {
                icon: DnsRoundedIcon,
                path: '/tasks/application-creation',
                name: intl.formatMessage({
                    id: 'Dashboard.tasksWorkflow.compactTasks.applicationCreation.name',
                    defaultMessage: 'Application Creation',
                }),
                count: allTasksSet.applicationCreation.length,
            },
            {
                icon: DnsRoundedIcon,
                path: '/tasks/application-deletion',
                name: intl.formatMessage({
                    id: 'Dashboard.tasksWorkflow.compactTasks.applicationDeletion.name',
                    defaultMessage: 'Application Deletion',
                }),
                count: allTasksSet.applicationDeletion.length,
            },
            {
                icon: PermMediaOutlinedIcon,
                path: '/tasks/subscription-creation',
                name: intl.formatMessage({
                    id: 'Dashboard.tasksWorkflow.compactTasks.subscriptionCreation.name',
                    defaultMessage: 'Subscription Creation',
                }),
                count: allTasksSet.subscriptionCreation.length,
            },
            {
                icon: PermMediaOutlinedIcon,
                path: '/tasks/subscription-deletion',
                name: intl.formatMessage({
                    id: 'Dashboard.tasksWorkflow.compactTasks.subscriptionDeletion.name',
                    defaultMessage: 'Subscription Deletion',
                }),
                count: allTasksSet.subscriptionDeletion.length,
            },
            {
                icon: PermMediaOutlinedIcon,
                path: '/tasks/subscription-update',
                name: intl.formatMessage({
                    id: 'Dashboard.tasksWorkflow.compactTasks.subscriptionUpdate.name',
                    defaultMessage: 'Subscription Update',
                }),
                count: allTasksSet.subscriptionUpdate.length,
            },
            {
                icon: PublicIcon,
                path: '/tasks/application-registration',
                name: intl.formatMessage({
                    id: 'Dashboard.tasksWorkflow.compactTasks.applicationRegistration.name',
                    defaultMessage: 'Application Registration',
                }),
                count: allTasksSet.registration.length,
            },
            {
                icon: SettingsEthernetIcon,
                path: '/tasks/api-state-change',
                name: intl.formatMessage({
                    id: 'Dashboard.tasksWorkflow.compactTasks.apiStateChange.name',
                    defaultMessage: 'API State Change',
                }),
                count: allTasksSet.stateChange.length,
            },
            {
                icon: SettingsEthernetIcon,
                path: '/tasks/api-product-state-change',
                name: intl.formatMessage({
                    id: 'Dashboard.tasksWorkflow.compactTasks.apiProductStateChange.name',
                    defaultMessage: 'API Product State Change',
                }),
                count: allTasksSet.productStateChange.length,
            },
        ];
        return (
            <Card className={classes.root} style={{ textAlign: 'left' }}>
                <CardContent>
                    <Box display='flex'>
                        <Box flexGrow={1}>
                            <Typography className={classes.title} gutterBottom>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.compactTasks.card.title'
                                    defaultMessage='Pending tasks'
                                />
                            </Typography>
                        </Box>
                        <Box>
                            <Typography className={classes.title} gutterBottom>
                                {getAllTaskCount()}
                            </Typography>
                        </Box>
                    </Box>

                    <Divider light />

                    <Box
                        display='flex'
                        flexWrap='wrap'
                        mt={2}
                        bgcolor='background.paper'

                    >
                        {compactTaskComponentDetails.map((c) => {
                            return getCompactTaskComponent(c.icon, c.path, c.name, c.count);
                        })}
                    </Box>
                </CardContent>
            </Card>
        );
    };

    // Approve/Reject button onClick handler
    const updateStatus = (referenceId, value) => {
        const body = {
            status: value,
        };
        restApi.updateWorkflow(referenceId, body)
            .then(() => {
                Alert.success(
                    <FormattedMessage
                        id='Dashboard.tasksWorkflow.card.task.update.success'
                        defaultMessage='Task status updated successfully'
                    />,
                );
            })
            .catch(() => {
                Alert.error(
                    <FormattedMessage
                        id='Dashboard.tasksWorkflow.card.task.update.failed'
                        defaultMessage='Task status updated failed'
                    />,
                );
            })
            .finally(() => {
                fetchAllWorkFlows();
            });
    };

    // Renders the approve/reject buttons with styles
    const getApproveRejectButtons = (referenceId) => {
        return (
            <Box>
                <Button
                    onClick={() => { updateStatus(referenceId, 'APPROVED'); }}
                    className={classes.approveButton}
                    variant='contained'
                    size='small'
                >
                    <FormattedMessage
                        id='Dashboard.tasksWorkflow.fewerTasks.card.task.accept'
                        defaultMessage='Accept'
                    />
                </Button>
                <Button
                    onClick={() => { updateStatus(referenceId, 'REJECTED'); }}
                    className={classes.rejectButton}
                    variant='contained'
                    size='small'
                >
                    <FormattedMessage
                        id='Dashboard.tasksWorkflow.fewerTasks.card.task.reject'
                        defaultMessage='Reject'
                    />
                </Button>
            </Box>
        );
    };

    // Fewer task component's application creation task element
    const getApplicationCreationFewerTaskComponent = () => {
        // Application Creation tasks related component generation
        return allTasksSet.applicationCreation.map((task) => {
            return (
                <Box display='flex' alignItems='center' mt={1}>
                    <Box flexGrow={1}>
                        <Typography variant='subtitle2'>
                            {task.properties.applicationName}
                        </Typography>
                        <Box display='flex'>
                            <Typography variant='body2'>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.application.createdBy.prefix'
                                    defaultMessage='Application Created by '
                                />
                            </Typography>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                &nbsp;
                                {task.properties.userName}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                {moment(task.createdTime).fromNow()}
                            </Typography>
                        </Box>
                    </Box>
                    {getApproveRejectButtons(task.referenceId)}
                </Box>
            );
        });
    };

    const getApplicationDeletionFewerTaskComponent = () => {
        // Application Creation tasks related component generation
        return allTasksSet.applicationDeletion.map((task) => {
            return (
                <Box display='flex' alignItems='center' mt={1}>
                    <Box flexGrow={1}>
                        <Typography variant='subtitle2'>
                            {task.properties.applicationName}
                        </Typography>
                        <Box display='flex'>
                            <Typography variant='body2'>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.application.deletedBy.prefix'
                                    defaultMessage='Application Deleted by '
                                />
                            </Typography>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                &nbsp;
                                {task.properties.userName}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                {moment(task.createdTime).fromNow()}
                            </Typography>
                        </Box>
                    </Box>
                    {getApproveRejectButtons(task.referenceId)}
                </Box>
            );
        });
    };

    // Fewer task component's user creation task element
    const getUserCreationFewerTaskComponent = () => {
        // User Creation tasks related component generation
        return allTasksSet.userCreation.map((task) => {
            return (
                <Box display='flex' alignItems='center' mt={1}>
                    <Box flexGrow={1}>
                        <Typography variant='subtitle2'>
                            {task.properties.tenantAwareUserName}
                        </Typography>
                        <Box display='flex'>
                            <Typography variant='body2'>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.user.createdOn.prefix'
                                    defaultMessage='User Created on '
                                />
                            </Typography>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                &nbsp;
                                {task.properties.tenantDomain}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                {moment(task.createdTime).fromNow()}
                            </Typography>
                        </Box>
                    </Box>
                    {getApproveRejectButtons(task.referenceId)}
                </Box>
            );
        });
    };

    // Fewer task component's subscription creation task element
    const getSubscriptionCreationFewerTaskComponent = () => {
        // Subscription Creation tasks related component generation
        return allTasksSet.subscriptionCreation.map((task) => {
            return (
                <Box display='flex' alignItems='center' mt={1}>
                    <Box flexGrow={1}>
                        <Typography variant='subtitle2'>
                            {task.properties.apiName + '-' + task.properties.apiVersion}
                        </Typography>
                        <Box display='flex'>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                {task.properties.applicationName + ','}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.subscription.subscribedBy'
                                    defaultMessage='Subscribed by'
                                />
                            </Typography>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                &nbsp;
                                {task.properties.subscriber}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                {moment(task.createdTime).fromNow()}
                            </Typography>
                        </Box>
                    </Box>
                    {getApproveRejectButtons(task.referenceId)}
                </Box>
            );
        });
    };

    // Fewer task component's subscription creation task element
    const getSubscriptionDeletionFewerTaskComponent = () => {
        // Subscription Update tasks related component generation
        return allTasksSet.subscriptionDeletion.map((task) => {
            return (
                <Box display='flex' alignItems='center' mt={1}>
                    <Box flexGrow={1}>
                        <Typography variant='subtitle2'>
                            {task.properties.apiName + '-' + task.properties.apiVersion}
                        </Typography>
                        <Box display='flex'>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                {task.properties.applicationName + ','}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.subscription.deletedBy'
                                    defaultMessage='Subscription Deleted by'
                                />
                            </Typography>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                &nbsp;
                                {task.properties.subscriber}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                {moment(task.createdTime).fromNow()}
                            </Typography>
                        </Box>
                    </Box>
                    {getApproveRejectButtons(task.referenceId)}
                </Box>
            );
        });
    };

    // Fewer task component's subscription creation task element
    const getSubscriptionUpdateFewerTaskComponent = () => {
        // Subscription Update tasks related component generation
        return allTasksSet.subscriptionUpdate.map((task) => {
            return (
                <Box display='flex' alignItems='center' mt={1}>
                    <Box flexGrow={1}>
                        <Typography variant='subtitle2'>
                            {task.properties.apiName + '-' + task.properties.apiVersion}
                        </Typography>
                        <Box display='flex'>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                {task.properties.applicationName + ','}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.subscription.subscribedBy'
                                    defaultMessage='Subscribed by'
                                />
                            </Typography>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                &nbsp;
                                {task.properties.subscriber}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                {moment(task.createdTime).fromNow()}
                            </Typography>
                        </Box>
                    </Box>
                    {getApproveRejectButtons(task.referenceId)}
                </Box>
            );
        });
    };

    // Fewer task component's registration creation task element
    const getRegistrationCreationFewerTaskComponent = () => {
        // Registration Creation tasks related component generation
        return allTasksSet.registration.map((task) => {
            let keyType;
            if (task.properties.keyType === 'PRODUCTION') {
                keyType = (
                    <FormattedMessage
                        id='Dashboard.tasksWorkflow.fewerTasks.card.registration.creation.keyType.Production'
                        defaultMessage='Production'
                    />
                );
            } else if (task.properties.keyType === 'SANDBOX') {
                keyType = (
                    <FormattedMessage
                        id='Dashboard.tasksWorkflow.fewerTasks.card.registration.creation.keyType.SandBox'
                        defaultMessage='SandBox'
                    />
                );
            } else {
                keyType = task.properties.keyType;
            }
            return (
                <Box display='flex' alignItems='center' mt={1}>
                    <Box flexGrow={1}>
                        <Typography variant='subtitle2'>
                            {task.properties.applicationName}
                        </Typography>
                        <Box display='flex'>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                {keyType}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.registration.key.generated.by'
                                    defaultMessage='Key generated by'
                                />
                            </Typography>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                &nbsp;
                                {task.properties.userName}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                {moment(task.createdTime).fromNow()}
                            </Typography>
                        </Box>
                    </Box>
                    {getApproveRejectButtons(task.referenceId)}
                </Box>
            );
        });
    };

    // Fewer task component's api state change task element
    const getStateChangeFewerTaskComponent = () => {
        // State Change tasks related component generation
        return allTasksSet.stateChange.map((task) => {
            return (
                <Box display='flex' alignItems='center' mt={1}>
                    <Box flexGrow={1}>
                        <Typography variant='subtitle2'>
                            {task.properties.apiName + '-' + task.properties.apiVersion}
                        </Typography>
                        <Box display='flex'>
                            <Typography variant='body2'>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.stateChangeAction.prefix'
                                    defaultMessage='State Change Action:'
                                />
                                &nbsp;
                            </Typography>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                {task.properties.action}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                {moment(task.createdTime).fromNow()}
                            </Typography>
                        </Box>
                    </Box>
                    {getApproveRejectButtons(task.referenceId)}
                </Box>
            );
        });
    };

    const getAPIProductStateChangeFewerTaskComponent = () => {
        // State Change tasks related component generation
        return allTasksSet.productStateChange.map((task) => {
            return (
                <Box display='flex' alignItems='center' mt={1}>
                    <Box flexGrow={1}>
                        <Typography variant='subtitle2'>
                            {task.properties.apiName}
                        </Typography>
                        <Box display='flex'>
                            <Typography variant='body2'>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.apiProduct.stateChangeAction.prefix'
                                    defaultMessage='State Change Action:'
                                />
                                &nbsp;
                            </Typography>
                            <Typography style={{ 'font-weight': 'bold' }} variant='body2'>
                                {task.properties.action}
                                &nbsp;
                            </Typography>
                            <Typography variant='body2'>
                                {moment(task.createdTime).fromNow()}
                            </Typography>
                        </Box>
                    </Box>
                    {getApproveRejectButtons(task.referenceId)}
                </Box>
            );
        });
    };

    // Component to be displayed when there are 4 or less remaining tasks
    // Renders some details of the task and approve/reject buttons
    const fewerTasksCard = () => {
        return (
            <Card className={classes.root} style={{ textAlign: 'left' }}>
                <CardContent>
                    <Box display='flex'>
                        <Box flexGrow={1}>
                            <Typography className={classes.title} gutterBottom>
                                <FormattedMessage
                                    id='Dashboard.tasksWorkflow.fewerTasks.card.title'
                                    defaultMessage='Pending tasks'
                                />
                            </Typography>
                        </Box>
                        <Box>
                            <Typography className={classes.title} gutterBottom>
                                {getAllTaskCount()}
                            </Typography>
                        </Box>
                    </Box>

                    <Divider light />
                    {getApplicationCreationFewerTaskComponent()}
                    {getApplicationDeletionFewerTaskComponent()}
                    {getUserCreationFewerTaskComponent()}
                    {getSubscriptionCreationFewerTaskComponent()}
                    {getSubscriptionDeletionFewerTaskComponent()}
                    {getSubscriptionUpdateFewerTaskComponent()}
                    {getRegistrationCreationFewerTaskComponent()}
                    {getStateChangeFewerTaskComponent()}
                    {getAPIProductStateChangeFewerTaskComponent()}
                </CardContent>
            </Card>
        );
    };

    // Render the card depending on the number of all remaining tasks
    const cnt = getAllTaskCount();
    if (cnt > 4) {
        return compactTasksCard();
    } else if (cnt > 0) {
        return fewerTasksCard();
    } else {
        return noTasksCard;
    }
}
