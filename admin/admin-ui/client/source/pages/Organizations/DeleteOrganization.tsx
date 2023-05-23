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

/**
 * Render delete dialog box.
 * @param {JSON} props component props.
 * @returns {TSX} Loading animation.
 */
import React, { useState } from 'react';
import DialogContentText from '@mui/material/DialogContentText';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import axios from "axios";
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import IconButton from '@mui/material/IconButton';
import CircularProgress from '@mui/material/CircularProgress';

interface DeleteProps {
    orgId: string;
    updateList: () => void;
}

const DeleteOrganization: React.FC<DeleteProps> = ({ orgId, updateList }) => {
    const [open, setOpen] = useState(false);
    const [saving, setSaving] = useState(false);

    const handleClickOpen = () => {
        setOpen(true);
    };

    const handleClose = () => {
        setOpen(false);
    };

    const saveTriggered = () => {
        setSaving(true);
        deleteOrganization();
        handleClose();
    };

    const deleteOrganization = () => {
        axios
            .delete(`/api/am/admin/organizations/${orgId}`, {
                withCredentials: true,
            })
            .catch((error) => {
                throw error.response.body.description;
            })
            .finally(() => {
                updateList();
                setSaving(false);
            });
    };

    

    return (
        <>
            <IconButton onClick={handleClickOpen} color='primary' size='small'>
                <DeleteForeverIcon fontSize='small' />
            </IconButton>

            <Dialog open={open} onClose={handleClose} aria-labelledby='form-dialog-title'>
                <DialogTitle id='form-dialog-title'>Delete Organization?</DialogTitle>
                <DialogContent>
                    <DialogContentText>Are you sure you want to delete this Organization?</DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>Cancel</Button>
                    <Button onClick={saveTriggered} color='primary' variant='contained' disabled={saving}>
                        {saving ? <CircularProgress size={16} /> : <>Delete</>}
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
};

export default DeleteOrganization;
