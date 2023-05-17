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

import AddCircleOutlineRoundedIcon from '@mui/icons-material/AddCircleOutlineRounded';
import EditIcon from '@mui/icons-material/Edit';
import TextField from '@mui/material/TextField';
import axios from 'axios';
import Alert from "components/Alert";
import PropTypes from 'prop-types';
import React, { useEffect, useState } from 'react';
import { FormattedMessage } from 'react-intl';
import FormDialogBase from '../../components/@extended/FormDialogBase';

export default function AddUpdateAPICategory({ id, nameProp, descriptionProp, updateList }) {
  // This component has been used to add API category when id is undefined and edit API category when id is defined
  const [APICategory, setAPICategory] = useState<{ name: string; description: string }>({ name: "", description: "" });

  useEffect(() => {
    if (id !== undefined) {
      setAPICategory({ name: nameProp, description: descriptionProp });
    }
  }, []);

  const hasErrors = (fieldName, value) => {
    let error;
    switch (fieldName) {
      case 'name':
        if (value === undefined) {
          error = false;
          break;
        }
        if (value === '') {
          error = 'Name is Empty';
        } else if (value.length > 255) {
          error = 'API Category name is too long';
        } else if (/\s/.test(value)) {
          error = 'Name contains spaces';
        } else if (/[!@#$%^&*(),?"{}[\]|<>\t\n]/i.test(value)) {
          error = 'Name field contains special characters';
        } else {
          error = false;
        }
        break;
      case 'description':
        if (value && value.length > 1024) {
          error = 'API Category description is too long';
        }
        break;
      default:
        break;
    }
    return error;
  };

  const getAllFormErrors = () => {
    let errorText = '';
    let NameErrors;
    let DescriptionErrors;

    if (APICategory.name === undefined) {
      NameErrors = hasErrors('name', '');
    } else {
      NameErrors = hasErrors('name', APICategory.name);
    }
    if (NameErrors) {
      errorText += NameErrors + '\n';
    }
    if (APICategory.description !== undefined) {
      DescriptionErrors = hasErrors('description', APICategory.description);
    }
    if (DescriptionErrors) {
      errorText += DescriptionErrors + '\n';
    }
    return errorText;
  };

  const formSaveCallback = () => {
    const formErrors = getAllFormErrors();
    if (formErrors !== '') {
      console.log(formErrors);
      Alert.error(formErrors);
      return false;
    } else {
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
