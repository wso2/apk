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
