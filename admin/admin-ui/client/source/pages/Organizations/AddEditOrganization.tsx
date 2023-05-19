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

import React, { useEffect, useState } from 'react';
import { FormattedMessage } from 'react-intl';
import { styled } from '@mui/material/styles';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import IconButton from '@mui/material/IconButton';
import TextField from '@mui/material/TextField';

import EditIcon from '@mui/icons-material/Edit';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';

import { MuiChipsInput } from 'mui-chips-input';

interface OrganizationData {
    id: string
    name: string;
    displayName: string;
    organizationClaimValue: string;
    enabled: boolean;
}

interface Props {
    datarow: any;
    updateList: (organization) => void;
}

/**
 * Render add & edit dialog boxes.
 * @param {JSON} props component props.
 * @returns {TSX} Loading animation.
 */
const AddEditOrganization: React.FC<Props> = ({
    datarow,
    updateList,
}) => {
    // const emptyDatarow = Object.keys(datarow).length === 0
    const [open, setOpen] = useState(false);
    const [saving, setSaving] = useState(false);
    const [organization, setOrganization] = useState<OrganizationData>({
        id: datarow ? datarow.id : '',
        name: datarow ? datarow.name : '',
        displayName: datarow ? datarow.displayName : '',
        organizationClaimValue: datarow ? datarow.organizationClaimValue : '',
        enabled: datarow ? datarow.enabled : true,
    });
    const [serviceNamespaces, setServiceNamespaces] = useState({ serviceNamespaces: datarow ? datarow.serviceNamespaces : [] });
    const [production, setProduction] = useState({ production: datarow ? datarow.production : [] });
    const [sandbox, setSandbox] = useState({ sandbox: datarow ? datarow.sandbox : [] });
    const [dialogTitle, setDialogTitle] = useState('Edit Organization');
    const [errors, setErrors] = useState<{
        name?: string;
        claimValue?: string;
        displayName?: string;
    }>({});

    const handleClickOpen = () => {
        setOpen(true);
    }

    const handleClose = () => {
        setOpen(false);
    };

    const handleSave = () => {
        const newErrors = validateOrganizationForm();
        setSaving(true);
        if (Object.keys(newErrors).length === 0) {
            updateList({ ...organization, ...serviceNamespaces, ...production, ...sandbox });
            setSaving(false);
        } else {
            setErrors(newErrors);
            setSaving(false);
            return false;
        }

        handleClose();
        return true;
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setOrganization({ ...organization, [e.target.name]: e.target.value });
        if (e.target.name === 'name') {
            setErrors((prevErrors) => ({ ...prevErrors, name: '' }));
        }
        if (e.target.name === 'displayName') {
            setErrors((prevErrors) => ({ ...prevErrors, displayName: '' }));
        }
        if (e.target.name === 'organizationClaimValue') {
            setErrors((prevErrors) => ({ ...prevErrors, claimValue: '' }));
        }
    };

    const handleSeriveNamespacesChipChange = (newValue) => {
        setServiceNamespaces({ ...serviceNamespaces, serviceNamespaces: newValue });
    }

    const handleProductionChipChange = (newValue) => {
        setProduction({ ...production, production: newValue })
    }

    const handleSandboxChipChange = (newValue) => {
        setSandbox({ ...sandbox, sandbox: newValue })
    }

    const { name: orgName, displayName: orgDisplayName, organizationClaimValue: orgClaimValue } = organization;
    const { serviceNamespaces: orgServiceNamespaces } = serviceNamespaces;
    const { production: orgProduction } = production;
    const { sandbox: orgSandbox } = sandbox;

    const validateOrganizationForm = (): { name?: string; orgClaimValue?: string } => {
        let newErrors: { name?: string; claimValue?: string; displayName?: string } = {};

        if (!orgName) {
            newErrors.name = 'Organization name should not be empty';
        }
        if (!orgClaimValue) {
            newErrors.claimValue = 'Organization claim value should not be empty';
        }
        if (!orgDisplayName) {
            newErrors.displayName = 'Organization display name should not be empty';
        }

        return newErrors;
    };

    const StyledChipsInput = styled(MuiChipsInput)(({ }) => ({
        '& .MuiChipsInput-Chip': {
            borderRadius: '25px',
        }
    }));

    useEffect(() => {
        if (!datarow) {
            setOrganization({
                id: '',
                name: '',
                displayName: '',
                organizationClaimValue: '',
                enabled: true,
            })
            setServiceNamespaces({ serviceNamespaces: [] });
            setProduction({ production: [] })
            setSandbox({ sandbox: [] })
            setDialogTitle('Add Organization');
        }
        setErrors((prevErrors) => ({ ...prevErrors, name: '' }));
        setErrors((prevErrors) => ({ ...prevErrors, displayName: '' }));
        setErrors((prevErrors) => ({ ...prevErrors, claimValue: '' }));
    }, [open]);

    return (
        <>
            {datarow && (
                <IconButton onClick={handleClickOpen} color='primary' size='small'>
                    <EditIcon fontSize='small' />
                </IconButton>
            )}

            {!datarow && (
                <Button onClick={handleClickOpen} variant='contained' style={{ margin: 'auto' }} startIcon={<AddCircleOutlineIcon />}>Add Organization</Button>
            )}

            <Dialog open={open} onClose={handleClose} aria-labelledby='form-dialog-title'>
                <DialogTitle id='form-dialog-title'>{dialogTitle}</DialogTitle>
                <DialogContent>
                    <>
                        <TextField
                            autoFocus
                            margin='dense'
                            name='name'
                            value={orgName}
                            onChange={handleInputChange}
                            fullWidth
                            variant='outlined'
                            label={(
                                <span>
                                    <FormattedMessage id='AdminPages.Organizations.AddEdit.form.name' defaultMessage='Name' />
                                    <span style={{ color: 'red' }}>*</span>
                                </span>
                            )}
                            helperText={errors?.name ? errors.name : 'Name of the Organization'}
                            error={!!errors?.name}
                            style={{ marginTop: '15px' }}
                        />
                        <TextField
                            margin='dense'
                            name='displayName'
                            value={orgDisplayName}
                            onChange={handleInputChange}
                            fullWidth
                            multiline
                            label={(
                                <span>
                                    <FormattedMessage id='AdminPages.Organizations.AddEdit.form.displayName' defaultMessage='Display Name' />
                                    <span style={{ color: 'red' }}>*</span>
                                </span>
                            )}
                            helperText={errors?.displayName ? errors.displayName : 'Display Name of the Organization'}
                            error={!!errors?.displayName}
                            variant='outlined'
                            style={{ marginTop: '15px' }}
                        />
                        <TextField
                            margin='dense'
                            name='organizationClaimValue'
                            value={orgClaimValue}
                            onChange={handleInputChange}
                            fullWidth
                            multiline
                            label={(
                                <span>
                                    <FormattedMessage id='AdminPages.Organizations.AddEdit.form.organizationClaimValue' defaultMessage='Organization Claim Value' />
                                    <span style={{ color: 'red' }}>*</span>
                                </span>
                            )}
                            helperText={errors?.claimValue ? errors.claimValue : 'Claim Value of the Organization'}
                            error={!!errors?.claimValue}
                            variant='outlined'
                            style={{ marginTop: '15px' }}
                        />
                        <StyledChipsInput
                            fullWidth
                            variant='outlined'
                            value={orgServiceNamespaces}
                            onChange={handleSeriveNamespacesChipChange}
                            hideClearAll
                            helperText={'Service Namespaces of the Organization'}
                            style={{ marginTop: '15px' }}
                        />
                        <StyledChipsInput
                            fullWidth
                            variant='outlined'
                            value={orgProduction}
                            onChange={handleProductionChipChange}
                            hideClearAll
                            helperText={'Production Endpoints of the Organization'}
                            style={{ marginTop: '15px' }}
                        />
                        <StyledChipsInput
                            fullWidth
                            variant='outlined'
                            value={orgSandbox}
                            onChange={handleSandboxChipChange}
                            hideClearAll
                            helperText={'Sandbox Endpoints of the Organization'}
                            style={{ marginTop: '15px' }}
                        />
                    </>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose} >
                        Cancel
                    </Button>
                    <Button
                        onClick={handleSave}
                        color='primary'
                        variant='contained'
                        disabled={saving}
                    >
                        {saving ? (<CircularProgress size={16} />) : (<>{'Save'}</>)}
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
}
export default AddEditOrganization;
