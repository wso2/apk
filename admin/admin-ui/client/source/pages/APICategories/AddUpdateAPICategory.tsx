import AddCircleOutlineRoundedIcon from '@mui/icons-material/AddCircleOutlineRounded';
import EditIcon from '@mui/icons-material/Edit';
import TextField from '@mui/material/TextField';
import axios from 'axios';
import PropTypes from 'prop-types';
import React, { useEffect, useState } from 'react';
import { FormattedMessage } from 'react-intl';
import FormDialogBase from '../../components/@extended/FormDialogBase';

export default function AddUpdateAPICategory({ id, nameProp, descriptionProp, updateList }) {
  // This component has been used to add API category when id is undefined and edit API category when id is defined
  const [APICategory, setAPICategory] = useState<{name: string; description: string}>({ name: "", description: "" });

  useEffect(() => {
    if (id !== undefined) {
      setAPICategory({ name: nameProp, description: descriptionProp });
    }
  }, []);

  const formSaveCallback = () => {
    if (id !== undefined) {
      axios.put('/api/am/admin/api-categories/' + id, { 'name': APICategory.name, 'description': APICategory.description }, {
        withCredentials: true,
      }).then(() => {
        return (
          <FormattedMessage
            id='AdminPages.ApiCategories.AddEdit.form.edit.successful'
            defaultMessage='API Category edited successfully.'
          />
        );
      }).catch((error) => {
        throw error.response.body.description;
      }).finally(() => {
        updateList();
      });
    } else {
      axios.post('/api/am/admin/api-categories/', { 'name': APICategory.name, 'description': APICategory.description }, {
        withCredentials: true,
      }).then(() => {
        return (
          <FormattedMessage
            id='AdminPages.ApiCategories.AddEdit.form.add.successful'
            defaultMessage='API Category added successfully.'
          />
        );
      }).catch((error) => {
        throw error.response.body.description;
      }).finally(() => {
        updateList();
      });
    }
  };

  const onChange = (e) => {
    setAPICategory({ ...APICategory, [e.target.name]: e.target.value });
  };

  return (
    <>
      <FormDialogBase
        title={id !== undefined ? 'Edit API Category' : 'Add API Category'}
        saveButtonText={id !== undefined ? 'Save' : 'Add'}
        icon={id !== undefined ? <EditIcon /> : <AddCircleOutlineRoundedIcon fontSize='large' />}
        formSaveCallback={formSaveCallback}
      >
        <>
          <TextField
            autoFocus
            margin='dense'
            name='name'
            value={APICategory.name}
            onChange={onChange}
            fullWidth
            variant='outlined'
            disabled={id !== undefined}
            label={(
              <span>
                <FormattedMessage id='AdminPages.ApiCategories.AddEdit.form.name' defaultMessage='Name' />
                <span>*</span>
              </span>
            )}
            helperText={'Name of the API category'}
          />
          <TextField
            margin='dense'
            name='description'
            value={APICategory.description}
            onChange={onChange}
            fullWidth
            variant='outlined'
            multiline
            label={(
              <span>
                <FormattedMessage id='AdminPages.ApiCategories.AddEdit.form.description' defaultMessage='Description' />
                <span>*</span>
              </span>
            )}
            helperText={'Description of the API category'}
          />
        </>
      </FormDialogBase>
    </>
  );
}

AddUpdateAPICategory.propTypes = {
  id: PropTypes.string,
  nameProp: PropTypes.string,
  descriptionProp: PropTypes.string,
  updateList: PropTypes.func.isRequired
};
