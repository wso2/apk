/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import DialogContentText from '@mui/material/DialogContentText';
import axios from 'axios';
import PropTypes from 'prop-types';
import React from 'react';
import { FormattedMessage } from 'react-intl';
import FormDialogBase from '../../components/@extended/FormDialogBase';

export default function DeleteAPICategory({ id, updateList }) {

    const formSaveCallback = () => {
        axios.delete('/api/am/admin/api-categories/' + id, {
            withCredentials: true,
        }).then(() => {
            return (
                <FormattedMessage
                    id='AdminPages.ApiCategories.Delete.form.delete.successful'
                    defaultMessage='API Category deleted successfully.'
                />
            );
        }).catch((error) => {
            throw error.response.body.description;
        }).finally(() => {
            updateList();
        });
    };

    return (
        <FormDialogBase
            title='Delete API category?'
            saveButtonText='Delete'
            icon={<DeleteForeverIcon />}
            formSaveCallback={formSaveCallback}
        >
            <DialogContentText>Are you sure you want to delete this API Category?</DialogContentText>
        </FormDialogBase>
    );
}

DeleteAPICategory.propTypes = {
    id: PropTypes.string.isRequired,
    updateList: PropTypes.func.isRequired
};
