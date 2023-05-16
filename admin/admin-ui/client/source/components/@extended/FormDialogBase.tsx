import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import IconButton from '@mui/material/IconButton';
import Alert from "components/Alert";
import PropTypes from 'prop-types';
import React, { useState } from 'react';

export default function FormDialogBase({
    title,
    children,
    icon,
    saveButtonText,
    formSaveCallback,
    dialogOpenCallback,
    triggerIconProps,
}) {
    const [open, setOpen] = useState<boolean>(false);
    const [saving, setSaving] = useState<boolean>(false);

    const handleClickOpen = () => {
        dialogOpenCallback();
        setOpen(true);
    };

    const handleClose = () => {
        setOpen(false);
    };

    const saveTriggerd = () => {
        const savedPromise = formSaveCallback();
        if (typeof savedPromise === 'function') {
            savedPromise(setOpen);
        } else if (savedPromise) {
            setSaving(true);
            savedPromise.then((data) => {
                Alert.success(data);
            }).catch((e) => {
                Alert.error(e);
            }).finally(() => {
                setSaving(false);
                handleClose();
            });
        }
    };

    return (
        <>
            {icon && (
                <IconButton {...triggerIconProps} onClick={handleClickOpen}>
                    {icon}
                </IconButton>
            )}

            <Dialog open={open} onClose={handleClose} aria-labelledby='form-dialog-title'>
                <DialogTitle id='form-dialog-title'>{title}</DialogTitle>
                <DialogContent>
                    {children}
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>
                        Cancel
                    </Button>
                    <Button
                        onClick={saveTriggerd}
                        color='primary'
                        variant='contained'
                        disabled={saving}
                        data-testid={saveButtonText + '-btn'}
                    >
                        {saving ? (<CircularProgress size={16} />) : (<>{saveButtonText}</>)}
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
}
FormDialogBase.defaultProps = {
    dialogOpenCallback: () => { },
    triggerButtonProps: {
        variant: 'contained',
        color: 'primary',
    },
    triggerIconProps: {
        color: 'primary',
        component: 'span',
    },
};

FormDialogBase.propTypes = {
    title: PropTypes.string.isRequired,
    children: PropTypes.element,
    icon: PropTypes.element.isRequired,
    saveButtonText: PropTypes.string.isRequired,
    dialogOpenCallback: PropTypes.func,
};
